# Receiver authentication benchmark (gha-ubuntu-24.04-run44820736544-1)

- CPU: AMD EPYC 7763 64-Core Processor (4 logical CPUs), linux/amd64, go1.25.6
- 20000 timed iterations per operation (per-scheme counts below), single goroutine, garbage collector enabled (GC percent 100)
- Run: 2026-11-20T21:29:53Z
- OS: Ubuntu 24.04.5 LTS (kernel 7.0.0-1017-azure); compiler gc; GOMAXPROCS 4; GOGC env ""; effective GC percent 100
- CPU features: aes, avx, avx2, bmi2, sha_ni
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
| ECDSA P-256 | primitive | 81.49 | 7.94 | 79.66 | 87.98 | 104.58 | 171.80 | 411.29 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.10 | 0.14 | 0.08 | 0.12 | 0.58 | 1.91 | 7.09 | 0.001 | 0.001 |
| ECDSA P-256 | digest_hit | 81.93 | 7.71 | 80.05 | 88.61 | 105.01 | 182.41 | 304.54 | 1.005 | 1.005 |
| ECDSA P-256 | cert_known | 82.66 | 9.19 | 80.51 | 89.34 | 110.99 | 192.05 | 445.94 | 1.014 | 1.011 |
| ECDSA P-256 | cert_new | 173.35 | 14.53 | 169.50 | 179.85 | 224.51 | 385.64 | 565.10 | 2.127 | 2.128 |
| ECDSA P-256 | digest_unknown | 0.11 | 0.15 | 0.08 | 0.14 | 0.91 | 1.98 | 4.17 | 0.001 | 0.001 |
| Falcon-512 | primitive | 57.93 | 10.40 | 56.31 | 58.68 | 82.28 | 232.49 | 424.35 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.11 | 0.16 | 0.08 | 0.14 | 0.98 | 2.03 | 3.54 | 0.002 | 0.001 |
| Falcon-512 | digest_hit | 58.45 | 9.74 | 56.88 | 59.05 | 82.15 | 206.32 | 358.73 | 1.009 | 1.010 |
| Falcon-512 | cert_known | 59.97 | 10.87 | 58.24 | 60.47 | 85.50 | 231.97 | 457.81 | 1.035 | 1.034 |
| Falcon-512 | cert_new | 118.94 | 22.14 | 113.28 | 126.29 | 205.54 | 328.14 | 1531.19 | 2.053 | 2.012 |
| Falcon-512 | digest_unknown | 0.12 | 0.14 | 0.09 | 0.15 | 0.79 | 1.97 | 3.06 | 0.002 | 0.002 |
| ML-DSA-44 | primitive | 115.52 | 21.66 | 110.29 | 121.19 | 225.95 | 355.07 | 554.39 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.10 | 0.11 | 0.08 | 0.10 | 0.43 | 1.66 | 3.72 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 115.55 | 21.24 | 110.26 | 121.30 | 223.30 | 363.85 | 453.01 | 1.000 | 1.000 |
| ML-DSA-44 | cert_known | 118.10 | 21.11 | 113.06 | 123.94 | 226.93 | 378.80 | 513.40 | 1.022 | 1.025 |
| ML-DSA-44 | cert_new | 237.40 | 34.52 | 227.64 | 248.31 | 429.87 | 580.10 | 950.56 | 2.055 | 2.064 |
| ML-DSA-44 | digest_unknown | 0.10 | 0.11 | 0.09 | 0.12 | 0.35 | 1.64 | 3.31 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 473.38 | 107.84 | 428.29 | 644.81 | 868.59 | 953.82 | 1034.74 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.10 | 0.12 | 0.08 | 0.12 | 0.43 | 1.78 | 2.55 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 476.81 | 104.57 | 432.71 | 644.02 | 853.66 | 933.19 | 949.07 | 1.007 | 1.010 |
| SLH-DSA-SHA2-128s | cert_known | 484.93 | 111.26 | 439.98 | 658.51 | 890.62 | 1000.47 | 1015.45 | 1.024 | 1.027 |
| SLH-DSA-SHA2-128s | cert_new | 984.66 | 173.03 | 890.85 | 1265.34 | 1499.24 | 1709.93 | 1846.46 | 2.080 | 2.080 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.10 | 0.13 | 0.09 | 0.11 | 0.38 | 1.72 | 4.29 | 0.000 | 0.000 |
