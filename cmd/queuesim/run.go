package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// options are the command-line settings.
type options struct {
	mode     string // fixed | measured
	arrival  string // legacy | certaware
	report   string // standard | extended
	dir      string
	tag      string
	armPrim  float64 // ms
	deadline float64 // ms
	configs  int
	slots    int
	seed     uint64
	counts   []int
	rp       radioParams
}

// caseOut is one simulated case.
type caseOut struct {
	res result
	ex  extended
}

// run executes the simulation and returns the report in the requested form
// (o.report), the extended report, and the per-case results.
func run(o *options) (string, string, []caseOut, error) {
	if o.arrival != "legacy" && o.arrival != "certaware" {
		return "", "", nil, fmt.Errorf("unknown -arrival %q", o.arrival)
	}
	if o.report != "standard" && o.report != "extended" {
		return "", "", nil, fmt.Errorf("unknown -report %q", o.report)
	}
	if o.arrival == "certaware" {
		o.report = "extended" // the legacy table has no radio accounting
		if o.rp.window < 1 || o.rp.window > 100 {
			return "", "", nil, errors.New("-sched-window-ms must be between 1 and 100 (one 100-ms period)")
		}
	}
	s := o.armPrim * 1e-3
	d := o.deadline * 1e-3
	var b strings.Builder
	var cs []caseOut

	sim := func(model string, n int, sc scenario) caseOut {
		if o.arrival == "certaware" && model == "cv2x" {
			r, e := simulateCertAware(n, o.slots, o.rp, sc, d, o.configs, 10, 10, o.seed)
			return caseOut{r, e}
		}
		r, e := simulate(model, n, o.slots, sc, d, o.configs, 10, 10, o.seed)
		return caseOut{r, e}
	}
	var std, ext string
	emit := func(title string) {
		var t1, t2 strings.Builder
		rs := make([]result, len(cs))
		for i, c := range cs {
			rs[i] = c.res
		}
		printTable(&t1, title, rs, d)
		printExtended(&t2, title, cs, d, *o)
		std, ext = b.String()+t1.String(), b.String()+t2.String()
	}

	switch o.mode {
	case "fixed":
		if o.arrival == "legacy" {
			cs = append(cs,
				sim("unconstrained", 100, fixedScenario("S", s, 1, 1)),
				sim("unconstrained", 200, fixedScenario("S", s, 1, 1)))
		}
		for _, sc := range []scenario{
			fixedScenario("S", s, 1, 1),
			fixedScenario("S/2S", s, 1, 2),
			fixedScenario("1.5S", s, 1.5, 1.5),
			fixedScenario("2S", s, 2, 2),
		} {
			for _, n := range o.counts {
				cs = append(cs, sim("cv2x", n, sc))
			}
		}
		title := fmt.Sprintf("Fixed work, S = %.3f ms (Falcon-512 primitive, ARMv8)", o.armPrim)
		if o.arrival == "certaware" {
			title += ", certificate-aware C-V2X scheduling"
		}
		emit(title)

	case "measured":
		if o.tag == "" {
			m, _ := filepath.Glob(filepath.Join(o.dir, "raw", "*_falcon512_primitive.txt"))
			if len(m) != 1 {
				return "", "", nil, fmt.Errorf("specify -tag: found %d benchmark runs", len(m))
			}
			o.tag = strings.TrimSuffix(filepath.Base(m[0]), "_falcon512_primitive.txt")
		}
		var loadErr error
		load := func(op string) []float64 {
			x, err := readSamples(filepath.Join(o.dir, "raw", o.tag+"_falcon512_"+op+".txt"))
			if err != nil && loadErr == nil {
				loadErr = err
			}
			return x
		}
		prim := load("primitive")
		ops := [nKinds][]float64{load("digest_hit"), load("cert_known"), load("cert_new"), load("digest_unknown")}
		if loadErr != nil {
			return "", "", nil, loadErr
		}
		names := [nKinds]string{"digest_hit", "cert_known", "cert_new", "digest_unknown"}

		// Median ratios give deterministic per-path work; scaled samples keep
		// the measured jitter.
		pm, pmed := mean(prim), median(prim)
		var kMed [nKinds]float64
		fmt.Fprintf(&b, "# Queue simulation from measured path times (%s)\n\n", o.tag)
		fmt.Fprintf(&b, "Measured Falcon-512 primitive: mean %.2f us, median %.2f us. Published ARMv8 primitive: %.3f ms.\n\n", pm*1e6, pmed*1e6, o.armPrim)
		b.WriteString("| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |\n|---|---:|---:|---:|\n")
		for k := 0; k < nKinds; k++ {
			kMed[k] = median(ops[k]) / pmed
			fmt.Fprintf(&b, "| %s | %.4f | %.4f | %.4f |\n", names[k], kMed[k], mean(ops[k])/pm, kMed[k]*o.armPrim)
		}
		fmt.Fprintf(&b, "\nC-V2X per-message budget T_sf/m = %.3f ms, i.e. %.3f x primitive.\n\n", subframe/float64(o.slots)*1e3, subframe/float64(o.slots)/s)

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
			for _, n := range o.counts {
				// Cold start: empty cache and empty queue at t = 0; no warm-up,
				// so the first 10 periods (1 s) are measured.
				cs = append(cs, sim("cv2x", n, sc))
			}
		}
		title := "C-V2X SPS, Falcon-512, measured receiver paths"
		if o.arrival == "certaware" {
			title += ", certificate-aware C-V2X scheduling"
		}
		emit(title)

	default:
		return "", "", nil, fmt.Errorf("unknown mode %q", o.mode)
	}
	if o.report == "extended" {
		return ext, ext, cs, nil
	}
	return std, ext, cs, nil
}

func pct(num, den uint64) float64 {
	if den == 0 {
		return 0
	}
	return float64(num) / float64(den)
}

// printExtended writes the accounting table with every denominator explicit.
func printExtended(w *strings.Builder, title string, cs []caseOut, deadline float64, o options) {
	fmt.Fprintf(w, "## %s\n\n", title)
	if o.arrival == "certaware" {
		fmt.Fprintf(w, "Radio model: %d subchannels per 1-ms subframe; digest-form message %d subchannels; certificate-carrying message %d subchannels and never shares its subframe; scheduling window %d ms; idealised deterministic earliest-fit scheduler (not a 3GPP Mode 4 / SAE J3161 MAC); no propagation loss, collisions, half-duplex loss or retransmissions.\n\n",
			o.rp.subchannels, o.rp.digestSub, o.rp.certSub, o.rp.window)
	} else {
		w.WriteString("Arrival model: legacy abstract receiver-arrival model; every generated message is assumed to reach the receiver, and the radio cost of certificate-carrying messages is not modelled.\n\n")
	}
	fmt.Fprintf(w, "Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = %.0f ms.\n\n", deadline*1e3)
	w.WriteString("| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \\| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |\n")
	w.WriteString("|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|\n")
	for _, c := range cs {
		r, e := c.res, c.ex
		if e.span == 0 { // legacy case with rho >= 1, not simulated
			fmt.Fprintf(w, "| %s | %s | %d | n/a | n/a | – | – | – | – | – | – | unstable (rho = %.3f, not simulated) | | | | | | | | |\n", r.model, r.scen, r.n, r.rho)
			continue
		}
		occ, budget, sched := "n/a", "n/a", "n/a"
		if e.radioModel {
			u := e.offeredSub / e.capacitySub
			occ = fmt.Sprintf("%.3f", u)
			yn := func(lim float64) string {
				if u <= lim {
					return "yes"
				}
				return "no"
			}
			budget = yn(0.30) + " / " + yn(0.65) + " / " + yn(1.00)
			if e.scheduled > 0 {
				sched = fmt.Sprintf("%.3f / %.0f", e.schedDelaySum/float64(e.scheduled)*1e3, e.schedDelayMax*1e3)
			}
		}
		onTime := pct(e.authenticated-e.missed, e.generated)
		util := e.busy / e.span
		if r.unstable {
			fmt.Fprintf(w, "| %s | %s | %d | %s | %s | %d | %d / %d | %d | %d | – | – | receiver unstable (utilisation %.3f) | | | | | | %s | %.3f | %d |\n",
				r.model, r.scen, r.n, occ, budget, e.generated, e.notSchedDig, e.notSchedCert, e.scheduled, e.discarded, util, sched, util, e.maxBacklog)
			continue
		}
		fmt.Fprintf(w, "| %s | %s | %d | %s | %s | %d | %d / %d | %d | %d | %d | %d | %.2e | %.6f | %.3f | %.3f | %.3f | %.3f | %s | %.3f | %d |\n",
			r.model, r.scen, r.n, occ, budget, e.generated, e.notSchedDig, e.notSchedCert, e.scheduled, e.discarded,
			e.authenticated, e.missed, r.pMiss, onTime, r.mean*1e3, e.p99*1e3, r.p999*1e3, r.max*1e3, sched, util, e.maxBacklog)
	}
	w.WriteString("\n")
}

// writeCSV writes one row per case with every count and its denominators.
func writeCSV(wr io.Writer, cs []caseOut, o options) error {
	var b strings.Builder
	b.WriteString("arrival_model,service_mode,tag,scenario,n,seed,configs,deadline_ms,subchannels,digest_subchannels,cert_subchannels,sched_window_ms,offered_occupancy,within_budget_0.30,within_budget_0.65,within_physical_1.00,generated,not_scheduled_digest,not_scheduled_cert,scheduled,unknown_signer_discarded,authenticated,authentication_deadline_missed,p_miss_given_authenticated,on_time_fraction_of_generated,discarded_fraction_of_generated,mean_ms,p99_ms,p999_ms,max_ms,mean_sched_delay_ms,max_sched_delay_ms,verifier_utilisation,max_backlog,rho_offered,receiver_unstable,configs_with_miss\n")
	for _, c := range cs {
		r, e := c.res, c.ex
		occ, b30, b65, b100 := "", "", "", ""
		msd, xsd := "", ""
		if e.radioModel {
			u := e.offeredSub / e.capacitySub
			occ = strconv.FormatFloat(u, 'f', 4, 64)
			b30, b65, b100 = strconv.FormatBool(u <= 0.30), strconv.FormatBool(u <= 0.65), strconv.FormatBool(u <= 1.00)
			if e.scheduled > 0 {
				msd = strconv.FormatFloat(e.schedDelaySum/float64(e.scheduled)*1e3, 'f', 4, 64)
				xsd = strconv.FormatFloat(e.schedDelayMax*1e3, 'f', 1, 64)
			}
		}
		util := 0.0
		if e.span > 0 {
			util = e.busy / e.span
		}
		model := r.model
		sub, ds, cert, win := "", "", "", ""
		if e.radioModel {
			sub, ds, cert, win = strconv.Itoa(o.rp.subchannels), strconv.Itoa(o.rp.digestSub), strconv.Itoa(o.rp.certSub), strconv.Itoa(o.rp.window)
		}
		stats := func(v float64) string {
			if r.unstable || e.authenticated == 0 {
				return ""
			}
			return strconv.FormatFloat(v*1e3, 'f', 4, 64)
		}
		fmt.Fprintf(&b, "%s,%s,%s,%q,%d,%d,%d,%.0f,%s,%s,%s,%s,%s,%s,%s,%s,%d,%d,%d,%d,%d,%d,%d,%.6e,%.6f,%.6f,%s,%s,%s,%s,%s,%s,%.4f,%d,%.4f,%t,%.4f\n",
			model, o.mode, o.tag, r.scen, r.n, o.seed, o.configs, o.deadline, sub, ds, cert, win, occ, b30, b65, b100,
			e.generated, e.notSchedDig, e.notSchedCert, e.scheduled, e.discarded, e.authenticated, e.missed,
			pctF(e.missed, e.authenticated), pctF(e.authenticated-e.missed, e.generated), pctF(e.discarded, e.generated),
			stats(r.mean), stats(e.p99), stats(r.p999), stats(r.max), msd, xsd, util, e.maxBacklog, r.rho, r.unstable, r.cfgMiss)
	}
	_, err := io.WriteString(wr, b.String())
	return err
}

func pctF(num, den uint64) float64 { return pct(num, den) }

func main() {
	var o options
	flag.StringVar(&o.mode, "mode", "measured", "service-time family: fixed | measured")
	flag.StringVar(&o.arrival, "arrival", "legacy", "arrival model: legacy (the original abstract models) | certaware (certificate-aware C-V2X scheduling)")
	flag.StringVar(&o.report, "report", "standard", "standard (the original table; legacy arrival only) | extended (full accounting)")
	flag.StringVar(&o.dir, "results", "results", "directory written by cmd/bench")
	flag.StringVar(&o.tag, "tag", "", "benchmark tag (default: only tag found)")
	flag.Float64Var(&o.armPrim, "arm-primitive-ms", 0.446, "published ARMv8 Falcon-512 primitive verification time")
	flag.Float64Var(&o.deadline, "deadline-ms", 10, "authentication deadline D")
	flag.IntVar(&o.configs, "configs", 10000, "random phase/resource configurations per case")
	flag.IntVar(&o.slots, "m", 2, "messages per C-V2X subframe (legacy model) and digest-form SPS resources per subframe")
	flag.Uint64Var(&o.seed, "seed", 7, "random seed")
	out := flag.String("out", "", "write the Markdown report to this file as well")
	csvOut := flag.String("csv", "", "write the per-case accounting as CSV to this file")
	extOut := flag.String("ext-out", "", "write the extended accounting report (Markdown) to this file")
	sendersFlag := flag.String("senders", "100,150,200", "comma-separated sender counts N")
	flag.IntVar(&o.rp.subchannels, "subchannels", defaultRadio.subchannels, "certaware: subchannels per subframe")
	flag.IntVar(&o.rp.digestSub, "digest-subch", defaultRadio.digestSub, "certaware: subchannels of a digest-form message")
	flag.IntVar(&o.rp.certSub, "cert-subch", defaultRadio.certSub, "certaware: subchannels of a certificate-carrying message")
	flag.IntVar(&o.rp.window, "sched-window-ms", defaultRadio.window, "certaware: a message not scheduled within this many 1-ms subframes of its generation is counted as not scheduled")
	flag.Parse()

	for _, f := range strings.Split(*sendersFlag, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil || n <= 0 {
			fmt.Fprintln(os.Stderr, "invalid -senders value:", f)
			os.Exit(1)
		}
		o.counts = append(o.counts, n)
	}
	rep, ext, cs, err := run(&o)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(rep)
	if *out != "" {
		if err := os.WriteFile(*out, []byte(rep), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if *extOut != "" {
		if err := os.WriteFile(*extOut, []byte(ext), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if *csvOut != "" {
		f, err := os.Create(*csvOut)
		if err == nil {
			err = writeCSV(f, cs, o)
			if cerr := f.Close(); err == nil {
				err = cerr
			}
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
