# Receiver authentication benchmark (gha-ubuntu-24.04-run44820736544-1-gcoff)

- CPU: AMD EPYC 7763 64-Core Processor (4 logical CPUs), linux/amd64, go1.25.6
- 20000 timed iterations per operation (per-scheme counts below), single goroutine, garbage collector disabled (only the explicit runtime.GC() before each timed loop)
- Run: 2026-11-20T21:30:37Z
- OS: Ubuntu 24.04.5 LTS (kernel 7.0.0-1017-azure); compiler gc; GOMAXPROCS 4; GOGC env ""; effective GC percent -1
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
| ECDSA P-256 | primitive | 82.19 | 6.77 | 79.80 | 88.80 | 109.13 | 146.41 | 233.28 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.12 | 0.06 | 0.13 | 0.14 | 0.30 | 0.97 | 2.29 | 0.002 | 0.002 |
| ECDSA P-256 | digest_hit | 81.21 | 4.73 | 79.88 | 82.80 | 102.28 | 134.99 | 175.42 | 0.988 | 1.001 |
| ECDSA P-256 | cert_known | 82.32 | 10.12 | 80.43 | 89.11 | 109.36 | 139.21 | 1227.16 | 1.002 | 1.008 |
| ECDSA P-256 | cert_new | 173.18 | 8.89 | 169.89 | 180.41 | 217.33 | 250.13 | 308.12 | 2.107 | 2.129 |
| ECDSA P-256 | digest_unknown | 0.11 | 0.05 | 0.09 | 0.14 | 0.25 | 0.56 | 1.97 | 0.001 | 0.001 |
| Falcon-512 | primitive | 57.41 | 5.02 | 56.37 | 57.65 | 78.50 | 117.88 | 254.84 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.11 | 0.05 | 0.12 | 0.14 | 0.25 | 0.54 | 1.81 | 0.002 | 0.002 |
| Falcon-512 | digest_hit | 57.94 | 5.13 | 56.96 | 58.15 | 79.00 | 106.29 | 447.38 | 1.009 | 1.011 |
| Falcon-512 | cert_known | 59.27 | 4.37 | 58.29 | 59.57 | 79.68 | 110.72 | 122.08 | 1.032 | 1.034 |
| Falcon-512 | cert_new | 118.88 | 9.27 | 117.09 | 126.17 | 167.82 | 200.62 | 231.31 | 2.071 | 2.077 |
| Falcon-512 | digest_unknown | 0.09 | 0.03 | 0.09 | 0.09 | 0.19 | 0.41 | 1.17 | 0.002 | 0.002 |
| ML-DSA-44 | primitive | 119.29 | 11.51 | 118.03 | 128.31 | 175.01 | 212.09 | 333.80 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.09 | 0.04 | 0.08 | 0.11 | 0.18 | 0.49 | 1.37 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 120.96 | 11.88 | 120.39 | 131.64 | 179.74 | 220.76 | 286.27 | 1.014 | 1.020 |
| ML-DSA-44 | cert_known | 116.63 | 9.24 | 113.60 | 124.31 | 166.49 | 203.73 | 250.24 | 0.978 | 0.962 |
| ML-DSA-44 | cert_new | 246.40 | 109.81 | 245.10 | 262.51 | 318.80 | 389.74 | 13535.15 | 2.065 | 2.077 |
| ML-DSA-44 | digest_unknown | 0.09 | 0.03 | 0.09 | 0.09 | 0.22 | 0.40 | 1.26 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 477.60 | 1908.55 | 430.34 | 458.24 | 613.95 | 653.74 | 85795.32 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.09 | 0.03 | 0.08 | 0.09 | 0.19 | 0.45 | 0.46 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 448.07 | 38.28 | 442.37 | 478.38 | 632.05 | 687.53 | 760.13 | 0.938 | 1.028 |
| SLH-DSA-SHA2-128s | cert_known | 463.21 | 44.54 | 450.91 | 500.69 | 671.23 | 746.49 | 799.82 | 0.970 | 1.048 |
| SLH-DSA-SHA2-128s | cert_new | 903.49 | 74.44 | 886.86 | 937.25 | 1304.12 | 1354.72 | 1393.67 | 1.892 | 2.061 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.10 | 0.04 | 0.09 | 0.09 | 0.21 | 0.52 | 0.76 | 0.000 | 0.000 |
