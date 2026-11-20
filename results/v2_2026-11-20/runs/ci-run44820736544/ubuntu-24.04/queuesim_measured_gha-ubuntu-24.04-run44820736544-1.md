# Queue simulation from measured path times (gha-ubuntu-24.04-run44820736544-1)

Measured Falcon-512 primitive: mean 57.93 us, median 56.31 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0101 | 1.0089 | 0.4505 |
| cert_known | 1.0343 | 1.0352 | 0.4613 |
| cert_new | 2.0119 | 2.0532 | 0.8973 |
| digest_unknown | 0.0015 | 0.0020 | 0.0007 |

C-V2X per-message budget T_sf/m = 0.500 ms, i.e. 1.121 x primitive.

## C-V2X SPS, Falcon-512, measured receiver paths

| Arrival model | Scenario | N | rho | Mean T (ms) | p99.9 T (ms) | Max T (ms) | Pr[T > 10 ms] | Configs with a miss | Discarded (unknown signer) |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| cv2x | steady, median work | 100 | 0.453 | 0.565 | 0.923 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 150 | 0.679 | 0.622 | 0.923 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 162 | 0.733 | 0.636 | 0.923 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 200 | 0.905 | 0.679 | 0.923 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 100 | 0.452 | 0.570 | 2.117 | 6.082 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 150 | 0.678 | 0.636 | 2.402 | 5.808 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 162 | 0.733 | 0.654 | 2.479 | 6.224 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 200 | 0.905 | 0.730 | 2.909 | 6.986 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 100 | 0.540 | 0.726 | 2.143 | 4.427 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 150 | 0.810 | 0.962 | 3.383 | 6.747 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 162 | 0.875 | 1.110 | 4.161 | 8.231 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 200 | 1.080 | unstable | | | | | |
| cv2x | no caching, measured jitter | 100 | 0.543 | 0.746 | 3.045 | 15.365 | 5.67e-05 | 0.0092 | 0.000 |
| cv2x | no caching, measured jitter | 150 | 0.815 | 1.027 | 4.848 | 22.091 | 2.21e-04 | 0.0149 | 0.000 |
| cv2x | no caching, measured jitter | 162 | 0.880 | 1.219 | 6.461 | 19.995 | 4.01e-04 | 0.0168 | 0.000 |
| cv2x | no caching, measured jitter | 200 | 1.086 | unstable | | | | | |
| cv2x | cold start, median work | 100 | 0.453 | 0.631 | 2.040 | 4.178 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 150 | 0.679 | 0.722 | 2.688 | 6.407 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 162 | 0.733 | 0.757 | 3.185 | 7.686 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 200 | 0.905 | 1.924 | 13.931 | 21.485 | 1.92e-02 | 0.4719 | 0.200 |
| cv2x | cold start, measured jitter | 100 | 0.452 | 0.641 | 2.602 | 13.566 | 2.15e-05 | 0.0048 | 0.200 |
| cv2x | cold start, measured jitter | 150 | 0.678 | 0.746 | 3.467 | 16.283 | 5.39e-05 | 0.0073 | 0.200 |
| cv2x | cold start, measured jitter | 162 | 0.733 | 0.791 | 4.125 | 21.028 | 8.14e-05 | 0.0079 | 0.200 |
| cv2x | cold start, measured jitter | 200 | 0.905 | 2.160 | 16.460 | 30.638 | 3.33e-02 | 0.6004 | 0.200 |

