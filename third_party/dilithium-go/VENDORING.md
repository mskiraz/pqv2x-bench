# Vendored copy of github.com/lattice-safe/dilithium-go v0.3.0

Upstream commit: 8e950de8d3d4402afdfd97ba0c538e5b109e559e (tag v0.3.0).

Why vendored: v0.3.0 does not build on amd64. `ntt_amd64.go` declares an
assembly function `nttAVX2_8` whose file `ntt_avx2_amd64.s` is not in the
repository, and the generator `gen/ntt_avx2.go` fails with "bad operands".

Local changes (two):

1. Removed `ntt_amd64.go`.
2. Removed the `//go:build !amd64` constraint from `ntt_generic.go`, so the
   portable Go NTT is used on every architecture.

No cryptographic code was changed. The library's own test suite, including
the NIST known-answer tests, passes on the patched copy (`go test .`).
