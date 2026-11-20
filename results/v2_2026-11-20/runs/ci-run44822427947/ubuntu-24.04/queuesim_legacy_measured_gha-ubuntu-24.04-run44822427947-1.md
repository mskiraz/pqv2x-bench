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

| Arrival model | Scenario | N | rho | Mean T (ms) | p99.9 T (ms) | Max T (ms) | Pr[T > 10 ms] | Configs with a miss | Discarded (unknown signer) |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| cv2x | steady, median work | 100 | 0.454 | 0.566 | 0.929 | 0.929 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 150 | 0.680 | 0.623 | 0.929 | 0.929 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 162 | 0.735 | 0.637 | 0.929 | 0.929 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 200 | 0.907 | 0.680 | 0.929 | 0.929 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 100 | 0.453 | 0.574 | 2.853 | 8.226 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 150 | 0.680 | 0.645 | 3.393 | 10.913 | 2.00e-08 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 162 | 0.734 | 0.665 | 3.527 | 9.352 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 200 | 0.907 | 0.766 | 4.308 | 10.736 | 9.50e-08 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 100 | 0.545 | 0.737 | 2.209 | 4.626 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 150 | 0.817 | 0.996 | 3.609 | 7.260 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 162 | 0.882 | 1.170 | 4.519 | 9.069 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 200 | 1.089 | unstable | | | | | |
| cv2x | no caching, measured jitter | 100 | 0.544 | 0.751 | 4.013 | 10.323 | 1.00e-08 | 0.0000 | 0.000 |
| cv2x | no caching, measured jitter | 150 | 0.816 | 1.047 | 5.733 | 13.511 | 5.97e-06 | 0.0012 | 0.000 |
| cv2x | no caching, measured jitter | 162 | 0.882 | 1.252 | 6.741 | 16.557 | 3.04e-05 | 0.0044 | 0.000 |
| cv2x | no caching, measured jitter | 200 | 1.088 | unstable | | | | | |
| cv2x | cold start, median work | 100 | 0.454 | 0.636 | 2.111 | 4.354 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 150 | 0.680 | 0.732 | 2.864 | 6.855 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 162 | 0.735 | 0.771 | 3.425 | 8.432 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 200 | 0.907 | 2.157 | 15.189 | 23.034 | 3.28e-02 | 0.6395 | 0.200 |
| cv2x | cold start, measured jitter | 100 | 0.453 | 0.645 | 3.323 | 10.472 | 3.75e-08 | 0.0000 | 0.200 |
| cv2x | cold start, measured jitter | 150 | 0.680 | 0.756 | 4.314 | 11.400 | 3.67e-07 | 0.0001 | 0.200 |
| cv2x | cold start, measured jitter | 162 | 0.734 | 0.805 | 4.939 | 12.204 | 2.17e-06 | 0.0004 | 0.200 |
| cv2x | cold start, measured jitter | 200 | 0.907 | 2.262 | 17.058 | 28.535 | 3.87e-02 | 0.6271 | 0.200 |

