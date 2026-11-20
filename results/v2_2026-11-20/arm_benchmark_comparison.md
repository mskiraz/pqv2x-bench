# Receiver-path benchmark across platforms

All times are **measurements** on the named platform (microseconds). Ratios are per-platform measurements of a path's time divided by the same platform's primitive verification time. Absolute times from different machines, operating systems or garbage-collector settings are not directly comparable; ratios are compared only to test whether they transfer between platforms.

## Runs

| Tag | Platform | CPU | OS | Go | GC | GOMAXPROCS | CPU features | Iterations (SLH-DSA) |
|---|---|---|---|---|---|---:|---|---|
| intel-i9-10980hk-run1 | windows/amd64 | Intel64 Family 6 Model 165 Stepping 2, GenuineIntel | Windows (version not recorded) | go1.25.6 | enabled (not recorded; original runs) | not recorded | not recorded | 20000 (2000) |
| intel-i9-10980hk-run2 | windows/amd64 | Intel64 Family 6 Model 165 Stepping 2, GenuineIntel | Windows (version not recorded) | go1.25.6 | enabled (not recorded; original runs) | not recorded | not recorded | 20000 (2000) |
| gha-ubuntu-24.04-run44822427947-1-gcoff | linux/amd64 | Intel(R) Xeon(R) 6975P-C | Ubuntu 24.04.5 LTS (kernel 7.0.0-1014-aws) | go1.25.6 | disabled (only the explicit runtime.GC() before each timed loop) | 4 | aes, avx, avx2, avx512f, bmi2, sha_ni | 20000 (2000) |
| gha-ubuntu-24.04-run44822427947-1 | linux/amd64 | Intel(R) Xeon(R) 6975P-C | Ubuntu 24.04.5 LTS (kernel 7.0.0-1014-aws) | go1.25.6 | enabled (GC percent 100) | 4 | aes, avx, avx2, avx512f, bmi2, sha_ni | 20000 (2000) |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | linux/arm64 | ARM implementer 0x41, part 0xd49 | Ubuntu 24.04.5 LTS (kernel 7.0.0-1017-azure) | go1.25.6 | disabled (only the explicit runtime.GC() before each timed loop) | 4 | aes, asimd, atomics, sha2, sha3, sha512 | 20000 (2000) |
| gha-ubuntu-24.04-arm-run44822427947-1 | linux/arm64 | ARM implementer 0x41, part 0xd49 | Ubuntu 24.04.5 LTS (kernel 7.0.0-1017-azure) | go1.25.6 | enabled (GC percent 100) | 4 | aes, asimd, atomics, sha2, sha3, sha512 | 20000 (2000) |

## Falcon-512: median ratio to the primitive (measured per platform)

| Tag | decode | digest_hit | cert_known | cert_new | digest_unknown | Δ cert_known vs intel-i9-10980hk-run1 | Δ cert_new vs intel-i9-10980hk-run1 |
|---|---:|---:|---:|---:|---:|---:|---:|
| intel-i9-10980hk-run1 | 0.0019 | 1.0174 | 1.0983 | 2.2330 | 0.0021 | +0.0000 | +0.0000 |
| intel-i9-10980hk-run2 | 0.0019 | 1.0174 | 1.0965 | 2.2326 | 0.0021 | -0.0017 | -0.0004 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | 0.0022 | 1.0137 | 1.0431 | 2.1179 | 0.0016 | -0.0552 | -0.1151 |
| gha-ubuntu-24.04-run44822427947-1 | 0.0015 | 1.0111 | 1.0410 | 2.0611 | 0.0016 | -0.0572 | -0.1719 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | 0.0023 | 1.0077 | 1.0330 | 2.1032 | 0.0023 | -0.0652 | -0.1298 |
| gha-ubuntu-24.04-arm-run44822427947-1 | 0.0022 | 1.0076 | 1.0353 | 2.0823 | 0.0025 | -0.0630 | -0.1507 |

## Falcon-512: mean ratio to the primitive (measured per platform)

| Tag | decode | digest_hit | cert_known | cert_new | digest_unknown | Δ cert_known vs intel-i9-10980hk-run1 | Δ cert_new vs intel-i9-10980hk-run1 |
|---|---:|---:|---:|---:|---:|---:|---:|
| intel-i9-10980hk-run1 | 0.0022 | 1.0220 | 1.0979 | 2.2397 | 0.0023 | +0.0000 | +0.0000 |
| intel-i9-10980hk-run2 | 0.0022 | 1.0169 | 1.0966 | 2.2429 | 0.0023 | -0.0013 | +0.0032 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | 0.0020 | 1.0118 | 1.0447 | 2.1077 | 0.0016 | -0.0532 | -0.1321 |
| gha-ubuntu-24.04-run44822427947-1 | 0.0020 | 1.0094 | 1.0440 | 2.0638 | 0.0021 | -0.0539 | -0.1759 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | 0.0025 | 1.0085 | 1.0346 | 2.1125 | 0.0023 | -0.0634 | -0.1272 |
| gha-ubuntu-24.04-arm-run44822427947-1 | 0.0027 | 1.0090 | 1.0375 | 2.0929 | 0.0029 | -0.0604 | -0.1468 |

## Falcon-512 path times scaled to the published ARMv8 primitive (platform-scaled estimates, not measurements)

Each entry is a platform's measured median ratio multiplied by 0.446 ms. The spread across rows is the uncertainty that the ratio-transfer assumption introduces.

| Tag | digest_hit (ms) | cert_known (ms) | cert_new (ms) | max(2 x digest_hit, cert_new) (ms) |
|---|---:|---:|---:|---:|
| intel-i9-10980hk-run1 | 0.4538 | 0.4898 | 0.9959 | 0.9959 |
| intel-i9-10980hk-run2 | 0.4538 | 0.4890 | 0.9957 | 0.9957 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | 0.4521 | 0.4652 | 0.9446 | 0.9446 |
| gha-ubuntu-24.04-run44822427947-1 | 0.4510 | 0.4643 | 0.9192 | 0.9192 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | 0.4495 | 0.4607 | 0.9380 | 0.9380 |
| gha-ubuntu-24.04-arm-run44822427947-1 | 0.4494 | 0.4617 | 0.9287 | 0.9287 |

## All schemes: median ratio to the primitive

| Scheme | Tag | digest_hit | cert_known | cert_new |
|---|---|---:|---:|---:|
| ECDSA P-256 | intel-i9-10980hk-run1 | 1.012 | 1.023 | 2.160 |
| ECDSA P-256 | intel-i9-10980hk-run2 | 1.012 | 1.024 | 2.160 |
| ECDSA P-256 | gha-ubuntu-24.04-run44822427947-1-gcoff | 0.997 | 1.002 | 2.129 |
| ECDSA P-256 | gha-ubuntu-24.04-run44822427947-1 | 1.006 | 1.010 | 2.135 |
| ECDSA P-256 | gha-ubuntu-24.04-arm-run44822427947-1-gcoff | 1.009 | 1.014 | 2.136 |
| ECDSA P-256 | gha-ubuntu-24.04-arm-run44822427947-1 | 1.008 | 1.014 | 2.132 |
| Falcon-512 | intel-i9-10980hk-run1 | 1.017 | 1.098 | 2.233 |
| Falcon-512 | intel-i9-10980hk-run2 | 1.017 | 1.097 | 2.233 |
| Falcon-512 | gha-ubuntu-24.04-run44822427947-1-gcoff | 1.014 | 1.043 | 2.118 |
| Falcon-512 | gha-ubuntu-24.04-run44822427947-1 | 1.011 | 1.041 | 2.061 |
| Falcon-512 | gha-ubuntu-24.04-arm-run44822427947-1-gcoff | 1.008 | 1.033 | 2.103 |
| Falcon-512 | gha-ubuntu-24.04-arm-run44822427947-1 | 1.008 | 1.035 | 2.082 |
| ML-DSA-44 | intel-i9-10980hk-run1 | 1.001 | 1.097 | 2.230 |
| ML-DSA-44 | intel-i9-10980hk-run2 | 1.002 | 1.099 | 2.234 |
| ML-DSA-44 | gha-ubuntu-24.04-run44822427947-1-gcoff | 1.025 | 0.980 | 2.144 |
| ML-DSA-44 | gha-ubuntu-24.04-run44822427947-1 | 1.005 | 1.037 | 2.092 |
| ML-DSA-44 | gha-ubuntu-24.04-arm-run44822427947-1-gcoff | 1.021 | 1.020 | 2.096 |
| ML-DSA-44 | gha-ubuntu-24.04-arm-run44822427947-1 | 1.004 | 1.026 | 2.071 |
| SLH-DSA-SHA2-128s | intel-i9-10980hk-run1 | 1.002 | 1.041 | 2.078 |
| SLH-DSA-SHA2-128s | intel-i9-10980hk-run2 | 1.002 | 1.036 | 2.029 |
| SLH-DSA-SHA2-128s | gha-ubuntu-24.04-run44822427947-1-gcoff | 1.039 | 1.057 | 2.095 |
| SLH-DSA-SHA2-128s | gha-ubuntu-24.04-run44822427947-1 | 1.018 | 1.018 | 2.054 |
| SLH-DSA-SHA2-128s | gha-ubuntu-24.04-arm-run44822427947-1-gcoff | 1.009 | 1.032 | 2.048 |
| SLH-DSA-SHA2-128s | gha-ubuntu-24.04-arm-run44822427947-1 | 0.988 | 1.017 | 2.056 |

## Primitive verification time per platform (measured, microseconds)

| Tag | ECDSA P-256 median | ECDSA P-256 mean | Falcon-512 median | Falcon-512 mean | ML-DSA-44 median | ML-DSA-44 mean | SLH-DSA-SHA2-128s median | SLH-DSA-SHA2-128s mean | ML-DSA-44 / Falcon-512 (median) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| intel-i9-10980hk-run1 | 91.00 | 92.52 | 63.10 | 64.55 | 124.20 | 128.59 | 1184.90 | 1234.01 | 1.97 |
| intel-i9-10980hk-run2 | 91.00 | 93.25 | 63.20 | 64.61 | 124.10 | 128.31 | 1206.70 | 1250.62 | 1.96 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | 61.92 | 63.30 | 38.40 | 39.05 | 73.15 | 73.90 | 322.79 | 325.96 | 1.91 |
| gha-ubuntu-24.04-run44822427947-1 | 61.52 | 63.20 | 38.46 | 39.61 | 68.78 | 72.37 | 332.32 | 361.72 | 1.79 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | 68.83 | 69.60 | 40.15 | 40.36 | 88.25 | 89.61 | 348.49 | 363.91 | 2.20 |
| gha-ubuntu-24.04-arm-run44822427947-1 | 68.91 | 69.49 | 40.15 | 40.77 | 88.12 | 90.51 | 410.23 | 417.47 | 2.20 |

## Falcon-512 distribution per operation (measured, microseconds)

| Tag | Operation | n | Mean | Std | Median | p90 | p99 | p99.9 | Max |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| intel-i9-10980hk-run1 | primitive | 20000 | 64.55 | 14.90 | 63.10 | 64.80 | 81.70 | 311.10 | 760.30 |
| intel-i9-10980hk-run1 | decode | 20000 | 0.14 | 0.18 | 0.12 | 0.13 | 0.57 | 2.84 | 6.07 |
| intel-i9-10980hk-run1 | digest_hit | 20000 | 65.97 | 13.91 | 64.20 | 67.00 | 101.40 | 264.70 | 616.90 |
| intel-i9-10980hk-run1 | cert_known | 20000 | 70.87 | 14.43 | 69.30 | 72.00 | 97.60 | 285.50 | 777.30 |
| intel-i9-10980hk-run1 | cert_new | 20000 | 144.57 | 25.38 | 140.90 | 145.90 | 232.40 | 524.60 | 837.00 |
| intel-i9-10980hk-run1 | digest_unknown | 20000 | 0.15 | 0.15 | 0.13 | 0.14 | 0.35 | 2.30 | 4.70 |
| intel-i9-10980hk-run2 | primitive | 20000 | 64.61 | 14.35 | 63.20 | 65.40 | 85.30 | 275.50 | 792.70 |
| intel-i9-10980hk-run2 | decode | 20000 | 0.14 | 0.18 | 0.12 | 0.13 | 0.65 | 2.88 | 4.81 |
| intel-i9-10980hk-run2 | digest_hit | 20000 | 65.70 | 14.75 | 64.30 | 66.10 | 90.60 | 303.60 | 741.70 |
| intel-i9-10980hk-run2 | cert_known | 20000 | 70.85 | 14.93 | 69.30 | 71.60 | 97.50 | 298.30 | 771.90 |
| intel-i9-10980hk-run2 | cert_new | 20000 | 144.90 | 25.57 | 141.10 | 145.90 | 236.80 | 517.00 | 845.00 |
| intel-i9-10980hk-run2 | digest_unknown | 20000 | 0.15 | 0.16 | 0.13 | 0.14 | 0.31 | 2.52 | 5.30 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | primitive | 20000 | 39.05 | 4.62 | 38.40 | 39.85 | 47.98 | 62.32 | 329.12 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | decode | 20000 | 0.08 | 0.04 | 0.09 | 0.10 | 0.21 | 0.34 | 2.05 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | digest_hit | 20000 | 39.51 | 2.39 | 38.92 | 40.38 | 48.21 | 57.16 | 145.35 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | cert_known | 20000 | 40.80 | 3.66 | 40.05 | 41.55 | 50.19 | 79.39 | 231.36 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | cert_new | 20000 | 82.30 | 5.63 | 81.32 | 87.28 | 95.63 | 135.05 | 275.21 |
| gha-ubuntu-24.04-run44822427947-1-gcoff | digest_unknown | 20000 | 0.06 | 0.03 | 0.06 | 0.06 | 0.15 | 0.24 | 2.72 |
| gha-ubuntu-24.04-run44822427947-1 | primitive | 20000 | 39.61 | 9.57 | 38.46 | 40.24 | 53.28 | 173.70 | 565.53 |
| gha-ubuntu-24.04-run44822427947-1 | decode | 20000 | 0.08 | 0.13 | 0.06 | 0.10 | 0.50 | 1.62 | 5.74 |
| gha-ubuntu-24.04-run44822427947-1 | digest_hit | 20000 | 39.98 | 7.91 | 38.88 | 40.59 | 53.29 | 160.00 | 346.45 |
| gha-ubuntu-24.04-run44822427947-1 | cert_known | 20000 | 41.35 | 10.42 | 40.03 | 42.07 | 54.61 | 212.69 | 519.17 |
| gha-ubuntu-24.04-run44822427947-1 | cert_new | 20000 | 81.74 | 16.46 | 79.26 | 85.79 | 133.74 | 327.14 | 604.63 |
| gha-ubuntu-24.04-run44822427947-1 | digest_unknown | 20000 | 0.08 | 0.11 | 0.06 | 0.11 | 0.34 | 1.68 | 3.96 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | primitive | 20000 | 40.36 | 1.95 | 40.15 | 40.68 | 43.90 | 77.97 | 98.68 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | decode | 20000 | 0.10 | 0.04 | 0.09 | 0.11 | 0.23 | 0.53 | 3.13 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | digest_hit | 20000 | 40.71 | 1.95 | 40.47 | 41.04 | 46.12 | 78.47 | 87.88 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | cert_known | 20000 | 41.76 | 2.29 | 41.48 | 42.12 | 47.94 | 84.10 | 155.99 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | cert_new | 20000 | 85.27 | 4.93 | 84.45 | 87.83 | 94.86 | 147.08 | 252.55 |
| gha-ubuntu-24.04-arm-run44822427947-1-gcoff | digest_unknown | 20000 | 0.09 | 0.02 | 0.09 | 0.10 | 0.17 | 0.22 | 0.78 |
| gha-ubuntu-24.04-arm-run44822427947-1 | primitive | 20000 | 40.77 | 6.63 | 40.15 | 40.92 | 47.34 | 162.55 | 268.24 |
| gha-ubuntu-24.04-arm-run44822427947-1 | decode | 20000 | 0.11 | 0.13 | 0.09 | 0.14 | 0.67 | 1.86 | 2.85 |
| gha-ubuntu-24.04-arm-run44822427947-1 | digest_hit | 20000 | 41.14 | 7.05 | 40.45 | 41.25 | 48.70 | 163.60 | 294.70 |
| gha-ubuntu-24.04-arm-run44822427947-1 | cert_known | 20000 | 42.30 | 7.37 | 41.56 | 42.46 | 50.71 | 171.45 | 265.22 |
| gha-ubuntu-24.04-arm-run44822427947-1 | cert_new | 20000 | 85.33 | 12.62 | 83.59 | 86.36 | 135.88 | 258.21 | 366.04 |
| gha-ubuntu-24.04-arm-run44822427947-1 | digest_unknown | 20000 | 0.12 | 0.12 | 0.10 | 0.15 | 0.34 | 1.77 | 3.01 |

Full per-scheme statistics for every run: arm_benchmark_comparison.csv.
