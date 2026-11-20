# Receiver authentication benchmark (gha-ubuntu-24.04-arm-run44822427947-1)

- CPU: ARM implementer 0x41, part 0xd49 (4 logical CPUs), linux/arm64, go1.25.6
- 20000 timed iterations per operation (per-scheme counts below), single goroutine, garbage collector enabled (GC percent 100)
- Run: 2026-11-20T21:41:07Z
- OS: Ubuntu 24.04.5 LTS (kernel 7.0.0-1017-azure); compiler gc; GOMAXPROCS 4; GOGC env ""; effective GC percent 100
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
| ECDSA P-256 | primitive | 69.49 | 5.19 | 68.91 | 70.53 | 75.24 | 117.85 | 460.82 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.10 | 0.12 | 0.09 | 0.12 | 0.40 | 1.83 | 3.30 | 0.001 | 0.001 |
| ECDSA P-256 | digest_hit | 70.03 | 5.45 | 69.44 | 70.75 | 75.36 | 113.12 | 367.40 | 1.008 | 1.008 |
| ECDSA P-256 | cert_known | 70.50 | 5.31 | 69.84 | 71.55 | 76.74 | 146.63 | 297.07 | 1.015 | 1.014 |
| ECDSA P-256 | cert_new | 148.17 | 9.90 | 146.90 | 149.39 | 170.66 | 316.95 | 515.32 | 2.132 | 2.132 |
| ECDSA P-256 | digest_unknown | 0.11 | 0.13 | 0.10 | 0.14 | 0.57 | 1.81 | 2.85 | 0.002 | 0.001 |
| Falcon-512 | primitive | 40.77 | 6.63 | 40.15 | 40.92 | 47.34 | 162.55 | 268.24 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.11 | 0.13 | 0.09 | 0.14 | 0.67 | 1.86 | 2.85 | 0.003 | 0.002 |
| Falcon-512 | digest_hit | 41.14 | 7.05 | 40.45 | 41.25 | 48.70 | 163.60 | 294.70 | 1.009 | 1.008 |
| Falcon-512 | cert_known | 42.30 | 7.37 | 41.56 | 42.46 | 50.71 | 171.45 | 265.22 | 1.038 | 1.035 |
| Falcon-512 | cert_new | 85.33 | 12.62 | 83.59 | 86.36 | 135.88 | 258.21 | 366.04 | 2.093 | 2.082 |
| Falcon-512 | digest_unknown | 0.12 | 0.12 | 0.10 | 0.15 | 0.34 | 1.77 | 3.01 | 0.003 | 0.003 |
| ML-DSA-44 | primitive | 90.51 | 15.18 | 88.12 | 92.57 | 165.53 | 292.02 | 361.26 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.10 | 0.11 | 0.09 | 0.12 | 0.29 | 1.77 | 5.03 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 91.14 | 15.82 | 88.44 | 93.08 | 169.99 | 304.55 | 450.98 | 1.007 | 1.004 |
| ML-DSA-44 | cert_known | 93.25 | 16.52 | 90.43 | 95.05 | 173.47 | 301.66 | 486.20 | 1.030 | 1.026 |
| ML-DSA-44 | cert_new | 187.73 | 25.42 | 182.49 | 190.11 | 339.79 | 443.78 | 900.65 | 2.074 | 2.071 |
| ML-DSA-44 | digest_unknown | 0.11 | 0.10 | 0.10 | 0.15 | 0.27 | 1.64 | 4.42 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 417.47 | 84.42 | 410.23 | 532.56 | 701.87 | 743.22 | 768.03 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.10 | 0.09 | 0.08 | 0.12 | 0.29 | 1.41 | 2.38 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 417.66 | 89.51 | 405.36 | 548.00 | 701.93 | 782.39 | 1310.03 | 1.000 | 0.988 |
| SLH-DSA-SHA2-128s | cert_known | 430.18 | 88.55 | 417.33 | 560.35 | 708.16 | 752.67 | 784.31 | 1.030 | 1.017 |
| SLH-DSA-SHA2-128s | cert_new | 853.43 | 139.54 | 843.48 | 1078.24 | 1183.57 | 1221.61 | 1246.68 | 2.044 | 2.056 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.11 | 0.11 | 0.10 | 0.13 | 0.36 | 1.69 | 2.43 | 0.000 | 0.000 |
