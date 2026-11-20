package main

import (
	"math"
	"math/rand/v2"
	"os"
	"testing"
)

func baseOpts() options {
	return options{
		mode: "fixed", arrival: "legacy", report: "standard", dir: "../../results",
		armPrim: 0.446, deadline: 10, configs: 2000, slots: 2, seed: 7,
		counts: []int{100, 150, 200}, rp: defaultRadio,
	}
}

// The legacy simulator must reproduce, byte for byte, the output the
// original code produced (golden files generated before any change).
func TestLegacyGolden(t *testing.T) {
	cases := []struct {
		golden string
		edit   func(*options)
	}{
		{"testdata/golden_fixed_c2000.md", func(o *options) {}},
		{"testdata/golden_measured_run1_c2000.md", func(o *options) {
			o.mode, o.tag, o.counts = "measured", "intel-i9-10980hk-run1", []int{100, 150, 162, 200}
		}},
	}
	for _, c := range cases {
		o := baseOpts()
		c.edit(&o)
		got, _, _, err := run(&o)
		if err != nil {
			t.Fatal(err)
		}
		want, err := os.ReadFile(c.golden)
		if err != nil {
			t.Fatal(err)
		}
		if got != string(want) {
			t.Errorf("%s: legacy output changed", c.golden)
		}
	}
}

func detScenario(name string, digest, cert float64, noCache, cold bool) scenario {
	sc := fixedScenario(name, 0.446e-3, digest, cert)
	sc.noCache, sc.coldStart = noCache, cold
	return sc
}

// checkSchedule verifies the scheduling invariants for one configuration.
func checkSchedule(t *testing.T, sch *scheduler, msgs []genMsg, placed []int, rp radioParams) {
	t.Helper()
	seen := make([]int, len(msgs))
	for k := 0; k < sch.nSub; k++ {
		if int(sch.occ[k]) > rp.subchannels {
			t.Fatalf("subframe %d: %d subchannels allocated", k, sch.occ[k])
		}
		sum, certs := 0, 0
		for j := 0; j < int(sch.count[k]); j++ {
			id := sch.slots[k*sch.maxPer+j]
			seen[id]++
			if msgs[id].cert {
				sum += rp.certSub
				certs++
			} else {
				sum += rp.digestSub
			}
		}
		if sum != int(sch.occ[k]) {
			t.Fatalf("subframe %d: occupancy %d, messages need %d", k, sch.occ[k], sum)
		}
		if certs > 0 && (sch.count[k] != 1 || int(sch.occ[k]) != rp.subchannels) {
			t.Fatalf("subframe %d: certificate message shares the subframe (%d messages, %d subchannels)", k, sch.count[k], sch.occ[k])
		}
	}
	for id, k := range placed {
		if k < 0 {
			if seen[id] != 0 {
				t.Fatalf("unscheduled message %d appears in a subframe", id)
			}
			continue
		}
		if seen[id] != 1 {
			t.Fatalf("message %d appears %d times", id, seen[id])
		}
		if k < msgs[id].gen {
			t.Fatalf("message %d scheduled in subframe %d before its generation subframe %d", id, k, msgs[id].gen)
		}
		if k >= msgs[id].gen+rp.window {
			t.Fatalf("message %d scheduled outside its window", id)
		}
	}
}

func TestSchedulingInvariants(t *testing.T) {
	rp := defaultRadio
	for _, n := range []int{100, 150, 200} {
		rng := rand.New(rand.NewPCG(1, 2))
		pool, pref, certOff, order := make([]int, 200), make([]int, n), make([]int, n), make([]int, n)
		nPer := 20
		sch := newScheduler(rp, nPer*100+rp.window)
		var msgs []genMsg
		var placed []int
		for c := 0; c < 200; c++ {
			configure(rng, n, 2, pool, pref, certOff, order)
			msgs, placed = scheduleConfig(sch, msgs, placed, order, pref, certOff, nPer, 100)
			checkSchedule(t, sch, msgs, placed, rp)
		}
	}
}

func TestOfferedLoadAndInfeasibility(t *testing.T) {
	want := map[int]float64{100: 0.52, 150: 0.78, 200: 1.04}
	for n, u := range want {
		_, ex := simulateCertAware(n, 2, defaultRadio, detScenario("S", 1, 1, false, false), 10e-3, 200, 10, 10, 7)
		if got := ex.offeredSub / ex.capacitySub; math.Abs(got-u) > 1e-12 {
			t.Errorf("N=%d: offered occupancy %.6f, want %.2f", n, got, u)
		}
		notSched := ex.notSchedDig + ex.notSchedCert
		if ex.generated != ex.scheduled+notSched {
			t.Errorf("N=%d: generated %d != scheduled %d + not scheduled %d", n, ex.generated, ex.scheduled, notSched)
		}
		if n == 100 && notSched != 0 {
			t.Errorf("N=100: %d messages not scheduled at 52%% offered load", notSched)
		}
		if n == 200 && notSched == 0 {
			t.Errorf("N=200: 104%% offered load scheduled without loss")
		}
	}
}

// With a certificate costing as much as a digest message, the
// certificate-aware schedule degenerates to the legacy slot model, so both
// simulators must agree.
func TestCertAwareReducesToLegacy(t *testing.T) {
	rp := defaultRadio
	rp.certSub = rp.digestSub
	for _, sc := range []scenario{
		detScenario("S", 1, 1, false, false),
		detScenario("S/2S", 1, 2, false, false),
		detScenario("cold", 1, 2, false, true),
	} {
		for _, n := range []int{100, 150, 200} {
			lr, lx := simulate("cv2x", n, 2, sc, 10e-3, 300, 10, 10, 7)
			cr, cx := simulateCertAware(n, 2, rp, sc, 10e-3, 300, 10, 10, 7)
			if lr.unstable {
				continue
			}
			if cx.notSchedDig+cx.notSchedCert != 0 {
				t.Fatalf("%s N=%d: messages not scheduled", sc.name, n)
			}
			if lx.authenticated != cx.authenticated || lx.discarded != cx.discarded || lx.missed != cx.missed || lx.generated != cx.generated {
				t.Errorf("%s N=%d: counts differ: legacy %+v cert-aware %+v", sc.name, n, lx, cx)
			}
			if math.Abs(lr.mean-cr.mean) > 1e-12 || math.Abs(lr.max-cr.max) > 1e-12 {
				t.Errorf("%s N=%d: delays differ: legacy %.9f/%.9f cert-aware %.9f/%.9f", sc.name, n, lr.mean, lr.max, cr.mean, cr.max)
			}
		}
	}
}

// Every case must satisfy the accounting identities, and unknown-signer
// discards must never enter the authentication or deadline statistics.
func TestAccountingIdentities(t *testing.T) {
	scens := []scenario{
		detScenario("S", 1, 1, false, false),
		detScenario("nocache", 1, 2.233, true, false),
		detScenario("cold", 1.017, 2.233, false, true),
	}
	for _, sc := range scens {
		for _, n := range []int{100, 150, 200} {
			for _, certaware := range []bool{false, true} {
				var r result
				var e extended
				if certaware {
					r, e = simulateCertAware(n, 2, defaultRadio, sc, 10e-3, 200, 10, 10, 7)
				} else {
					r, e = simulate("cv2x", n, 2, sc, 10e-3, 200, 10, 10, 7)
					if r.unstable {
						continue
					}
				}
				if e.generated != e.scheduled+e.notSchedDig+e.notSchedCert {
					t.Errorf("%s N=%d cert-aware=%t: generated != scheduled + not scheduled", sc.name, n, certaware)
				}
				if e.scheduled != e.discarded+e.authenticated {
					t.Errorf("%s N=%d cert-aware=%t: scheduled %d != discarded %d + authenticated %d", sc.name, n, certaware, e.scheduled, e.discarded, e.authenticated)
				}
				if e.missed > e.authenticated {
					t.Errorf("%s N=%d: more misses than authenticated messages", sc.name, n)
				}
				if !sc.coldStart && e.discarded != 0 {
					t.Errorf("%s N=%d: discards outside cold start", sc.name, n)
				}
			}
		}
	}
}

// The original cold-start experiment discards about 20% of the messages of
// the first second: a sender's digest-form messages before its first
// certificate-carrying message (offset uniform in 0..4, so 2 of 10 messages
// on average). This is a property of the transmission schedule, not a
// measured packet-loss rate.
func TestColdStartDiscardBaseline(t *testing.T) {
	for _, n := range []int{100, 150, 162, 200} {
		r, e := simulate("cv2x", n, 2, detScenario("cold", 1.017, 2.233, false, true), 10e-3, 2000, 10, 10, 7)
		if math.Abs(r.discarded-0.2) > 0.005 {
			t.Errorf("N=%d: discarded fraction %.4f, want about 0.200", n, r.discarded)
		}
		if got := float64(e.discarded) / float64(e.generated); got != r.discarded {
			t.Errorf("N=%d: discard fraction denominator is not the generated count", n)
		}
	}
}

// Discards and deadline misses are separate: with a verifier far too slow,
// every authenticated message misses the deadline, and the discarded ones
// are neither authenticated nor counted as misses.
func TestDiscardsAreNotMisses(t *testing.T) {
	slow := scenario{name: "slow", svc: func(int, *rand.Rand) float64 { return 20e-3 }, meanWork: 20e-3, coldStart: true}
	for _, certaware := range []bool{false, true} {
		var r result
		var e extended
		if certaware {
			r, e = simulateCertAware(100, 2, defaultRadio, slow, 10e-3, 20, 10, 10, 7)
		} else {
			r, e = simulate("cv2x", 100, 2, slow, 10e-3, 20, 10, 10, 7)
		}
		if e.discarded == 0 || e.authenticated == 0 {
			t.Fatalf("cert-aware=%t: expected discards and authentications", certaware)
		}
		if e.missed != e.authenticated || r.pMiss != 1 {
			t.Errorf("cert-aware=%t: missed %d of %d authenticated", certaware, e.missed, e.authenticated)
		}
		if e.missed+e.discarded != e.scheduled {
			t.Errorf("cert-aware=%t: discards counted as misses or lost", certaware)
		}
	}
}

// The receiver must charge each path its own cost: an unknown signer's
// digest message the digest_unknown time, a first certificate the cert_new
// time, a repeated certificate the cert_known time, a cached digest the
// digest_hit time.
func TestServiceCostPerPath(t *testing.T) {
	const n = 100
	for _, certaware := range []bool{false, true} {
		for _, cold := range []bool{false, true} {
			for _, noCache := range []bool{false, true} {
				if cold && noCache {
					continue
				}
				var counts [nKinds]int
				sc := scenario{name: "count", meanWork: 0.1e-3, coldStart: cold, noCache: noCache,
					svc: func(kind int, _ *rand.Rand) float64 { counts[kind]++; return 0.1e-3 }}
				configs := 50
				if certaware {
					simulateCertAware(n, 2, defaultRadio, sc, 10e-3, configs, 10, 10, 7)
				} else {
					simulate("cv2x", n, 2, sc, 10e-3, configs, 10, 10, 7)
				}
				periods := 20
				if cold {
					periods = 10
				}
				certs := n * configs * periods / certEvery
				total := n * configs * periods
				switch {
				case cold:
					// Each sender: first certificate new, second known.
					if counts[kCertNew] != n*configs || counts[kCertKnown] != certs-n*configs || counts[kDigestUnknown] == 0 {
						t.Errorf("cert-aware=%t cold: kinds %v", certaware, counts)
					}
				case noCache:
					if counts[kCertNew] != certs || counts[kCertKnown] != 0 || counts[kDigestUnknown] != 0 {
						t.Errorf("cert-aware=%t no cache: kinds %v", certaware, counts)
					}
				default:
					if counts[kCertNew] != 0 || counts[kCertKnown] != certs || counts[kDigestUnknown] != 0 {
						t.Errorf("cert-aware=%t steady: kinds %v", certaware, counts)
					}
				}
				sum := 0
				for _, c := range counts {
					sum += c
				}
				if sum != total {
					t.Errorf("cert-aware=%t: %d service calls for %d messages", certaware, sum, total)
				}
			}
		}
	}
}

// Under the default certificate-aware model a subframe carries at most two
// digest-form messages or one certificate-carrying message, so with median
// Falcon-512 work scaled to 0.446 ms (each subframe's work below 1 ms) the
// receiver never carries work into the next subframe.
func TestCertAwareMedianWorkBound(t *testing.T) {
	k := [nKinds]float64{1.0174, 1.0983, 2.2330, 0.0021}
	det := func(kind int, _ *rand.Rand) float64 { return k[kind] * 0.446e-3 }
	for _, sc := range []scenario{
		{name: "steady", svc: det, meanWork: 0.5e-3},
		{name: "nocache", svc: det, meanWork: 0.5e-3, noCache: true},
		{name: "cold", svc: det, meanWork: 0.5e-3, coldStart: true},
	} {
		for _, n := range []int{100, 150, 200} {
			r, _ := simulateCertAware(n, 2, defaultRadio, sc, 10e-3, 100, 10, 10, 7)
			bound := math.Max(2*k[kDigestHit], k[kCertNew]) * 0.446e-3
			if r.max > bound+1e-12 {
				t.Errorf("%s N=%d: max delay %.6f ms exceeds per-subframe work bound %.6f ms", sc.name, n, r.max*1e3, bound*1e3)
			}
		}
	}
}
