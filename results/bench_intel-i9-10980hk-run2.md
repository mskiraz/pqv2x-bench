# Receiver authentication benchmark (intel-i9-10980hk-run2)

- CPU: Intel64 Family 6 Model 165 Stepping 2, GenuineIntel (16 logical CPUs), windows/amd64, go1.25.6
- 20000 timed iterations per operation, single goroutine, garbage collector enabled
- Run: 2026-11-19T15:27:04Z

## Encoded sizes (bytes)

| Scheme | Public key | Signature | Certificate | Digest SPDU | Certificate SPDU |
|---|---:|---:|---:|---:|---:|
| ECDSA P-256 | 33 | 64 | 144 | 132 | 270 |
| Falcon-512 | 897 | 666 | 1615 | 736 | 2345 |
| ML-DSA-44 | 1312 | 2420 | 3784 | 2490 | 6268 |
| SLH-DSA-SHA2-128s | 32 | 7856 | 7938 | 7926 | 15858 |

## Time per message (microseconds)

| Scheme | Operation | Mean | Std | Median | p99 | p99.9 | Max | Mean / primitive | Median / primitive |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| ECDSA P-256 | primitive | 93.25 | 13.91 | 91.00 | 155.00 | 230.70 | 928.10 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.13 | 0.15 | 0.12 | 0.36 | 2.61 | 5.21 | 0.001 | 0.001 |
| ECDSA P-256 | digest_hit | 94.14 | 12.06 | 92.10 | 151.40 | 227.70 | 780.70 | 1.010 | 1.012 |
| ECDSA P-256 | cert_known | 95.97 | 15.97 | 93.20 | 159.30 | 257.20 | 854.50 | 1.029 | 1.024 |
| ECDSA P-256 | cert_new | 200.40 | 25.23 | 196.60 | 285.50 | 592.40 | 1070.40 | 2.149 | 2.160 |
| ECDSA P-256 | digest_unknown | 0.15 | 0.17 | 0.13 | 0.56 | 2.82 | 5.11 | 0.002 | 0.001 |
| Falcon-512 | primitive | 64.61 | 14.35 | 63.20 | 85.30 | 275.50 | 792.70 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.14 | 0.18 | 0.12 | 0.65 | 2.88 | 4.81 | 0.002 | 0.002 |
| Falcon-512 | digest_hit | 65.70 | 14.75 | 64.30 | 90.60 | 303.60 | 741.70 | 1.017 | 1.017 |
| Falcon-512 | cert_known | 70.85 | 14.93 | 69.30 | 97.50 | 298.30 | 771.90 | 1.097 | 1.097 |
| Falcon-512 | cert_new | 144.90 | 25.57 | 141.10 | 236.80 | 517.00 | 845.00 | 2.243 | 2.233 |
| Falcon-512 | digest_unknown | 0.15 | 0.16 | 0.13 | 0.31 | 2.52 | 5.30 | 0.002 | 0.002 |
| ML-DSA-44 | primitive | 128.31 | 26.58 | 124.10 | 234.60 | 492.40 | 1001.30 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.13 | 0.12 | 0.12 | 0.25 | 1.95 | 4.29 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 128.66 | 25.47 | 124.40 | 238.40 | 492.00 | 825.00 | 1.003 | 1.002 |
| ML-DSA-44 | cert_known | 140.55 | 26.50 | 136.40 | 244.50 | 512.30 | 847.40 | 1.095 | 1.099 |
| ML-DSA-44 | cert_new | 285.54 | 42.95 | 277.30 | 496.40 | 781.00 | 1357.40 | 2.225 | 2.234 |
| ML-DSA-44 | digest_unknown | 0.14 | 0.12 | 0.13 | 0.25 | 2.11 | 4.25 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 1250.62 | 147.86 | 1206.70 | 1750.40 | 1924.30 | 2163.60 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.14 | 0.16 | 0.12 | 0.26 | 2.52 | 3.51 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 1257.43 | 156.91 | 1209.00 | 1796.00 | 1932.10 | 2225.40 | 1.005 | 1.002 |
| SLH-DSA-SHA2-128s | cert_known | 1283.42 | 151.18 | 1250.30 | 1794.90 | 2196.10 | 2361.00 | 1.026 | 1.036 |
| SLH-DSA-SHA2-128s | cert_new | 2546.31 | 222.51 | 2448.30 | 3140.70 | 4337.90 | 4579.20 | 2.036 | 2.029 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.15 | 0.18 | 0.13 | 0.89 | 2.61 | 3.88 | 0.000 | 0.000 |
