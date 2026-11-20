# Queue simulation from measured path times (gha-ubuntu-24.04-arm-run44820736544-1)

Measured Falcon-512 primitive: mean 40.80 us, median 40.13 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0080 | 1.0082 | 0.4496 |
| cert_known | 1.0353 | 1.0354 | 0.4617 |
| cert_new | 2.0843 | 2.0981 | 0.9296 |
| digest_unknown | 0.0024 | 0.0029 | 0.0011 |

C-V2X per-message budget T_sf/m = 0.500 ms, i.e. 1.121 x primitive.

## C-V2X SPS, Falcon-512, measured receiver paths

| Arrival model | Scenario | N | rho | Mean T (ms) | p99.9 T (ms) | Max T (ms) | Pr[T > 10 ms] | Configs with a miss | Discarded (unknown signer) |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| cv2x | steady, median work | 100 | 0.452 | 0.564 | 0.924 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 150 | 0.678 | 0.621 | 0.924 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 162 | 0.732 | 0.635 | 0.924 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 200 | 0.904 | 0.678 | 0.924 | 0.923 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 100 | 0.452 | 0.569 | 2.220 | 5.046 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 150 | 0.678 | 0.635 | 2.450 | 7.141 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 162 | 0.732 | 0.652 | 2.513 | 6.638 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 200 | 0.904 | 0.727 | 2.865 | 6.903 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 100 | 0.546 | 0.740 | 2.245 | 4.715 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 150 | 0.818 | 1.008 | 3.716 | 7.473 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 162 | 0.884 | 1.190 | 4.657 | 9.369 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 200 | 1.091 | unstable | | | | | |
| cv2x | no caching, measured jitter | 100 | 0.547 | 0.754 | 3.213 | 13.061 | 2.93e-05 | 0.0105 | 0.000 |
| cv2x | no caching, measured jitter | 150 | 0.820 | 1.054 | 5.000 | 17.778 | 1.10e-04 | 0.0144 | 0.000 |
| cv2x | no caching, measured jitter | 162 | 0.886 | 1.268 | 6.490 | 18.967 | 2.36e-04 | 0.0166 | 0.000 |
| cv2x | no caching, measured jitter | 200 | 1.094 | unstable | | | | | |
| cv2x | cold start, median work | 100 | 0.452 | 0.636 | 2.138 | 4.437 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 150 | 0.678 | 0.733 | 2.927 | 7.039 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 162 | 0.732 | 0.774 | 3.526 | 8.705 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 200 | 0.904 | 2.189 | 15.528 | 23.554 | 3.60e-02 | 0.6716 | 0.200 |
| cv2x | cold start, measured jitter | 100 | 0.452 | 0.644 | 2.722 | 12.486 | 1.36e-05 | 0.0055 | 0.200 |
| cv2x | cold start, measured jitter | 150 | 0.678 | 0.752 | 3.628 | 14.811 | 2.82e-05 | 0.0076 | 0.200 |
| cv2x | cold start, measured jitter | 162 | 0.732 | 0.800 | 4.340 | 14.437 | 4.63e-05 | 0.0084 | 0.200 |
| cv2x | cold start, measured jitter | 200 | 0.904 | 2.338 | 17.205 | 31.674 | 4.52e-02 | 0.7067 | 0.200 |

