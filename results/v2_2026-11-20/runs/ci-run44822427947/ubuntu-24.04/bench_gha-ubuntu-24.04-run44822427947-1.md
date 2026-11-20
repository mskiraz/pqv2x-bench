# Receiver authentication benchmark (gha-ubuntu-24.04-run44822427947-1)

- CPU: Intel(R) Xeon(R) 6975P-C (4 logical CPUs), linux/amd64, go1.25.6
- 20000 timed iterations per operation (per-scheme counts below), single goroutine, garbage collector enabled (GC percent 100)
- Run: 2026-11-20T21:41:12Z
- OS: Ubuntu 24.04.5 LTS (kernel 7.0.0-1014-aws); compiler gc; GOMAXPROCS 4; GOGC env ""; effective GC percent 100
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
| ECDSA P-256 | primitive | 63.20 | 8.46 | 61.52 | 66.08 | 81.33 | 199.21 | 459.45 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.07 | 0.11 | 0.05 | 0.08 | 0.34 | 1.52 | 5.11 | 0.001 | 0.001 |
| ECDSA P-256 | digest_hit | 63.36 | 6.65 | 61.90 | 66.22 | 76.49 | 154.58 | 373.75 | 1.003 | 1.006 |
| ECDSA P-256 | cert_known | 63.63 | 8.58 | 62.16 | 66.37 | 74.14 | 188.21 | 463.80 | 1.007 | 1.010 |
| ECDSA P-256 | cert_new | 134.45 | 13.43 | 131.36 | 139.19 | 181.88 | 349.96 | 572.41 | 2.127 | 2.135 |
| ECDSA P-256 | digest_unknown | 0.08 | 0.11 | 0.06 | 0.10 | 0.39 | 1.58 | 5.14 | 0.001 | 0.001 |
| Falcon-512 | primitive | 39.61 | 9.57 | 38.46 | 40.24 | 53.28 | 173.70 | 565.53 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.08 | 0.13 | 0.06 | 0.10 | 0.50 | 1.62 | 5.74 | 0.002 | 0.001 |
| Falcon-512 | digest_hit | 39.98 | 7.91 | 38.88 | 40.59 | 53.29 | 160.00 | 346.45 | 1.009 | 1.011 |
| Falcon-512 | cert_known | 41.35 | 10.42 | 40.03 | 42.07 | 54.61 | 212.69 | 519.17 | 1.044 | 1.041 |
| Falcon-512 | cert_new | 81.74 | 16.46 | 79.26 | 85.79 | 133.74 | 327.14 | 604.63 | 2.064 | 2.061 |
| Falcon-512 | digest_unknown | 0.08 | 0.11 | 0.06 | 0.11 | 0.34 | 1.68 | 3.96 | 0.002 | 0.002 |
| ML-DSA-44 | primitive | 72.37 | 20.03 | 68.78 | 76.44 | 155.07 | 351.18 | 620.35 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.07 | 0.09 | 0.06 | 0.08 | 0.28 | 1.45 | 4.38 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 73.43 | 19.15 | 69.15 | 78.52 | 151.47 | 313.08 | 600.38 | 1.015 | 1.005 |
| ML-DSA-44 | cert_known | 74.77 | 17.87 | 71.31 | 79.22 | 152.13 | 307.49 | 471.24 | 1.033 | 1.037 |
| ML-DSA-44 | cert_new | 150.81 | 30.00 | 143.86 | 155.21 | 301.35 | 480.88 | 934.76 | 2.084 | 2.092 |
| ML-DSA-44 | digest_unknown | 0.08 | 0.13 | 0.06 | 0.12 | 0.29 | 1.66 | 9.55 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 361.72 | 86.75 | 332.32 | 492.56 | 673.38 | 1046.20 | 1078.45 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.07 | 0.08 | 0.06 | 0.07 | 0.27 | 1.07 | 2.04 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 368.61 | 93.08 | 338.25 | 493.05 | 709.99 | 1074.64 | 1296.12 | 1.019 | 1.018 |
| SLH-DSA-SHA2-128s | cert_known | 367.62 | 86.69 | 338.32 | 496.36 | 707.38 | 934.71 | 964.18 | 1.016 | 1.018 |
| SLH-DSA-SHA2-128s | cert_new | 744.60 | 140.44 | 682.57 | 938.63 | 1263.53 | 1602.60 | 2319.88 | 2.059 | 2.054 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.07 | 0.09 | 0.06 | 0.08 | 0.19 | 1.42 | 2.52 | 0.000 | 0.000 |
