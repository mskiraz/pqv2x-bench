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
