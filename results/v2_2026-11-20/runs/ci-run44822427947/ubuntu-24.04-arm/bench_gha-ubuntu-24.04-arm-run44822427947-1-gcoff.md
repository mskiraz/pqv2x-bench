# Receiver authentication benchmark (gha-ubuntu-24.04-arm-run44822427947-1-gcoff)

- CPU: ARM implementer 0x41, part 0xd49 (4 logical CPUs), linux/arm64, go1.25.6
- 20000 timed iterations per operation (per-scheme counts below), single goroutine, garbage collector disabled (only the explicit runtime.GC() before each timed loop)
- Run: 2026-11-20T21:41:44Z
- OS: Ubuntu 24.04.5 LTS (kernel 7.0.0-1017-azure); compiler gc; GOMAXPROCS 4; GOGC env ""; effective GC percent -1
- CPU features: aes, asimd, atomics, sha2, sha3, sha512
- Build settings: -buildmode=exe -compiler=gc CGO_CFLAGS= CGO_CPPFLAGS= CGO_CXXFLAGS= CGO_ENABLED=1 CGO_LDFLAGS= GOARCH=arm64 GOARM64=v8.0 GOOS=linux
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
| ECDSA P-256 | primitive | 69.60 | 3.56 | 68.83 | 71.12 | 75.74 | 111.40 | 308.33 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.10 | 0.06 | 0.10 | 0.11 | 0.21 | 0.88 | 2.43 | 0.001 | 0.001 |
| ECDSA P-256 | digest_hit | 69.84 | 2.33 | 69.42 | 70.51 | 74.51 | 105.65 | 134.60 | 1.003 | 1.009 |
| ECDSA P-256 | cert_known | 70.24 | 2.45 | 69.78 | 70.93 | 76.17 | 106.68 | 148.32 | 1.009 | 1.014 |
| ECDSA P-256 | cert_new | 148.05 | 4.25 | 147.00 | 150.24 | 159.04 | 195.77 | 300.87 | 2.127 | 2.136 |
| ECDSA P-256 | digest_unknown | 0.10 | 0.04 | 0.10 | 0.12 | 0.17 | 0.56 | 1.63 | 0.001 | 0.001 |
| Falcon-512 | primitive | 40.36 | 1.95 | 40.15 | 40.68 | 43.90 | 77.97 | 98.68 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.10 | 0.04 | 0.09 | 0.11 | 0.23 | 0.53 | 3.13 | 0.002 | 0.002 |
| Falcon-512 | digest_hit | 40.71 | 1.95 | 40.47 | 41.04 | 46.12 | 78.47 | 87.88 | 1.008 | 1.008 |
| Falcon-512 | cert_known | 41.76 | 2.29 | 41.48 | 42.12 | 47.94 | 84.10 | 155.99 | 1.035 | 1.033 |
| Falcon-512 | cert_new | 85.27 | 4.93 | 84.45 | 87.83 | 94.86 | 147.08 | 252.55 | 2.113 | 2.103 |
| Falcon-512 | digest_unknown | 0.09 | 0.02 | 0.09 | 0.10 | 0.17 | 0.22 | 0.78 | 0.002 | 0.002 |
| ML-DSA-44 | primitive | 89.61 | 5.80 | 88.25 | 93.21 | 117.20 | 148.35 | 261.75 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.08 | 0.02 | 0.08 | 0.09 | 0.12 | 0.25 | 2.10 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 91.19 | 5.62 | 90.14 | 95.82 | 107.26 | 149.28 | 286.79 | 1.018 | 1.021 |
| ML-DSA-44 | cert_known | 91.06 | 4.77 | 90.00 | 94.08 | 103.41 | 149.11 | 237.06 | 1.016 | 1.020 |
| ML-DSA-44 | cert_new | 187.02 | 37.50 | 185.01 | 193.25 | 235.78 | 269.10 | 5317.04 | 2.087 | 2.096 |
| ML-DSA-44 | digest_unknown | 0.09 | 0.02 | 0.09 | 0.10 | 0.19 | 0.29 | 1.57 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 363.91 | 576.01 | 348.49 | 374.34 | 400.32 | 432.92 | 26103.77 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.09 | 0.02 | 0.08 | 0.09 | 0.20 | 0.28 | 0.37 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 353.53 | 18.38 | 351.59 | 377.16 | 398.33 | 426.59 | 452.37 | 0.971 | 1.009 |
| SLH-DSA-SHA2-128s | cert_known | 358.87 | 16.63 | 359.71 | 377.16 | 407.22 | 439.32 | 440.36 | 0.986 | 1.032 |
| SLH-DSA-SHA2-128s | cert_new | 714.53 | 18.23 | 713.83 | 736.82 | 764.91 | 780.71 | 802.41 | 1.963 | 2.048 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.10 | 0.02 | 0.09 | 0.10 | 0.22 | 0.28 | 0.30 | 0.000 | 0.000 |
