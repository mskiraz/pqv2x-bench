# Receiver authentication benchmark (gha-ubuntu-24.04-arm-run44820736544-1-gcoff)

- CPU: ARM implementer 0x41, part 0xd49 (4 logical CPUs), linux/arm64, go1.25.6
- 20000 timed iterations per operation (per-scheme counts below), single goroutine, garbage collector disabled (only the explicit runtime.GC() before each timed loop)
- Run: 2026-11-20T21:30:22Z
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
| ECDSA P-256 | primitive | 69.61 | 3.80 | 68.81 | 71.12 | 75.92 | 113.15 | 290.82 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.10 | 0.05 | 0.10 | 0.11 | 0.18 | 0.82 | 2.46 | 0.001 | 0.001 |
| ECDSA P-256 | digest_hit | 69.78 | 2.97 | 69.36 | 70.38 | 73.52 | 104.50 | 211.26 | 1.003 | 1.008 |
| ECDSA P-256 | cert_known | 70.21 | 2.49 | 69.75 | 70.87 | 75.59 | 108.62 | 193.70 | 1.009 | 1.014 |
| ECDSA P-256 | cert_new | 147.93 | 4.39 | 146.94 | 150.05 | 158.60 | 194.55 | 298.23 | 2.125 | 2.135 |
| ECDSA P-256 | digest_unknown | 0.10 | 0.04 | 0.09 | 0.11 | 0.18 | 0.40 | 3.10 | 0.001 | 0.001 |
| Falcon-512 | primitive | 40.41 | 2.60 | 40.09 | 40.69 | 47.46 | 82.23 | 171.09 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.10 | 0.04 | 0.09 | 0.11 | 0.21 | 0.44 | 1.65 | 0.002 | 0.002 |
| Falcon-512 | digest_hit | 40.68 | 2.20 | 40.44 | 41.01 | 45.91 | 83.38 | 139.54 | 1.007 | 1.009 |
| Falcon-512 | cert_known | 41.78 | 2.67 | 41.50 | 42.14 | 47.65 | 83.07 | 187.31 | 1.034 | 1.035 |
| Falcon-512 | cert_new | 85.20 | 4.69 | 84.38 | 87.62 | 94.70 | 145.89 | 244.75 | 2.108 | 2.105 |
| Falcon-512 | digest_unknown | 0.09 | 0.02 | 0.09 | 0.10 | 0.17 | 0.23 | 0.83 | 0.002 | 0.002 |
| ML-DSA-44 | primitive | 89.86 | 17.06 | 88.26 | 93.50 | 118.75 | 159.52 | 2315.22 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.08 | 0.02 | 0.08 | 0.09 | 0.12 | 0.24 | 0.89 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 90.92 | 5.77 | 89.99 | 95.37 | 106.39 | 150.34 | 288.65 | 1.012 | 1.020 |
| ML-DSA-44 | cert_known | 90.83 | 4.57 | 89.84 | 93.38 | 102.84 | 143.65 | 238.14 | 1.011 | 1.018 |
| ML-DSA-44 | cert_new | 185.87 | 9.29 | 184.10 | 191.62 | 233.36 | 267.11 | 542.63 | 2.068 | 2.086 |
| ML-DSA-44 | digest_unknown | 0.09 | 0.02 | 0.09 | 0.10 | 0.18 | 0.29 | 0.94 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 357.34 | 572.73 | 342.29 | 363.47 | 391.30 | 449.04 | 25953.26 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.09 | 0.03 | 0.09 | 0.09 | 0.20 | 0.29 | 0.87 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 347.30 | 18.15 | 345.20 | 367.02 | 410.74 | 440.45 | 457.70 | 0.972 | 1.009 |
| SLH-DSA-SHA2-128s | cert_known | 355.48 | 18.46 | 350.65 | 379.89 | 407.12 | 437.15 | 445.32 | 0.995 | 1.024 |
| SLH-DSA-SHA2-128s | cert_new | 707.37 | 31.40 | 701.32 | 748.30 | 787.48 | 1028.60 | 1267.01 | 1.980 | 2.049 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.09 | 0.02 | 0.09 | 0.10 | 0.20 | 0.27 | 0.67 | 0.000 | 0.000 |
