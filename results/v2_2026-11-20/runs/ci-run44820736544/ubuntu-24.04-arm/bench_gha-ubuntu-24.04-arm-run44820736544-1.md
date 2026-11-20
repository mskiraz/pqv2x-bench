# Receiver authentication benchmark (gha-ubuntu-24.04-arm-run44820736544-1)

- CPU: ARM implementer 0x41, part 0xd49 (4 logical CPUs), linux/arm64, go1.25.6
- 20000 timed iterations per operation (per-scheme counts below), single goroutine, garbage collector enabled (GC percent 100)
- Run: 2026-11-20T21:29:45Z
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
| ECDSA P-256 | primitive | 69.52 | 4.81 | 68.89 | 70.70 | 75.90 | 123.78 | 282.01 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.10 | 0.12 | 0.08 | 0.13 | 0.41 | 1.76 | 4.37 | 0.001 | 0.001 |
| ECDSA P-256 | digest_hit | 70.04 | 4.56 | 69.44 | 70.93 | 76.13 | 149.00 | 255.08 | 1.008 | 1.008 |
| ECDSA P-256 | cert_known | 70.48 | 5.85 | 69.83 | 71.28 | 76.77 | 161.23 | 362.51 | 1.014 | 1.014 |
| ECDSA P-256 | cert_new | 148.05 | 9.70 | 146.79 | 149.25 | 168.98 | 335.78 | 481.89 | 2.130 | 2.131 |
| ECDSA P-256 | digest_unknown | 0.11 | 0.12 | 0.09 | 0.12 | 0.52 | 1.90 | 3.39 | 0.002 | 0.001 |
| Falcon-512 | primitive | 40.80 | 7.42 | 40.13 | 41.13 | 47.73 | 154.03 | 536.08 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.11 | 0.13 | 0.09 | 0.12 | 0.57 | 1.80 | 3.57 | 0.003 | 0.002 |
| Falcon-512 | digest_hit | 41.13 | 6.52 | 40.45 | 41.37 | 48.60 | 159.51 | 231.71 | 1.008 | 1.008 |
| Falcon-512 | cert_known | 42.24 | 7.12 | 41.54 | 42.49 | 49.99 | 163.61 | 279.27 | 1.035 | 1.035 |
| Falcon-512 | cert_new | 85.59 | 14.30 | 83.64 | 86.54 | 143.51 | 263.38 | 940.59 | 2.098 | 2.084 |
| Falcon-512 | digest_unknown | 0.12 | 0.12 | 0.10 | 0.15 | 0.35 | 1.79 | 3.23 | 0.003 | 0.002 |
| ML-DSA-44 | primitive | 90.16 | 17.47 | 87.69 | 91.65 | 160.76 | 289.57 | 1315.62 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.10 | 0.10 | 0.09 | 0.12 | 0.27 | 1.68 | 2.93 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 90.65 | 15.30 | 88.26 | 92.39 | 160.70 | 298.38 | 475.12 | 1.005 | 1.006 |
| ML-DSA-44 | cert_known | 93.12 | 25.00 | 90.23 | 94.75 | 166.76 | 312.52 | 2185.72 | 1.033 | 1.029 |
| ML-DSA-44 | cert_new | 186.71 | 24.12 | 181.98 | 188.45 | 332.56 | 442.63 | 533.10 | 2.071 | 2.075 |
| ML-DSA-44 | digest_unknown | 0.11 | 0.09 | 0.10 | 0.16 | 0.26 | 1.64 | 2.96 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 411.66 | 81.86 | 400.43 | 533.22 | 692.86 | 755.88 | 814.89 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.10 | 0.09 | 0.09 | 0.10 | 0.35 | 1.34 | 2.21 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 406.31 | 83.67 | 392.82 | 536.38 | 686.11 | 733.91 | 751.03 | 0.987 | 0.981 |
| SLH-DSA-SHA2-128s | cert_known | 415.15 | 81.19 | 402.77 | 537.34 | 681.88 | 717.68 | 735.42 | 1.008 | 1.006 |
| SLH-DSA-SHA2-128s | cert_new | 828.43 | 131.65 | 812.69 | 1043.70 | 1133.89 | 1169.35 | 1207.26 | 2.012 | 2.030 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.11 | 0.11 | 0.09 | 0.16 | 0.31 | 1.67 | 2.23 | 0.000 | 0.000 |
