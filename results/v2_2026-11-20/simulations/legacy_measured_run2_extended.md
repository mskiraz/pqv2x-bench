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

Arrival model: legacy abstract receiver-arrival model; every generated message is assumed to reach the receiver, and the radio cost of certificate-carrying messages is not modelled.

Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = 10 ms.

| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |
|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|
| cv2x | steady, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.575 | 0.943 | 0.979 | 0.978 | n/a | 0.461 | 2 |
| cv2x | steady, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.633 | 0.979 | 0.979 | 0.978 | n/a | 0.691 | 2 |
| cv2x | steady, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.647 | 0.979 | 0.979 | 0.978 | n/a | 0.747 | 2 |
| cv2x | steady, median work | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 0 | 0.00e+00 | 1.000000 | 0.691 | 0.979 | 0.979 | 0.978 | n/a | 0.922 | 2 |
| cv2x | steady, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 1 | 1.00e-08 | 1.000000 | 0.585 | 1.281 | 3.043 | 10.240 | n/a | 0.461 | 13 |
| cv2x | steady, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 15 | 1.00e-07 | 1.000000 | 0.660 | 1.848 | 3.741 | 11.981 | n/a | 0.691 | 18 |
| cv2x | steady, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 45 | 2.78e-07 | 1.000000 | 0.683 | 2.048 | 4.036 | 16.256 | n/a | 0.746 | 27 |
| cv2x | steady, measured jitter | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 310 | 1.55e-06 | 0.999998 | 0.815 | 3.106 | 5.143 | 12.587 | n/a | 0.921 | 26 |
| cv2x | no caching, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.777 | 1.992 | 2.693 | 5.323 | n/a | 0.562 | 9 |
| cv2x | no caching, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 1.145 | 3.336 | 4.588 | 9.149 | n/a | 0.843 | 16 |
| cv2x | no caching, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 466 | 2.88e-06 | 0.999997 | 1.462 | 4.491 | 6.128 | 12.068 | n/a | 0.911 | 22 |
| cv2x | no caching, median work | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.124, not simulated) | | | | | | | | |
| cv2x | no caching, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 7 | 7.00e-08 | 1.000000 | 0.797 | 2.387 | 4.145 | 11.095 | n/a | 0.563 | 14 |
| cv2x | no caching, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 1794 | 1.20e-05 | 0.999988 | 1.222 | 4.227 | 6.258 | 15.264 | n/a | 0.844 | 23 |
| cv2x | no caching, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 24860 | 1.53e-04 | 0.999847 | 1.601 | 5.652 | 8.062 | 20.412 | n/a | 0.912 | 33 |
| cv2x | no caching, measured jitter | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.126, not simulated) | | | | | | | | |
| cv2x | cold start, median work | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.658 | 1.899 | 2.349 | 5.218 | n/a | 0.421 | 8 |
| cv2x | cold start, median work | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.774 | 2.349 | 3.508 | 8.612 | n/a | 0.631 | 14 |
| cv2x | cold start, median work | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 48 | 3.70e-07 | 0.799968 | 0.834 | 2.919 | 4.545 | 11.236 | n/a | 0.682 | 19 |
| cv2x | cold start, median work | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 39994256 | 160005744 | 19025597 | 1.19e-01 | 0.704901 | 3.446 | 16.705 | 20.041 | 28.607 | n/a | 0.842 | 58 |
| cv2x | cold start, measured jitter | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 19999903 | 80000097 | 0 | 0.00e+00 | 0.800001 | 0.670 | 1.965 | 3.543 | 9.977 | n/a | 0.421 | 12 |
| cv2x | cold start, measured jitter | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 29997779 | 120002221 | 181 | 1.51e-06 | 0.800014 | 0.808 | 2.875 | 4.855 | 14.052 | n/a | 0.632 | 21 |
| cv2x | cold start, measured jitter | 162 | n/a | n/a | 162000000 | 0 / 0 | 162000000 | 32406775 | 129593225 | 2069 | 1.60e-05 | 0.799945 | 0.881 | 3.544 | 5.796 | 16.391 | n/a | 0.682 | 24 |
| cv2x | cold start, measured jitter | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 39991870 | 160008130 | 20694233 | 1.29e-01 | 0.696569 | 3.684 | 18.033 | 22.350 | 34.451 | n/a | 0.843 | 69 |

