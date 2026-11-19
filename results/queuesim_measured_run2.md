# Queue simulation from measured path times (intel-i9-10980hk-run2)

Measured Falcon-512 primitive: mean 64.61 us, median 63.20 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0174 | 1.0169 | 0.4538 |
| cert_known | 1.0965 | 1.0966 | 0.4890 |
| cert_new | 2.2326 | 2.2429 | 0.9957 |
| digest_unknown | 0.0021 | 0.0023 | 0.0009 |

C-V2X per-message budget T_sf/m = 0.500 ms, i.e. 1.121 x primitive.

## C-V2X SPS, Falcon-512, measured receiver paths

| Arrival model | Scenario | N | rho | Mean T (ms) | p99.9 T (ms) | Max T (ms) | Pr[T > 10 ms] | Configs with a miss | Discarded (unknown signer) |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| cv2x | steady, median work | 100 | 0.461 | 0.575 | 0.979 | 0.978 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 150 | 0.691 | 0.633 | 0.979 | 0.978 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, median work | 200 | 0.922 | 0.691 | 0.979 | 0.978 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 100 | 0.461 | 0.585 | 3.043 | 10.240 | 1.00e-08 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 150 | 0.691 | 0.660 | 3.741 | 11.981 | 1.00e-07 | 0.0000 | 0.000 |
| cv2x | steady, measured jitter | 200 | 0.921 | 0.815 | 5.143 | 12.587 | 1.55e-06 | 0.0002 | 0.000 |
| cv2x | no caching, median work | 100 | 0.562 | 0.777 | 2.693 | 5.323 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 150 | 0.843 | 1.145 | 4.588 | 9.149 | 0.00e+00 | 0.0000 | 0.000 |
| cv2x | no caching, median work | 200 | 1.124 | unstable | | | | | |
| cv2x | no caching, measured jitter | 100 | 0.563 | 0.797 | 4.145 | 11.095 | 7.00e-08 | 0.0000 | 0.000 |
| cv2x | no caching, measured jitter | 150 | 0.844 | 1.222 | 6.258 | 15.264 | 1.20e-05 | 0.0022 | 0.000 |
| cv2x | no caching, measured jitter | 200 | 1.126 | unstable | | | | | |
| cv2x | cold start, median work | 100 | 0.461 | 0.658 | 2.349 | 5.218 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 150 | 0.691 | 0.774 | 3.508 | 8.612 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, median work | 200 | 0.922 | 3.446 | 20.041 | 28.607 | 1.19e-01 | 0.9718 | 0.200 |
| cv2x | cold start, measured jitter | 100 | 0.461 | 0.670 | 3.543 | 9.977 | 0.00e+00 | 0.0000 | 0.200 |
| cv2x | cold start, measured jitter | 150 | 0.691 | 0.808 | 4.855 | 14.052 | 1.51e-06 | 0.0002 | 0.200 |
| cv2x | cold start, measured jitter | 200 | 0.921 | 3.684 | 22.350 | 34.451 | 1.29e-01 | 0.9616 | 0.200 |

