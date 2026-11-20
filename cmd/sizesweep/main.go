// Command sizesweep measures how the COER-encoded size of a signed IEEE 1609.2
// message depends on the length of its application payload, for every
// signature scheme, in digest form and in certificate-carrying form.
//
// Every size comes from the encoder in internal/ieee1609. For each scheme the
// program:
//
//  1. builds an issuing authority and one pseudonym certificate;
//  2. encodes the message for every payload length from 0 to the limit (2481
//     bytes by default) with a signature of the scheme's exact length, and
//     finds the largest payload whose message fits and the first that does
//     not (the encoder adds a second length byte once the payload reaches
//     128 bytes, so the search does not assume one byte per byte);
//  3. for each reported payload length (the sweep points, 0, the largest
//     fitting payload and the first exceeding one) signs a real message,
//     checks that its length equals the scanned size, and checks that the
//     receiver accepts it.
//
// Only the application payload varies. Certificate extensions, region and
// assurance fields, and other protocol options are not modelled.
package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"pqv2xbench/internal/fleet"
	"pqv2xbench/internal/ieee1609"
	"pqv2xbench/internal/receiver"
	"pqv2xbench/internal/sigs"
)

// Row is one encoded message.
type Row struct {
	Scheme    string `json:"scheme"`
	Format    string `json:"format"` // digest | certificate
	Payload   int    `json:"payload_bytes"`
	Size      int    `json:"spdu_bytes"`
	Remaining int    `json:"remaining_bytes"` // limit - size; negative when over
	Fits      bool   `json:"fits"`
	Signed    bool   `json:"signed_and_verified"`
}

// Summary is the per-scheme, per-format result of the exact search.
type Summary struct {
	Scheme        string `json:"scheme"`
	Format        string `json:"format"`
	Certificate   int    `json:"certificate_bytes"`
	SizeAt39      int    `json:"spdu_bytes_at_39"`
	SizeAt0       int    `json:"spdu_bytes_at_0"`
	MaxPayload    int    `json:"max_fitting_payload_bytes"` // -1: none fits
	SizeAtMax     int    `json:"spdu_bytes_at_max"`
	FirstExceeds  int    `json:"first_exceeding_payload_bytes"`
	SizeAtFirstEx int    `json:"spdu_bytes_at_first_exceeding"`
	Monotonic     bool   `json:"size_nondecreasing_over_scan"`
}

type Report struct {
	Limit   int       `json:"limit_bytes"`
	Note    string    `json:"note"`
	Points  []int     `json:"sweep_points"`
	Summary []Summary `json:"summary"`
	Rows    []Row     `json:"rows"`
}

func encodeSize(f *fleet.Fleet, payloadLen int, carry bool, sig []byte) int {
	s := f.Senders[0]
	m := &ieee1609.SPDU{
		Payload:   make([]byte, payloadLen),
		Psid:      ieee1609.PsidBSM,
		GenTime:   fleet.NowMicros,
		SigAlg:    f.Scheme.Alg(),
		Signature: sig,
	}
	if carry {
		m.Cert = s.Cert
	} else {
		m.Digest = s.ID
	}
	return len(m.Encode())
}

func signedSize(f *fleet.Fleet, payloadLen int, carry bool) (int, error) {
	p := make([]byte, payloadLen)
	if _, err := rand.Read(p); err != nil {
		return 0, err
	}
	msg, err := f.SignMessage(0, p, carry)
	if err != nil {
		return 0, err
	}
	rx := receiver.New(fleet.NowMicros, fleet.MaxAgeMicros)
	if err := rx.Trust(f.Issuer); err != nil {
		return 0, err
	}
	if !carry {
		if err := rx.Trust(f.Senders[0].Cert); err != nil {
			return 0, err
		}
	}
	if o, _ := rx.Process(msg.Raw); o != receiver.Accepted {
		return 0, fmt.Errorf("signed %d-byte-payload message not accepted (outcome %d)", payloadLen, o)
	}
	return len(msg.Raw), nil
}

func main() {
	out := flag.String("out", "results/sizesweep", "output directory")
	limit := flag.Int("limit", 2481, "single-subframe size limit in bytes")
	pts := flag.String("points", "39,64,100,128,160,170,174,175,176,200,256,400", "payload lengths to report")
	only := flag.String("schemes", "", "comma-separated subset: ecdsa,falcon,mldsa,slhdsa")
	flag.Parse()

	var points []int
	for _, s := range strings.Split(*pts, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || v < 0 {
			fmt.Fprintln(os.Stderr, "invalid -points value:", s)
			os.Exit(1)
		}
		points = append(points, v)
	}
	rep := Report{Limit: *limit, Points: points,
		Note: "Sizes are COER encoder output for the implemented IEEE 1609.2-2022 subset with hypothetical PQ extension additions; only the application payload varies; the issuing authority uses the same scheme as the pseudonym certificate; region, assurance level and other optional certificate fields are not present."}

	for _, sc := range sigs.All() {
		if *only != "" {
			keep := false
			for _, k := range strings.Split(*only, ",") {
				name := strings.NewReplacer(" ", "", "-", "").Replace(strings.ToLower(sc.Name()))
				if strings.Contains(name, strings.TrimSpace(k)) {
					keep = true
				}
			}
			if !keep {
				continue
			}
		}
		fmt.Fprintf(os.Stderr, "%s...\n", sc.Name())
		f, err := fleet.NewWithPayload(sc, 1, 0, fleet.BSMPayloadSize)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		ref, err := f.SignMessage(0, make([]byte, fleet.BSMPayloadSize), false)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, carry := range []bool{false, true} {
			format := "digest"
			if carry {
				format = "certificate"
			}
			// Exact scan over every payload length up to the limit.
			sizes := make([]int, *limit+2)
			mono := true
			for p := range sizes {
				sizes[p] = encodeSize(f, p, carry, ref.Signature)
				if p > 0 && sizes[p] < sizes[p-1] {
					mono = false
				}
			}
			sm := Summary{Scheme: sc.Name(), Format: format, Certificate: len(f.Senders[0].Cert.Raw),
				SizeAt39: sizes[fleet.BSMPayloadSize], SizeAt0: sizes[0], MaxPayload: -1, FirstExceeds: -1, Monotonic: mono}
			for p, s := range sizes {
				if s <= *limit {
					sm.MaxPayload, sm.SizeAtMax = p, s
				}
			}
			for p, s := range sizes {
				if s > *limit && p > sm.MaxPayload {
					sm.FirstExceeds, sm.SizeAtFirstEx = p, s
					break
				}
			}
			rep.Summary = append(rep.Summary, sm)

			// Real signed messages at every reported length.
			report := map[int]bool{0: true}
			for _, p := range points {
				report[p] = true
			}
			if sm.MaxPayload >= 0 {
				report[sm.MaxPayload] = true
			}
			if sm.FirstExceeds >= 0 {
				report[sm.FirstExceeds] = true
			}
			var lens []int
			for p := range report {
				lens = append(lens, p)
			}
			sort.Ints(lens)
			for _, p := range lens {
				scanned := encodeSize(f, p, carry, ref.Signature)
				got, err := signedSize(f, p, carry)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if got != scanned {
					fmt.Fprintf(os.Stderr, "%s %s payload %d: signed size %d differs from scanned size %d\n", sc.Name(), format, p, got, scanned)
					os.Exit(1)
				}
				rep.Rows = append(rep.Rows, Row{Scheme: sc.Name(), Format: format, Payload: p, Size: got,
					Remaining: *limit - got, Fits: got <= *limit, Signed: true})
			}
		}
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	js, _ := json.MarshalIndent(rep, "", "  ")
	write(filepath.Join(*out, "payload_sensitivity.json"), string(js))
	write(filepath.Join(*out, "payload_sensitivity.csv"), csv(rep))
	md := markdown(rep)
	write(filepath.Join(*out, "payload_sensitivity.md"), md)
	fmt.Print(md)
}

func write(path, s string) {
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func csv(r Report) string {
	var b strings.Builder
	b.WriteString("scheme,format,payload_bytes,spdu_bytes,limit_bytes,remaining_bytes,fits,signed_and_verified\n")
	for _, x := range r.Rows {
		fmt.Fprintf(&b, "%s,%s,%d,%d,%d,%d,%t,%t\n", x.Scheme, x.Format, x.Payload, x.Size, r.Limit, x.Remaining, x.Fits, x.Signed)
	}
	b.WriteString("\nscheme,format,certificate_bytes,spdu_bytes_at_0,spdu_bytes_at_39,max_fitting_payload_bytes,spdu_bytes_at_max,first_exceeding_payload_bytes,spdu_bytes_at_first_exceeding,size_nondecreasing\n")
	for _, s := range r.Summary {
		fmt.Fprintf(&b, "%s,%s,%d,%d,%d,%d,%d,%d,%d,%t\n", s.Scheme, s.Format, s.Certificate, s.SizeAt0, s.SizeAt39, s.MaxPayload, s.SizeAtMax, s.FirstExceeds, s.SizeAtFirstEx, s.Monotonic)
	}
	return b.String()
}

func markdown(r Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# COER payload-size sensitivity (limit %d B)\n\n%s\n\n", r.Limit, r.Note)
	b.WriteString("## Largest application payload that fits\n\n| Scheme | Format | Certificate (B) | SPDU at 39 B payload | Max payload (B) | SPDU at max | First exceeding payload (B) | SPDU there |\n|---|---|---:|---:|---:|---:|---:|---:|\n")
	for _, s := range r.Summary {
		maxP, atMax := "none", "–"
		if s.MaxPayload >= 0 {
			maxP, atMax = strconv.Itoa(s.MaxPayload), strconv.Itoa(s.SizeAtMax)
		}
		first, atFirst := strconv.Itoa(s.FirstExceeds), strconv.Itoa(s.SizeAtFirstEx)
		if s.MaxPayload < 0 {
			first, atFirst = "0", strconv.Itoa(s.SizeAt0)
		}
		fmt.Fprintf(&b, "| %s | %s | %d | %d | %s | %s | %s | %s |\n", s.Scheme, s.Format, s.Certificate, s.SizeAt39, maxP, atMax, first, atFirst)
	}
	// Sweep table: one column per scheme and format.
	type key struct {
		scheme, format string
		p              int
	}
	cell := map[key]Row{}
	var cols []key
	seen := map[[2]string]bool{}
	for _, x := range r.Rows {
		cell[key{x.Scheme, x.Format, x.Payload}] = x
		k := [2]string{x.Scheme, x.Format}
		if !seen[k] {
			seen[k] = true
			cols = append(cols, key{x.Scheme, x.Format, 0})
		}
	}
	b.WriteString("\n## Encoded SPDU size (bytes) by payload length; remaining bytes to the limit in parentheses, negative when over\n\n| Payload (B) |")
	for _, c := range cols {
		fmt.Fprintf(&b, " %s %s |", c.scheme, c.format)
	}
	b.WriteString("\n|---:|")
	for range cols {
		b.WriteString("---:|")
	}
	b.WriteString("\n")
	for _, p := range r.Points {
		fmt.Fprintf(&b, "| %d |", p)
		for _, c := range cols {
			x := cell[key{c.scheme, c.format, p}]
			fmt.Fprintf(&b, " %d (%+d) |", x.Size, x.Remaining)
		}
		b.WriteString("\n")
	}
	b.WriteString("\nEvery listed size was produced by signing a real message with that payload and checking that the receiver accepts it.\n")
	return b.String()
}
