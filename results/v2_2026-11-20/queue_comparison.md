# Receiver queue: legacy abstract C-V2X model vs certificate-aware scheduling

All values are **simulation outputs** (no radio propagation, collisions, half-duplex loss or retransmissions). Service times are either fixed multiples of the published ARMv8 Falcon-512 primitive time S = 0.446 ms (`fixed`) or measured path ratios scaled to S (`measured`, a platform-scaled estimate). Both models use the same random configurations (seed and draw order are shared).

- **Legacy (`cv2x`)**: up to two messages per 1-ms subframe whatever their form; every generated message is delivered.
- **Certificate-aware (`cv2x-cert`)**: digest-form message 4 of 10 subchannels, certificate-carrying message all 10 (never shares its subframe); idealised earliest-fit scheduling within 100 ms; messages that do not fit are counted as not scheduled.

Denominators: Pr[miss | auth] = deadline missed / authenticated. On-time / generated = authenticated within the deadline / generated. Discarded = unknown-signer discards / generated.

| Service times | Benchmark tag | Scenario | N | Legacy: max T (ms) | Legacy: Pr[miss \| auth] | Legacy: on-time / generated | Cert-aware: offered occupancy | Cert-aware: not scheduled (digest / cert) | Cert-aware: discarded / generated | Cert-aware: max T (ms) | Cert-aware: p99.9 T (ms) | Cert-aware: Pr[miss \| auth] | Cert-aware: on-time / generated | Cert-aware: max scheduling delay (ms) |
|---|---|---|---:|---:|---:|---:|---:|---|---:|---:|---:|---:|---:|---:|
| fixed |  | S | 100 | 0.8920 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.8920 | 0.8920 | 0.000000e+00 | 1.000000 | 6.0 |
| fixed |  | S | 150 | 0.8920 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.8920 | 0.8920 | 0.000000e+00 | 1.000000 | 11.0 |
| fixed |  | S | 200 | 0.8920 | 0.000000e+00 | 1.000000 | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.8920 | 0.8930 | 0.000000e+00 | 0.873339 | 99.0 |
| fixed |  | S/2S | 100 | 4.3660 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.8920 | 0.8920 | 0.000000e+00 | 1.000000 | 6.0 |
| fixed |  | S/2S | 150 | 6.5440 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.8920 | 0.8920 | 0.000000e+00 | 1.000000 | 11.0 |
| fixed |  | S/2S | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.8920 | 0.8930 | 0.000000e+00 | 0.873339 | 99.0 |
| fixed |  | 1.5S | 100 | 4.7320 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 4.4010 | 2.3670 | 0.000000e+00 | 1.000000 | 6.0 |
| fixed |  | 1.5S | 150 | unstable | – | – | 0.7800 | 0 / 0 | 0.0000 | unstable | unstable | unstable | 0.843162 | 11.0 |
| fixed |  | 1.5S | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | unstable | unstable | unstable | 0.000000 | 99.0 |
| fixed |  | 2S | 100 | 12.8440 | 5.610000e-05 | 0.999944 | 0.5200 | 0 / 0 | 0.0000 | 12.5200 | 7.8160 | 2.992000e-05 | 0.999970 | 6.0 |
| fixed |  | 2S | 150 | unstable | – | – | 0.7800 | 0 / 0 | 0.0000 | unstable | unstable | unstable | 0.000000 | 11.0 |
| fixed |  | 2S | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | unstable | unstable | unstable | 0.000000 | 99.0 |
| measured | intel-i9-10980hk-run1 | steady, median work | 100 | 0.9796 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.9075 | 0.9080 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | intel-i9-10980hk-run1 | steady, median work | 150 | 0.9796 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.9075 | 0.9080 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | intel-i9-10980hk-run1 | steady, median work | 162 | 0.9796 | 0.000000e+00 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 0.9075 | 0.9080 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | intel-i9-10980hk-run1 | steady, median work | 200 | 0.9796 | 0.000000e+00 | 1.000000 | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.9075 | 0.9080 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | intel-i9-10980hk-run1 | steady, measured jitter | 100 | 9.7415 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 8.0969 | 2.6890 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | intel-i9-10980hk-run1 | steady, measured jitter | 150 | 9.0233 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 8.9873 | 3.3820 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | intel-i9-10980hk-run1 | steady, measured jitter | 162 | 9.0007 | 0.000000e+00 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 9.2770 | 3.6120 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | intel-i9-10980hk-run1 | steady, measured jitter | 200 | 12.3097 | 6.450000e-07 | 0.999999 | 1.0400 | 10664208 / 14667896 | 0.0000 | 9.4038 | 3.7690 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | intel-i9-10980hk-run1 | no caching, median work | 100 | 5.3244 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.9959 | 0.9960 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | intel-i9-10980hk-run1 | no caching, median work | 150 | 9.1538 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.9959 | 0.9960 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | intel-i9-10980hk-run1 | no caching, median work | 162 | 12.0749 | 2.913580e-06 | 0.999997 | 0.8424 | 0 / 0 | 0.0000 | 0.9959 | 0.9960 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | intel-i9-10980hk-run1 | no caching, median work | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.9959 | 0.9960 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | intel-i9-10980hk-run1 | no caching, measured jitter | 100 | 10.9562 | 4.000000e-08 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 9.0776 | 3.7670 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | intel-i9-10980hk-run1 | no caching, measured jitter | 150 | 14.9252 | 1.450000e-05 | 0.999985 | 0.7800 | 0 / 0 | 0.0000 | 11.1684 | 4.7740 | 4.066667e-07 | 1.000000 | 11.0 |
| measured | intel-i9-10980hk-run1 | no caching, measured jitter | 162 | 17.1498 | 1.453519e-04 | 0.999855 | 0.8424 | 0 / 0 | 0.0000 | 13.0936 | 5.2410 | 2.358025e-06 | 0.999998 | 16.0 |
| measured | intel-i9-10980hk-run1 | no caching, measured jitter | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | 11.9034 | 5.1600 | 7.614450e-07 | 0.873339 | 99.0 |
| measured | intel-i9-10980hk-run1 | cold start, median work | 100 | 5.2197 | 0.000000e+00 | 0.800052 | 0.5200 | 0 / 0 | 0.1999 | 0.9959 | 0.9960 | 0.000000e+00 | 0.800052 | 6.0 |
| measured | intel-i9-10980hk-run1 | cold start, median work | 150 | 8.6155 | 0.000000e+00 | 0.800016 | 0.7800 | 0 / 0 | 0.2000 | 0.9959 | 0.9960 | 0.000000e+00 | 0.800016 | 11.0 |
| measured | intel-i9-10980hk-run1 | cold start, median work | 162 | 11.2431 | 3.858176e-07 | 0.799968 | 0.8424 | 0 / 0 | 0.2000 | 0.9959 | 0.9960 | 0.000000e+00 | 0.799969 | 16.0 |
| measured | intel-i9-10980hk-run1 | cold start, median work | 200 | 28.6196 | 1.194261e-01 | 0.704484 | 1.0400 | 5415614 / 7434715 | 0.2018 | 0.9959 | 0.9960 | 0.000000e+00 | 0.733973 | 99.0 |
| measured | intel-i9-10980hk-run1 | cold start, measured jitter | 100 | 10.7332 | 1.249998e-08 | 0.800001 | 0.5200 | 0 / 0 | 0.1999 | 10.5669 | 3.2350 | 2.499838e-08 | 0.800052 | 6.0 |
| measured | intel-i9-10980hk-run1 | cold start, measured jitter | 150 | 14.5247 | 2.308291e-06 | 0.800013 | 0.7800 | 0 / 0 | 0.2000 | 9.7051 | 3.8800 | 0.000000e+00 | 0.800016 | 11.0 |
| measured | intel-i9-10980hk-run1 | cold start, measured jitter | 162 | 16.4755 | 1.315655e-05 | 0.799948 | 0.8424 | 0 / 0 | 0.2000 | 11.4291 | 4.1270 | 4.784138e-07 | 0.799968 | 16.0 |
| measured | intel-i9-10980hk-run1 | cold start, measured jitter | 200 | 33.5952 | 1.407385e-01 | 0.687444 | 1.0400 | 5415614 / 7434715 | 0.2018 | 11.9266 | 4.3030 | 3.201752e-07 | 0.733973 | 99.0 |
| measured | intel-i9-10980hk-run2 | steady, median work | 100 | 0.9781 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.9075 | 0.9080 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | intel-i9-10980hk-run2 | steady, median work | 150 | 0.9781 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.9075 | 0.9080 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | intel-i9-10980hk-run2 | steady, median work | 162 | 0.9781 | 0.000000e+00 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 0.9075 | 0.9080 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | intel-i9-10980hk-run2 | steady, median work | 200 | 0.9781 | 0.000000e+00 | 1.000000 | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.9075 | 0.9080 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | intel-i9-10980hk-run2 | steady, measured jitter | 100 | 10.2404 | 1.000000e-08 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 9.6802 | 2.9840 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | intel-i9-10980hk-run2 | steady, measured jitter | 150 | 11.9810 | 1.000000e-07 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 9.7271 | 3.6570 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | intel-i9-10980hk-run2 | steady, measured jitter | 162 | 16.2556 | 2.777778e-07 | 1.000000 | 0.8424 | 0 / 0 | 0.0000 | 9.7306 | 3.9700 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | intel-i9-10980hk-run2 | steady, measured jitter | 200 | 12.5866 | 1.550000e-06 | 0.999998 | 1.0400 | 10664208 / 14667896 | 0.0000 | 11.2665 | 4.2850 | 1.259533e-07 | 0.873339 | 99.0 |
| measured | intel-i9-10980hk-run2 | no caching, median work | 100 | 5.3229 | 0.000000e+00 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 0.9957 | 0.9960 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | intel-i9-10980hk-run2 | no caching, median work | 150 | 9.1494 | 0.000000e+00 | 1.000000 | 0.7800 | 0 / 0 | 0.0000 | 0.9957 | 0.9960 | 0.000000e+00 | 1.000000 | 11.0 |
| measured | intel-i9-10980hk-run2 | no caching, median work | 162 | 12.0680 | 2.876543e-06 | 0.999997 | 0.8424 | 0 / 0 | 0.0000 | 0.9957 | 0.9960 | 0.000000e+00 | 1.000000 | 16.0 |
| measured | intel-i9-10980hk-run2 | no caching, median work | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | 0.9957 | 0.9960 | 0.000000e+00 | 0.873339 | 99.0 |
| measured | intel-i9-10980hk-run2 | no caching, measured jitter | 100 | 11.0954 | 7.000000e-08 | 1.000000 | 0.5200 | 0 / 0 | 0.0000 | 9.6802 | 3.7170 | 0.000000e+00 | 1.000000 | 6.0 |
| measured | intel-i9-10980hk-run2 | no caching, measured jitter | 150 | 15.2641 | 1.196000e-05 | 0.999988 | 0.7800 | 0 / 0 | 0.0000 | 11.4735 | 4.8800 | 7.866667e-07 | 0.999999 | 11.0 |
| measured | intel-i9-10980hk-run2 | no caching, measured jitter | 162 | 20.4122 | 1.534568e-04 | 0.999847 | 0.8424 | 0 / 0 | 0.0000 | 12.4789 | 5.3410 | 2.327160e-06 | 0.999998 | 16.0 |
| measured | intel-i9-10980hk-run2 | no caching, measured jitter | 200 | unstable | – | – | 1.0400 | 10664208 / 14667896 | 0.0000 | 14.8287 | 5.4010 | 4.179360e-06 | 0.873336 | 99.0 |
| measured | intel-i9-10980hk-run2 | cold start, median work | 100 | 5.2177 | 0.000000e+00 | 0.800052 | 0.5200 | 0 / 0 | 0.1999 | 0.9957 | 0.9960 | 0.000000e+00 | 0.800052 | 6.0 |
| measured | intel-i9-10980hk-run2 | cold start, median work | 150 | 8.6116 | 0.000000e+00 | 0.800016 | 0.7800 | 0 / 0 | 0.2000 | 0.9957 | 0.9960 | 0.000000e+00 | 0.800016 | 11.0 |
| measured | intel-i9-10980hk-run2 | cold start, median work | 162 | 11.2359 | 3.703849e-07 | 0.799968 | 0.8424 | 0 / 0 | 0.2000 | 0.9957 | 0.9960 | 0.000000e+00 | 0.799969 | 16.0 |
| measured | intel-i9-10980hk-run2 | cold start, median work | 200 | 28.6068 | 1.189057e-01 | 0.704901 | 1.0400 | 5415614 / 7434715 | 0.2018 | 0.9957 | 0.9960 | 0.000000e+00 | 0.733973 | 99.0 |
| measured | intel-i9-10980hk-run2 | cold start, measured jitter | 100 | 9.9765 | 0.000000e+00 | 0.800001 | 0.5200 | 0 / 0 | 0.1999 | 9.6802 | 3.2920 | 0.000000e+00 | 0.800052 | 6.0 |
| measured | intel-i9-10980hk-run2 | cold start, measured jitter | 150 | 14.0520 | 1.508305e-06 | 0.800014 | 0.7800 | 0 / 0 | 0.2000 | 10.1144 | 4.0300 | 2.499950e-08 | 0.800016 | 11.0 |
| measured | intel-i9-10980hk-run2 | cold start, measured jitter | 162 | 16.3907 | 1.596534e-05 | 0.799945 | 0.8424 | 0 / 0 | 0.2000 | 10.3698 | 4.3570 | 2.314905e-08 | 0.799969 | 16.0 |
| measured | intel-i9-10980hk-run2 | cold start, measured jitter | 200 | 34.4514 | 1.293324e-01 | 0.696569 | 1.0400 | 5415614 / 7434715 | 0.2018 | 11.3648 | 4.5980 | 4.223588e-07 | 0.733973 | 99.0 |

Full accounting (generated, scheduled, discarded, authenticated, missed, utilisation, backlog) for every row: queue_comparison.csv.
