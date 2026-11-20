# Queue simulation from measured path times (gha-ubuntu-24.04-run44822427947-1)

Measured Falcon-512 primitive: mean 39.61 us, median 38.46 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0111 | 1.0094 | 0.4510 |
| cert_known | 1.0410 | 1.0440 | 0.4643 |
| cert_new | 2.0611 | 2.0638 | 0.9192 |
| digest_unknown | 0.0016 | 0.0021 | 0.0007 |

C-V2X per-message budget T_sf/m = 0.500 ms, i.e. 1.121 x primitive.

## C-V2X SPS, Falcon-512, measured receiver paths, certificate-aware C-V2X scheduling

Radio model: 10 subchannels per 1-ms subframe; digest-form message 4 subchannels; certificate-carrying message 10 subchannels and never shares its subframe; scheduling window 100 ms; idealised deterministic earliest-fit scheduler (not a 3GPP Mode 4 / SAE J3161 MAC); no propagation loss, collisions, half-duplex loss or retransmissions.

Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = 10 ms.

| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |
|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|
| cv2x-cert | steady, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.545 | 0.902 | 0.902 | 0.902 | 0.214 / 6 | 0.454 | 2 |
| cv2x-cert | steady, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.610 | 0.902 | 0.902 | 0.902 | 0.947 / 11 | 0.680 | 2 |
| cv2x-cert | steady, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.627 | 0.902 | 0.902 | 0.902 | 2.144 / 16 | 0.735 | 2 |
| cv2x-cert | steady, median work | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.646 | 0.902 | 0.902 | 0.902 | 98.809 / 99 | 0.791 | 2 |
| cv2x-cert | steady, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.552 | 1.146 | 2.703 | 8.470 | 0.214 / 6 | 0.453 | 14 |
| cv2x-cert | steady, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.630 | 1.634 | 3.291 | 9.939 | 0.947 / 11 | 0.680 | 17 |
| cv2x-cert | steady, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.652 | 1.772 | 3.434 | 8.926 | 2.144 / 16 | 0.734 | 16 |
| cv2x-cert | steady, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.679 | 1.919 | 3.575 | 9.036 | 98.809 / 99 | 0.790 | 16 |
| cv2x-cert | no caching, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.636 | 0.920 | 0.920 | 0.919 | 0.214 / 6 | 0.545 | 2 |
| cv2x-cert | no caching, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.701 | 0.920 | 0.920 | 0.919 | 0.947 / 11 | 0.817 | 2 |
| cv2x-cert | no caching, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.718 | 0.920 | 0.920 | 0.919 | 2.144 / 16 | 0.882 | 2 |
| cv2x-cert | no caching, median work | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.712 | 0.920 | 0.920 | 0.919 | 98.809 / 99 | 0.906 | 2 |
| cv2x-cert | no caching, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 4 | 4.00e-08 | 1.000000 | 0.655 | 1.661 | 3.761 | 11.145 | 0.214 / 6 | 0.544 | 13 |
| cv2x-cert | no caching, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 42 | 2.80e-07 | 1.000000 | 0.764 | 2.526 | 4.700 | 11.528 | 0.947 / 11 | 0.816 | 20 |
| cv2x-cert | no caching, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 179 | 1.10e-06 | 0.999999 | 0.809 | 2.925 | 5.040 | 12.507 | 2.144 / 16 | 0.882 | 22 |
| cv2x-cert | no caching, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 225 | 1.29e-06 | 0.873338 | 0.808 | 2.929 | 4.925 | 14.996 | 98.809 / 99 | 0.905 | 27 |
| cv2x-cert | cold start, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.590 | 0.920 | 0.920 | 0.919 | 0.213 / 6 | 0.409 | 2 |
| cv2x-cert | cold start, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.645 | 0.920 | 0.920 | 0.919 | 0.940 / 11 | 0.614 | 2 |
| cv2x-cert | cold start, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 0 | 0.00e+00 | 0.799969 | 0.660 | 0.920 | 0.920 | 0.919 | 2.083 / 16 | 0.663 | 2 |
| cv2x-cert | cold start, median work | 200 | 1.040 | no / no / no | 200000000 | 5415614 / 7434715 | 187149671 | 40355054 | 146794617 | 0 | 0.00e+00 | 0.733973 | 0.675 | 0.920 | 0.920 | 0.919 | 72.836 / 99 | 0.757 | 2 |
| cv2x-cert | cold start, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.600 | 1.332 | 3.105 | 9.856 | 0.213 / 6 | 0.409 | 13 |
| cv2x-cert | cold start, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 10 | 8.33e-08 | 0.800016 | 0.671 | 1.843 | 3.789 | 11.018 | 0.940 / 11 | 0.613 | 18 |
| cv2x-cert | cold start, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 1 | 7.72e-09 | 0.799969 | 0.694 | 2.002 | 3.959 | 10.309 | 2.083 / 16 | 0.662 | 18 |
| cv2x-cert | cold start, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 5415614 / 7434715 | 187149671 | 40355054 | 146794617 | 10 | 6.81e-08 | 0.733973 | 0.716 | 2.136 | 4.079 | 11.805 | 72.836 / 99 | 0.757 | 20 |

