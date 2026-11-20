# Queue simulation from measured path times (gha-ubuntu-24.04-arm-run44822427947-1)

Measured Falcon-512 primitive: mean 40.77 us, median 40.15 us. Published ARMv8 primitive: 0.446 ms.

| Path | Median / primitive median | Mean / primitive mean | Scaled median (ms) |
|---|---:|---:|---:|
| digest_hit | 1.0076 | 1.0090 | 0.4494 |
| cert_known | 1.0353 | 1.0375 | 0.4617 |
| cert_new | 2.0823 | 2.0929 | 0.9287 |
| digest_unknown | 0.0025 | 0.0029 | 0.0011 |

C-V2X per-message budget T_sf/m = 0.500 ms, i.e. 1.121 x primitive.

## C-V2X SPS, Falcon-512, measured receiver paths, certificate-aware C-V2X scheduling

Radio model: 10 subchannels per 1-ms subframe; digest-form message 4 subchannels; certificate-carrying message 10 subchannels and never shares its subframe; scheduling window 100 ms; idealised deterministic earliest-fit scheduler (not a 3GPP Mode 4 / SAE J3161 MAC); no propagation loss, collisions, half-duplex loss or retransmissions.

Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = 10 ms.

| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |
|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|
| cv2x-cert | steady, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.543 | 0.899 | 0.899 | 0.899 | 0.214 / 6 | 0.452 | 2 |
| cv2x-cert | steady, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.607 | 0.899 | 0.899 | 0.899 | 0.947 / 11 | 0.678 | 2 |
| cv2x-cert | steady, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.625 | 0.899 | 0.899 | 0.899 | 2.144 / 16 | 0.732 | 2 |
| cv2x-cert | steady, median work | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.643 | 0.899 | 0.899 | 0.899 | 98.809 / 99 | 0.788 | 2 |
| cv2x-cert | steady, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.549 | 1.048 | 2.269 | 5.648 | 0.214 / 6 | 0.453 | 10 |
| cv2x-cert | steady, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.623 | 1.448 | 2.602 | 6.515 | 0.947 / 11 | 0.679 | 12 |
| cv2x-cert | steady, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.644 | 1.555 | 2.666 | 6.456 | 2.144 / 16 | 0.733 | 11 |
| cv2x-cert | steady, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.669 | 1.705 | 2.826 | 6.507 | 98.809 / 99 | 0.789 | 12 |
| cv2x-cert | no caching, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.636 | 0.929 | 0.929 | 0.929 | 0.214 / 6 | 0.545 | 2 |
| cv2x-cert | no caching, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.701 | 0.929 | 0.929 | 0.929 | 0.947 / 11 | 0.818 | 2 |
| cv2x-cert | no caching, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.718 | 0.929 | 0.929 | 0.929 | 2.144 / 16 | 0.883 | 2 |
| cv2x-cert | no caching, median work | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.711 | 0.929 | 0.929 | 0.929 | 98.809 / 99 | 0.906 | 2 |
| cv2x-cert | no caching, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.652 | 1.524 | 2.792 | 6.799 | 0.214 / 6 | 0.547 | 10 |
| cv2x-cert | no caching, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.748 | 2.123 | 3.270 | 7.993 | 0.947 / 11 | 0.820 | 14 |
| cv2x-cert | no caching, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 0 | 162000000 | 0 | 0.00e+00 | 1.000000 | 0.783 | 2.300 | 3.449 | 7.751 | 2.144 / 16 | 0.886 | 15 |
| cv2x-cert | no caching, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.781 | 2.312 | 3.445 | 7.768 | 98.809 / 99 | 0.909 | 16 |
| cv2x-cert | cold start, median work | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.589 | 0.929 | 0.929 | 0.929 | 0.213 / 6 | 0.409 | 2 |
| cv2x-cert | cold start, median work | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.645 | 0.929 | 0.929 | 0.929 | 0.940 / 11 | 0.613 | 2 |
| cv2x-cert | cold start, median work | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 0 | 0.00e+00 | 0.799969 | 0.659 | 0.929 | 0.929 | 0.929 | 2.083 / 16 | 0.662 | 2 |
| cv2x-cert | cold start, median work | 200 | 1.040 | no / no / no | 200000000 | 5415614 / 7434715 | 187149671 | 40355054 | 146794617 | 0 | 0.00e+00 | 0.733973 | 0.674 | 0.929 | 0.929 | 0.929 | 72.836 / 99 | 0.757 | 2 |
| cv2x-cert | cold start, measured jitter | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 19994828 | 80005172 | 0 | 0.00e+00 | 0.800052 | 0.598 | 1.210 | 2.541 | 6.737 | 0.213 / 6 | 0.410 | 9 |
| cv2x-cert | cold start, measured jitter | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 29997582 | 120002418 | 0 | 0.00e+00 | 0.800016 | 0.664 | 1.627 | 2.823 | 6.742 | 0.940 / 11 | 0.615 | 13 |
| cv2x-cert | cold start, measured jitter | 162 | 0.842 | no / no / yes | 162000000 | 0 / 0 | 162000000 | 32405068 | 129594932 | 0 | 0.00e+00 | 0.799969 | 0.684 | 1.763 | 2.947 | 6.764 | 2.083 / 16 | 0.664 | 14 |
| cv2x-cert | cold start, measured jitter | 200 | 1.040 | no / no / no | 200000000 | 5415614 / 7434715 | 187149671 | 40355054 | 146794617 | 0 | 0.00e+00 | 0.733973 | 0.705 | 1.866 | 3.032 | 8.296 | 72.836 / 99 | 0.759 | 16 |

