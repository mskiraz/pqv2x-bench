// Package fleet builds the test material for the benchmark: one issuing
// authority, a set of pseudonym certificates it signs, and signed Basic
// Safety Messages in digest and certificate-carrying form.
package fleet

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"pqv2xbench/internal/ieee1609"
	"pqv2xbench/internal/sigs"
)

// BSMPayloadSize is the size of the SAE J2735 BSM Part I payload used by
// Abrar et al. (39 bytes).
const BSMPayloadSize = 39

// Epoch values used for all generated material (arbitrary but consistent).
const (
	CertStart    = 700_000_000 // Time32 seconds
	CertHours    = 168         // one-week pseudonym certificate
	NowMicros    = (CertStart + 3600) * 1_000_000
	MaxAgeMicros = 5_000_000
	PeriodMicros = 100_000 // 10 Hz
)

// Sender is one vehicle pseudonym.
type Sender struct {
	Cert     *ieee1609.Certificate
	CertHash [32]byte
	ID       [8]byte
	key      sigs.PrivateKey
}

// Message is one encoded SPDU with the fields a primitive-only benchmark
// needs.
type Message struct {
	Raw       []byte
	Sender    int
	SignInput [32]byte // digest the message signature covers
	Signature []byte
}

// Fleet is the generated material for one signature scheme.
type Fleet struct {
	Scheme  sigs.Scheme
	Issuer  *ieee1609.Certificate
	Senders []*Sender
	Digest  []Message // signer given as HashedId8
	Carry   []Message // signer certificate carried
}

func randBytes(b []byte) {
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}

func newCert(sc sigs.Scheme, issuerID [8]byte, pk []byte) *ieee1609.Certificate {
	c := &ieee1609.Certificate{
		IssuerID:      issuerID,
		CrlSeries:     1,
		Start:         CertStart,
		DurationHours: CertHours,
		Psid:          ieee1609.PsidBSM,
		KeyAlg:        sc.Alg(),
		VerifyKey:     pk,
		SigAlg:        sc.Alg(),
	}
	var ic [2]byte
	randBytes(ic[:])
	c.ICert = binary.BigEndian.Uint16(ic[:])
	randBytes(c.LinkageValue[:])
	randBytes(c.CracaID[:])
	return c
}

// New generates an issuer, nSenders pseudonym certificates, and perSender
// digest and certificate-carrying messages for each sender. The issuing
// authority signs with the same scheme as the pseudonyms. Messages carry a
// BSMPayloadSize-byte payload.
func New(sc sigs.Scheme, nSenders, perSender int) (*Fleet, error) {
	return NewWithPayload(sc, nSenders, perSender, BSMPayloadSize)
}

// NewWithPayload is New with a configurable application-payload length in
// bytes. It is used by the payload-size sensitivity analysis; the benchmark
// uses New.
func NewWithPayload(sc sigs.Scheme, nSenders, perSender, payloadSize int) (*Fleet, error) {
	if payloadSize < 0 {
		return nil, fmt.Errorf("negative payload size %d", payloadSize)
	}
	f := &Fleet{Scheme: sc}

	// Issuing authority: self-signed, signer hash is the hash of the empty
	// string as IEEE 1609.2 specifies for self-signed certificates.
	isk, ipk, err := sc.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("issuer keygen: %w", err)
	}
	var self [8]byte
	issuer := newCert(sc, self, ipk)
	emptyHash := ieee1609.HashCert(nil)
	in := ieee1609.SignInput(issuer.EncodeTBS(), emptyHash)
	if issuer.Signature, err = isk.Sign(in[:]); err != nil {
		return nil, fmt.Errorf("issuer self-sign: %w", err)
	}
	issuer.Encode()
	issuerHash := ieee1609.HashCert(issuer.Raw)
	f.Issuer = issuer

	payload := make([]byte, payloadSize)
	for i := 0; i < nSenders; i++ {
		sk, pk, err := sc.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("sender keygen: %w", err)
		}
		c := newCert(sc, ieee1609.HashedID8(issuerHash), pk)
		in := ieee1609.SignInput(c.EncodeTBS(), issuerHash)
		if c.Signature, err = isk.Sign(in[:]); err != nil {
			return nil, fmt.Errorf("certificate sign: %w", err)
		}
		c.Encode()
		s := &Sender{Cert: c, CertHash: ieee1609.HashCert(c.Raw), key: sk}
		s.ID = ieee1609.HashedID8(s.CertHash)
		f.Senders = append(f.Senders, s)

		for j := 0; j < perSender; j++ {
			for _, carry := range []bool{false, true} {
				randBytes(payload)
				m := &ieee1609.SPDU{
					Payload: append([]byte(nil), payload...),
					Psid:    ieee1609.PsidBSM,
					GenTime: NowMicros - uint64(j)*PeriodMicros/uint64(perSender+1),
					SigAlg:  sc.Alg(),
				}
				if carry {
					m.Cert = c
				} else {
					m.Digest = s.ID
				}
				in := ieee1609.SignInput(m.EncodeTBS(), s.CertHash)
				if m.Signature, err = sk.Sign(in[:]); err != nil {
					return nil, fmt.Errorf("message sign: %w", err)
				}
				msg := Message{Raw: m.Encode(), Sender: i, SignInput: in, Signature: m.Signature}
				if carry {
					f.Carry = append(f.Carry, msg)
				} else {
					f.Digest = append(f.Digest, msg)
				}
			}
		}
	}
	return f, nil
}

// SignMessage signs and encodes one message from sender i with the given
// payload, in digest form (carry false) or certificate-carrying form.
func (f *Fleet) SignMessage(i int, payload []byte, carry bool) (Message, error) {
	if i < 0 || i >= len(f.Senders) {
		return Message{}, fmt.Errorf("no sender %d", i)
	}
	s := f.Senders[i]
	m := &ieee1609.SPDU{
		Payload: append([]byte(nil), payload...),
		Psid:    ieee1609.PsidBSM,
		GenTime: NowMicros,
		SigAlg:  f.Scheme.Alg(),
	}
	if carry {
		m.Cert = s.Cert
	} else {
		m.Digest = s.ID
	}
	in := ieee1609.SignInput(m.EncodeTBS(), s.CertHash)
	sig, err := s.key.Sign(in[:])
	if err != nil {
		return Message{}, fmt.Errorf("message sign: %w", err)
	}
	m.Signature = sig
	return Message{Raw: m.Encode(), Sender: i, SignInput: in, Signature: sig}, nil
}
