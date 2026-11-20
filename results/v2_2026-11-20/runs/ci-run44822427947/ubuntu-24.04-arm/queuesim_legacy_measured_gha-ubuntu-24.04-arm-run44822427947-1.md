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

| Arrival model | Scenario | N | rho | Mean T (ms) | p99.9 T (ms) | Max T (ms) | Pr[T > 10 ms] | Configs with a miss | Discarded (unknown signer) |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| cv2x | steady, median work | 100 | 0.452 | 0.564 | 0.924 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 150 | 0.678 | 0.621 | 0.924 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 162 | 0.732 | 0.635 | 0.924 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 200 | 0.904 | 0.678 | 0.924 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 100 | 0.453 | 0.571 | 2.343 | 6.212 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 150 | 0.679 | 0.637 | 2.652 | 6.336 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 162 | 0.733 | 0.655 | 2.728 | 6.882 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 200 | 0.905 | 0.736 | 3.216 | 7.868 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 100 | 0.545 | 0.740 | 2.239 | 4.706 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 150 | 0.818 | 1.006 | 3.705 | 7.449 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 162 | 0.883 | 1.186 | 4.639 | 9.325 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 200 | 1.090 | unstable | | | | | |
| cv2x | no caching, measured jitter | 100 | 0.547 | 0.752 | 3.216 | 9.559 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, measured jitter | 150 | 0.820 | 1.050 | 4.636 | 12.151 | 2.27e-07 | 0.0000 | 0.000 |
| cv2x | no caching, measured jitter | 162 | 0.886 | 1.261 | 5.703 | 12.610 | 1.85e-06 | 0.0004 | 0.000 |
| cv2x | no caching, measured jitter | 200 | 1.093 | unstable | | | | | |
| cv2x | cold start, median work | 100 | 0.452 | 0.636 | 2.135 | 4.430 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 150 | 0.678 | 0.733 | 2.921 | 7.014 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 162 | 0.732 | 0.773 | 3.511 | 8.667 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 200 | 0.904 | 2.172 | 15.453 | 23.470 | 3.49e-02 | 0.6659 | 0.200 |
| cv2x | cold start, measured jitter | 100 | 0.453 | 0.644 | 2.771 | 7.410 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, measured jitter | 150 | 0.679 | 0.753 | 3.591 | 10.173 | 1.67e-08 | 0.0000 | 0.200 |
| cv2x | cold start, measured jitter | 162 | 0.733 | 0.801 | 4.237 | 11.164 | 4.63e-08 | 0.0000 | 0.200 |
| cv2x | cold start, measured jitter | 200 | 0.905 | 2.334 | 16.878 | 27.192 | 4.43e-02 | 0.7015 | 0.200 |

