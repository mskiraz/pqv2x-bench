// Package sigs wraps the signature schemes compared in the paper behind one
// interface. Post-quantum schemes come from github.com/lattice-safe; ECDSA
// P-256 comes from the Go standard library.
//
// Every scheme signs and verifies a 32-byte digest, mirroring IEEE 1609.2,
// where the signature input is a SHA-256 hash of the to-be-signed data and
// the signer certificate.
package sigs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"math/big"

	dilithium "github.com/lattice-safe/dilithium-go"
	falconnative "github.com/lattice-safe/falcon-go/native"
	slhdsa "github.com/lattice-safe/lattice-slh-dsa-go"
)

// Alg identifies a scheme. The value is also the CHOICE index used for the
// scheme in the (extended) IEEE 1609.2 PublicVerificationKey and Signature
// types: index 0 is the standard ecdsaNistP256 alternative; indices 5-7 are
// hypothetical extension additions placed after the five alternatives of
// IEEE 1609.2-2022.
type Alg uint8

const (
	ECDSAP256   Alg = 0
	Falcon512   Alg = 5
	MLDSA44     Alg = 6
	SLHDSA128s  Alg = 7
	falconLogN      = 9
	falconSigTy     = falconnative.SigPadded
)

// PublicKey is a parsed verification key, as a receiver would cache it.
type PublicKey interface {
	Verify(digest, sig []byte) bool
}

// PrivateKey signs 32-byte digests.
type PrivateKey interface {
	Sign(digest []byte) ([]byte, error)
}

// Scheme describes one signature scheme.
type Scheme interface {
	Alg() Alg
	Name() string
	PublicKeySize() int // bytes as carried in a certificate
	SignatureSize() int // bytes as carried in a message or certificate
	GenerateKey() (PrivateKey, []byte, error)
	ParsePublicKey(pk []byte) (PublicKey, error)
}

// ByAlg returns the scheme for an algorithm identifier.
func ByAlg(a Alg) (Scheme, bool) {
	switch a {
	case ECDSAP256:
		return ecdsaScheme{}, true
	case Falcon512:
		return falconScheme{}, true
	case MLDSA44:
		return mldsaScheme{}, true
	case SLHDSA128s:
		return slhScheme{}, true
	}
	return nil, false
}

// All lists the schemes in the order used in the paper.
func All() []Scheme {
	return []Scheme{ecdsaScheme{}, falconScheme{}, mldsaScheme{}, slhScheme{}}
}

var errKey = errors.New("sigs: malformed key")

// ---------------------------------------------------------------- ECDSA P-256

type ecdsaScheme struct{}

func (ecdsaScheme) Alg() Alg           { return ECDSAP256 }
func (ecdsaScheme) Name() string       { return "ECDSA P-256" }
func (ecdsaScheme) PublicKeySize() int { return 33 } // compressed point
func (ecdsaScheme) SignatureSize() int { return 64 } // r || s

type ecdsaPriv struct{ k *ecdsa.PrivateKey }
type ecdsaPub struct{ k ecdsa.PublicKey }

func (ecdsaScheme) GenerateKey() (PrivateKey, []byte, error) {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return ecdsaPriv{k}, elliptic.MarshalCompressed(elliptic.P256(), k.X, k.Y), nil
}

func (ecdsaScheme) ParsePublicKey(pk []byte) (PublicKey, error) {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), pk)
	if x == nil {
		return nil, errKey
	}
	return &ecdsaPub{ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}}, nil
}

func (p ecdsaPriv) Sign(digest []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, p.k, digest)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 64)
	r.FillBytes(out[:32])
	s.FillBytes(out[32:])
	return out, nil
}

func (p *ecdsaPub) Verify(digest, sig []byte) bool {
	if len(sig) != 64 {
		return false
	}
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:])
	return ecdsa.Verify(&p.k, digest, r, s)
}

// ----------------------------------------------- Falcon-512 (padded, 666 B)

type falconScheme struct{}

func (falconScheme) Alg() Alg           { return Falcon512 }
func (falconScheme) Name() string       { return "Falcon-512" }
func (falconScheme) PublicKeySize() int { return falconnative.PubKeySize(falconLogN) }
func (falconScheme) SignatureSize() int { return falconnative.SigPaddedSize(falconLogN) }

type falconPriv struct {
	sk  []byte
	rng falconnative.Shake256Context
}
type falconPub struct{ pk []byte }

func (falconScheme) GenerateKey() (PrivateKey, []byte, error) {
	p := &falconPriv{sk: make([]byte, falconnative.PrivKeySize(falconLogN))}
	if rc := falconnative.Shake256InitPRNGFromSystem(&p.rng); rc != 0 {
		return nil, nil, errors.New("sigs: falcon rng init failed")
	}
	pk := make([]byte, falconnative.PubKeySize(falconLogN))
	if rc := falconnative.FalconKeygenMake(&p.rng, falconLogN, p.sk, pk, nil); rc != 0 {
		return nil, nil, errors.New("sigs: falcon keygen failed")
	}
	return p, pk, nil
}

func (s falconScheme) ParsePublicKey(pk []byte) (PublicKey, error) {
	if len(pk) != s.PublicKeySize() {
		return nil, errKey
	}
	return &falconPub{pk: append([]byte(nil), pk...)}, nil
}

func (p *falconPriv) Sign(digest []byte) ([]byte, error) {
	sig := make([]byte, falconnative.SigPaddedSize(falconLogN))
	n := len(sig)
	if rc := falconnative.FalconSignDyn(&p.rng, sig, &n, falconSigTy, p.sk, digest, nil); rc != 0 {
		return nil, errors.New("sigs: falcon sign failed")
	}
	return sig[:n], nil
}

func (p *falconPub) Verify(digest, sig []byte) bool {
	return falconnative.FalconVerify(sig, falconSigTy, p.pk, digest, nil) == 0
}

// ------------------------------------------------------------------ ML-DSA-44

type mldsaScheme struct{}

const mldsaMode = dilithium.Dilithium2

func (mldsaScheme) Alg() Alg           { return MLDSA44 }
func (mldsaScheme) Name() string       { return "ML-DSA-44" }
func (mldsaScheme) PublicKeySize() int { return mldsaMode.PublicKeyBytes() }
func (mldsaScheme) SignatureSize() int { return mldsaMode.SignatureBytes() }

type mldsaPriv struct{ kp *dilithium.DilithiumKeyPair }
type mldsaPub struct{ pk []byte }

func (mldsaScheme) GenerateKey() (PrivateKey, []byte, error) {
	kp, err := dilithium.Generate(mldsaMode)
	if err != nil {
		return nil, nil, err
	}
	return mldsaPriv{kp}, append([]byte(nil), kp.PublicKey()...), nil
}

func (s mldsaScheme) ParsePublicKey(pk []byte) (PublicKey, error) {
	if len(pk) != s.PublicKeySize() {
		return nil, errKey
	}
	return &mldsaPub{pk: append([]byte(nil), pk...)}, nil
}

func (p mldsaPriv) Sign(digest []byte) ([]byte, error) {
	sig, err := p.kp.Sign(digest, nil)
	if err != nil {
		return nil, err
	}
	return append([]byte(nil), sig.AsBytes()...), nil
}

func (p *mldsaPub) Verify(digest, sig []byte) bool {
	if len(sig) != mldsaMode.SignatureBytes() {
		return false
	}
	return dilithium.Verify(mldsaMode, sig, digest, nil, p.pk)
}

// --------------------------------------------------------- SLH-DSA-SHA2-128s

type slhScheme struct{}

var slhMode = slhdsa.SlhDsaSha2_128s

func (slhScheme) Alg() Alg           { return SLHDSA128s }
func (slhScheme) Name() string       { return "SLH-DSA-SHA2-128s" }
func (slhScheme) PublicKeySize() int { return slhMode.PkBytes() }
func (slhScheme) SignatureSize() int { return slhMode.SigBytes() }

type slhPriv struct{ sk *slhdsa.PrivateKey }
type slhPub struct{ pk *slhdsa.PublicKey }

func (slhScheme) GenerateKey() (PrivateKey, []byte, error) {
	sk, pk, err := slhdsa.GenerateKey(slhMode)
	if err != nil {
		return nil, nil, err
	}
	return slhPriv{sk}, append([]byte(nil), pk.Pk...), nil
}

func (s slhScheme) ParsePublicKey(pk []byte) (PublicKey, error) {
	if len(pk) != s.PublicKeySize() {
		return nil, errKey
	}
	return &slhPub{&slhdsa.PublicKey{Mode: slhMode, Pk: append([]byte(nil), pk...)}}, nil
}

func (p slhPriv) Sign(digest []byte) ([]byte, error) { return p.sk.Sign(digest) }

func (p *slhPub) Verify(digest, sig []byte) bool { return p.pk.Verify(digest, sig) }
