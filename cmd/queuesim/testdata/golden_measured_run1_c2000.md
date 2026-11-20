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

| Arrival model | Scenario | N | rho | Mean T (ms) | p99.9 T (ms) | Max T (ms) | Pr[T > 10 ms] | Configs with a miss | Discarded (unknown signer) |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| cv2x | steady, median work | 100 | 0.461 | 0.576 | 0.980 | 0.980 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 150 | 0.691 | 0.634 | 0.980 | 0.980 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 162 | 0.747 | 0.647 | 0.980 | 0.980 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 200 | 0.922 | 0.691 | 0.980 | 0.980 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 100 | 0.463 | 0.587 | 2.711 | 5.852 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 150 | 0.694 | 0.661 | 3.432 | 6.145 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 162 | 0.749 | 0.683 | 3.680 | 7.469 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 200 | 0.925 | 0.808 | 4.485 | 7.278 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 100 | 0.562 | 0.776 | 2.514 | 4.690 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 150 | 0.843 | 1.148 | 4.749 | 7.523 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 162 | 0.911 | 1.467 | 6.225 | 10.683 | 3.70e-06 | 0.0005 | 0.000 |
| cv2x | no caching, median work | 200 | 1.124 | unstable | | | | | |
| cv2x | no caching, measured jitter | 100 | 0.564 | 0.799 | 4.089 | 8.663 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, measured jitter | 150 | 0.847 | 1.224 | 6.185 | 11.957 | 1.17e-05 | 0.0030 | 0.000 |
| cv2x | no caching, measured jitter | 162 | 0.914 | 1.617 | 8.098 | 12.170 | 1.16e-04 | 0.0185 | 0.000 |
| cv2x | no caching, measured jitter | 200 | 1.129 | unstable | | | | | |
| cv2x | cold start, median work | 100 | 0.461 | 0.658 | 2.345 | 4.320 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 150 | 0.691 | 0.775 | 3.565 | 6.998 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 162 | 0.747 | 0.835 | 4.548 | 8.275 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 200 | 0.922 | 3.450 | 20.190 | 25.813 | 1.19e-01 | 0.9760 | 0.200 |
| cv2x | cold start, measured jitter | 100 | 0.463 | 0.672 | 3.453 | 8.623 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, measured jitter | 150 | 0.694 | 0.810 | 4.719 | 10.586 | 4.17e-07 | 0.0005 | 0.200 |
| cv2x | cold start, measured jitter | 162 | 0.749 | 0.882 | 5.818 | 10.430 | 2.31e-06 | 0.0020 | 0.200 |
| cv2x | cold start, measured jitter | 200 | 0.925 | 3.922 | 22.960 | 31.526 | 1.44e-01 | 0.9765 | 0.200 |

