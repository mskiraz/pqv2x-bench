# Queue simulation from measured path times (intel-i9-10980hk-run1)

Measured Falcon-512 primitive: mean 64.55 us, median 63.10 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0174 | 1.0220 | 0.4538 |
| cert_known | 1.0983 | 1.0979 | 0.4898 |
| cert_new | 2.2330 | 2.2397 | 0.9959 |
| digest_unknown | 0.0021 | 0.0023 | 0.0009 |

C-V2X per-message budget T_sf/m = 0.500 ms, i.e. 1.121 x primitive.

## C-V2X SPS, Falcon-512, measured receiver paths

Arrival model: legacy abstract receiver-arrival model; every generated message is assumed to reach the receiver, and the radio cost of certificate-carrying messages is not modelled.

Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = 10 ms.

| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |
|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|
| cv2x | steady, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.576 | 0.944 | 0.980 | 0.980 | n/a | 0.461 | 2 |
| cv2x | steady, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.634 | 0.980 | 0.980 | 0.980 | n/a | 0.691 | 2 |
| cv2x | steady, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.647 | 0.980 | 0.980 | 0.980 | n/a | 0.747 | 2 |
| cv2x | steady, median work | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 0 | 0.00e+00 | 1.000000 | 0.691 | 0.980 | 0.980 | 0.980 | n/a | 0.922 | 2 |
| cv2x | steady, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.587 | 1.274 | 2.788 | 9.742 | n/a | 0.463 | 12 |
| cv2x | steady, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.661 | 1.730 | 3.474 | 9.023 | n/a | 0.694 | 15 |
| cv2x | steady, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.683 | 1.881 | 3.688 | 9.001 | n/a | 0.749 | 17 |
| cv2x | steady, measured jitter | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 129 | 6.45e-07 | 0.999999 | 0.809 | 2.906 | 4.590 | 12.310 | n/a | 0.925 | 25 |
| cv2x | no caching, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.777 | 1.992 | 2.694 | 5.324 | n/a | 0.562 | 9 |
| cv2x | no caching, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 1.146 | 3.337 | 4.590 | 9.154 | n/a | 0.843 | 16 |
| cv2x | no caching, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 472 | 2.91e-06 | 0.999997 | 1.463 | 4.494 | 6.133 | 12.075 | n/a | 0.911 | 22 |
| cv2x | no caching, median work | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.124, not simulated) | | | | | | | | |
| cv2x | no caching, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 4 | 4.00e-08 | 1.000000 | 0.798 | 2.354 | 4.109 | 10.956 | n/a | 0.564 | 13 |
| cv2x | no caching, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 2175 | 1.45e-05 | 0.999985 | 1.226 | 4.182 | 6.267 | 14.925 | n/a | 0.847 | 24 |
| cv2x | no caching, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 23547 | 1.45e-04 | 0.999855 | 1.615 | 5.645 | 8.074 | 17.150 | n/a | 0.914 | 29 |
| cv2x | no caching, measured jitter | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.129, not simulated) | | | | | | | | |
| cv2x | cold start, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.658 | 1.900 | 2.350 | 5.220 | n/a | 0.421 | 8 |
| cv2x | cold start, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.775 | 2.350 | 3.510 | 8.616 | n/a | 0.632 | 14 |
| cv2x | cold start, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 50 | 3.86e-07 | 0.799968 | 0.834 | 2.922 | 4.548 | 11.243 | n/a | 0.682 | 19 |
| cv2x | cold start, median work | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 39994256 | 160005744 | 19108857 | 1.19e-01 | 0.704484 | 3.456 | 16.719 | 20.055 | 28.620 | n/a | 0.842 | 58 |
| cv2x | cold start, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 19999903 | 80000097 | 1 | 1.25e-08 | 0.800001 | 0.671 | 1.956 | 3.424 | 10.733 | n/a | 0.423 | 13 |
| cv2x | cold start, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 29997779 | 120002221 | 277 | 2.31e-06 | 0.800013 | 0.809 | 2.810 | 4.705 | 14.525 | n/a | 0.634 | 23 |
| cv2x | cold start, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 32406775 | 129593225 | 1705 | 1.32e-05 | 0.799948 | 0.883 | 3.507 | 5.793 | 16.476 | n/a | 0.684 | 27 |
| cv2x | cold start, measured jitter | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 39991870 | 160008130 | 22519312 | 1.41e-01 | 0.687444 | 3.870 | 18.320 | 22.510 | 33.595 | n/a | 0.845 | 68 |

