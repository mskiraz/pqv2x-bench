# Queue simulation from measured path times (gha-ubuntu-24.04-arm-run44822427947-1)

Measured Falcon-512 primitive: mean 40.77 us, median 40.15 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0076 | 1.0090 | 0.4494 |
| cert_known | 1.0353 | 1.0375 | 0.4617 |
| cert_new | 2.0823 | 2.0929 | 0.9287 |
| digest_unknown | 0.0025 | 0.0029 | 0.0011 |

C-V2X per-message budget T_sf/m = 0.500 ms, i.e. 1.121 x primitive.

## C-V2X SPS, Falcon-512, measured receiver paths

Arrival model: legacy abstract receiver-arrival model; every generated message is assumed to reach the receiver, and the radio cost of certificate-carrying messages is not modelled.

Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = 10 ms.

| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |
|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|
| cv2x | steady, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.564 | 0.912 | 0.924 | 0.923 | n/a | 0.452 | 2 |
| cv2x | steady, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.621 | 0.924 | 0.924 | 0.923 | n/a | 0.678 | 2 |
| cv2x | steady, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.635 | 0.924 | 0.924 | 0.923 | n/a | 0.732 | 2 |
| cv2x | steady, median work | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 0 | 0.00e+00 | 1.000000 | 0.678 | 0.924 | 0.924 | 0.923 | n/a | 0.904 | 2 |
| cv2x | steady, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.571 | 1.099 | 2.343 | 6.212 | n/a | 0.453 | 9 |
| cv2x | steady, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.637 | 1.488 | 2.652 | 6.336 | n/a | 0.679 | 12 |
| cv2x | steady, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.655 | 1.598 | 2.728 | 6.882 | n/a | 0.733 | 13 |
| cv2x | steady, measured jitter | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 0 | 0.00e+00 | 1.000000 | 0.736 | 2.152 | 3.216 | 7.868 | n/a | 0.905 | 16 |
| cv2x | no caching, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.740 | 1.858 | 2.239 | 4.706 | n/a | 0.545 | 8 |
| cv2x | no caching, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 1.006 | 2.718 | 3.705 | 7.449 | n/a | 0.818 | 13 |
| cv2x | no caching, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 1.186 | 3.403 | 4.639 | 9.325 | n/a | 0.883 | 16 |
| cv2x | no caching, median work | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.090, not simulated) | | | | | | | | |
| cv2x | no caching, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.752 | 2.075 | 3.216 | 9.559 | n/a | 0.547 | 15 |
| cv2x | no caching, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 34 | 2.27e-07 | 1.000000 | 1.050 | 3.228 | 4.636 | 12.151 | n/a | 0.820 | 18 |
| cv2x | no caching, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 299 | 1.85e-06 | 0.999998 | 1.261 | 4.035 | 5.703 | 12.610 | n/a | 0.886 | 22 |
| cv2x | no caching, measured jitter | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.093, not simulated) | | | | | | | | |
| cv2x | cold start, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.636 | 1.686 | 2.135 | 4.430 | n/a | 0.409 | 7 |
| cv2x | cold start, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.733 | 2.034 | 2.921 | 7.014 | n/a | 0.613 | 12 |
| cv2x | cold start, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 0 | 0.00e+00 | 0.799969 | 0.773 | 2.340 | 3.511 | 8.667 | n/a | 0.662 | 15 |
| cv2x | cold start, median work | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 39994256 | 160005744 | 5586277 | 3.49e-02 | 0.772097 | 2.172 | 12.395 | 15.453 | 23.470 | n/a | 0.818 | 47 |
| cv2x | cold start, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 19999903 | 80000097 | 0 | 0.00e+00 | 0.800001 | 0.644 | 1.825 | 2.771 | 7.410 | n/a | 0.410 | 10 |
| cv2x | cold start, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 29997779 | 120002221 | 2 | 1.67e-08 | 0.800015 | 0.753 | 2.310 | 3.591 | 10.173 | n/a | 0.615 | 15 |
| cv2x | cold start, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 32406775 | 129593225 | 6 | 4.63e-08 | 0.799958 | 0.801 | 2.686 | 4.237 | 11.164 | n/a | 0.664 | 18 |
| cv2x | cold start, measured jitter | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 39991870 | 160008130 | 7082138 | 4.43e-02 | 0.764630 | 2.334 | 13.292 | 16.878 | 27.192 | n/a | 0.820 | 55 |

