# Receiver authentication benchmark (gha-ubuntu-24.04-run44822427947-1-gcoff)

- CPU: Intel(R) Xeon(R) 6975P-C (4 logical CPUs), linux/amd64, go1.25.6
- 20000 timed iterations per operation (per-scheme counts below), single goroutine, garbage collector disabled (only the explicit runtime.GC() before each timed loop)
- Run: 2026-11-20T21:41:44Z
- OS: Ubuntu 24.04.5 LTS (kernel 7.0.0-1014-aws); compiler gc; GOMAXPROCS 4; GOGC env ""; effective GC percent -1
- CPU features: aes, avx, avx2, avx512f, bmi2, sha_ni
- Build settings: -buildmode=exe -compiler=gc CGO_CFLAGS= CGO_CPPFLAGS= CGO_CXXFLAGS= CGO_ENABLED=1 CGO_LDFLAGS= GOAMD64=v1 GOARCH=amd64 GOOS=linux
- Modules: github.com/lattice-safe/dilithium-go@v0.3.0 => ./third_party/dilithium-go@(devel); github.com/lattice-safe/falcon-go@v0.1.1; github.com/lattice-safe/lattice-slh-dsa-go@v0.2.0; golang.org/x/crypto@v0.51.0; golang.org/x/sys@v0.44.0

## Implementations and sample counts

| Scheme | Timed iterations | Warm-up | Fleet (senders x messages per form) | Implementation |
|---|---:|---:|---|---|
| ECDSA P-256 | 20000 | 1000 | 50 x 4 | Go standard library crypto/ecdsa; the public key is parsed once into an *ecdsa.PublicKey when cached |
| Falcon-512 | 20000 | 1000 | 50 x 4 | github.com/lattice-safe/falcon-go v0.1.1, native package, padded 666-byte signatures; verification is pure Go (the module's assembly files serve signing only); the cached PublicKey holds only the encoded key, so every Verify call decodes the 897-byte public key, computes SHAKE256 over it, converts it to NTT form and decodes the signature |
| ML-DSA-44 | 20000 | 1000 | 50 x 4 | github.com/lattice-safe/dilithium-go v0.3.0, vendored in third_party/dilithium-go; the NTT is scalar Go on every architecture (the upstream AVX2 assembly is missing and the vendored nttAVX2_8 is a scalar Go loop), so no SIMD is used; ML-DSA with empty context over the 32-byte digest |
| SLH-DSA-SHA2-128s | 2000 | 50 | 4 x 2 | github.com/lattice-safe/lattice-slh-dsa-go v0.2.0, pure Go without assembly; hashing through the standard library crypto/sha256 and crypto/sha512 |

## Encoded sizes (bytes)

| Scheme | Public key | Signature | Certificate | Digest SPDU | Certificate SPDU |
|---|---:|---:|---:|---:|---:|
| ECDSA P-256 | 33 | 64 | 144 | 132 | 270 |
| Falcon-512 | 897 | 666 | 1615 | 736 | 2345 |
| ML-DSA-44 | 1312 | 2420 | 3784 | 2490 | 6268 |
| SLH-DSA-SHA2-128s | 32 | 7856 | 7938 | 7926 | 15858 |

## Time per message (microseconds)

| Scheme | Operation | Mean | Std | Median | p90 | p99 | p99.9 | Max | Mean / primitive | Median / primitive |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| ECDSA P-256 | primitive | 63.30 | 5.25 | 61.92 | 66.99 | 80.34 | 102.25 | 376.53 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.08 | 0.05 | 0.08 | 0.09 | 0.16 | 0.55 | 1.83 | 0.001 | 0.001 |
| ECDSA P-256 | digest_hit | 62.89 | 3.53 | 61.74 | 65.66 | 72.21 | 84.36 | 236.67 | 0.994 | 0.997 |
| ECDSA P-256 | cert_known | 63.27 | 4.97 | 62.03 | 66.34 | 72.41 | 102.66 | 267.39 | 1.000 | 1.002 |
| ECDSA P-256 | cert_new | 133.76 | 6.01 | 131.85 | 139.46 | 147.63 | 191.32 | 473.25 | 2.113 | 2.129 |
| ECDSA P-256 | digest_unknown | 0.08 | 0.05 | 0.06 | 0.09 | 0.16 | 0.38 | 3.28 | 0.001 | 0.001 |
| Falcon-512 | primitive | 39.05 | 4.62 | 38.40 | 39.85 | 47.98 | 62.32 | 329.12 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.08 | 0.04 | 0.09 | 0.10 | 0.21 | 0.34 | 2.05 | 0.002 | 0.002 |
| Falcon-512 | digest_hit | 39.51 | 2.39 | 38.92 | 40.38 | 48.21 | 57.16 | 145.35 | 1.012 | 1.014 |
| Falcon-512 | cert_known | 40.80 | 3.66 | 40.05 | 41.55 | 50.19 | 79.39 | 231.36 | 1.045 | 1.043 |
| Falcon-512 | cert_new | 82.30 | 5.63 | 81.32 | 87.28 | 95.63 | 135.05 | 275.21 | 2.108 | 2.118 |
| Falcon-512 | digest_unknown | 0.06 | 0.03 | 0.06 | 0.06 | 0.15 | 0.24 | 2.72 | 0.002 | 0.002 |
| ML-DSA-44 | primitive | 73.90 | 12.81 | 73.15 | 79.24 | 92.09 | 173.45 | 1519.82 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.06 | 0.02 | 0.06 | 0.06 | 0.12 | 0.24 | 1.07 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 75.48 | 7.97 | 75.00 | 82.67 | 96.84 | 146.83 | 362.96 | 1.021 | 1.025 |
| ML-DSA-44 | cert_known | 73.28 | 4.39 | 71.69 | 79.24 | 90.29 | 111.73 | 169.01 | 0.992 | 0.980 |
| ML-DSA-44 | cert_new | 157.32 | 43.50 | 156.80 | 170.21 | 188.59 | 266.98 | 6026.59 | 2.129 | 2.144 |
| ML-DSA-44 | digest_unknown | 0.06 | 0.02 | 0.06 | 0.06 | 0.19 | 0.29 | 0.79 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 325.96 | 19.58 | 322.79 | 347.11 | 372.10 | 513.22 | 536.70 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.06 | 0.02 | 0.06 | 0.06 | 0.18 | 0.28 | 0.47 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 339.74 | 26.74 | 335.31 | 368.77 | 440.28 | 508.66 | 527.35 | 1.042 | 1.039 |
| SLH-DSA-SHA2-128s | cert_known | 344.91 | 24.70 | 341.03 | 365.02 | 431.62 | 596.54 | 824.44 | 1.058 | 1.057 |
| SLH-DSA-SHA2-128s | cert_new | 683.08 | 41.36 | 676.36 | 716.93 | 856.75 | 1154.33 | 1159.90 | 2.096 | 2.095 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.06 | 0.02 | 0.06 | 0.06 | 0.17 | 0.24 | 0.39 | 0.000 | 0.000 |
