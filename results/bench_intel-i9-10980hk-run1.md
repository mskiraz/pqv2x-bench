# Receiver authentication benchmark (intel-i9-10980hk-run1)

- CPU: Intel64 Family 6 Model 165 Stepping 2, GenuineIntel (16 logical CPUs), windows/amd64, go1.25.6
- 20000 timed iterations per operation, single goroutine, garbage collector enabled
- Run: 2026-11-19T15:25:49Z

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
| ECDSA P-256 | primitive | 92.52 | 12.56 | 91.00 | 115.80 | 230.70 | 803.20 | 1.000 | 1.000 |
| ECDSA P-256 | decode | 0.13 | 0.15 | 0.12 | 0.42 | 2.37 | 4.76 | 0.001 | 0.001 |
| ECDSA P-256 | digest_hit | 93.70 | 11.64 | 92.10 | 120.00 | 238.10 | 800.60 | 1.013 | 1.012 |
| ECDSA P-256 | cert_known | 95.06 | 14.63 | 93.10 | 156.30 | 247.20 | 803.40 | 1.027 | 1.023 |
| ECDSA P-256 | cert_new | 200.05 | 21.66 | 196.60 | 274.00 | 540.30 | 1008.70 | 2.162 | 2.160 |
| ECDSA P-256 | digest_unknown | 0.15 | 0.18 | 0.13 | 0.75 | 2.68 | 10.76 | 0.002 | 0.001 |
| Falcon-512 | primitive | 64.55 | 14.90 | 63.10 | 81.70 | 311.10 | 760.30 | 1.000 | 1.000 |
| Falcon-512 | decode | 0.14 | 0.18 | 0.12 | 0.57 | 2.84 | 6.07 | 0.002 | 0.002 |
| Falcon-512 | digest_hit | 65.97 | 13.91 | 64.20 | 101.40 | 264.70 | 616.90 | 1.022 | 1.017 |
| Falcon-512 | cert_known | 70.87 | 14.43 | 69.30 | 97.60 | 285.50 | 777.30 | 1.098 | 1.098 |
| Falcon-512 | cert_new | 144.57 | 25.38 | 140.90 | 232.40 | 524.60 | 837.00 | 2.240 | 2.233 |
| Falcon-512 | digest_unknown | 0.15 | 0.15 | 0.13 | 0.35 | 2.30 | 4.70 | 0.002 | 0.002 |
| ML-DSA-44 | primitive | 128.59 | 26.91 | 124.20 | 231.00 | 501.60 | 833.40 | 1.000 | 1.000 |
| ML-DSA-44 | decode | 0.13 | 0.13 | 0.12 | 0.26 | 2.23 | 4.00 | 0.001 | 0.001 |
| ML-DSA-44 | digest_hit | 128.77 | 26.78 | 124.30 | 236.40 | 484.20 | 907.50 | 1.001 | 1.001 |
| ML-DSA-44 | cert_known | 140.68 | 28.13 | 136.20 | 243.60 | 543.80 | 870.00 | 1.094 | 1.097 |
| ML-DSA-44 | cert_new | 285.19 | 41.86 | 277.00 | 496.30 | 743.60 | 993.70 | 2.218 | 2.230 |
| ML-DSA-44 | digest_unknown | 0.15 | 0.15 | 0.13 | 0.28 | 2.06 | 11.22 | 0.001 | 0.001 |
| SLH-DSA-SHA2-128s | primitive | 1234.01 | 144.55 | 1184.90 | 1720.80 | 2076.10 | 2209.30 | 1.000 | 1.000 |
| SLH-DSA-SHA2-128s | decode | 0.13 | 0.14 | 0.12 | 0.30 | 2.23 | 2.46 | 0.000 | 0.000 |
| SLH-DSA-SHA2-128s | digest_hit | 1239.87 | 153.97 | 1187.00 | 1789.80 | 2214.80 | 2268.30 | 1.005 | 1.002 |
| SLH-DSA-SHA2-128s | cert_known | 1274.04 | 137.17 | 1233.40 | 1766.90 | 1864.00 | 2118.90 | 1.032 | 1.041 |
| SLH-DSA-SHA2-128s | cert_new | 2553.78 | 230.83 | 2461.90 | 3218.60 | 4206.70 | 4350.30 | 2.070 | 2.078 |
| SLH-DSA-SHA2-128s | digest_unknown | 0.15 | 0.15 | 0.13 | 0.27 | 2.43 | 3.62 | 0.000 | 0.000 |
