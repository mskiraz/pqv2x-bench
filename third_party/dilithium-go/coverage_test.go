package dilithium

import (
	"errors"
	"fmt"
	"testing"
)

// TestInvalidModeCoverage tests all parameter functions in params.go with an invalid Mode to hit fallback branches.
func TestInvalidModeCoverage(t *testing.T) {
	invalidMode := Mode(99)

	if invalidMode.ModeTag() != 0 {
		t.Errorf("expected 0 for invalid ModeTag")
	}
	if invalidMode.K() != 0 {
		t.Errorf("expected 0 for invalid K")
	}
	if invalidMode.L() != 0 {
		t.Errorf("expected 0 for invalid L")
	}
	if invalidMode.Eta() != 0 {
		t.Errorf("expected 0 for invalid Eta")
	}
	if invalidMode.Tau() != 0 {
		t.Errorf("expected 0 for invalid Tau")
	}
	if invalidMode.Beta() != 0 {
		t.Errorf("expected 0 for invalid Beta")
	}
	if invalidMode.Gamma1() != 0 {
		t.Errorf("expected 0 for invalid Gamma1")
	}
	if invalidMode.Gamma2() != 0 {
		t.Errorf("expected 0 for invalid Gamma2")
	}
	if invalidMode.Omega() != 0 {
		t.Errorf("expected 0 for invalid Omega")
	}
	if invalidMode.Ctildebytes() != 0 {
		t.Errorf("expected 0 for invalid Ctildebytes")
	}
	if invalidMode.PolyZPackedBytes() != 0 {
		t.Errorf("expected 0 for invalid PolyZPackedBytes")
	}
	if invalidMode.PolyW1PackedBytes() != 0 {
		t.Errorf("expected 0 for invalid PolyW1PackedBytes")
	}
	if invalidMode.PolyEtaPackedBytes() != 0 {
		t.Errorf("expected 0 for invalid PolyEtaPackedBytes")
	}
	if invalidMode.HashOID() != nil {
		t.Errorf("expected nil for invalid HashOID")
	}
	if invalidMode.FipsName() != "" {
		t.Errorf("expected empty string for invalid FipsName")
	}
}

// TestNttDispatchCoverage tests both hasAVX2 = true and hasAVX2 = false dispatch paths in ntt_dispatch.go.
func TestNttDispatchCoverage(t *testing.T) {
	origAVX2 := hasAVX2

	var poly Poly
	for i := 0; i < N; i++ {
		poly.Coeffs[i] = int32(i * 100)
	}

	// Test hasAVX2 = true
	hasAVX2 = true
	polyCopy1 := poly
	NTT(&polyCopy1.Coeffs)
	InvNTT(&polyCopy1.Coeffs)

	// Test hasAVX2 = false
	hasAVX2 = false
	polyCopy2 := poly
	NTT(&polyCopy2.Coeffs)
	InvNTT(&polyCopy2.Coeffs)

	hasAVX2 = origAVX2
}

// TestRandomErrorCoverage tests error handling when randRead fails.
func TestRandomErrorCoverage(t *testing.T) {
	origRandRead := randRead
	defer func() { randRead = origRandRead }()

	randRead = func(b []byte) (int, error) {
		return 0, errors.New("entropy error")
	}

	_, err := Generate(Dilithium2)
	if !errors.Is(err, ErrRandomError) {
		t.Fatalf("expected ErrRandomError on Generate, got: %v", err)
	}

	kp := &DilithiumKeyPair{mode: Dilithium2, privkey: make([]byte, Dilithium2.SecretKeyBytes()), pubkey: make([]byte, Dilithium2.PublicKeyBytes())}

	_, err = kp.Sign([]byte("msg"), []byte(""))
	if !errors.Is(err, ErrRandomError) {
		t.Fatalf("expected ErrRandomError on Sign, got: %v", err)
	}

	_, err = kp.SignPrehash([]byte("msg"), []byte(""))
	if !errors.Is(err, ErrRandomError) {
		t.Fatalf("expected ErrRandomError on SignPrehash, got: %v", err)
	}
}

type mockSqueezer struct {
	call int
}

func (m *mockSqueezer) Squeeze(out []byte) {
	if m.call == 0 {
		for i := range out {
			out[i] = 0xFF // All rejected
		}
	} else {
		for i := range out {
			out[i] = 0 // Valid
		}
	}
	m.call++
}

// TestPolyUniformMultiBlockCoverage tests the for ctr < N loop in PolyUniform.
func TestPolyUniformMultiBlockCoverage(t *testing.T) {
	var poly Poly
	mock := &mockSqueezer{}
	polyUniformStream(&poly, mock)
	if mock.call < 2 {
		t.Fatalf("expected multi-block squeezing in polyUniformStream")
	}
}

// TestUseHintCoverage tests all branch paths of useHint in rounding.go.
func TestUseHintCoverage(t *testing.T) {
	// Mode with gamma2 == (Q-1)/88 (Dilithium2)
	mode2 := Dilithium2

	// hint = false
	_ = useHint(mode2, 100, false)

	// gamma2 == (Q-1)/88 with hint = true:
	// 1. a0 > 0, a1 == 43 -> returns 0
	res1 := useHint(mode2, 43*2*95232+100, true)
	if res1 != 0 {
		t.Errorf("expected 0 for a1=43, a0>0, got %d", res1)
	}

	// 2. a0 > 0, a1 != 43 -> returns a1 + 1
	res2 := useHint(mode2, 2*2*95232+100, true)
	if res2 != 2+1 {
		t.Errorf("expected 3 for a1=2, a0>0, got %d", res2)
	}

	// 3. a0 <= 0, a1 == 0 -> returns 43
	res3 := useHint(mode2, -100, true)
	if res3 != 43 {
		t.Errorf("expected 43 for a1=0, a0<=0, got %d", res3)
	}

	// 4. a0 <= 0, a1 != 0 -> returns a1 - 1
	res4 := useHint(mode2, 2*2*95232-100, true)
	if res4 != 2-1 {
		t.Errorf("expected 1 for a1=2, a0<=0, got %d", res4)
	}

	// Mode with gamma2 == (Q-1)/32 (Dilithium3)
	mode3 := Dilithium3

	// hint = false
	_ = useHint(mode3, 100, false)

	// gamma2 == (Q-1)/32 with hint = true:
	// a0 > 0 -> (a1 + 1) & 15
	res5 := useHint(mode3, 2*2*261888+100, true)
	if res5 != (2+1)&15 {
		t.Errorf("expected 3 for mode3 a0>0, got %d", res5)
	}

	// a0 <= 0 -> (a1 - 1) & 15
	res6 := useHint(mode3, 2*2*261888-100, true)
	if res6 != (2-1)&15 {
		t.Errorf("expected 1 for mode3 a0<=0, got %d", res6)
	}
}

// TestUnpackSigBranchCoverage tests hint validation error branches in UnpackSig.
func TestUnpackSigBranchCoverage(t *testing.T) {
	mode := Dilithium2
	kp, _ := Generate(mode)
	msg := []byte("test")
	sig, _ := kp.Sign(msg, []byte{})
	sigBytes := sig.AsBytes()

	ctilde := mode.Ctildebytes()
	zPacked := mode.PolyZPackedBytes()
	l := mode.L()
	k := mode.K()
	omega := mode.Omega()
	hStart := ctilde + l*zPacked

	c := make([]byte, ctilde)
	z := NewPolyVecL()
	h := NewPolyVecK()

	// 0. Bad length check
	if !UnpackSig(mode, c, z, h, make([]byte, 10)) {
		t.Errorf("expected UnpackSig to return true for bad signature length")
	}

	// 1. end < idx branch
	corruptEnd := make([]byte, len(sigBytes))
	copy(corruptEnd, sigBytes)
	corruptEnd[hStart+omega+0] = 5
	corruptEnd[hStart+omega+1] = 3 // end < idx (3 < 5)
	if !UnpackSig(mode, c, z, h, corruptEnd) {
		t.Errorf("expected UnpackSig to return true for end < idx")
	}

	// 2. Non-strictly-increasing hint indices: sig[hStart+j] <= sig[hStart+j-1] when j > idx
	corruptOrder := make([]byte, len(sigBytes))
	copy(corruptOrder, sigBytes)
	corruptOrder[hStart+0] = 10
	corruptOrder[hStart+1] = 5 // Decreasing indices in first polynomial (j=1 > idx=0, 5 <= 10)
	corruptOrder[hStart+omega+0] = 2
	for i := 1; i < k; i++ {
		corruptOrder[hStart+omega+i] = 2
	}
	if !UnpackSig(mode, c, z, h, corruptOrder) {
		t.Errorf("expected UnpackSig to return true for non-increasing hint indices")
	}

	// 3. Trailing non-zero byte in hint array (sig[hStart+j] != 0 for j >= idx)
	corruptTrailing := make([]byte, len(sigBytes))
	copy(corruptTrailing, sigBytes)
	for i := 0; i < k; i++ {
		corruptTrailing[hStart+omega+i] = 0 // idx = 0
	}
	corruptTrailing[hStart+0] = 1 // non-zero at index 0 when idx = 0
	if !UnpackSig(mode, c, z, h, corruptTrailing) {
		t.Errorf("expected UnpackSig to return true for non-zero trailing hint byte")
	}
}

// TestFromPublicKeyCoverage tests edge cases in FromPublicKey.
func TestFromPublicKeyCoverage(t *testing.T) {
	// Mode tag valid, but length wrong (not 1 or 1+pkLen, e.g. 5 bytes)
	badLen := []byte{0x02, 1, 2, 3, 4}
	_, _, err := FromPublicKey(badLen)
	if !errors.Is(err, ErrFormatError) {
		t.Fatalf("expected ErrFormatError for wrong length in FromPublicKey")
	}
}

// TestSignSignatureInternalRejectionBranches tests internal signing rejection loop branches.
func TestSignSignatureInternalRejectionBranches(t *testing.T) {
	mode := Dilithium2
	var seed [RNDBYTES]byte

	_, sk := Keypair(mode, &seed)

	s1Offset := 2*SEEDBYTES + TRBYTES
	s2Offset := s1Offset + mode.L()*mode.PolyEtaPackedBytes()

	// 1. Corrupt s1 to trigger z norm rejection
	corruptSk1 := make([]byte, len(sk))
	copy(corruptSk1, sk)
	for i := s1Offset; i < s1Offset+mode.L()*mode.PolyEtaPackedBytes(); i++ {
		corruptSk1[i] = 0xAA
	}
	sig1 := make([]byte, mode.SignatureBytes())
	SignSignatureInternal(mode, sig1, []byte("m"), []byte{0, 0}, &seed, corruptSk1)

	// 2. Corrupt s2 to trigger w0 norm rejection
	corruptSk2 := make([]byte, len(sk))
	copy(corruptSk2, sk)
	for i := s2Offset; i < s2Offset+mode.K()*mode.PolyEtaPackedBytes(); i++ {
		corruptSk2[i] = 0x55
	}
	sig2 := make([]byte, mode.SignatureBytes())
	SignSignatureInternal(mode, sig2, []byte("m"), []byte{0, 0}, &seed, corruptSk2)

	// 3. Keep s1 and s2 intact, but set t0 bytes to 0x00 (all +4096) across multiple seeds/messages
	// to ensure PolyVecKChkNorm(h, gamma2) (line 193 continue) is triggered on at least one iteration.
	for _, m := range []Mode{Dilithium2, Dilithium3, Dilithium5} {
		for i := 0; i < 50; i++ {
			var s [RNDBYTES]byte
			s[0] = byte(i)
			_, skM := Keypair(m, &s)
			t0OffM := 2*SEEDBYTES + TRBYTES + m.L()*m.PolyEtaPackedBytes() + m.K()*m.PolyEtaPackedBytes()

			corruptSkT0 := make([]byte, len(skM))
			copy(corruptSkT0, skM)
			for j := 0; j < m.K()*POLYT0_PACKEDBYTES; j++ {
				corruptSkT0[t0OffM+j] = 0x00
			}
			sigM := make([]byte, m.SignatureBytes())
			msg := []byte(fmt.Sprintf("test-msg-%d", i))
			SignSignatureInternal(m, sigM, msg, []byte{0, 0}, &s, corruptSkT0)
		}
	}
}
