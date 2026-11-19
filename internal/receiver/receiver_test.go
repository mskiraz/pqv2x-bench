package receiver_test

import (
	"bytes"
	"testing"

	"pqv2xbench/internal/fleet"
	"pqv2xbench/internal/ieee1609"
	"pqv2xbench/internal/receiver"
	"pqv2xbench/internal/sigs"
)

func build(t *testing.T, sc sigs.Scheme) *fleet.Fleet {
	t.Helper()
	f, err := fleet.New(sc, 2, 2)
	if err != nil {
		t.Fatalf("%s: %v", sc.Name(), err)
	}
	return f
}

func TestRoundTripAndPaths(t *testing.T) {
	for _, sc := range sigs.All() {
		if testing.Short() && sc.Alg() == sigs.SLHDSA128s {
			continue
		}
		t.Run(sc.Name(), func(t *testing.T) {
			f := build(t, sc)

			// Encoding round trip.
			c, err := ieee1609.DecodeCertificate(f.Senders[0].Cert.Raw)
			if err != nil {
				t.Fatalf("decode certificate: %v", err)
			}
			if !bytes.Equal(c.VerifyKey, f.Senders[0].Cert.VerifyKey) ||
				!bytes.Equal(c.Signature, f.Senders[0].Cert.Signature) ||
				!bytes.Equal(c.TBS, f.Senders[0].Cert.TBS) {
				t.Fatal("certificate round trip mismatch")
			}
			if len(c.VerifyKey) != sc.PublicKeySize() || len(c.Signature) != sc.SignatureSize() {
				t.Fatalf("unexpected key/signature sizes %d/%d", len(c.VerifyKey), len(c.Signature))
			}
			s, err := ieee1609.DecodeSPDU(f.Carry[0].Raw)
			if err != nil || s.Cert == nil || len(s.Payload) != fleet.BSMPayloadSize {
				t.Fatalf("decode SPDU: %v", err)
			}

			rx := receiver.New(fleet.NowMicros, fleet.MaxAgeMicros)
			if err := rx.Trust(f.Issuer); err != nil {
				t.Fatal(err)
			}

			// Digest from an unknown signer cannot be verified.
			if o, _ := rx.Process(f.Digest[0].Raw); o != receiver.UnknownSigner {
				t.Fatalf("unknown signer: got %d", o)
			}
			// First certificate-carrying message validates the certificate.
			if o, p := rx.Process(f.Carry[0].Raw); o != receiver.Accepted || p != receiver.PathCertNew {
				t.Fatalf("new certificate: got %d/%d", o, p)
			}
			// Second one finds it cached.
			if o, p := rx.Process(f.Carry[1].Raw); o != receiver.Accepted || p != receiver.PathCertKnown {
				t.Fatalf("known certificate: got %d/%d", o, p)
			}
			// Digest now resolves.
			if o, p := rx.Process(f.Digest[0].Raw); o != receiver.Accepted || p != receiver.PathDigestHit {
				t.Fatalf("digest hit: got %d/%d", o, p)
			}

			// A flipped payload bit must be rejected.
			bad := append([]byte(nil), f.Digest[1].Raw...)
			bad[10] ^= 0x01
			if o, _ := rx.Process(bad); o != receiver.BadSignature {
				t.Fatalf("tampered payload: got %d", o)
			}
			// A flipped bit in a carried, uncached certificate must be rejected.
			rx2 := receiver.New(fleet.NowMicros, fleet.MaxAgeMicros)
			_ = rx2.Trust(f.Issuer)
			badc := append([]byte(nil), f.Carry[0].Raw...)
			badc[len(badc)-sc.SignatureSize()-sc.SignatureSize()-40] ^= 0x01
			if o, _ := rx2.Process(badc); o == receiver.Accepted {
				t.Fatal("tampered certificate accepted")
			}
			// Stale message.
			rx.Now += 10 * fleet.MaxAgeMicros
			if o, _ := rx.Process(f.Digest[0].Raw); o != receiver.Stale {
				t.Fatalf("stale: got %d", o)
			}
		})
	}
}

// TestRejections checks that the receiver refuses messages it must not
// accept, and reports the reason it should.
func TestRejections(t *testing.T) {
	for _, sc := range sigs.All() {
		if testing.Short() && sc.Alg() == sigs.SLHDSA128s {
			continue
		}
		t.Run(sc.Name(), func(t *testing.T) {
			f := build(t, sc)
			other := build(t, sc) // a second, unrelated issuing authority

			warm := func() *receiver.Receiver {
				rx := receiver.New(fleet.NowMicros, fleet.MaxAgeMicros)
				if err := rx.Trust(f.Issuer); err != nil {
					t.Fatal(err)
				}
				for _, s := range f.Senders {
					if err := rx.Trust(s.Cert); err != nil {
						t.Fatal(err)
					}
				}
				return rx
			}
			digest, carry := f.Digest[0].Raw, f.Carry[0].Raw

			// Truncated and over-long encodings are malformed.
			if o, _ := warm().Process(digest[:len(digest)-1]); o != receiver.Malformed {
				t.Fatalf("truncated message: got %d", o)
			}
			if o, _ := warm().Process(append(append([]byte(nil), digest...), 0)); o != receiver.Malformed {
				t.Fatalf("trailing byte: got %d", o)
			}
			if o, _ := warm().Process(nil); o != receiver.Malformed {
				t.Fatalf("empty message: got %d", o)
			}

			// A certificate issued by an authority the receiver does not know.
			if o, _ := warm().Process(other.Carry[0].Raw); o != receiver.UnknownIssuer {
				t.Fatalf("unknown issuer: got %d", o)
			}
			// A digest that refers to that authority's certificate.
			if o, _ := warm().Process(other.Digest[0].Raw); o != receiver.UnknownSigner {
				t.Fatalf("foreign digest: got %d", o)
			}

			// A message whose PSID differs from the certificate's permission.
			// Byte 48 is the PSID value: 3 outer bytes, 4 payload header bytes,
			// 39 payload bytes, header preamble, PSID length.
			psid := append([]byte(nil), digest...)
			if psid[48] != 0x20 {
				t.Fatalf("unexpected PSID byte %#x", psid[48])
			}
			psid[48] = 0x21
			if o, _ := warm().Process(psid); o != receiver.BadCertificate {
				t.Fatalf("PSID mismatch: got %d", o)
			}

			// A flipped signature bit.
			sig := append([]byte(nil), digest...)
			sig[len(sig)-5] ^= 0x10
			if o, _ := warm().Process(sig); o != receiver.BadSignature {
				t.Fatalf("tampered signature: got %d", o)
			}

			// A valid message replayed with another sender's digest.
			swap := append([]byte(nil), digest...)
			var otherID [8]byte
			for _, s := range f.Senders {
				if s.ID != f.Senders[f.Digest[0].Sender].ID {
					otherID = s.ID
				}
			}
			// 3 outer bytes, 54 bytes of to-be-signed data, 1 signer choice byte.
			if swap[57] != 0x80 {
				t.Fatalf("unexpected signer choice byte %#x", swap[57])
			}
			copy(swap[58:66], otherID[:])
			if o, _ := warm().Process(swap); o != receiver.BadSignature {
				t.Fatalf("swapped signer digest: got %d", o)
			}

			// Certificates that have expired are refused, whether cached or
			// carried. The acceptance window is widened so that the message
			// itself is not rejected as stale first.
			late := uint64(fleet.NowMicros) + uint64(fleet.CertHours+2)*3600*1_000_000
			huge := uint64(1) << 62
			expired := receiver.New(late, huge)
			_ = expired.Trust(f.Issuer)
			for _, s := range f.Senders {
				_ = expired.Trust(s.Cert)
			}
			if o, _ := expired.Process(digest); o != receiver.BadCertificate {
				t.Fatalf("expired cached certificate: got %d", o)
			}
			fresh := receiver.New(late, huge)
			_ = fresh.Trust(f.Issuer)
			if o, _ := fresh.Process(carry); o != receiver.BadCertificate {
				t.Fatalf("expired carried certificate: got %d", o)
			}

			// After all rejections, a valid message is still accepted.
			if o, _ := warm().Process(digest); o != receiver.Accepted {
				t.Fatalf("valid message after rejections: got %d", o)
			}
		})
	}
}
