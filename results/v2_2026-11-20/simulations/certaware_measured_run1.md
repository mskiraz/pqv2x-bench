# Queue simulation from measured path times (intel-i9-10980hk-run1)

Measured Falcon-512 primitive: mean 64.55 us, median 63.10 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0174 | 1.0220 | 0.4538 |
| cert_known | 1.0983 | 1.0979 | 0.4898 |
| cert_new | 2.2330 | 2.2397 | 0.9959 |
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
| cv2x-cert | steady, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.564 | 1.230 | 2.689 | 8.097 | 0.214 / 6 | 0.463 | 10 |
| cv2x-cert | steady, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.644 | 1.695 | 3.382 | 8.987 | 0.947 / 11 | 0.694 | 15 |
| cv2x-cert | steady, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.669 | 1.829 | 3.612 | 9.277 | 2.144 / 16 | 0.749 | 18 |
| cv2x-cert | steady, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.697 | 2.013 | 3.769 | 9.404 | 98.809 / 99 | 0.805 | 18 |
| cv2x-cert | no caching, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.654 | 0.996 | 0.996 | 0.996 | 0.214 / 6 | 0.562 | 2 |
| cv2x-cert | no caching, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.719 | 0.996 | 0.996 | 0.996 | 0.947 / 11 | 0.843 | 2 |
| cv2x-cert | no caching, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.737 | 0.996 | 0.996 | 0.996 | 2.144 / 16 | 0.911 | 2 |
| cv2x-cert | no caching, median work | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.726 | 0.996 | 0.996 | 0.996 | 98.809 / 99 | 0.930 | 2 |
| cv2x-cert | no caching, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.681 | 1.773 | 3.767 | 9.078 | 0.214 / 6 | 0.564 | 13 |
| cv2x-cert | no caching, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 61 | 4.07e-07 | 1.000000 | 0.815 | 2.846 | 4.774 | 11.168 | 0.947 / 11 | 0.847 | 21 |
| cv2x-cert | no caching, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 382 | 2.36e-06 | 0.999998 | 0.888 | 3.383 | 5.241 | 13.094 | 2.144 / 16 | 0.914 | 24 |
| cv2x-cert | no caching, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 133 | 7.61e-07 | 0.873339 | 0.884 | 3.365 | 5.160 | 11.903 | 98.809 / 99 | 0.934 | 22 |
| cv2x-cert | cold start, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.605 | 0.996 | 0.996 | 0.996 | 0.213 / 6 | 0.421 | 2 |
| cv2x-cert | cold start, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.661 | 0.996 | 0.996 | 0.996 | 0.940 / 11 | 0.632 | 2 |
| cv2x-cert | cold start, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 0 | 0.00e+00 | 0.799969 | 0.676 | 0.996 | 0.996 | 0.996 | 2.083 / 16 | 0.682 | 2 |
| cv2x-cert | cold start, median work | 200 | 1.040 | no / no / no | 200000000 | 5415614 / 7434715 | 187149671 | 40355054 | 146794617 | 0 | 0.00e+00 | 0.733973 | 0.691 | 0.996 | 0.996 | 0.996 | 72.836 / 99 | 0.779 | 2 |
| cv2x-cert | cold start, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 2 | 2.50e-08 | 0.800052 | 0.619 | 1.457 | 3.235 | 10.567 | 0.213 / 6 | 0.423 | 13 |
| cv2x-cert | cold start, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.697 | 1.951 | 3.880 | 9.705 | 0.940 / 11 | 0.634 | 16 |
| cv2x-cert | cold start, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 62 | 4.78e-07 | 0.799968 | 0.724 | 2.158 | 4.127 | 11.429 | 2.083 / 16 | 0.684 | 21 |
| cv2x-cert | cold start, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 5415614 / 7434715 | 187149671 | 40355054 | 146794617 | 47 | 3.20e-07 | 0.733973 | 0.751 | 2.349 | 4.303 | 11.927 | 72.836 / 99 | 0.782 | 23 |

