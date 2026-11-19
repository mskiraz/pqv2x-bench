// Command bench measures receiver-side authentication of IEEE 1609.2 signed
// messages for ECDSA P-256, Falcon-512, ML-DSA-44 and SLH-DSA-SHA2-128s.
//
// For each scheme it reports the encoded certificate and message sizes and
// the per-message time of six operations:
//
//	primitive       signature verification only, key already parsed
//	decode          COER decoding of a digest-form message
//	digest_hit      full path, signer given as HashedId8, certificate cached
//	cert_known      full path, certificate carried and already cached
//	cert_new        full path, certificate carried and validated against
//	                its issuer (one extra signature verification)
//	digest_unknown  digest-form message whose certificate is not cached
//	                (decoded, then discarded)
//
// Per-operation samples are written to <out>/raw so that cmd/queuesim can
// replay them.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"pqv2xbench/internal/fleet"
	"pqv2xbench/internal/hrtime"
	"pqv2xbench/internal/ieee1609"
	"pqv2xbench/internal/receiver"
	"pqv2xbench/internal/sigs"
)

// Stats summarises one operation; times are in microseconds.
type Stats struct {
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

// Sizes are encoded sizes in bytes.
type Sizes struct {
	PublicKey   int `json:"public_key"`
	Signature   int `json:"signature"`
	Certificate int `json:"certificate"`
	DigestSPDU  int `json:"digest_spdu"`
	CertSPDU    int `json:"certificate_spdu"`
}

type SchemeResult struct {
	Name  string           `json:"name"`
	Sizes Sizes            `json:"sizes"`
	Ops   map[string]Stats `json:"ops"`
}

type Report struct {
	Tag       string         `json:"tag"`
	Time      string         `json:"time_utc"`
	GoVersion string         `json:"go_version"`
	GOOS      string         `json:"goos"`
	GOARCH    string         `json:"goarch"`
	CPU       string         `json:"cpu"`
	NumCPU    int            `json:"num_cpu"`
	Iters     int            `json:"iterations"`
	Schemes   []SchemeResult `json:"schemes"`
}

// smallBatch is the batch size for sub-microsecond operations.
const smallBatch = 64

var opOrder = []string{"primitive", "decode", "digest_hit", "cert_known", "cert_new", "digest_unknown"}

func cpuName() string {
	if b, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, l := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(l, "model name") || strings.HasPrefix(l, "Model") {
				if i := strings.Index(l, ":"); i >= 0 {
					return strings.TrimSpace(l[i+1:])
				}
			}
		}
	}
	if v := os.Getenv("PROCESSOR_IDENTIFIER"); v != "" {
		return v
	}
	return "unknown"
}

func summarize(ns []int64) Stats {
	s := append([]int64(nil), ns...)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	var sum, sq float64
	for _, v := range s {
		sum += float64(v)
	}
	mean := sum / float64(len(s))
	for _, v := range s {
		d := float64(v) - mean
		sq += d * d
	}
	q := func(p float64) float64 {
		i := int(math.Ceil(p*float64(len(s)))) - 1
		if i < 0 {
			i = 0
		}
		return float64(s[i]) / 1e3
	}
	return Stats{
		N: len(s), Mean: mean / 1e3, Std: math.Sqrt(sq/float64(len(s))) / 1e3,
		Min: float64(s[0]) / 1e3, P50: q(0.50), P90: q(0.90), P99: q(0.99), P999: q(0.999),
		Max: float64(s[len(s)-1]) / 1e3,
	}
}

func writeRaw(dir, tag, scheme, op string, ns []int64) error {
	var b strings.Builder
	for _, v := range ns {
		fmt.Fprintf(&b, "%d\n", v)
	}
	name := fmt.Sprintf("%s_%s_%s.txt", tag, slug(scheme), op)
	return os.WriteFile(filepath.Join(dir, name), []byte(b.String()), 0o644)
}

func slug(s string) string {
	s = strings.ToLower(s)
	return strings.NewReplacer(" ", "", "-", "").Replace(s)
}

// measure times op after warm untimed-result iterations and returns iters
// samples in nanoseconds per operation. Each sample times batch consecutive
// calls and divides by batch; batch > 1 is used for operations that are
// shorter than the clock resolution. prep runs untimed before each call when
// batch is 1.
func measure(iters, warm, batch int, prep func(i int), op func(i int) bool) []int64 {
	for i := 0; i < warm; i++ {
		if prep != nil {
			prep(i)
		}
		if !op(i) {
			panic("benchmark operation returned an unexpected outcome during warm-up")
		}
	}
	runtime.GC()
	out := make([]int64, iters)
	k := 0
	for i := 0; i < iters; i++ {
		if prep != nil {
			prep(k)
		}
		ok := true
		t0 := hrtime.Now()
		for j := 0; j < batch; j++ {
			ok = op(k) && ok
			k++
		}
		out[i] = (hrtime.Now() - t0) / int64(batch)
		if !ok {
			panic("benchmark operation returned an unexpected outcome")
		}
	}
	return out
}

func runScheme(sc sigs.Scheme, senders, per, iters, warm int, rawDir, tag string) (SchemeResult, error) {
	f, err := fleet.New(sc, senders, per)
	if err != nil {
		return SchemeResult{}, err
	}
	res := SchemeResult{Name: sc.Name(), Ops: map[string]Stats{}}
	res.Sizes = Sizes{
		PublicKey: sc.PublicKeySize(), Signature: sc.SignatureSize(),
		Certificate: len(f.Senders[0].Cert.Raw),
		DigestSPDU:  len(f.Digest[0].Raw), CertSPDU: len(f.Carry[0].Raw),
	}

	// Parsed keys for the primitive-only measurement.
	keys := make([]sigs.PublicKey, len(f.Senders))
	for i, s := range f.Senders {
		if keys[i], err = sc.ParsePublicKey(s.Cert.VerifyKey); err != nil {
			return res, err
		}
	}
	// Receiver with every sender certificate cached.
	warmRx := receiver.New(fleet.NowMicros, fleet.MaxAgeMicros)
	if err := warmRx.Trust(f.Issuer); err != nil {
		return res, err
	}
	for _, s := range f.Senders {
		if err := warmRx.Trust(s.Cert); err != nil {
			return res, err
		}
	}
	// Receiver that knows only the issuing authority.
	coldRx := receiver.New(fleet.NowMicros, fleet.MaxAgeMicros)
	if err := coldRx.Trust(f.Issuer); err != nil {
		return res, err
	}

	nd, nc := len(f.Digest), len(f.Carry)
	samples := map[string][]int64{}

	samples["primitive"] = measure(iters, warm, 1, nil, func(i int) bool {
		m := &f.Digest[i%nd]
		return keys[m.Sender].Verify(m.SignInput[:], m.Signature)
	})
	samples["decode"] = measure(iters, warm, smallBatch, nil, func(i int) bool {
		_, err := ieee1609.DecodeSPDU(f.Digest[i%nd].Raw)
		return err == nil
	})
	samples["digest_hit"] = measure(iters, warm, 1, nil, func(i int) bool {
		o, p := warmRx.Process(f.Digest[i%nd].Raw)
		return o == receiver.Accepted && p == receiver.PathDigestHit
	})
	samples["cert_known"] = measure(iters, warm, 1, nil, func(i int) bool {
		o, p := warmRx.Process(f.Carry[i%nc].Raw)
		return o == receiver.Accepted && p == receiver.PathCertKnown
	})
	samples["cert_new"] = measure(iters, warm, 1,
		func(i int) { coldRx.Forget(f.Senders[f.Carry[i%nc].Sender].ID) },
		func(i int) bool {
			o, p := coldRx.Process(f.Carry[i%nc].Raw)
			return o == receiver.Accepted && p == receiver.PathCertNew
		})
	// The cold receiver now holds every sender certificate; rebuild it so
	// that digest-form messages find no certificate.
	coldRx = receiver.New(fleet.NowMicros, fleet.MaxAgeMicros)
	if err := coldRx.Trust(f.Issuer); err != nil {
		return res, err
	}
	samples["digest_unknown"] = measure(iters, warm, smallBatch, nil, func(i int) bool {
		o, _ := coldRx.Process(f.Digest[i%nd].Raw)
		return o == receiver.UnknownSigner
	})

	prim := summarize(samples["primitive"])
	for _, op := range opOrder {
		st := summarize(samples[op])
		st.Ratio = st.Mean / prim.Mean
		st.RatioM = st.P50 / prim.P50
		res.Ops[op] = st
		if err := writeRaw(rawDir, tag, sc.Name(), op, samples[op]); err != nil {
			return res, err
		}
	}
	return res, nil
}

func markdown(r Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Receiver authentication benchmark (%s)\n\n", r.Tag)
	fmt.Fprintf(&b, "- CPU: %s (%d logical CPUs), %s/%s, %s\n", r.CPU, r.NumCPU, r.GOOS, r.GOARCH, r.GoVersion)
	fmt.Fprintf(&b, "- %d timed iterations per operation, single goroutine, garbage collector enabled\n- Run: %s\n\n", r.Iters, r.Time)
	b.WriteString("## Encoded sizes (bytes)\n\n| Scheme | Public key | Signature | Certificate | Digest SPDU | Certificate SPDU |\n|---|---:|---:|---:|---:|---:|\n")
	for _, s := range r.Schemes {
		z := s.Sizes
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d | %d |\n", s.Name, z.PublicKey, z.Signature, z.Certificate, z.DigestSPDU, z.CertSPDU)
	}
	b.WriteString("\n## Time per message (microseconds)\n\n| Scheme | Operation | Mean | Std | Median | p99 | p99.9 | Max | Mean / primitive | Median / primitive |\n|---|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, s := range r.Schemes {
		for _, op := range opOrder {
			t := s.Ops[op]
			fmt.Fprintf(&b, "| %s | %s | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f | %.3f | %.3f |\n",
				s.Name, op, t.Mean, t.Std, t.P50, t.P99, t.P999, t.Max, t.Ratio, t.RatioM)
		}
	}
	return b.String()
}

func main() {
	iters := flag.Int("iters", 20000, "timed iterations per operation")
	slhIters := flag.Int("slh-iters", 2000, "timed iterations per operation for SLH-DSA")
	warm := flag.Int("warmup", 1000, "untimed warm-up iterations per operation")
	senders := flag.Int("senders", 50, "pseudonym certificates per scheme")
	per := flag.Int("per", 4, "messages of each form per sender")
	out := flag.String("out", "results", "output directory")
	tag := flag.String("tag", runtime.GOOS+"-"+runtime.GOARCH, "label for this run")
	only := flag.String("schemes", "", "comma-separated subset: ecdsa,falcon,mldsa,slhdsa")
	flag.Parse()

	runtime.LockOSThread()
	rawDir := filepath.Join(*out, "raw")
	if err := os.MkdirAll(rawDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rep := Report{
		Tag: *tag, Time: time.Now().UTC().Format(time.RFC3339), GoVersion: runtime.Version(),
		GOOS: runtime.GOOS, GOARCH: runtime.GOARCH, CPU: cpuName(), NumCPU: runtime.NumCPU(), Iters: *iters,
	}
	want := func(sc sigs.Scheme) bool {
		if *only == "" {
			return true
		}
		for _, k := range strings.Split(*only, ",") {
			if strings.Contains(slug(sc.Name()), strings.TrimSpace(k)) {
				return true
			}
		}
		return false
	}
	for _, sc := range sigs.All() {
		if !want(sc) {
			continue
		}
		n, ns, np, w := *iters, *senders, *per, *warm
		if sc.Alg() == sigs.SLHDSA128s { // signing is slow; keep the fleet small
			n, ns, np, w = *slhIters, 4, 2, 50
		}
		fmt.Fprintf(os.Stderr, "%s: generating fleet and measuring...\n", sc.Name())
		r, err := runScheme(sc, ns, np, n, w, rawDir, *tag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		rep.Schemes = append(rep.Schemes, r)
	}
	js, _ := json.MarshalIndent(rep, "", "  ")
	if err := os.WriteFile(filepath.Join(*out, "bench_"+*tag+".json"), js, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	md := markdown(rep)
	if err := os.WriteFile(filepath.Join(*out, "bench_"+*tag+".md"), []byte(md), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(md)
}
