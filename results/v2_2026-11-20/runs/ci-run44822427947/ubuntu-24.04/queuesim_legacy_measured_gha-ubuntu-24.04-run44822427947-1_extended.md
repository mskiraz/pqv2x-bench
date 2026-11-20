# Queue simulation from measured path times (gha-ubuntu-24.04-run44822427947-1)

Measured Falcon-512 primitive: mean 39.61 us, median 38.46 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0111 | 1.0094 | 0.4510 |
| cert_known | 1.0410 | 1.0440 | 0.4643 |
| cert_new | 2.0611 | 2.0638 | 0.9192 |
| digest_unknown | 0.0016 | 0.0021 | 0.0007 |

C-V2X per-message budget T_sf/m = 0.500 ms, i.e. 1.121 x primitive.

## C-V2X SPS, Falcon-512, measured receiver paths

Arrival model: legacy abstract receiver-arrival model; every generated message is assumed to reach the receiver, and the radio cost of certificate-carrying messages is not modelled.

Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = 10 ms.

| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |
|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|
| cv2x | steady, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.566 | 0.916 | 0.929 | 0.929 | n/a | 0.454 | 2 |
| cv2x | steady, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.623 | 0.929 | 0.929 | 0.929 | n/a | 0.680 | 2 |
| cv2x | steady, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.637 | 0.929 | 0.929 | 0.929 | n/a | 0.735 | 2 |
| cv2x | steady, median work | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 0 | 0.00e+00 | 1.000000 | 0.680 | 0.929 | 0.929 | 0.929 | n/a | 0.907 | 2 |
| cv2x | steady, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.574 | 1.201 | 2.853 | 8.226 | n/a | 0.453 | 12 |
| cv2x | steady, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 3 | 2.00e-08 | 1.000000 | 0.645 | 1.695 | 3.393 | 10.913 | n/a | 0.680 | 16 |
| cv2x | steady, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.665 | 1.843 | 3.527 | 9.352 | n/a | 0.734 | 18 |
| cv2x | steady, measured jitter | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 19 | 9.50e-08 | 1.000000 | 0.766 | 2.671 | 4.308 | 10.736 | n/a | 0.906 | 22 |
| cv2x | no caching, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.737 | 1.839 | 2.209 | 4.626 | n/a | 0.545 | 8 |
| cv2x | no caching, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.996 | 2.673 | 3.609 | 7.260 | n/a | 0.817 | 12 |
| cv2x | no caching, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 1.170 | 3.318 | 4.519 | 9.069 | n/a | 0.882 | 16 |
| cv2x | no caching, median work | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.089, not simulated) | | | | | | | | |
| cv2x | no caching, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 1 | 1.00e-08 | 1.000000 | 0.751 | 2.119 | 4.013 | 10.323 | n/a | 0.544 | 14 |
| cv2x | no caching, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 896 | 5.97e-06 | 0.999994 | 1.047 | 3.488 | 5.733 | 13.511 | n/a | 0.816 | 23 |
| cv2x | no caching, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 4920 | 3.04e-05 | 0.999970 | 1.252 | 4.374 | 6.741 | 16.557 | n/a | 0.882 | 29 |
| cv2x | no caching, measured jitter | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.088, not simulated) | | | | | | | | |
| cv2x | cold start, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.636 | 1.661 | 2.111 | 4.354 | n/a | 0.409 | 7 |
| cv2x | cold start, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.732 | 2.013 | 2.864 | 6.855 | n/a | 0.614 | 12 |
| cv2x | cold start, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 0 | 0.00e+00 | 0.799969 | 0.771 | 2.285 | 3.425 | 8.432 | n/a | 0.663 | 15 |
| cv2x | cold start, median work | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 39994256 | 160005744 | 5241625 | 3.28e-02 | 0.773821 | 2.157 | 12.196 | 15.189 | 23.034 | n/a | 0.818 | 47 |
| cv2x | cold start, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 19999903 | 80000097 | 3 | 3.75e-08 | 0.800001 | 0.645 | 1.803 | 3.323 | 10.472 | n/a | 0.409 | 14 |
| cv2x | cold start, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 29997779 | 120002221 | 44 | 3.67e-07 | 0.800015 | 0.756 | 2.427 | 4.314 | 11.400 | n/a | 0.613 | 19 |
| cv2x | cold start, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 32406775 | 129593225 | 281 | 2.17e-06 | 0.799956 | 0.805 | 2.831 | 4.939 | 12.204 | n/a | 0.662 | 21 |
| cv2x | cold start, measured jitter | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 39991870 | 160008130 | 6194898 | 3.87e-02 | 0.769066 | 2.262 | 13.075 | 17.058 | 28.535 | n/a | 0.818 | 58 |

