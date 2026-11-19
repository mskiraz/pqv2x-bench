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
| cv2x | steady, median work | 200 | 0.922 | 0.691 | 0.980 | 0.980 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 100 | 0.463 | 0.587 | 2.788 | 9.742 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 150 | 0.694 | 0.661 | 3.474 | 9.023 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 200 | 0.925 | 0.809 | 4.590 | 12.310 | 6.45e-07 | 0.0001 | 0.000 |
| cv2x | no caching, median work | 100 | 0.562 | 0.777 | 2.694 | 5.324 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 150 | 0.843 | 1.146 | 4.590 | 9.154 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 200 | 1.124 | unstable | | | | | |
| cv2x | no caching, measured jitter | 100 | 0.564 | 0.798 | 4.109 | 10.956 | 4.00e-08 | 0.0000 | 0.000 |
| cv2x | no caching, measured jitter | 150 | 0.847 | 1.226 | 6.267 | 14.925 | 1.45e-05 | 0.0024 | 0.000 |
| cv2x | no caching, measured jitter | 200 | 1.129 | unstable | | | | | |
| cv2x | cold start, median work | 100 | 0.461 | 0.658 | 2.350 | 5.220 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 150 | 0.691 | 0.775 | 3.510 | 8.616 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 200 | 0.922 | 3.456 | 20.055 | 28.620 | 1.19e-01 | 0.9719 | 0.200 |
| cv2x | cold start, measured jitter | 100 | 0.463 | 0.671 | 3.424 | 10.733 | 1.25e-08 | 0.0000 | 0.200 |
| cv2x | cold start, measured jitter | 150 | 0.694 | 0.809 | 4.705 | 14.525 | 2.31e-06 | 0.0003 | 0.200 |
| cv2x | cold start, measured jitter | 200 | 0.925 | 3.870 | 22.510 | 33.595 | 1.41e-01 | 0.9723 | 0.200 |

