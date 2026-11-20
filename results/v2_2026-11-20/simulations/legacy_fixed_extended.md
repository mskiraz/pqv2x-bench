## Fixed work, S = 0.446 ms (Falcon-512 primitive, ARMv8)

Arrival model: legacy abstract receiver-arrival model; every generated message is assumed to reach the receiver, and the radio cost of certificate-carrying messages is not modelled.

Denominators: Pr[miss | auth] = deadline missed / authenticated; On-time / generated = (authenticated - deadline missed) / generated; Discarded = unknown-signer discards / generated; offered occupancy = subchannel demand of generated messages / available subchannel-slots. Deadline D = 10 ms.

| Model | Scenario | N | Offered occupancy | Within 30% / 65% / 100% | Generated | Not scheduled (digest / cert) | Scheduled | Unknown-signer discarded | Authenticated | Deadline missed | Pr[miss \| auth] | On-time / generated | Mean T (ms) | p99 T (ms) | p99.9 T (ms) | Max T (ms) | Mean / max scheduling delay (ms) | Verifier utilisation | Max backlog |
|---|---|---:|---:|---|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|---:|
| unconstrained | S | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.620 | 1.672 | 2.325 | 4.994 | n/a | 0.446 | 12 |
| unconstrained | S | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 14890 | 7.44e-05 | 0.999926 | 1.853 | 6.104 | 8.102 | 13.285 | n/a | 0.892 | 30 |
| cv2x | S | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.557 | 0.892 | 0.892 | 0.892 | n/a | 0.446 | 2 |
| cv2x | S | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.613 | 0.892 | 0.892 | 0.892 | n/a | 0.669 | 2 |
| cv2x | S | 200 | n/a | n/a | 200000000 | 0 / 0 | 200000000 | 0 | 200000000 | 0 | 0.00e+00 | 1.000000 | 0.669 | 0.892 | 0.892 | 0.892 | n/a | 0.892 | 2 |
| cv2x | S/2S | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.719 | 1.784 | 2.122 | 4.366 | n/a | 0.535 | 8 |
| cv2x | S/2S | 150 | n/a | n/a | 150000000 | 0 / 0 | 150000000 | 0 | 150000000 | 0 | 0.00e+00 | 1.000000 | 0.941 | 2.460 | 3.272 | 6.544 | n/a | 0.803 | 12 |
| cv2x | S/2S | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.070, not simulated) | | | | | | | | |
| cv2x | 1.5S | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 0 | 0.00e+00 | 1.000000 | 0.972 | 2.029 | 2.698 | 4.732 | n/a | 0.669 | 8 |
| cv2x | 1.5S | 150 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.004, not simulated) | | | | | | | | |
| cv2x | 1.5S | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.338, not simulated) | | | | | | | | |
| cv2x | 2S | 100 | n/a | n/a | 100000000 | 0 / 0 | 100000000 | 0 | 100000000 | 5610 | 5.61e-05 | 0.999944 | 2.232 | 6.248 | 8.112 | 12.844 | n/a | 0.892 | 15 |
| cv2x | 2S | 150 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.338, not simulated) | | | | | | | | |
| cv2x | 2S | 200 | n/a | n/a | – | – | – | – | – | – | unstable (rho = 1.784, not simulated) | | | | | | | | |

