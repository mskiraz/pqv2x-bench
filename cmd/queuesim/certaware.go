package main

// Certificate-aware C-V2X arrival model ("cv2x-cert").
//
// The legacy "cv2x" model gives every sender one of 100 x m message slots and
// lets every message, digest or certificate-carrying, use that slot. Under
// the SAE J3161 configuration analysed in the paper, however, a Falcon-512
// digest-form message occupies 4 of the 10 subchannels of a 1-ms subframe
// and a certificate-carrying message occupies all 10. This file adds an
// arrival model that respects those radio costs.
//
// Scheduling policy (idealised, deterministic, centralised; it is NOT an
// implementation of the 3GPP Mode 4 / SAE J3161 MAC):
//
//   - Senders, their preferred subframes, their certificate offsets and the
//     tie order within a subframe are drawn exactly as in the legacy model,
//     from a random stream seeded identically, so both models see the same
//     random configurations.
//   - Each sender generates one message per 100-ms period at the start of its
//     preferred subframe; one message in five carries the certificate.
//   - Messages are placed in generation order (ties in the shuffled sender
//     order) into the earliest subframe at or after their generation subframe
//     that has room: a digest-form message needs digestSub free subchannels,
//     and a certificate-carrying message needs certSub free subchannels.
//     With the defaults (4 and 10 of 10) two digest-form messages can share
//     a subframe, and a certificate-carrying message needs an empty subframe,
//     which it then fills completely.
//   - A message that finds no room within window subframes of its generation
//     subframe is not scheduled. It is counted, never silently delivered.
//   - With window <= 100 subframes a message is placed before its sender's
//     next message is generated, so per-sender order is preserved.
//   - A scheduled message reaches the receiver at the end of its subframe.
//     There is no propagation loss, no collision, no half-duplex loss and no
//     retransmission; every sender is in range.

import (
	"math"
	"math/rand/v2"
	"sort"
)

// radioParams describes the certificate-aware subchannel model.
type radioParams struct {
	subchannels int // per subframe
	digestSub   int // subchannels of a digest-form message
	certSub     int // subchannels of a certificate-carrying message
	window      int // scheduling window in subframes
}

var defaultRadio = radioParams{subchannels: 10, digestSub: 4, certSub: 10, window: 100}

// genMsg is one generated message.
type genMsg struct {
	sender int
	period int
	gen    int // generation subframe (absolute index)
	cert   bool
}

// nextFree is a union-find "next usable index" structure: find(i) returns
// the smallest j >= i that has not been closed.
type nextFree []int32

func newNextFree(n int) nextFree {
	f := make(nextFree, n+1)
	for i := range f {
		f[i] = int32(i)
	}
	return f
}

func (f nextFree) reset() {
	for i := range f {
		f[i] = int32(i)
	}
}

func (f nextFree) find(i int) int {
	r := i
	for int(f[r]) != r {
		r = int(f[r])
	}
	for int(f[i]) != r { // path compression
		n := int(f[i])
		f[i] = int32(r)
		i = n
	}
	return r
}

func (f nextFree) close(i int) { f[i] = int32(i + 1) }

// scheduler places messages into subframes under radioParams.
type scheduler struct {
	rp      radioParams
	nSub    int
	occ     []uint8 // allocated subchannels per subframe
	hasCert []bool
	slots   []int32 // message indices per subframe, maxPer each
	count   []uint8
	maxPer  int
	dFree   nextFree // subframes that can still take a digest message
	cFree   nextFree // subframes that can still take a certificate message
}

func newScheduler(rp radioParams, nSub int) *scheduler {
	maxPer := rp.subchannels / rp.digestSub
	if m := rp.subchannels / rp.certSub; m > maxPer {
		maxPer = m
	}
	return &scheduler{
		rp: rp, nSub: nSub, maxPer: maxPer,
		occ: make([]uint8, nSub), hasCert: make([]bool, nSub),
		slots: make([]int32, nSub*maxPer), count: make([]uint8, nSub),
		dFree: newNextFree(nSub), cFree: newNextFree(nSub),
	}
}

func (s *scheduler) reset() {
	clear(s.occ)
	clear(s.hasCert)
	clear(s.count)
	s.dFree.reset()
	s.cFree.reset()
}

// fits reports whether a message of the given kind can join subframe k.
// Only capacity is checked: with certSub equal to the number of subchannels
// a certificate-carrying message needs an empty subframe and leaves no room
// for anything else, so it never shares its subframe.
func (s *scheduler) fits(k int, cert bool) bool {
	need := s.rp.digestSub
	if cert {
		need = s.rp.certSub
	}
	return int(s.occ[k])+need <= s.rp.subchannels
}

// place schedules message id generated in subframe gen and returns the
// subframe, or -1 if it does not fit within the window or the horizon.
func (s *scheduler) place(id int, gen int, cert bool) int {
	free := s.dFree
	if cert {
		free = s.cFree
	}
	k := free.find(gen)
	for k < s.nSub && k < gen+s.rp.window && !s.fits(k, cert) {
		free.close(k) // stale entry: subframe filled by the other kind
		k = free.find(k)
	}
	if k >= s.nSub || k >= gen+s.rp.window {
		return -1
	}
	if cert {
		s.occ[k] += uint8(s.rp.certSub)
		s.hasCert[k] = true
	} else {
		s.occ[k] += uint8(s.rp.digestSub)
	}
	s.slots[k*s.maxPer+int(s.count[k])] = int32(id)
	s.count[k]++
	if !s.fits(k, false) {
		s.dFree.close(k)
	}
	if !s.fits(k, true) {
		s.cFree.close(k)
	}
	return k
}

// extended holds the accounting of one simulated case. Counts cover the
// messages generated in the measured periods.
type extended struct {
	generated     uint64
	notSchedDig   uint64 // digest-form messages not scheduled
	notSchedCert  uint64 // certificate-carrying messages not scheduled
	scheduled     uint64
	discarded     uint64 // unknown signer at the receiver
	authenticated uint64
	missed        uint64 // authenticated after the deadline
	busy          float64
	span          float64 // measured time across all configurations
	maxBacklog    int
	offeredSub    float64 // subchannel demand of generated messages
	carriedSub    float64 // subchannels of scheduled messages
	capacitySub   float64 // subchannels available in the measured periods
	schedDelaySum float64
	schedDelayMax float64
	p99           float64
	radioModel    bool // whether radio accounting applies
}

// backlog tracks the number of messages in a FCFS system.
type backlog struct {
	done []float64 // completion times, nondecreasing
	head int
}

func (b *backlog) reset() { b.done, b.head = b.done[:0], 0 }

// arrive records a message arriving at t that completes at c and returns
// the number of messages in the system, including this one.
func (b *backlog) arrive(t, c float64) int {
	for b.head < len(b.done) && b.done[b.head] <= t {
		b.head++
	}
	if b.head > 1024 && b.head*2 > len(b.done) {
		b.done = append(b.done[:0], b.done[b.head:]...)
		b.head = 0
	}
	b.done = append(b.done, c)
	return len(b.done) - b.head
}

// configure draws one random configuration exactly as the legacy simulate
// function does (same calls, same order) for the cv2x model.
func configure(rng *rand.Rand, n, m int, pool, pref, certOff, order []int) {
	nRes := len(pool)
	for i := range pool {
		pool[i] = i
	}
	for i := 0; i < n; i++ {
		j := i + rng.IntN(nRes-i)
		pool[i], pool[j] = pool[j], pool[i]
		pref[i] = pool[i] / m
	}
	for i := range order {
		order[i] = i
		certOff[i] = rng.IntN(certEvery)
	}
	rng.Shuffle(n, func(a, b int) { order[a], order[b] = order[b], order[a] })
	sort.SliceStable(order, func(a, b int) bool { return pref[order[a]] < pref[order[b]] })
}

// scheduleConfig generates the messages of one configuration in generation
// order and places them; placed[id] is the subframe of msgs[id] or -1.
func scheduleConfig(sch *scheduler, msgs []genMsg, placed []int, order, pref, certOff []int, nPer, perPeriod int) ([]genMsg, []int) {
	sch.reset()
	msgs, placed = msgs[:0], placed[:0]
	for p := 0; p < nPer; p++ {
		for _, i := range order {
			g := genMsg{sender: i, period: p, gen: p*perPeriod + pref[i], cert: p%certEvery == certOff[i]}
			msgs = append(msgs, g)
			placed = append(placed, sch.place(len(msgs)-1, g.gen, g.cert))
		}
	}
	return msgs, placed
}

// simulateCertAware runs the certificate-aware model. The configuration
// stream is seeded like the legacy model's; service-time draws use a
// separate stream, so the configurations are the legacy ones whatever the
// service-time family.
func simulateCertAware(n, m int, rp radioParams, sc scenario, deadline float64, configs, warm, meas int, seed uint64) (result, extended) {
	res := result{model: "cv2x-cert", scen: sc.name, n: n, rho: float64(n) / period * sc.meanWork}
	ex := extended{radioModel: true}
	if sc.coldStart {
		warm = 0
	}
	perPeriod := int(math.Round(period / subframe))
	if n > perPeriod*m {
		panic("more senders than SPS resources")
	}
	cfgRng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	svcRng := rand.New(rand.NewPCG(seed^0x5bd1e995, seed+0x632be59bd9b4e019))
	h := &hist{bins: make([]uint32, nBins), limit: deadline}
	pool := make([]int, perPeriod*m)
	pref := make([]int, n)
	certOff := make([]int, n)
	order := make([]int, n)
	known := make([]bool, n)
	nPer := warm + meas
	nSub := nPer*perPeriod + rp.window
	sch := newScheduler(rp, nSub)
	msgs := make([]genMsg, 0, n*nPer)
	placed := make([]int, 0, n*nPer)
	var bl backlog
	var cfgMiss uint64

	for c := 0; c < configs; c++ {
		configure(cfgRng, n, m, pool, pref, certOff, order)
		for i := range known {
			known[i] = !sc.coldStart
		}
		msgs, placed = scheduleConfig(sch, msgs, placed, order, pref, certOff, nPer, perPeriod)
		for id, g := range msgs {
			if g.period < warm {
				continue
			}
			k := placed[id]
			ex.generated++
			cost := rp.digestSub
			if g.cert {
				cost = rp.certSub
			}
			ex.offeredSub += float64(cost)
			switch {
			case k < 0 && g.cert:
				ex.notSchedCert++
			case k < 0:
				ex.notSchedDig++
			default:
				ex.scheduled++
				ex.carriedSub += float64(cost)
				d := float64(k-g.gen) * subframe
				ex.schedDelaySum += d
				if d > ex.schedDelayMax {
					ex.schedDelayMax = d
				}
			}
		}
		ex.capacitySub += float64(rp.subchannels * meas * perPeriod)
		ex.span += float64(meas) * period

		// Receiver: serial FCFS verifier, arrivals at the end of each subframe
		// in placement order.
		free := 0.0
		bl.reset()
		missBefore := h.over
		for k := 0; k < nSub; k++ {
			for j := 0; j < int(sch.count[k]); j++ {
				g := msgs[sch.slots[k*sch.maxPer+j]]
				arr := float64(k+1) * subframe
				i := g.sender
				kind := kDigestHit
				switch {
				case g.cert && (sc.noCache || !known[i]):
					kind = kCertNew
					known[i] = true
				case g.cert:
					kind = kCertKnown
				case !known[i]:
					kind = kDigestUnknown
				}
				start := arr
				if free > start {
					start = free
				}
				svc := sc.svc(kind, svcRng)
				free = start + svc
				q := bl.arrive(arr, free)
				if g.period >= warm {
					ex.busy += svc
					if q > ex.maxBacklog {
						ex.maxBacklog = q
					}
					if kind == kDigestUnknown {
						ex.discarded++
						continue
					}
					h.add(free - arr)
				}
			}
		}
		if h.over > missBefore {
			cfgMiss++
		}
	}
	ex.authenticated = h.n
	ex.missed = h.over
	if h.n > 0 {
		res.mean = h.sum / float64(h.n)
		res.p999 = h.quantile(0.999)
		ex.p99 = h.quantile(0.99)
		res.max = h.max
		res.pMiss = float64(h.over) / float64(h.n)
	}
	res.cfgMiss = float64(cfgMiss) / float64(configs)
	if ex.generated > 0 {
		res.discarded = float64(ex.discarded) / float64(ex.generated)
	}
	// The receiver is unstable when the carried verification work reaches
	// its capacity; the delay statistics then depend on the horizon.
	if ex.busy/ex.span >= 1 {
		res.unstable = true
	}
	return res, ex
}
