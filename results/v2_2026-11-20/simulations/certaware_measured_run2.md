# Queue simulation from measured path times (intel-i9-10980hk-run2)

Measured Falcon-512 primitive: mean 64.61 us, median 63.20 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0174 | 1.0169 | 0.4538 |
| cert_known | 1.0965 | 1.0966 | 0.4890 |
| cert_new | 2.2326 | 2.2429 | 0.9957 |
| digest_unknown | 0.0021 | 0.0023 | 0.0009 |

C-V2X per-message budget T_sf/m = 0.500 ms, i.e. 1.121 x primitive.

## C-V2X SPS, Falcon-512, measured receiver paths, certificate-aware C-V2X scheduling

Radio model: 10 subchannels per 1-ms subframe; digest-form message 4 subchannels; certificate-carrying message 10 subchannels and never shares its subframe; scheduling window 100 ms; idealised deterministic earliest-fit scheduler (not a 3GPP Mode 4 / SAE J3161 MAC); no propagation loss, collisions, half-duplex loss or retransmissions.

Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = 10 ms.

| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |
|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|
| cv2x-cert | steady, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.553 | 0.908 | 0.908 | 0.908 | 0.214 / 6 | 0.461 | 2 |
| cv2x-cert | steady, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.618 | 0.908 | 0.908 | 0.908 | 0.947 / 11 | 0.691 | 2 |
| cv2x-cert | steady, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.635 | 0.908 | 0.908 | 0.908 | 2.144 / 16 | 0.747 | 2 |
| cv2x-cert | steady, median work | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.653 | 0.908 | 0.908 | 0.908 | 98.809 / 99 | 0.802 | 2 |
| cv2x-cert | steady, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.562 | 1.215 | 2.984 | 9.680 | 0.214 / 6 | 0.461 | 13 |
| cv2x-cert | steady, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.643 | 1.781 | 3.657 | 9.727 | 0.947 / 11 | 0.691 | 17 |
| cv2x-cert | steady, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.668 | 1.977 | 3.970 | 9.731 | 2.144 / 16 | 0.746 | 19 |
| cv2x-cert | steady, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 22 | 1.26e-07 | 0.873339 | 0.697 | 2.228 | 4.285 | 11.266 | 98.809 / 99 | 0.801 | 20 |
| cv2x-cert | no caching, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.654 | 0.996 | 0.996 | 0.996 | 0.214 / 6 | 0.562 | 2 |
| cv2x-cert | no caching, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.719 | 0.996 | 0.996 | 0.996 | 0.947 / 11 | 0.843 | 2 |
| cv2x-cert | no caching, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.737 | 0.996 | 0.996 | 0.996 | 2.144 / 16 | 0.911 | 2 |
| cv2x-cert | no caching, median work | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.726 | 0.996 | 0.996 | 0.996 | 98.809 / 99 | 0.930 | 2 |
| cv2x-cert | no caching, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.680 | 1.841 | 3.717 | 9.680 | 0.214 / 6 | 0.563 | 13 |
| cv2x-cert | no caching, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 118 | 7.87e-07 | 0.999999 | 0.817 | 2.963 | 4.880 | 11.473 | 0.947 / 11 | 0.844 | 20 |
| cv2x-cert | no caching, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 377 | 2.33e-06 | 0.999998 | 0.892 | 3.418 | 5.341 | 12.479 | 2.144 / 16 | 0.912 | 22 |
| cv2x-cert | no caching, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 730 | 4.18e-06 | 0.873336 | 0.889 | 3.463 | 5.401 | 14.829 | 98.809 / 99 | 0.931 | 28 |
| cv2x-cert | cold start, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.605 | 0.996 | 0.996 | 0.996 | 0.213 / 6 | 0.421 | 2 |
| cv2x-cert | cold start, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.661 | 0.996 | 0.996 | 0.996 | 0.940 / 11 | 0.631 | 2 |
| cv2x-cert | cold start, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 0 | 0.00e+00 | 0.799969 | 0.676 | 0.996 | 0.996 | 0.996 | 2.083 / 16 | 0.682 | 2 |
| cv2x-cert | cold start, median work | 200 | 1.040 | no / no / no | 200000000 | 5415614 / 7434715 | 187149671 | 40355054 | 146794617 | 0 | 0.00e+00 | 0.733973 | 0.691 | 0.996 | 0.996 | 0.996 | 72.836 / 99 | 0.779 | 2 |
| cv2x-cert | cold start, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.618 | 1.515 | 3.292 | 9.680 | 0.213 / 6 | 0.421 | 13 |
| cv2x-cert | cold start, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 3 | 2.50e-08 | 0.800016 | 0.697 | 2.090 | 4.030 | 10.114 | 0.940 / 11 | 0.632 | 17 |
| cv2x-cert | cold start, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 3 | 2.31e-08 | 0.799969 | 0.724 | 2.348 | 4.357 | 10.370 | 2.083 / 16 | 0.682 | 18 |
| cv2x-cert | cold start, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 5415614 / 7434715 | 187149671 | 40355054 | 146794617 | 62 | 4.22e-07 | 0.733973 | 0.752 | 2.548 | 4.598 | 11.365 | 72.836 / 99 | 0.780 | 22 |

