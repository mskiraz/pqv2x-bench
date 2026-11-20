# Receiver queue: legacy abstract C-V2X model vs certificate-aware scheduling

All values are **simulation outputs** (no radio propagation, collisions, half-duplex loss or retransmissions). Service times are either fixed multiples of the published ARMv8 Falcon-512 primitive time S = 0.446 ms (`fixed`) or measured path ratios scaled to S (`measured`, a platform-scaled estimate). Both models use the same random configurations (seed and draw order are shared).

- **Legacy (`cv2x`)**: up to two messages per 1-ms subframe whatever their form; every generated message is delivered.
- **Certificate-aware (`cv2x-cert`)**: digest-form message 4 of 10 subchannels, certificate-carrying message all 10 (never shares its subframe); idealised earliest-fit scheduling within 100 ms; messages that do not fit are counted as not scheduled.

Denominators: Pr[miss | auth] = deadline missed / authenticated. On-time / generated = authenticated within the deadline / generated. Discarded = unknown-signer discards / generated.

| Service times | Benchmark tag | Scenario | N | Legacy: max T (ms) | Legacy: Pr[miss \| auth] | Legacy: on-time / generated | Cert-aware: offered occupancy | Cert-aware: not scheduled (digest / cert) | Cert-aware: discarded / generated | Cert-aware: max T (ms) | Cert-aware: p99.9 T (ms) | Cert-aware: Pr[miss \| auth] | Cert-aware: on-time / generated | Cert-aware: max scheduling delay (ms) |
|---|---|---|---:|---:|---:|---:|---:|---|---:|---:|---:|---:|---:|---:|
| measured | gha-ubuntu-24.04-run44822427947-1 | steady, median work | 100 | 0.9286 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.9019 | 0.9020 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | steady, median work | 150 | 0.9286 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.9019 | 0.9020 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | steady, median work | 162 | 0.9286 | 0.000000e+00 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 0.9019 | 0.9020 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | steady, median work | 200 | 0.9286 | 0.000000e+00 | 1.000000 | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.9019 | 0.9020 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | steady, measured jitter | 100 | 8.2261 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 8.4701 | 2.7030 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | steady, measured jitter | 150 | 10.9134 | 2.000000e-08 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 9.9390 | 3.2910 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | steady, measured jitter | 162 | 9.3522 | 0.000000e+00 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 8.9265 | 3.4340 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | steady, measured jitter | 200 | 10.7359 | 9.500000e-08 | 1.000000 | 1.0400 | 10664208 / 14667896 | 0.0000 | 9.0360 | 3.5750 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | no caching, median work | 100 | 4.6260 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.9192 | 0.9200 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | no caching, median work | 150 | 7.2597 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.9192 | 0.9200 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | no caching, median work | 162 | 9.0691 | 0.000000e+00 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 0.9192 | 0.9200 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | no caching, median work | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.9192 | 0.9200 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | no caching, measured jitter | 100 | 10.3230 | 1.000000e-08 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 11.1450 | 3.7610 | 4.000000e-08 | 1.000000 | 6.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | no caching, measured jitter | 150 | 13.5107 | 5.973333e-06 | 0.999994 | 0.7800 | 0 / 0 | 0.0000 | 11.5282 | 4.7000 | 2.800000e-07 | 1.000000 | 11.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | no caching, measured jitter | 162 | 16.5574 | 3.037037e-05 | 0.999970 | 0.8424 | 0 / 0 | 0.0000 | 12.5070 | 5.0400 | 1.104938e-06 | 0.999999 | 16.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | no caching, measured jitter | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | 14.9958 | 4.9250 | 1.288159e-06 | 0.873338 | 99.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | cold start, median work | 100 | 4.3539 | 0.000000e+00 | 0.800052 | 0.5200 | 0 / 0 | 0.1999 | 0.9192 | 0.9200 | 0.000000e+00 | 0.800052 | 6.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | cold start, median work | 150 | 6.8548 | 0.000000e+00 | 0.800016 | 0.7800 | 0 / 0 | 0.2000 | 0.9192 | 0.9200 | 0.000000e+00 | 0.800016 | 11.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | cold start, median work | 162 | 8.4316 | 0.000000e+00 | 0.799969 | 0.8424 | 0 / 0 | 0.2000 | 0.9192 | 0.9200 | 0.000000e+00 | 0.799969 | 16.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | cold start, median work | 200 | 23.0341 | 3.275898e-02 | 0.773821 | 1.0400 | 5415614 / 7434715 | 0.2018 | 0.9192 | 0.9200 | 0.000000e+00 | 0.733973 | 99.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | cold start, measured jitter | 100 | 10.4719 | 3.749995e-08 | 0.800001 | 0.5200 | 0 / 0 | 0.1999 | 9.8562 | 3.1050 | 0.000000e+00 | 0.800052 | 6.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | cold start, measured jitter | 150 | 11.4005 | 3.666599e-07 | 0.800015 | 0.7800 | 0 / 0 | 0.2000 | 11.0175 | 3.7890 | 8.333165e-08 | 0.800016 | 11.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | cold start, measured jitter | 162 | 12.2041 | 2.168323e-06 | 0.799956 | 0.8424 | 0 / 0 | 0.2000 | 10.3086 | 3.9590 | 7.716351e-09 | 0.799969 | 16.0 |
| measured | gha-ubuntu-24.04-run44822427947-1 | cold start, measured jitter | 200 | 28.5352 | 3.871615e-02 | 0.769066 | 1.0400 | 5415614 / 7434715 | 0.2018 | 11.8046 | 4.0790 | 6.812239e-08 | 0.733973 | 99.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | steady, median work | 100 | 0.9235 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.8988 | 0.8990 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | steady, median work | 150 | 0.9235 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.8988 | 0.8990 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | steady, median work | 162 | 0.9235 | 0.000000e+00 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 0.8988 | 0.8990 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | steady, median work | 200 | 0.9235 | 0.000000e+00 | 1.000000 | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.8988 | 0.8990 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | steady, measured jitter | 100 | 6.2124 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 5.6484 | 2.2690 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | steady, measured jitter | 150 | 6.3365 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 6.5154 | 2.6020 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | steady, measured jitter | 162 | 6.8822 | 0.000000e+00 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 6.4562 | 2.6660 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | steady, measured jitter | 200 | 7.8679 | 0.000000e+00 | 1.000000 | 1.0400 | 10664208 / 14667896 | 0.0000 | 6.5067 | 2.8260 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | no caching, median work | 100 | 4.7065 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.9287 | 0.9290 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | no caching, median work | 150 | 7.4490 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.9287 | 0.9290 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | no caching, median work | 162 | 9.3248 | 0.000000e+00 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 0.9287 | 0.9290 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | no caching, median work | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.9287 | 0.9290 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | no caching, measured jitter | 100 | 9.5591 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 6.7991 | 2.7920 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | no caching, measured jitter | 150 | 12.1513 | 2.266667e-07 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 7.9933 | 3.2700 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | no caching, measured jitter | 162 | 12.6101 | 1.845679e-06 | 0.999998 | 0.8424 | 0 / 0 | 0.0000 | 7.7514 | 3.4490 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | no caching, measured jitter | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | 7.7676 | 3.4450 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | cold start, median work | 100 | 4.4297 | 0.000000e+00 | 0.800052 | 0.5200 | 0 / 0 | 0.1999 | 0.9287 | 0.9290 | 0.000000e+00 | 0.800052 | 6.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | cold start, median work | 150 | 7.0141 | 0.000000e+00 | 0.800016 | 0.7800 | 0 / 0 | 0.2000 | 0.9287 | 0.9290 | 0.000000e+00 | 0.800016 | 11.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | cold start, median work | 162 | 8.6668 | 0.000000e+00 | 0.799969 | 0.8424 | 0 / 0 | 0.2000 | 0.9287 | 0.9290 | 0.000000e+00 | 0.799969 | 16.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | cold start, median work | 200 | 23.4699 | 3.491298e-02 | 0.772097 | 1.0400 | 5415614 / 7434715 | 0.2018 | 0.9287 | 0.9290 | 0.000000e+00 | 0.733973 | 99.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | cold start, measured jitter | 100 | 7.4103 | 0.000000e+00 | 0.800001 | 0.5200 | 0 / 0 | 0.1999 | 6.7366 | 2.5410 | 0.000000e+00 | 0.800052 | 6.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | cold start, measured jitter | 150 | 10.1732 | 1.666636e-08 | 0.800015 | 0.7800 | 0 / 0 | 0.2000 | 6.7423 | 2.8230 | 0.000000e+00 | 0.800016 | 11.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | cold start, measured jitter | 162 | 11.1640 | 4.629872e-08 | 0.799958 | 0.8424 | 0 / 0 | 0.2000 | 6.7642 | 2.9470 | 0.000000e+00 | 0.799969 | 16.0 |
| measured | gha-ubuntu-24.04-arm-run44822427947-1 | cold start, measured jitter | 200 | 27.1922 | 4.426111e-02 | 0.764630 | 1.0400 | 5415614 / 7434715 | 0.2018 | 8.2963 | 3.0320 | 0.000000e+00 | 0.733973 | 99.0 |

Full accounting (generated, scheduled, discarded, authenticated, missed, utilisation, backlog) for every row: queue_comparison_runner_ratios.csv.
