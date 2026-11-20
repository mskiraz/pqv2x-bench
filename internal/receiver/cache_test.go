package receiver_test

import (
	"testing"

	"pqv2xbench/internal/fleet"
	"pqv2xbench/internal/receiver"
	"pqv2xbench/internal/sigs"
)

// TestColdStartCache follows one receiver from an empty cache: digests from
// a sender are rejected until that sender's certificate has been received
// and validated, after which they verify; other senders stay unknown; an
// invalid certificate is not cached and does not unlock its digests.
func TestColdStartCache(t *testing.T) {
	for _, sc := range sigs.All() {
		if testing.Short() && sc.Alg() == sigs.SLHDSA128s {
			continue
		}
		t.Run(sc.Name(), func(t *testing.T) {
			f, err := fleet.New(sc, 2, 2)
			if err != nil {
				t.Fatal(err)
			}
			byS := func(ms []fleet.Message, s int) []fleet.Message {
				var out []fleet.Message
				for _, m := range ms {
					if m.Sender == s {
						out = append(out, m)
					}
				}
				return out
			}
			d0, c0 := byS(f.Digest, 0), byS(f.Carry, 0)
			d1, c1 := byS(f.Digest, 1), byS(f.Carry, 1)

			rx := receiver.New(fleet.NowMicros, fleet.MaxAgeMicros)
			if err := rx.Trust(f.Issuer); err != nil {
				t.Fatal(err)
			}
			// Empty cache apart from the trust anchor.
			if rx.Len() != 1 {
				t.Fatalf("cold cache holds %d entries, want 1 (issuer)", rx.Len())
			}
			// Unknown digest rejected before any certificate is received.
			for _, m := range d0 {
				if o, p := rx.Process(m.Raw); o != receiver.UnknownSigner || p != receiver.PathNone {
					t.Fatalf("digest before certificate: got %d/%d", o, p)
				}
			}
			if rx.Len() != 1 {
				t.Fatal("rejected digest changed the cache")
			}

			// A tampered certificate (issuer signature broken) is refused and
			// not cached; sender 1 stays unknown.
			bad := append([]byte(nil), c1[0].Raw...)
			bad[len(bad)-sc.SignatureSize()-sc.SignatureSize()/2-4] ^= 0x01 // middle of the certificate's issuer signature
			if o, _ := rx.Process(bad); o == receiver.Accepted {
				t.Fatal("tampered certificate accepted")
			}
			if rx.Len() != 1 {
				t.Fatal("invalid certificate was cached")
			}
			if o, _ := rx.Process(d1[0].Raw); o != receiver.UnknownSigner {
				t.Fatalf("digest after invalid certificate: got %d", o)
			}

			// First valid certificate from sender 0 is validated and cached.
			if o, p := rx.Process(c0[0].Raw); o != receiver.Accepted || p != receiver.PathCertNew {
				t.Fatalf("first certificate: got %d/%d", o, p)
			}
			if rx.Len() != 2 {
				t.Fatalf("cache holds %d entries after first certificate, want 2", rx.Len())
			}
			// Its digests now verify; sender 1's still do not.
			for _, m := range d0 {
				if o, p := rx.Process(m.Raw); o != receiver.Accepted || p != receiver.PathDigestHit {
					t.Fatalf("digest after certificate: got %d/%d", o, p)
				}
			}
			if o, _ := rx.Process(d1[0].Raw); o != receiver.UnknownSigner {
				t.Fatalf("other sender's digest: got %d", o)
			}
			// A repeated certificate takes the cached path.
			if o, p := rx.Process(c0[1].Raw); o != receiver.Accepted || p != receiver.PathCertKnown {
				t.Fatalf("repeated certificate: got %d/%d", o, p)
			}
			// Sender 1's valid certificate unlocks sender 1.
			if o, p := rx.Process(c1[0].Raw); o != receiver.Accepted || p != receiver.PathCertNew {
				t.Fatalf("second sender's certificate: got %d/%d", o, p)
			}
			if o, p := rx.Process(d1[0].Raw); o != receiver.Accepted || p != receiver.PathDigestHit {
				t.Fatalf("second sender's digest: got %d/%d", o, p)
			}

			// Forgetting a certificate returns its sender to unknown.
			rx.Forget(f.Senders[0].ID)
			if o, _ := rx.Process(d0[0].Raw); o != receiver.UnknownSigner {
				t.Fatalf("digest after Forget: got %d", o)
			}
		})
	}
}

// TestUnrecognisedCertificateAlgorithm: a carried certificate whose key uses
// an algorithm the receiver does not support is not accepted and not cached.
func TestUnrecognisedCertificateAlgorithm(t *testing.T) {
	f, err := fleet.New(sigs.All()[0], 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	rx := receiver.New(fleet.NowMicros, fleet.MaxAgeMicros)
	if err := rx.Trust(f.Issuer); err != nil {
		t.Fatal(err)
	}
	raw := append([]byte(nil), f.Carry[0].Raw...)
	// The verification-key CHOICE of an ECDSA certificate is "0x80 0x80
	// 0x82|0x83" followed by the 32-byte x coordinate; turn the curve
	// alternative into an unsupported index.
	found := false
	for i := 0; i+2 < len(raw); i++ {
		if raw[i] == 0x80 && raw[i+1] == 0x80 && (raw[i+2] == 0x82 || raw[i+2] == 0x83) {
			raw[i+1] = 0x84 // ecdsaBrainpoolP256r1, not supported here
			found = true
			break
		}
	}
	if !found {
		t.Fatal("verification key not found in encoding")
	}
	if o, _ := rx.Process(raw); o == receiver.Accepted {
		t.Fatal("certificate with an unsupported key algorithm accepted")
	}
	if rx.Len() != 1 {
		t.Fatal("unsupported certificate was cached")
	}
}
