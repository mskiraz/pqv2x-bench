// Command queuesim evaluates the receiver authentication queue under periodic
// V2X traffic (10 Hz senders, one in five messages carrying the certificate).
//
// Two arrival models:
//
//	unconstrained  independent uniform phases, no channel serialisation
//	               (the N*D/D/1 queue)
//	cv2x           SAE J3161 C-V2X Mode 4: each sender holds a distinct SPS
//	               resource out of 100 subframes x m slots; messages reach
//	               the verifier at the end of their subframe
//
// Two families of service times:
//
//	-mode fixed     multiples of the published ARMv8 Falcon-512 primitive
//	                time (reproduces the fixed-work table of the paper)
//	-mode measured  per-path times measured by cmd/bench, rescaled so that
//	                the measured primitive mean equals the published ARMv8
//	                primitive time
//
// The verifier is a single FCFS server.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	period    = 0.100 // s, 10 Hz
	subframe  = 1e-3  // s
	certEvery = 5     // SAE J2945/1
	binWidth  = 1e-6  // s, histogram resolution
	nBins     = 400_000
)

// kind of message from the receiver's point of view
const (
	kDigestHit = iota
	kCertKnown
	kCertNew
	kDigestUnknown
	nKinds
)

// service returns a service time in seconds for a message kind.
type service func(kind int, rng *rand.Rand) float64

type scenario struct {
	name      string
	svc       service
	meanWork  float64 // mean service per message in steady state, for rho
	noCache   bool    // every certificate-carrying message is validated anew
	coldStart bool    // all senders unknown at t = 0; transient is measured
}

type result struct {
	model, scen string
	n           int
	rho         float64
	mean, p999  float64
	max         float64
	pMiss       float64
	cfgMiss     float64 // share of configurations with at least one miss
	discarded   float64 // share of messages discarded as unknown signer
	unstable    bool
}

type hist struct {
	bins  []uint32
	n     uint64
	sum   float64
	max   float64
	over  uint64
	limit float64
}

func (h *hist) add(t float64) {
	i := int(t / binWidth)
	if i >= nBins {
		i = nBins - 1
	}
	h.bins[i]++
	h.n++
	h.sum += t
	if t > h.max {
		h.max = t
	}
	if t > h.limit {
		h.over++
	}
}

func (h *hist) quantile(p float64) float64 {
	target := uint64(math.Ceil(p * float64(h.n)))
	var c uint64
	for i, v := range h.bins {
		c += uint64(v)
		if c >= target {
			return (float64(i) + 1) * binWidth
		}
	}
	return h.max
}

// simulate runs configs independent phase/resource configurations.
func simulate(model string, n, m int, sc scenario, deadline float64, configs, warm, meas int, seed uint64) result {
	res := result{model: model, scen: sc.name, n: n, rho: float64(n) / period * sc.meanWork}
	if res.rho >= 1 && !sc.coldStart {
		res.unstable = true
		return res
	}
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	h := &hist{bins: make([]uint32, nBins), limit: deadline}
	phase := make([]float64, n)
	certOff := make([]int, n)
	known := make([]bool, n)
	order := make([]int, n)
	nRes := int(math.Round(period/subframe)) * m
	pool := make([]int, nRes)
	var cfgMiss, discarded, total uint64

	if sc.coldStart {
		warm = 0
	}
	for c := 0; c < configs; c++ {
		switch model {
		case "unconstrained":
			for i := range phase {
				phase[i] = rng.Float64() * period
			}
		case "cv2x":
			for i := range pool {
				pool[i] = i
			}
			for i := 0; i < n; i++ { // partial Fisher-Yates: n distinct resources
				j := i + rng.IntN(nRes-i)
				pool[i], pool[j] = pool[j], pool[i]
				phase[i] = float64(pool[i]/m+1) * subframe
			}
		default:
			panic("unknown model " + model)
		}
		for i := range order {
			order[i] = i
			certOff[i] = rng.IntN(certEvery)
			known[i] = !sc.coldStart
		}
		rng.Shuffle(n, func(a, b int) { order[a], order[b] = order[b], order[a] }) // random order within a batch
		sort.SliceStable(order, func(a, b int) bool { return phase[order[a]] < phase[order[b]] })

		free := 0.0 // time at which the verifier becomes idle
		missBefore := h.over
		for p := 0; p < warm+meas; p++ {
			base := float64(p) * period
			for _, i := range order {
				arr := base + phase[i]
				carries := p%certEvery == certOff[i]
				kind := kDigestHit
				switch {
				case carries && (sc.noCache || !known[i]):
					kind = kCertNew
					known[i] = true
				case carries:
					kind = kCertKnown
				case !known[i]:
					kind = kDigestUnknown
				}
				start := arr
				if free > start {
					start = free
				}
				free = start + sc.svc(kind, rng)
				if p >= warm {
					total++
					if kind == kDigestUnknown {
						discarded++ // not authenticated; excluded from delay statistics
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
	res.mean = h.sum / float64(h.n)
	res.p999 = h.quantile(0.999)
	res.max = h.max
	res.pMiss = float64(h.over) / float64(h.n)
	res.cfgMiss = float64(cfgMiss) / float64(configs)
	res.discarded = float64(discarded) / float64(total)
	return res
}

func readSamples(path string) ([]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []float64
	s := bufio.NewScanner(f)
	for s.Scan() {
		v, err := strconv.ParseFloat(strings.TrimSpace(s.Text()), 64)
		if err != nil {
			return nil, err
		}
		out = append(out, v*1e-9)
	}
	return out, s.Err()
}

func mean(x []float64) float64 {
	var s float64
	for _, v := range x {
		s += v
	}
	return s / float64(len(x))
}

func median(x []float64) float64 {
	s := append([]float64(nil), x...)
	sort.Float64s(s)
	return s[len(s)/2]
}

func fixedScenario(name string, s, kDigest, kCert float64) scenario {
	return scenario{
		name: name,
		svc: func(kind int, _ *rand.Rand) float64 {
			if kind == kCertKnown || kind == kCertNew {
				return kCert * s
			}
			return kDigest * s
		},
		meanWork: (0.8*kDigest + 0.2*kCert) * s,
	}
}

func printTable(w *strings.Builder, title string, rs []result, deadline float64) {
	fmt.Fprintf(w, "## %s\n\n", title)
	fmt.Fprintf(w, "| Arrival model | Scenario | N | rho | Mean T (ms) | p99.9 T (ms) | Max T (ms) | Pr[T > %.0f ms] | Configs with a miss | Discarded (unknown signer) |\n", deadline*1e3)
	w.WriteString("|---|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, r := range rs {
		if r.unstable {
			fmt.Fprintf(w, "| %s | %s | %d | %.3f | unstable | | | | | |\n", r.model, r.scen, r.n, r.rho)
			continue
		}
		fmt.Fprintf(w, "| %s | %s | %d | %.3f | %.3f | %.3f | %.3f | %.2e | %.4f | %.3f |\n",
			r.model, r.scen, r.n, r.rho, r.mean*1e3, r.p999*1e3, r.max*1e3, r.pMiss, r.cfgMiss, r.discarded)
	}
	w.WriteString("\n")
}

func main() {
	mode := flag.String("mode", "measured", "fixed | measured")
	dir := flag.String("results", "results", "directory written by cmd/bench")
	tag := flag.String("tag", "", "benchmark tag (default: only tag found)")
	armPrim := flag.Float64("arm-primitive-ms", 0.446, "published ARMv8 Falcon-512 primitive verification time")
	deadline := flag.Float64("deadline-ms", 10, "authentication deadline D")
	configs := flag.Int("configs", 10000, "random phase/resource configurations per case")
	slots := flag.Int("m", 2, "messages per C-V2X subframe")
	seed := flag.Uint64("seed", 7, "random seed")
	out := flag.String("out", "", "write the Markdown report to this file as well")
	sendersFlag := flag.String("senders", "100,150,200", "comma-separated sender counts N")
	flag.Parse()

	var counts []int
	for _, f := range strings.Split(*sendersFlag, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil || n <= 0 {
			fmt.Fprintln(os.Stderr, "invalid -senders value:", f)
			os.Exit(1)
		}
		counts = append(counts, n)
	}

	s := *armPrim * 1e-3
	d := *deadline * 1e-3
	var b strings.Builder
	var rs []result

	switch *mode {
	case "fixed":
		rs = append(rs,
			simulate("unconstrained", 100, *slots, fixedScenario("S", s, 1, 1), d, *configs, 10, 10, *seed),
			simulate("unconstrained", 200, *slots, fixedScenario("S", s, 1, 1), d, *configs, 10, 10, *seed))
		for _, sc := range []scenario{
			fixedScenario("S", s, 1, 1),
			fixedScenario("S/2S", s, 1, 2),
			fixedScenario("1.5S", s, 1.5, 1.5),
			fixedScenario("2S", s, 2, 2),
		} {
			for _, n := range counts {
				rs = append(rs, simulate("cv2x", n, *slots, sc, d, *configs, 10, 10, *seed))
			}
		}
		printTable(&b, fmt.Sprintf("Fixed work, S = %.3f ms (Falcon-512 primitive, ARMv8)", *armPrim), rs, d)

	case "measured":
		if *tag == "" {
			m, _ := filepath.Glob(filepath.Join(*dir, "raw", "*_falcon512_primitive.txt"))
			if len(m) != 1 {
				fmt.Fprintln(os.Stderr, "specify -tag: found", len(m), "benchmark runs")
				os.Exit(1)
			}
			*tag = strings.TrimSuffix(filepath.Base(m[0]), "_falcon512_primitive.txt")
		}
		load := func(op string) []float64 {
			x, err := readSamples(filepath.Join(*dir, "raw", *tag+"_falcon512_"+op+".txt"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return x
		}
		prim := load("primitive")
		ops := [nKinds][]float64{load("digest_hit"), load("cert_known"), load("cert_new"), load("digest_unknown")}
		names := [nKinds]string{"digest_hit", "cert_known", "cert_new", "digest_unknown"}

		// Median ratios give deterministic per-path work; scaled samples keep
		// the measured jitter.
		pm, pmed := mean(prim), median(prim)
		var kMed [nKinds]float64
		fmt.Fprintf(&b, "# Queue simulation from measured path times (%s)\n\n", *tag)
		fmt.Fprintf(&b, "Measured Falcon-512 primitive: mean %.2f us, median %.2f us. Published ARMv8 primitive: %.3f ms.\n\n", pm*1e6, pmed*1e6, *armPrim)
		b.WriteString("| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |\n|---|---:|---:|---:|\n")
		for k := 0; k < nKinds; k++ {
			kMed[k] = median(ops[k]) / pmed
			fmt.Fprintf(&b, "| %s | %.4f | %.4f | %.4f |\n", names[k], kMed[k], mean(ops[k])/pm, kMed[k]**armPrim)
		}
		fmt.Fprintf(&b, "\nC-V2X per-message budget T_sf/m = %.3f ms, i.e. %.3f x primitive.\n\n", subframe/float64(*slots)*1e3, subframe/float64(*slots)/s)

		det := func(kind int, _ *rand.Rand) float64 { return kMed[kind] * s }
		scale := s / pm
		emp := func(kind int, rng *rand.Rand) float64 { return ops[kind][rng.IntN(len(ops[kind]))] * scale }
		steadyDet := (0.8*kMed[kDigestHit] + 0.2*kMed[kCertKnown]) * s
		noCacheDet := (0.8*kMed[kDigestHit] + 0.2*kMed[kCertNew]) * s
		steadyEmp := (0.8*mean(ops[kDigestHit]) + 0.2*mean(ops[kCertKnown])) * scale
		noCacheEmp := (0.8*mean(ops[kDigestHit]) + 0.2*mean(ops[kCertNew])) * scale

		scens := []scenario{
			{name: "steady, median work", svc: det, meanWork: steadyDet},
			{name: "steady, measured jitter", svc: emp, meanWork: steadyEmp},
			{name: "no caching, median work", svc: det, meanWork: noCacheDet, noCache: true},
			{name: "no caching, measured jitter", svc: emp, meanWork: noCacheEmp, noCache: true},
			{name: "cold start, median work", svc: det, meanWork: steadyDet, coldStart: true},
			{name: "cold start, measured jitter", svc: emp, meanWork: steadyEmp, coldStart: true},
		}
		for _, sc := range scens {
			for _, n := range counts {
				// Cold start: empty cache and empty queue at t = 0, 20 periods measured.
				rs = append(rs, simulate("cv2x", n, *slots, sc, d, *configs, 10, 10, *seed))
			}
		}
		printTable(&b, "C-V2X SPS, Falcon-512, measured receiver paths", rs, d)

	default:
		fmt.Fprintln(os.Stderr, "unknown mode", *mode)
		os.Exit(1)
	}

	fmt.Print(b.String())
	if *out != "" {
		if err := os.WriteFile(*out, []byte(b.String()), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
