// Package receiver implements the receiver-side authentication path for a
// signed IEEE 1609.2 message: decode, signer resolution through a
// certificate cache, certificate validation on a cache miss, relevance and
// permission checks, and signature verification.
package receiver

import (
	"pqv2xbench/internal/ieee1609"
	"pqv2xbench/internal/sigs"
)

// Outcome is the result of processing one message.
type Outcome uint8

const (
	Accepted Outcome = iota
	Malformed
	Stale          // generation time outside the accepted window
	UnknownSigner  // digest refers to a certificate not in the cache
	UnknownIssuer  // carried certificate chains to an unknown issuer
	BadCertificate // certificate signature, validity, or permissions invalid
	BadSignature
)

// Path identifies which authentication path a message took.
type Path uint8

const (
	PathNone      Path = iota
	PathDigestHit      // signer given as digest, certificate cached
	PathCertKnown      // certificate carried and already cached
	PathCertNew        // certificate carried and validated against its issuer
)

type entry struct {
	cert *ieee1609.Certificate
	hash [32]byte
	key  sigs.PublicKey
}

// Receiver holds the certificate cache and acceptance parameters.
type Receiver struct {
	cache map[[8]byte]*entry
	// Now is the receiver clock in microseconds since the 1609.2 epoch.
	Now uint64
	// MaxAge is the accepted message age in microseconds.
	MaxAge uint64
}

// New returns a receiver with an empty cache.
func New(now, maxAge uint64) *Receiver {
	return &Receiver{cache: make(map[[8]byte]*entry), Now: now, MaxAge: maxAge}
}

// Trust installs a certificate without validating it (trust anchors and
// issuing-authority certificates provisioned out of band).
func (rx *Receiver) Trust(c *ieee1609.Certificate) error {
	e, err := newEntry(c)
	if err != nil {
		return err
	}
	rx.cache[ieee1609.HashedID8(e.hash)] = e
	return nil
}

// Forget removes a certificate from the cache.
func (rx *Receiver) Forget(id [8]byte) { delete(rx.cache, id) }

// Len returns the number of cached certificates.
func (rx *Receiver) Len() int { return len(rx.cache) }

func newEntry(c *ieee1609.Certificate) (*entry, error) {
	sc, ok := sigs.ByAlg(c.KeyAlg)
	if !ok {
		return nil, ieee1609.ErrUnsupported
	}
	key, err := sc.ParsePublicKey(c.VerifyKey)
	if err != nil {
		return nil, err
	}
	return &entry{cert: c, hash: ieee1609.HashCert(c.Raw), key: key}, nil
}

func (rx *Receiver) validAt(c *ieee1609.Certificate) bool {
	nowSec := rx.Now / 1_000_000
	start := uint64(c.Start)
	return nowSec >= start && nowSec <= start+uint64(c.DurationHours)*3600
}

// Process authenticates one encoded SPDU.
func (rx *Receiver) Process(raw []byte) (Outcome, Path) {
	s, err := ieee1609.DecodeSPDU(raw)
	if err != nil {
		return Malformed, PathNone
	}
	// Relevance: generation time within the accepted window.
	if s.GenTime > rx.Now+rx.MaxAge || s.GenTime+rx.MaxAge < rx.Now {
		return Stale, PathNone
	}

	var signer *entry
	path := PathDigestHit
	if s.Cert == nil {
		e, ok := rx.cache[s.Digest]
		if !ok {
			return UnknownSigner, PathNone
		}
		signer = e
	} else {
		h := ieee1609.HashCert(s.Cert.Raw)
		id := ieee1609.HashedID8(h)
		if e, ok := rx.cache[id]; ok && e.hash == h {
			signer, path = e, PathCertKnown
		} else {
			issuer, ok := rx.cache[s.Cert.IssuerID]
			if !ok {
				return UnknownIssuer, PathNone
			}
			if s.Cert.SigAlg != issuer.cert.KeyAlg || !rx.validAt(issuer.cert) {
				return BadCertificate, PathNone
			}
			in := ieee1609.SignInput(s.Cert.TBS, issuer.hash)
			if !issuer.key.Verify(in[:], s.Cert.Signature) {
				return BadCertificate, PathNone
			}
			// Cache a private copy: the decoded certificate refers to the
			// message buffer, which the caller may reuse.
			own, err := ieee1609.DecodeCertificate(append([]byte(nil), s.Cert.Raw...))
			if err != nil {
				return BadCertificate, PathNone
			}
			ne, err := newEntry(own)
			if err != nil {
				return BadCertificate, PathNone
			}
			rx.cache[id] = ne
			signer, path = ne, PathCertNew
		}
	}

	// Consistency: validity period, application permission, algorithm.
	if !rx.validAt(signer.cert) || signer.cert.Psid != s.Psid || s.SigAlg != signer.cert.KeyAlg {
		return BadCertificate, path
	}
	in := ieee1609.SignInput(s.TBS, signer.hash)
	if !signer.key.Verify(in[:], s.Signature) {
		return BadSignature, path
	}
	return Accepted, path
}
