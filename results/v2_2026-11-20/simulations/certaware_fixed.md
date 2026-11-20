## Fixed work, S = 0.446 ms (Falcon-512 primitive, ARMv8), certificate-aware C-V2X scheduling

Radio model: 10 subchannels per 1-ms subframe; digest-form message 4 subchannels; certificate-carrying message 10 subchannels and never shares its subframe; scheduling window 100 ms; idealised deterministic earliest-fit scheduler (not a 3GPP Mode 4 / SAE J3161 MAC); no propagation loss, collisions, half-duplex loss or retransmissions.

Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = 10 ms.

| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |
|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|
| cv2x-cert | S | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.537 | 0.892 | 0.892 | 0.892 | 0.214 / 6 | 0.446 | 2 |
| cv2x-cert | S | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.600 | 0.892 | 0.892 | 0.892 | 0.947 / 11 | 0.669 | 2 |
| cv2x-cert | S | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.637 | 0.893 | 0.893 | 0.892 | 98.809 / 99 | 0.779 | 2 |
| cv2x-cert | S/2S | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.626 | 0.892 | 0.892 | 0.892 | 0.214 / 6 | 0.535 | 2 |
| cv2x-cert | S/2S | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.690 | 0.892 | 0.892 | 0.892 | 0.947 / 11 | 0.803 | 2 |
| cv2x-cert | S/2S | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | 174667896 | 0 | 0.00e+00 | 0.873339 | 0.701 | 0.893 | 0.893 | 0.892 | 98.809 / 99 | 0.892 | 2 |
| cv2x-cert | 1.5S | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.911 | 2.015 | 2.367 | 4.401 | 0.214 / 6 | 0.669 | 7 |
| cv2x-cert | 1.5S | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | – | – | receiver unstable (utilisation 1.004) | | | | | | 0.947 / 11 | 1.004 | 28 |
| cv2x-cert | 1.5S | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | – | – | receiver unstable (utilisation 1.169) | | | | | | 98.809 / 99 | 1.169 | 517 |
| cv2x-cert | 2S | 100 | 0.520 | no / yes / yes | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 2992 | 2.99e-05 | 0.999970 | 2.079 | 5.976 | 7.816 | 12.520 | 0.214 / 6 | 0.892 | 15 |
| cv2x-cert | 2S | 150 | 0.780 | no / no / yes | 150000000 | 0 / 0 | 150000000 | 0 | – | – | receiver unstable (utilisation 1.338) | | | | | | 0.947 / 11 | 1.338 | 765 |
| cv2x-cert | 2S | 200 | 1.040 | no / no / no | 200000000 | 10664208 / 14667896 | 174667896 | 0 | – | – | receiver unstable (utilisation 1.558) | | | | | | 98.809 / 99 | 1.558 | 1299 |

