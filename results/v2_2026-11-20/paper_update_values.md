# Values and wording for the manuscript update

Manuscript: *Post-Quantum V2X Broadcast Authentication: A Feasibility
Framework for Secure Migration* (Main.tex as of 2026-11-19). Section and
table numbers refer to the compiled PDF. The manuscript itself has not been
edited.

Every value below is copied from a file in this directory. Each value is one
of four kinds, labelled throughout:

- **[M] measurement**: a time or size observed on a named platform;
- **[E] platform-scaled estimate**: a ratio measured on one platform
  multiplied by the published ARMv8 Falcon-512 primitive time 0.446 ms
  (Twardokus et al.); not a measurement on ARMv8 automotive hardware;
- **[A] analytical calculation**: arithmetic on stated assumptions;
- **[S] simulation**: output of `cmd/queuesim` (receiver-only; no
  propagation, collisions, half-duplex loss or retransmissions).

Sources: `arm_benchmark_comparison.md` (benchmarks), `queue_comparison.md`
and `simulations/` (queues, 100,000 random configurations per case, seed 7,
deadline 10 ms), `payload_sensitivity.md` (sizes). CI run used for the new
benchmark numbers: GitHub Actions run 44822427947, Go 1.25.6.

---

## 1. Abstract

Current: "measures the complete authentication path at a median of
1.02–1.10 times the signature verification time when certificates are cached
and 2.23 times for a new certificate."

- These are Falcon-512 values on the x86 laptop. [M]
- On native ARM64 Linux (Neoverse-N2 server core, GitHub-hosted) the Falcon-512
  median ratios are 1.008 (digest, cached), 1.035 (certificate, cached) and
  2.082 (new certificate). On x86-64 Linux (Intel Xeon 6975P-C runner) they
  are 1.011, 1.041 and 2.061. [M]

Suggested: "For Falcon-512, the complete authentication path takes a median
of 1.01–1.10 times the signature verification time when certificates are
cached and 2.06–2.23 times for a new certificate, across an x86 laptop and
native x86-64 and ARM64 Linux hosts."

Current: "Receiver-only queue simulations with measured timing variability
show low authentication delay for cached certificates, while
certificate-aware radio scheduling remains to be validated."

Suggested: "A certificate-aware scheduling model, in which a
certificate-carrying Falcon-512 message occupies a whole subframe, bounds the
median-work receiver delay below 1 ms for 100–162 transmitters in all cache
states; at 200 transmitters the offered load is 104% of the subchannel
capacity and 12.7% of messages cannot be scheduled."

Current: "Falcon-512 authentication does not queue if every per-message check
finishes within 0.5 ms."

Suggested (certificate-aware form): "Falcon-512 authentication does not queue
if every subframe's work finishes within the 1-ms subframe: two digest-form
checks, or one certificate-carrying check." Values [E]: 2 x 0.4538 = 0.908 ms
and 0.996 ms (laptop ratios); 2 x 0.4494 = 0.899 ms and 0.929 ms (ARM64
ratios).

## 2. Section V-B (radio resources) and V-C (mixed deployment)

No value changes. For consistency with the repository [A]:

| Assumption | Subchannels per message | Vehicles within 65% | Within 100% |
|---|---:|---:|---:|
| Digest-form only | 4 | 162 | 250 |
| 1 in 5 carries the certificate in one subframe | 5.2 | 125 | 192 |

The paper's 125 (65%) and 57 (30%) remain correct. Offered occupancy in the
certificate-aware simulation equals 0.052 N exactly: 0.52 (N = 100), 0.78
(150), 0.8424 (162), 1.04 (200). [S, matches A]

## 3. Section V-D, Table VIII (fixed work, S = 0.446 ms, N = 100)

Current rows come from the legacy abstract model (at most two messages per
subframe whatever their form). Keep them, labelled "abstract model", and add
the certificate-aware column [S]:

| Authentication work | Legacy max T (ms) | Legacy Pr[T > 10 ms] | Cert-aware max T (ms) | Cert-aware Pr[T > 10 ms] |
|---|---:|---:|---:|---:|
| S | 0.89 | 0 | 0.89 | 0 |
| S/2S | 4.37 | 0 | 0.89 | 0 |
| 2S | 12.84 | 5.6 x 10^-5 | 12.52 | 3.0 x 10^-5 |

(Exact: legacy 0.8920, 4.3660, 12.8440, 5.61e-05; cert-aware 0.8920, 0.8920,
12.5200, 2.992e-05. At N = 100 no message is left unscheduled.)

Suggested sentence: "When the radio cost of certificates is modelled, a
certificate-carrying message occupies a subframe alone, so its extra
verification no longer delays a second message in the same subframe: with
S/2S work the maximum delay falls from 4.37 ms to 0.89 ms."

## 4. Section VI-A (implementation and method)

Add [M, from the library sources of the pinned versions]:

- "The Falcon-512 library keeps only the encoded public key; each
  verification decodes the 897-byte key, hashes it with SHAKE256 and converts
  it to NTT form, so the measured primitive includes this per-call work."
- "The ML-DSA-44 implementation uses scalar Go on every architecture (no
  SIMD)."
- "SLH-DSA's SHA-2 calls use the CPU's SHA instructions where present, which
  explains its large cross-platform differences (median primitive 1184.9 us
  on the laptop, 332.3 us on the Xeon runner, 410.2 us on the Neoverse-N2
  runner)." [M]
- "Measurements were repeated natively on x86-64 and ARM64 Linux GitHub-hosted
  runners (4 vCPUs, Go 1.25.6); these are server cores, not automotive
  processors."

## 5. Section VI-B, Table IX (sizes) — unchanged values, new sentence

Table IX values reproduce exactly (regression-tested): Falcon-512
certificate 1615 B, digest message 736 B, certificate message 2345 B;
ML-DSA-44 digest message 2490 B. [M]

Add a payload-sensitivity sentence [M, encoder output, signed and verified
messages]:

"With the implemented encoding profile, a Falcon-512 message carrying a
Falcon-512-signed certificate fits the 2481-B transport block for application
payloads up to 174 B (2481 B at 174 B; 2482 B at 175 B); the digest form fits
up to 1782 B. An ML-DSA-44 digest-form message fits only for payloads up to
30 B. The COER length determinant of the payload grows at 128 and 256 bytes,
so the margin is not a constant byte offset."

Caveat to state: only the payload varies; region, assurance level and other
optional certificate fields are not modelled.

## 6. Section VI-C, Table X (timings) — add platform rows or a note

Laptop values unchanged. Native Linux medians, microseconds [M]
(`arm_benchmark_comparison.md`, GC enabled):

| Scheme | Platform | Signature only | Digest, cached | Certificate, cached | Certificate, new |
|---|---|---:|---:|---:|---:|
| Falcon-512 | x86-64 Linux (Xeon 6975P-C) | 38.46 | ratio 1.011 | ratio 1.041 | ratio 2.061 |
| Falcon-512 | ARM64 Linux (Neoverse-N2) | 40.15 | ratio 1.008 | ratio 1.035 | ratio 2.082 |
| ML-DSA-44 | x86-64 Linux | 68.78 | 1.005 | 1.037 | 2.092 |
| ML-DSA-44 | ARM64 Linux | 88.12 | 1.004 | 1.026 | 2.071 |
| ECDSA P-256 | x86-64 Linux | 61.52 | 1.006 | 1.010 | 2.135 |
| ECDSA P-256 | ARM64 Linux | 68.91 | 1.008 | 1.014 | 2.132 |

(Full means, medians, p90, p99, p99.9, maxima and standard deviations are in
`arm_benchmark_comparison.csv`.)

Suggested sentence: "On native ARM64, Falcon-512 verification in the same Go
library is faster than ML-DSA-44 verification (median 40.2 us vs 88.1 us),
whereas the published ARMv8 measurements of Table VII show the opposite order
(0.446 ms vs 0.189 ms); relative verification cost therefore depends on the
implementation, not only on the instruction set."

## 7. Section VI-D (receiver capacity and latency)

Current: "This gives estimated authentication times of 0.454 ms for a cached
digest message, 0.490 ms for a cached full-certificate message, and 0.996 ms
for a new certificate."

These remain correct as laptop-ratio estimates [E]. Add the ARM64-ratio
estimates [E]: 0.449 ms, 0.462 ms and 0.929 ms (x86-64 Linux ratios: 0.451,
0.464, 0.919 ms).

Suggested: "Using ratios measured natively on ARM64 instead gives 0.449, 0.462
and 0.929 ms; the laptop ratios are the most conservative of the three
platforms measured."

Current: "Processing both without accumulating a queue would require every
authentication to finish within 0.5 ms. The estimated cached-message times
are below this limit, but new-certificate authentication exceeds it."

Suggested: "Under certificate-aware scheduling a certificate-carrying message
occupies a subframe alone, so the sufficient condition becomes that each
subframe's work fits in 1 ms: two cached digest checks (0.908 ms) or one
certificate check (at most 0.996 ms) [E]. Both hold for all three ratio
sources."

## 8. Table XI (Falcon-512 receiver simulation, median work, laptop ratios)

Current Table XI is the legacy abstract model. Suggested replacement with both
models, run 1 ratios [S] (`queue_comparison.md`):

| Scenario | N | Legacy max T (ms) | Legacy misses (of authenticated) | Cert-aware offered occupancy | Cert-aware not scheduled (of generated) | Cert-aware max T (ms) | Cert-aware misses | Cert-aware on time / generated |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| Certificates cached | 100 | 0.98 | 0 | 0.52 | 0 | 0.91 | 0 | 1.000 |
| No certificate caching | 100 | 5.32 | 0 | 0.52 | 0 | 1.00 | 0 | 1.000 |
| Cold start | 100 | 5.22 | 0 | 0.52 | 0 | 1.00 | 0 | 0.800 |
| Certificates cached | 200 | 0.98 | 0 | 1.04 | 12.67% | 0.91 | 0 | 0.873 |
| No certificate caching | 200 | unstable | – | 1.04 | 12.67% | 1.00 | 0 | 0.873 |
| Cold start | 200 | 28.62 | 11.9% | 1.04 | 6.43% | 1.00 | 0 | 0.734 |

Exact values: cert-aware max T 0.9075 (cached) and 0.9959 ms (no caching,
cold start). N = 200 steady: 10,664,208 digest-form and 14,667,896
certificate-carrying messages of 200,000,000 generated are not scheduled
(6.7% of digest-form, 36.7% of certificate-carrying messages). Cold start
N = 200: 5,415,614 + 7,434,715 of 200,000,000 not scheduled; unknown-signer
discards 40,355,054 (20.18% of generated). On-time / generated =
(authenticated − missed) / generated. Run 2 ratios give the same
certificate-aware values to three decimals.

Current: "At 200 vehicles, repeated certificate validation makes the receiver
unable to keep up, while a cold start ... causes approximately 12% of
authenticated messages to miss the deadline."

Suggested: "In the abstract arrival model, which lets two
certificate-carrying messages share a subframe, 200 vehicles overload the
receiver without caching and cause 11.9% deadline misses at cold start. When
certificate-carrying messages occupy a whole subframe, the receiver never
queues across subframes with median work (maximum delay 1.0 ms in every
scenario up to N = 200), and the binding constraint becomes the radio: at 200
vehicles the offered load is 104% of the subchannel capacity and 12.7% of the
messages (36.7% of the certificate-carrying ones) cannot be scheduled even by
an ideal scheduler. The 11.9% cold-start miss rate is therefore an artefact of
the abstract model."

Add (cold start): "About 20% of the messages in the first second are
discarded because their sender's certificate has not yet arrived (on average
2 of a sender's 10 messages precede its first certificate). This follows from
the transmission schedule and is not a measured packet-loss rate." [S]
(Exact: 0.1999–0.2018 of generated across N; 19,994,828 of 100,000,000 at
N = 100.)

Current: "additional tests with measured timing variability show that
occasional deadline misses remain possible."

Values [S, laptop jitter scaled to 0.446 ms] for this sentence:

- Certificate-aware, measured jitter: maximum delay 8.1–14.8 ms across all
  scenarios and N; Pr[miss | authenticated] at most 4.2 x 10^-6 (no caching,
  N = 200, run 2); 0 for the steady state at N <= 162.
- Legacy, measured jitter: maximum 9.0–34.5 ms; up to 1.5 x 10^-4 (no caching,
  N = 162) and 14.1% (cold start, N = 200).
- The tails come from laptop execution-time outliers (Windows, garbage
  collector), which the simulation rescales to the ARMv8 time.

## 9. Section VI-E (limitations) and VIII (conclusion)

Current: "The ARMv8 estimates assume that the measured processing overhead
scales similarly across processors."

Suggested: "Native measurements on x86-64 and ARM64 server cores give
Falcon-512 path ratios within 0.01 of the laptop's for cached digests
(1.008–1.011 vs 1.017), but lower for the certificate paths: 1.035–1.041 vs
1.098 for a cached certificate and 2.061–2.082 vs 2.233 for a new
certificate. This is consistent with the certificate-hashing share depending
on SHA hardware (both runners report SHA instructions; the laptop run did not
record CPU features). The ratios are therefore
not exactly portable; the laptop ratios used in the paper are the most
conservative of those measured. Automotive-processor measurements remain
necessary."

Add: "The certificate-aware arrival model uses an idealised, deterministic
earliest-fit scheduler; it is not a 3GPP Mode 4 or SAE J3161 MAC."

## 10. Claims that must change

1. "approximately 12% of authenticated messages miss the deadline" at 200
   vehicles (Sec. VI-D) — artefact of the abstract arrival model; replace as
   in item 8.
2. "No certificate caching, 200: Unstable" (Table XI) — unstable only in the
   abstract model; with certificate-aware scheduling the receiver keeps up
   (max 1.0 ms) but 12.7% of messages cannot be scheduled.
3. "The ARMv8 estimates assume that the measured processing overhead scales
   similarly across processors" — now measured: partially supported (item 9).
4. Abstract ratio range 1.02–1.10 / 2.23 — label as Falcon-512 on the laptop
   or widen to the measured range (item 1).

## 11. Still requiring automotive hardware

- Complete receiver path (all six operations) on an automotive-class ARM SoC
  (for example Cortex-A53/A72/A76 class, the class used by Twardokus et al.),
  to replace the [E] values.
- The same with an automotive HSM or V2X accelerator, if verification is
  offloaded.
- Execution-time tails under a real-time OS, to replace the laptop jitter in
  the "measured jitter" simulations.
