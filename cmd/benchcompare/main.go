// Command benchcompare puts several cmd/bench runs side by side: their
// environments, every scheme's per-operation statistics, and the
// receiver-path-to-primitive ratios that the paper transfers to an ARMv8
// primitive time. It writes a long-form CSV and a Markdown summary.
//
// Usage:
//
//	benchcompare -ref <tag> -out <prefix> run1.json run2.json ...
//
// The output separates measured values (per platform) from platform-scaled
// estimates (a measured ratio multiplied by the published ARMv8 primitive
// time), which are not measurements.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
)

type stats struct {
	N      int     `json:"n"`
	Mean   float64 `json:"mean_us"`
	Std    float64 `json:"std_us"`
	Min    float64 `json:"min_us"`
	P50    float64 `json:"p50_us"`
	P90    float64 `json:"p90_us"`
	P99    float64 `json:"p99_us"`
	P999   float64 `json:"p999_us"`
	Max    float64 `json:"max_us"`
	Ratio  float64 `json:"ratio_to_primitive_mean"`
	RatioM float64 `json:"ratio_to_primitive_median"`
}

type scheme struct {
	Name       string           `json:"name"`
	Iterations int              `json:"iterations"`
	Ops        map[string]stats `json:"ops"`
}

type env struct {
	CPU         string            `json:"cpu"`
	CPUFeatures []string          `json:"cpu_features"`
	OSVersion   string            `json:"os_version"`
	Kernel      string            `json:"kernel"`
	GOMAXPROCS  int               `json:"gomaxprocs"`
	GCPercent   int               `json:"gc_percent"`
	BuildSet    map[string]string `json:"build_settings"`
}

type report struct {
	Tag       string   `json:"tag"`
	Time      string   `json:"time_utc"`
	GoVersion string   `json:"go_version"`
	GOOS      string   `json:"goos"`
	GOARCH    string   `json:"goarch"`
	CPU       string   `json:"cpu"`
	NumCPU    int      `json:"num_cpu"`
	Iters     int      `json:"iterations"`
	GCMode    string   `json:"gc_mode"`
	Env       *env     `json:"environment"`
	Schemes   []scheme `json:"schemes"`
	file      string
}

var ops = []string{"primitive", "decode", "digest_hit", "cert_known", "cert_new", "digest_unknown"}

func (r *report) gc() string {
	if r.GCMode == "" {
		return "enabled (not recorded; original runs)"
	}
	return r.GCMode
}

func (r *report) find(name string) *scheme {
	for i := range r.Schemes {
		if r.Schemes[i].Name == name {
			return &r.Schemes[i]
		}
	}
	return nil
}

func main() {
	ref := flag.String("ref", "", "tag of the reference run for differences")
	out := flag.String("out", "arm_benchmark_comparison", "output path prefix (.csv and .md are appended)")
	armPrim := flag.Float64("arm-primitive-ms", 0.446, "published ARMv8 Falcon-512 primitive time used for platform-scaled estimates")
	flag.Parse()

	var rs []*report
	for _, f := range flag.Args() {
		b, err := os.ReadFile(f)
		if err != nil {
			fail(err)
		}
		r := &report{file: f}
		if err := json.Unmarshal(b, r); err != nil {
			fail(fmt.Errorf("%s: %w", f, err))
		}
		rs = append(rs, r)
	}
	if len(rs) == 0 {
		fail(fmt.Errorf("no input files"))
	}
	var refR *report
	for _, r := range rs {
		if r.Tag == *ref {
			refR = r
		}
	}
	if refR == nil {
		refR = rs[0]
	}

	// CSV: one row per run, scheme and operation.
	var c strings.Builder
	c.WriteString("tag,source_file,goos,goarch,cpu,go_version,gc_mode,scheme,operation,n,mean_us,std_us,min_us,p50_us,p90_us,p99_us,p999_us,max_us,ratio_to_primitive_mean,ratio_to_primitive_median,ratio_median_minus_reference\n")
	for _, r := range rs {
		for _, s := range r.Schemes {
			rs0 := refR.find(s.Name)
			for _, op := range ops {
				t, ok := s.Ops[op]
				if !ok {
					continue
				}
				diff := ""
				if rs0 != nil {
					if t0, ok := rs0.Ops[op]; ok {
						diff = fmt.Sprintf("%.5f", t.RatioM-t0.RatioM)
					}
				}
				fmt.Fprintf(&c, "%s,%s,%s,%s,%q,%s,%q,%s,%s,%d,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.5f,%.5f,%s\n",
					r.Tag, r.file, r.GOOS, r.GOARCH, r.CPU, r.GoVersion, r.gc(), s.Name, op, t.N,
					t.Mean, t.Std, t.Min, t.P50, t.P90, t.P99, t.P999, t.Max, t.Ratio, t.RatioM, diff)
			}
		}
	}
	if err := os.WriteFile(*out+".csv", []byte(c.String()), 0o644); err != nil {
		fail(err)
	}

	// Markdown.
	var m strings.Builder
	m.WriteString("# Receiver-path benchmark across platforms\n\n")
	m.WriteString("All times are **measurements** on the named platform (microseconds). Ratios are per-platform measurements of a path's time divided by the same platform's primitive verification time. Absolute times from different machines, operating systems or garbage-collector settings are not directly comparable; ratios are compared only to test whether they transfer between platforms.\n\n")
	m.WriteString("## Runs\n\n| Tag | Platform | CPU | OS | Go | GC | GOMAXPROCS | CPU features | Iterations (SLH-DSA) |\n|---|---|---|---|---|---|---:|---|---|\n")
	for _, r := range rs {
		osv, feat, gmp := "not recorded", "not recorded", "not recorded"
		if r.Env != nil {
			osv = r.Env.OSVersion
			if r.Env.Kernel != "" {
				osv += " (kernel " + r.Env.Kernel + ")"
			}
			feat = strings.Join(r.Env.CPUFeatures, ", ")
			gmp = fmt.Sprint(r.Env.GOMAXPROCS)
		} else if r.GOOS == "windows" {
			osv = "Windows (version not recorded)"
		}
		slh := "–"
		if s := r.find("SLH-DSA-SHA2-128s"); s != nil {
			slh = fmt.Sprint(s.Ops["primitive"].N)
		}
		fmt.Fprintf(&m, "| %s | %s/%s | %s | %s | %s | %s | %s | %s | %d (%s) |\n", r.Tag, r.GOOS, r.GOARCH, r.CPU, osv, r.GoVersion, r.gc(), gmp, feat, r.Iters, slh)
	}

	names := []string{}
	seen := map[string]bool{}
	for _, r := range rs {
		for _, s := range r.Schemes {
			if !seen[s.Name] {
				seen[s.Name] = true
				names = append(names, s.Name)
			}
		}
	}

	m.WriteString("\n## Falcon-512: median ratio to the primitive (measured per platform)\n\n")
	ratioTable(&m, rs, refR, "Falcon-512", true)
	m.WriteString("\n## Falcon-512: mean ratio to the primitive (measured per platform)\n\n")
	ratioTable(&m, rs, refR, "Falcon-512", false)

	m.WriteString("\n## Falcon-512 path times scaled to the published ARMv8 primitive (platform-scaled estimates, not measurements)\n\n")
	fmt.Fprintf(&m, "Each entry is a platform's measured median ratio multiplied by %.3f ms. The spread across rows is the uncertainty that the ratio-transfer assumption introduces.\n\n", *armPrim)
	m.WriteString("| Tag | digest_hit (ms) | cert_known (ms) | cert_new (ms) | max(2 x digest_hit, cert_new) (ms) |\n|---|---:|---:|---:|---:|\n")
	for _, r := range rs {
		s := r.find("Falcon-512")
		if s == nil {
			continue
		}
		dh, ck, cn := s.Ops["digest_hit"].RatioM**armPrim, s.Ops["cert_known"].RatioM**armPrim, s.Ops["cert_new"].RatioM**armPrim
		fmt.Fprintf(&m, "| %s | %.4f | %.4f | %.4f | %.4f |\n", r.Tag, dh, ck, cn, math.Max(2*dh, cn))
	}

	m.WriteString("\n## All schemes: median ratio to the primitive\n\n| Scheme | Tag | digest_hit | cert_known | cert_new |\n|---|---|---:|---:|---:|\n")
	for _, n := range names {
		for _, r := range rs {
			if s := r.find(n); s != nil {
				fmt.Fprintf(&m, "| %s | %s | %.3f | %.3f | %.3f |\n", n, r.Tag, s.Ops["digest_hit"].RatioM, s.Ops["cert_known"].RatioM, s.Ops["cert_new"].RatioM)
			}
		}
	}

	m.WriteString("\n## Primitive verification time per platform (measured, microseconds)\n\n| Tag |")
	for _, n := range names {
		fmt.Fprintf(&m, " %s median | %s mean |", n, n)
	}
	m.WriteString(" ML-DSA-44 / Falcon-512 (median) |\n|---|")
	for range names {
		m.WriteString("---:|---:|")
	}
	m.WriteString("---:|\n")
	for _, r := range rs {
		fmt.Fprintf(&m, "| %s |", r.Tag)
		for _, n := range names {
			if s := r.find(n); s != nil {
				fmt.Fprintf(&m, " %.2f | %.2f |", s.Ops["primitive"].P50, s.Ops["primitive"].Mean)
			} else {
				m.WriteString(" – | – |")
			}
		}
		f, d := r.find("Falcon-512"), r.find("ML-DSA-44")
		if f != nil && d != nil {
			fmt.Fprintf(&m, " %.2f |\n", d.Ops["primitive"].P50/f.Ops["primitive"].P50)
		} else {
			m.WriteString(" – |\n")
		}
	}

	m.WriteString("\n## Falcon-512 distribution per operation (measured, microseconds)\n\n| Tag | Operation | n | Mean | Std | Median | p90 | p99 | p99.9 | Max |\n|---|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, r := range rs {
		if s := r.find("Falcon-512"); s != nil {
			for _, op := range ops {
				t := s.Ops[op]
				fmt.Fprintf(&m, "| %s | %s | %d | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f |\n", r.Tag, op, t.N, t.Mean, t.Std, t.P50, t.P90, t.P99, t.P999, t.Max)
			}
		}
	}
	fmt.Fprintf(&m, "\nFull per-scheme statistics for every run: %s.csv.\n", baseName(*out))
	if err := os.WriteFile(*out+".md", []byte(m.String()), 0o644); err != nil {
		fail(err)
	}
	fmt.Print(m.String())
}

func ratioTable(m *strings.Builder, rs []*report, refR *report, name string, median bool) {
	paths := []string{"decode", "digest_hit", "cert_known", "cert_new", "digest_unknown"}
	m.WriteString("| Tag |")
	for _, p := range paths {
		fmt.Fprintf(m, " %s |", p)
	}
	fmt.Fprintf(m, " Δ cert_known vs %s | Δ cert_new vs %s |\n|---|", refR.Tag, refR.Tag)
	for range paths {
		m.WriteString("---:|")
	}
	m.WriteString("---:|---:|\n")
	get := func(s *scheme, op string) float64 {
		if median {
			return s.Ops[op].RatioM
		}
		return s.Ops[op].Ratio
	}
	r0 := refR.find(name)
	for _, r := range rs {
		s := r.find(name)
		if s == nil {
			continue
		}
		fmt.Fprintf(m, "| %s |", r.Tag)
		for _, p := range paths {
			fmt.Fprintf(m, " %.4f |", get(s, p))
		}
		if r0 != nil {
			fmt.Fprintf(m, " %+.4f | %+.4f |\n", get(s, "cert_known")-get(r0, "cert_known"), get(s, "cert_new")-get(r0, "cert_new"))
		} else {
			m.WriteString(" | |\n")
		}
	}
}

func baseName(p string) string {
	i := strings.LastIndexAny(p, `/\`)
	return p[i+1:]
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
