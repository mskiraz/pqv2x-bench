// Command queuecompare joins cmd/queuesim CSV outputs (-csv) and puts the
// legacy abstract C-V2X model ("cv2x") next to the certificate-aware model
// ("cv2x-cert") for the same service-time source, scenario and sender count.
//
// Usage:
//
//	queuecompare -out <prefix> legacy.csv certaware.csv ...
//
// Every row of every input is copied to <prefix>.csv; <prefix>.md holds the
// side-by-side summary. All values are simulation outputs.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type row map[string]string

func load(path string) ([]string, []row, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	recs, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	if len(recs) < 1 {
		return nil, nil, fmt.Errorf("%s: empty", path)
	}
	head := recs[0]
	var rs []row
	for _, rec := range recs[1:] {
		r := row{"source_file": path}
		for i, h := range head {
			r[h] = rec[i]
		}
		rs = append(rs, r)
	}
	return head, rs, nil
}

func num(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func key(r row) string {
	return r["service_mode"] + "|" + r["tag"] + "|" + r["scenario"] + "|" + r["n"]
}

func main() {
	out := flag.String("out", "queue_comparison", "output path prefix (.csv and .md are appended)")
	flag.Parse()
	var head []string
	var all []row
	for _, p := range flag.Args() {
		h, rs, err := load(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if head == nil {
			head = h
		}
		all = append(all, rs...)
	}
	if len(all) == 0 {
		fmt.Fprintln(os.Stderr, "no input rows")
		os.Exit(1)
	}

	// Combined CSV.
	var c strings.Builder
	w := csv.NewWriter(&c)
	_ = w.Write(append([]string{"source_file"}, head...))
	for _, r := range all {
		rec := []string{r["source_file"]}
		for _, h := range head {
			rec = append(rec, r[h])
		}
		_ = w.Write(rec)
	}
	w.Flush()
	if err := os.WriteFile(*out+".csv", []byte(c.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	legacy := map[string]row{}
	var order []string
	cert := map[string]row{}
	for _, r := range all {
		switch r["arrival_model"] {
		case "cv2x":
			k := key(r)
			if _, ok := legacy[k]; !ok {
				order = append(order, k)
			}
			legacy[k] = r
		case "cv2x-cert":
			cert[key(r)] = r
		}
	}

	var m strings.Builder
	m.WriteString("# Receiver queue: legacy abstract C-V2X model vs certificate-aware scheduling\n\n")
	m.WriteString("All values are **simulation outputs** (no radio propagation, collisions, half-duplex loss or retransmissions). Service times are either fixed multiples of the published ARMv8 Falcon-512 primitive time S = 0.446 ms (`fixed`) or measured path ratios scaled to S (`measured`, a platform-scaled estimate). Both models use the same random configurations (seed and draw order are shared).\n\n")
	m.WriteString("- **Legacy (`cv2x`)**: up to two messages per 1-ms subframe whatever their form; every generated message is delivered.\n")
	m.WriteString("- **Certificate-aware (`cv2x-cert`)**: digest-form message 4 of 10 subchannels, certificate-carrying message all 10 (never shares its subframe); idealised earliest-fit scheduling within 100 ms; messages that do not fit are counted as not scheduled.\n\n")
	m.WriteString("Denominators: Pr[miss | auth] = deadline missed / authenticated. On-time / generated = authenticated within the deadline / generated. Discarded = unknown-signer discards / generated.\n\n")
	m.WriteString("| Service times | Benchmark tag | Scenario | N | Legacy: max T (ms) | Legacy: Pr[miss \\| auth] | Legacy: on-time / generated | Cert-aware: offered occupancy | Cert-aware: not scheduled (digest / cert) | Cert-aware: discarded / generated | Cert-aware: max T (ms) | Cert-aware: p99.9 T (ms) | Cert-aware: Pr[miss \\| auth] | Cert-aware: on-time / generated | Cert-aware: max scheduling delay (ms) |\n")
	m.WriteString("|---|---|---|---:|---:|---:|---:|---:|---|---:|---:|---:|---:|---:|---:|\n")
	cell := func(r row, f string, unstableText string) string {
		if r == nil {
			return "–"
		}
		if r["receiver_unstable"] == "true" && unstableText != "" {
			return unstableText
		}
		if r[f] == "" {
			return "–"
		}
		return r[f]
	}
	for _, k := range order {
		l, ce := legacy[k], cert[k]
		lmax := cell(l, "max_ms", "unstable")
		lp := cell(l, "p_miss_given_authenticated", "unstable")
		lon := cell(l, "on_time_fraction_of_generated", "unstable")
		if l["receiver_unstable"] == "true" {
			lmax, lp, lon = "unstable", "–", "–"
		}
		occ, ns, disc, cmax, c999, cp, con, csd := "–", "–", "–", "–", "–", "–", "–", "–"
		if ce != nil {
			occ = ce["offered_occupancy"]
			ns = ce["not_scheduled_digest"] + " / " + ce["not_scheduled_cert"]
			disc = fmt.Sprintf("%.4f", num(ce["discarded_fraction_of_generated"]))
			cmax = cell(ce, "max_ms", "unstable")
			c999 = cell(ce, "p999_ms", "unstable")
			cp = cell(ce, "p_miss_given_authenticated", "unstable")
			con = cell(ce, "on_time_fraction_of_generated", "")
			csd = cell(ce, "max_sched_delay_ms", "")
		}
		fmt.Fprintf(&m, "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			l["service_mode"], l["tag"], l["scenario"], l["n"], lmax, lp, lon, occ, ns, disc, cmax, c999, cp, con, csd)
	}
	m.WriteString("\nFull accounting (generated, scheduled, discarded, authenticated, missed, utilisation, backlog) for every row: " + baseName(*out) + ".csv.\n")
	if err := os.WriteFile(*out+".md", []byte(m.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(m.String())
}

func baseName(p string) string {
	i := strings.LastIndexAny(p, `/\`)
	return p[i+1:]
}
