package fleet_test

import (
	"testing"

	"pqv2xbench/internal/fleet"
	"pqv2xbench/internal/sigs"
)

// The encoded sizes published with the paper (39-byte payload) must not
// change.
func TestPublishedSizes(t *testing.T) {
	want := map[sigs.Alg][3]int{ // certificate, digest SPDU, certificate SPDU
		sigs.ECDSAP256:  {144, 132, 270},
		sigs.Falcon512:  {1615, 736, 2345},
		sigs.MLDSA44:    {3784, 2490, 6268},
		sigs.SLHDSA128s: {7938, 7926, 15858},
	}
	for _, sc := range sigs.All() {
		if testing.Short() && sc.Alg() == sigs.SLHDSA128s {
			continue
		}
		f, err := fleet.New(sc, 1, 1)
		if err != nil {
			t.Fatal(err)
		}
		got := [3]int{len(f.Senders[0].Cert.Raw), len(f.Digest[0].Raw), len(f.Carry[0].Raw)}
		if got != want[sc.Alg()] {
			t.Errorf("%s: sizes %v, want %v", sc.Name(), got, want[sc.Alg()])
		}
	}
}

// Boundary payload lengths against the 2481-byte single-subframe limit. The
// payload's COER length determinant grows from one to two bytes at 128 and
// to three bytes at 256, so the boundaries are not offset by a constant.
func TestPayloadBoundaries(t *testing.T) {
	cases := []struct {
		alg     sigs.Alg
		carry   bool
		fit     int // largest payload that fits
		fitSize int
	}{
		{sigs.Falcon512, true, 174, 2481},
		{sigs.Falcon512, false, 1782, 2481},
		{sigs.MLDSA44, false, 30, 2481},
		{sigs.ECDSAP256, true, 2248, 2481},
	}
	for _, c := range cases {
		sc, _ := sigs.ByAlg(c.alg)
		f, err := fleet.NewWithPayload(sc, 1, 0, fleet.BSMPayloadSize)
		if err != nil {
			t.Fatal(err)
		}
		at := func(p int) int {
			m, err := f.SignMessage(0, make([]byte, p), c.carry)
			if err != nil {
				t.Fatal(err)
			}
			return len(m.Raw)
		}
		if got := at(c.fit); got != c.fitSize {
			t.Errorf("%s carry=%t payload %d: %d bytes, want %d", sc.Name(), c.carry, c.fit, got, c.fitSize)
		}
		if got := at(c.fit + 1); got <= 2481 {
			t.Errorf("%s carry=%t payload %d: %d bytes, want more than 2481", sc.Name(), c.carry, c.fit+1, got)
		}
	}
	// Length-determinant steps: 127 -> 128 adds two bytes, 255 -> 256 adds two.
	sc, _ := sigs.ByAlg(sigs.Falcon512)
	f, err := fleet.NewWithPayload(sc, 1, 0, fleet.BSMPayloadSize)
	if err != nil {
		t.Fatal(err)
	}
	size := func(p int) int {
		m, err := f.SignMessage(0, make([]byte, p), false)
		if err != nil {
			t.Fatal(err)
		}
		return len(m.Raw)
	}
	if d := size(128) - size(127); d != 2 {
		t.Errorf("127->128 adds %d bytes, want 2", d)
	}
	if d := size(256) - size(255); d != 2 {
		t.Errorf("255->256 adds %d bytes, want 2", d)
	}
	if size(39) != 736 {
		t.Errorf("SignMessage at 39 bytes gives %d, want 736", size(39))
	}
}
