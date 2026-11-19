// Package ieee1609 encodes and decodes the IEEE 1609.2 structures a receiver
// handles for a signed Basic Safety Message: the secured message
// (Ieee1609Dot2Data / SignedData) and an explicit pseudonym certificate.
//
// The encoding follows the COER rules for the IEEE 1609.2-2022 ASN.1 types
// with the options an SAE J2945/1 BSM uses (psid and generationTime in the
// header; signer given as HashedId8 or as one certificate; pseudonym
// certificate identified by linkage data, with one application permission).
//
// IEEE 1609.2 defines no post-quantum algorithms. Post-quantum keys and
// signatures are therefore carried as hypothetical extension additions to the
// PublicVerificationKey and Signature CHOICE types (indices 5-7, see package
// sigs). COER encodes an extension addition as an open type, i.e. with a
// length determinant, which is what the sizes reported here include.
package ieee1609

import (
	"crypto/sha256"
	"errors"

	"pqv2xbench/internal/coer"
	"pqv2xbench/internal/sigs"
)

// PsidBSM is the PSID of the vehicle-to-vehicle safety application.
const PsidBSM = 0x20

var (
	ErrUnsupported = errors.New("ieee1609: unsupported structure or option")
	ErrTrailing    = errors.New("ieee1609: trailing bytes")
)

// Certificate is an explicit IEEE 1609.2 certificate.
type Certificate struct {
	IssuerID      [8]byte // HashedId8 of the issuer certificate
	ICert         uint16  // linkage data
	LinkageValue  [9]byte
	CracaID       [3]byte
	CrlSeries     uint16
	Start         uint32 // Time32, seconds since the 1609.2 epoch
	DurationHours uint16
	Psid          uint64
	KeyAlg        sigs.Alg
	VerifyKey     []byte
	SigAlg        sigs.Alg
	Signature     []byte

	Raw []byte // complete encoding (set by Encode and Decode)
	TBS []byte // ToBeSignedCertificate encoding
}

// SPDU is a signed secured protocol data unit carrying an unsecured payload.
type SPDU struct {
	Payload   []byte
	Psid      uint64
	GenTime   uint64       // Time64, microseconds since the 1609.2 epoch
	Digest    [8]byte      // signer HashedId8, valid when Cert == nil
	Cert      *Certificate // signer certificate, when carried
	SigAlg    sigs.Alg
	Signature []byte

	TBS []byte // ToBeSignedData encoding
}

// HashCert returns SHA-256 of a certificate encoding.
func HashCert(raw []byte) [32]byte { return sha256.Sum256(raw) }

// HashedID8 returns the low-order eight bytes of a certificate hash.
func HashedID8(h [32]byte) (id [8]byte) {
	copy(id[:], h[24:])
	return
}

// SignInput computes the IEEE 1609.2 signature input for explicit
// certificates: SHA-256( SHA-256(tbs) || SHA-256(signer certificate) ).
func SignInput(tbs []byte, signerHash [32]byte) [32]byte {
	var buf [64]byte
	h := sha256.Sum256(tbs)
	copy(buf[:32], h[:])
	copy(buf[32:], signerHash[:])
	return sha256.Sum256(buf[:])
}

func putSignature(w *coer.Writer, alg sigs.Alg, sig []byte) {
	if alg == sigs.ECDSAP256 {
		w.Choice(0)     // ecdsaNistP256Signature
		w.Choice(0)     // rSig: x-only
		w.Raw(sig[:32]) // r
		w.Raw(sig[32:]) // sSig
		return
	}
	w.Choice(int(alg)) // extension addition, open type
	w.Var(sig)
}

func getSignature(r *coer.Reader) (sigs.Alg, []byte) {
	alg := sigs.Alg(r.Choice())
	if alg == sigs.ECDSAP256 {
		if r.Choice() != 0 {
			r.Fail()
		}
		return alg, r.Raw(64) // r || s are contiguous
	}
	return alg, r.Var()
}

func putVerifyKey(w *coer.Writer, alg sigs.Alg, pk []byte) {
	w.Choice(0) // VerificationKeyIndicator: verificationKey
	if alg == sigs.ECDSAP256 {
		w.Choice(0)                // ecdsaNistP256
		w.Choice(2 + int(pk[0]&1)) // compressed-y-0 / compressed-y-1
		w.Raw(pk[1:33])
		return
	}
	w.Choice(int(alg))
	w.Var(pk)
}

func getVerifyKey(r *coer.Reader) (sigs.Alg, []byte) {
	if r.Choice() != 0 {
		r.Fail()
	}
	alg := sigs.Alg(r.Choice())
	if alg == sigs.ECDSAP256 {
		form := r.Choice()
		x := r.Raw(32)
		if form != 2 && form != 3 || x == nil {
			r.Fail()
			return alg, nil
		}
		pk := make([]byte, 33)
		pk[0] = byte(form) // 0x02 or 0x03, the SEC 1 compressed prefix
		copy(pk[1:], x)
		return alg, pk
	}
	return alg, r.Var()
}

// EncodeTBS encodes the ToBeSignedCertificate and stores it in c.TBS.
func (c *Certificate) EncodeTBS() []byte {
	w := coer.Writer{B: make([]byte, 0, 80+len(c.VerifyKey))}
	w.U8(0x10)  // preamble: only appPermissions present
	w.Choice(0) // CertificateId: linkageData
	w.U8(0x00)  // LinkageData preamble: no group linkage value
	w.U16(c.ICert)
	w.Raw(c.LinkageValue[:])
	w.Raw(c.CracaID[:])
	w.U16(c.CrlSeries)
	w.U32(c.Start)
	w.Choice(4) // Duration: hours
	w.U16(c.DurationHours)
	w.Count(1) // SequenceOfPsidSsp
	w.U8(0x00) // PsidSsp preamble: no ssp
	w.UintUnbounded(c.Psid)
	putVerifyKey(&w, c.KeyAlg, c.VerifyKey)
	c.TBS = w.B
	return w.B
}

// Encode encodes the certificate (c.Signature must be set) and stores the
// result in c.Raw.
func (c *Certificate) Encode() []byte {
	tbs := c.EncodeTBS()
	w := coer.Writer{B: make([]byte, 0, 16+len(tbs)+len(c.Signature))}
	w.U8(0x80)  // CertificateBase preamble: signature present
	w.U8(3)     // version
	w.U8(0)     // type: explicit
	w.Choice(0) // IssuerIdentifier: sha256AndDigest
	w.Raw(c.IssuerID[:])
	start := len(w.B)
	w.Raw(tbs)
	putSignature(&w, c.SigAlg, c.Signature)
	c.Raw = w.B
	c.TBS = w.B[start : start+len(tbs)]
	return w.B
}

func decodeCertificate(r *coer.Reader, buf []byte) *Certificate {
	c := &Certificate{}
	begin := r.Offset()
	if r.U8() != 0x80 || r.U8() != 3 || r.U8() != 0 || r.Choice() != 0 {
		r.Fail()
		return nil
	}
	copy(c.IssuerID[:], r.Raw(8))
	tbsStart := r.Offset()
	if r.U8() != 0x10 || r.Choice() != 0 || r.U8() != 0x00 {
		r.Fail()
		return nil
	}
	c.ICert = r.U16()
	copy(c.LinkageValue[:], r.Raw(9))
	copy(c.CracaID[:], r.Raw(3))
	c.CrlSeries = r.U16()
	c.Start = r.U32()
	if r.Choice() != 4 {
		r.Fail()
		return nil
	}
	c.DurationHours = r.U16()
	if r.Count() != 1 || r.U8() != 0x00 {
		r.Fail()
		return nil
	}
	c.Psid = r.UintUnbounded()
	c.KeyAlg, c.VerifyKey = getVerifyKey(r)
	tbsEnd := r.Offset()
	c.SigAlg, c.Signature = getSignature(r)
	if r.Err() != nil {
		return nil
	}
	c.TBS = buf[tbsStart:tbsEnd]
	c.Raw = buf[begin:r.Offset()]
	return c
}

// DecodeCertificate decodes a certificate that occupies all of b. The result
// refers to b; b must not be modified afterwards.
func DecodeCertificate(b []byte) (*Certificate, error) {
	r := coer.NewReader(b)
	c := decodeCertificate(r, b)
	if err := r.Err(); err != nil {
		return nil, err
	}
	if r.Left() != 0 {
		return nil, ErrTrailing
	}
	return c, nil
}

// EncodeTBS encodes the ToBeSignedData and stores it in s.TBS.
func (s *SPDU) EncodeTBS() []byte {
	w := coer.Writer{B: make([]byte, 0, 24+len(s.Payload))}
	w.U8(0x40)  // SignedDataPayload preamble: data present
	w.U8(3)     // inner Ieee1609Dot2Data: protocolVersion
	w.Choice(0) // unsecuredData
	w.Var(s.Payload)
	w.U8(0x40) // HeaderInfo preamble: generationTime present
	w.UintUnbounded(s.Psid)
	w.U64(s.GenTime)
	s.TBS = w.B
	return w.B
}

// Encode encodes the SPDU; s.Signature must be set.
func (s *SPDU) Encode() []byte {
	tbs := s.EncodeTBS()
	n := 16 + len(tbs) + len(s.Signature)
	if s.Cert != nil {
		n += len(s.Cert.Raw)
	}
	w := coer.Writer{B: make([]byte, 0, n)}
	w.U8(3)     // protocolVersion
	w.Choice(1) // signedData
	w.U8(0)     // hashId: sha256
	w.Raw(tbs)
	if s.Cert == nil {
		w.Choice(0) // SignerIdentifier: digest
		w.Raw(s.Digest[:])
	} else {
		w.Choice(1) // SignerIdentifier: certificate
		w.Count(1)
		w.Raw(s.Cert.Raw)
	}
	putSignature(&w, s.SigAlg, s.Signature)
	return w.B
}

// DecodeSPDU decodes a signed SPDU. The result refers to b.
func DecodeSPDU(b []byte) (*SPDU, error) {
	r := coer.NewReader(b)
	s := &SPDU{}
	if r.U8() != 3 || r.Choice() != 1 || r.U8() != 0 {
		return nil, ErrUnsupported
	}
	tbsStart := r.Offset()
	if r.U8() != 0x40 || r.U8() != 3 || r.Choice() != 0 {
		return nil, ErrUnsupported
	}
	s.Payload = r.Var()
	if r.U8() != 0x40 {
		return nil, ErrUnsupported
	}
	s.Psid = r.UintUnbounded()
	s.GenTime = r.U64()
	tbsEnd := r.Offset()
	switch r.Choice() {
	case 0:
		copy(s.Digest[:], r.Raw(8))
	case 1:
		if r.Count() != 1 {
			return nil, ErrUnsupported
		}
		s.Cert = decodeCertificate(r, b)
	default:
		return nil, ErrUnsupported
	}
	s.SigAlg, s.Signature = getSignature(r)
	if err := r.Err(); err != nil {
		return nil, err
	}
	if r.Left() != 0 {
		return nil, ErrTrailing
	}
	s.TBS = b[tbsStart:tbsEnd]
	return s, nil
}
