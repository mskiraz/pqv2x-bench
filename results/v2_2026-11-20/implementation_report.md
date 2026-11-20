# Implementation report: strengthening the PQ-V2X benchmark (version 2)

Date: 2026-11-20. Branch `bench-v2-methodology`. Code commits `a4a20d5`
(benchmark metadata and CI) and `0a262ef` (simulation, size sweep, tests); the
results and documentation commit follows them. The four-constraint framework
of the paper is unchanged; this work adds measurements and a second arrival
model. The manuscript was not edited (see `paper_update_values.md`).

## 1. Status by task

| Task | Status |
|---|---|
| 1. Native ARM64 receiver-path benchmarks | **Done.** Native ARM64 and x86-64 Linux measurements (GC on and GC off) from two CI runs; laptop results untouched |
| 2. Certificate-aware C-V2X scheduling | **Done.** New `-arrival certaware` model; legacy model preserved and byte-identical |
| 3. Cold-start and cache accounting | **Done.** Full accounting in both models; new unit and regression tests |
| 4. COER payload-size sensitivity | **Done.** `cmd/sizesweep`, exact search, signed and verified messages |
| 5. Documentation and reproducibility | **Done.** README, source comments, workflow; tests for invariants, boundaries, identities and published outputs |

## 2. Where every number came from

| Output | Kind | Produced by |
|---|---|---|
| `arm_benchmark_comparison.{csv,md}` | measurements (per platform); platform-scaled estimates in one marked table | `cmd/benchcompare`, CI run 44822427947, job `report` |
| `runs/ci-run44822427947/…` | measurements, raw samples, environment | `cmd/bench`, CI run 44822427947 (commit `0a262ef`) |
| `runs/ci-run44820736544/…` | measurements, raw samples, environment | `cmd/bench`, CI run 44820736544 (commit `a4a20d5`; bench code identical) |
| `queue_comparison.{csv,md}`, `simulations/` | simulations | `cmd/queuesim`, `cmd/queuecompare`, CI run 44822427947, job `analysis` |
| `simulations/queue_comparison_runner_ratios.*` | simulations with ratios measured on the runners | same run, jobs `bench` and `report` |
| `payload_sensitivity.{csv,md,json}` | measurements (encoder output) | `cmd/sizesweep`, same run, job `analysis` |
| `../*.md`, `../*.json`, `../raw/` | original published results | unchanged; regression-tested |

Every new number in the reports was produced by code in this repository, in
GitHub Actions, from a pushed commit. Nothing was estimated by hand except
the [A]-labelled arithmetic in `paper_update_values.md`, which is stated with
its formula.

### Fixed experimental parameters

| Parameter | Value |
|---|---|
| Go toolchain | 1.25.6 (`GOTOOLCHAIN=go1.25.6`, `setup-go 1.25.6`), same as the laptop runs |
| Benchmark | 20,000 timed iterations (SLH-DSA 2,000), 1,000 warm-up (SLH-DSA 50), 50 senders x 4 messages per form (SLH-DSA 4 x 2), one locked OS thread, GOMAXPROCS 4 (runner default) |
| Garbage collector | enabled, GOGC 100 (original methodology); separate sensitivity runs with `-gogc -1` |
| Queue simulation | 100,000 random configurations per case, seed 7, 10 warm-up + 10 measured periods (cold start: 0 + 10), deadline 10 ms, 10 Hz, one in five messages carries the certificate, S = 0.446 ms |
| Certificate-aware radio | 10 subchannels per 1-ms subframe, digest-form 4, certificate-carrying 10, scheduling window 100 ms |
| Size sweep | limit 2481 B, payloads 0–2482 scanned, sweep points 39, 64, 100, 128, 160, 170, 174, 175, 176, 200, 256, 400 |
| Runners | `ubuntu-24.04` (run 1: AMD EPYC 7763; run 2: Intel Xeon 6975P-C), `ubuntu-24.04-arm` (Arm Neoverse-N2, implementer 0x41 part 0xd49), Ubuntu 24.04.5, 4 vCPUs |

## 3. Task 1 — native ARM64 measurements

### What was changed

- `cmd/bench` now records, in JSON and Markdown, the CPU model and relevant
  features (x86: AVX, AVX2, AVX-512F, AES, BMI2, SHA-NI from
  `/proc/cpuinfo`; ARM64: ASIMD, AES, SHA2, SHA3, SHA512, atomics), OS and
  kernel, Go version, compiler, build settings (`GOAMD64`, `GOARM64`,
  `CGO_ENABLED`, VCS revision where stamped), GOMAXPROCS, GOGC and the
  effective GC percentage, module versions, and per scheme the iteration
  count, warm-up, fleet size and an implementation description
  (`internal/hostinfo`). It also writes a CSV. No hostname or user data.
- `-gogc` flag for a separate GC sensitivity run; default unchanged.
- The `primitive` description no longer claims "key already parsed".
- Workflow: pinned toolchain, environment capture (`uname`, `lscpu`,
  `go env`), unique run tags, GC-off run, artifacts always uploaded, plus the
  `analysis` and `report` jobs.

### Implementation facts established from the library sources

- **Falcon-512** (`falcon-go` v0.1.1, `native`): `ParsePublicKey` stores the
  encoded bytes only. `FalconVerify` → `fndsa.Verify` decodes the 897-byte key
  (`modq_decode`), computes SHAKE256 over it (`hash_verifying_key`, computed
  even though the raw-Falcon hash mode does not use it), converts it to NTT
  form, decodes the signature, and computes hash-to-point over nonce and
  message. All of this is inside every `primitive` sample. Verification code
  is pure Go; the module's `.s` files implement floating-point helpers used
  by signing.
- **ML-DSA-44** (`dilithium-go` v0.3.0, vendored): no SIMD on any
  architecture. On AVX2 CPUs `hasAVX2` is true and `NTT` calls `nttAVX2_8`,
  but in the vendored copy that function is the scalar Go loop from
  `ntt_generic.go` (the upstream assembly is absent); `InvNTT` is scalar in
  both branches. The NIST KAT tests of the library pass on the copy.
- **SLH-DSA-SHA2-128s** (`lattice-slh-dsa-go` v0.2.0): pure Go, hashing via
  `crypto/sha256`/`crypto/sha512`, which use SHA instructions where present.
- **ECDSA P-256**: standard library; key parsed once when cached.

### Results (measurements; `arm_benchmark_comparison.md`)

Falcon-512 median ratio to the primitive, GC enabled:

| Platform | digest_hit | cert_known | cert_new |
|---|---:|---:|---:|
| Laptop i9-10980HK, Windows (published run 1) | 1.0174 | 1.0983 | 2.2330 |
| x86-64 Linux, Xeon 6975P-C (CI run 2) | 1.0111 | 1.0410 | 2.0611 |
| x86-64 Linux, EPYC 7763 (CI run 1) | 1.010 | 1.034 | 2.012 |
| ARM64 Linux, Neoverse-N2 (CI run 2) | 1.0076 | 1.0353 | 2.0823 |
| ARM64 Linux, Neoverse-N2 (CI run 1) | 1.008 | 1.035 | 2.084 |

Median primitive times (us): Falcon-512 63.10 / 38.46 / 40.15, ML-DSA-44
124.20 / 68.78 / 88.12, ECDSA 91.00 / 61.52 / 68.91, SLH-DSA 1184.9 / 332.3 /
410.2 (laptop / Xeon / Neoverse-N2).

**Is the x86-to-ARM ratio assumption supported?** Partly.

- The cached-digest ratio transfers: 1.008–1.017 everywhere.
- The certificate paths do not transfer exactly: the laptop's `cert_known`
  (1.098) and `cert_new` (2.233) ratios exceed the native ARM64 ones (1.035,
  2.082) and the x86-64 Linux ones. The laptop values are the most
  conservative of the three; using the ARM64 ratios lowers the
  platform-scaled estimates from 0.454 / 0.490 / 0.996 ms to 0.449 / 0.462 /
  0.929 ms.
- A likely cause is SHA-256 hardware (certificate hashing): both runners
  report SHA instructions, and SLH-DSA, which is almost all SHA-256, runs
  2.9–3.6 times faster on them. The laptop run did not record CPU features, so
  this explanation is not verified for the laptop.
- The Neoverse-N2 is a server core. No measurement here is on automotive
  hardware, so the transfer to the ARMv8 platform of the published 0.446 ms
  remains an assumption.
- The relative cost of Falcon-512 and ML-DSA-44 depends on the
  implementation: in this Go code Falcon-512 is faster on ARM64 (40.2 vs
  88.1 us), the opposite of the published ARMv8 order (0.446 vs 0.189 ms).

**GC sensitivity** (separate condition, not a replacement): with the collector
disabled the Falcon-512 median ratios are 1.0137 / 1.0431 / 2.1179 (x86-64)
and 1.0077 / 1.0330 / 2.1032 (ARM64); the Falcon-512 primitive median changes
by less than 0.2% (38.46 vs 38.40 us, 40.15 vs 40.15 us), while other schemes
move more (ML-DSA-44 on x86-64: 68.78 vs 73.15 us), so GC-on and GC-off runs
are reported as separate conditions. Laptop GC-off runs were not possible (section 8).

**x86 runner variability.** GitHub assigned different x86 hosts to the two CI
runs (EPYC 7763, then Xeon 6975P-C); their Falcon `cert_new` ratios differ
(2.012 vs 2.061). Hosted runners are not a controlled platform.

## 4. Task 2 — certificate-aware scheduling

### Model (`cmd/queuesim/certaware.go`)

- Configurations (sender subframes, certificate offsets, tie order) are drawn
  by the same calls, in the same order, from an identically seeded stream as
  the legacy `cv2x` model, so both models simulate the same configurations.
  Service-time jitter uses a separate stream.
- Each message is generated at the start of its sender's subframe and
  placed, in generation order, into the earliest subframe at or after that
  one with enough free subchannels (digest 4, certificate 10, of 10), within a
  100-ms window; otherwise it is counted as not scheduled. Union-find
  structures find the next usable subframe for each message kind.
- A message arrives at the end of its subframe; one FCFS verifier as before.
- Idealised, deterministic, centralised; no collisions, half-duplex loss,
  propagation loss or retransmissions; not a 3GPP/SAE MAC. A message is never
  scheduled before its generation subframe, and with a window of at most one
  period per-sender order is preserved.

### Invariants and tests (`cmd/queuesim/queuesim_test.go`)

- Every subframe carries at most 10 subchannels; its occupancy equals the sum
  of its messages' costs; a subframe with a certificate-carrying message
  holds exactly that message and 10 subchannels; no message precedes its
  generation subframe or leaves its window (200 configurations each for
  N = 100, 150, 200).
- Offered occupancy is exactly 0.52, 0.78 and 1.04 for N = 100, 150 and 200;
  N = 100 schedules everything; N = 200 leaves messages unscheduled.
- With certificates costing as much as digest messages the model reduces to
  the legacy one: identical counts and delays.
- With median Falcon-512 work the delay never exceeds max(2 x digest,
  cert_new) work.

### Results (simulation; `queue_comparison.md`)

- With median work, the certificate-aware receiver's maximum delay is
  0.9075 ms (steady) and 0.9959 ms (no caching, cold start) for every N up to
  200, against 0.98–28.62 ms in the legacy model. Each subframe's work (two
  digest checks, or one certificate check) fits in 1 ms, so no work carries
  over.
- The binding constraint becomes the radio. N = 100, 150 and 162 schedule
  every message (offered occupancy 0.52, 0.78, 0.8424; maximum scheduling
  delay 6, 11 and 16 ms). N = 200 offers 1.04 and leaves 12.67% of the
  generated messages unscheduled (6.7% of the digest-form and 36.7% of the
  certificate-carrying ones); under overload the greedy FCFS scheduler
  defers almost every message to the end of its window (mean scheduling
  delay 98.8 ms), which is a symptom of infeasibility, not a design
  recommendation.
- The legacy model's cold-start result at N = 200 (11.9% deadline misses)
  does not survive: it relied on two certificate-carrying messages sharing a
  subframe.
- With laptop jitter, certificate-aware maximum delays are 8.1–14.8 ms and
  Pr[miss | authenticated] is at most 4.2 x 10^-6.
- With ratios measured on the runners instead of the laptop, the
  certificate-aware maximum delay is 0.90–0.93 ms; the legacy cold start at
  N = 200 drops from 11.9% to 3.3–3.5% misses
  (`simulations/queue_comparison_runner_ratios.md`), another sign of how
  sensitive the legacy result is.

## 5. Task 3 — accounting

Both models now report `generated`, `scheduled`, `not_scheduled_digest`,
`not_scheduled_cert`, `unknown_signer_discarded`, `authenticated`,
`authentication_deadline_missed`, `p_miss_given_authenticated` (denominator:
authenticated), `on_time_fraction_of_generated` (denominator: generated),
`discarded_fraction_of_generated`, p99, utilisation, maximum backlog,
scheduling delay and budget checks (`-report extended`, `-ext-out`, `-csv`).
The original table format is unchanged.

Checks:

- Tests enforce generated = scheduled + not scheduled, scheduled = discarded
  + authenticated, missed <= authenticated, and no discards outside cold
  start, for both models.
- With a verifier far too slow, every authenticated message misses and the
  discarded ones are neither authenticated nor misses.
- Path costs: each sender's first certificate is charged `cert_new`, the
  second `cert_known`; unknown digests are charged `digest_unknown`; without
  caching every certificate is `cert_new`. Counts are checked exactly for
  both models.
- Cold-start discard baseline: 0.1999–0.2018 of generated messages in the
  first second (legacy and certificate-aware), regression-tested at 0.200 ±
  0.005. This follows from the uniform certificate offset (2 of 10 messages
  on average); it is not a packet-loss rate.
- Receiver-level tests (`internal/receiver/cache_test.go`): empty cache
  except the issuer; unknown digest rejected before the certificate; first
  valid certificate validated and cached and its digests then verify; other
  senders stay unknown; a tampered certificate is refused and not cached; a
  certificate with an unsupported key algorithm is refused and not cached;
  `Forget` returns a sender to unknown.

## 6. Task 4 — payload sensitivity

`fleet.NewWithPayload` and `Fleet.SignMessage` were added; `fleet.New`
behaves as before. `cmd/sizesweep` encodes every payload length 0–2482 with
the real encoder (signature of the exact length), finds the largest fitting
payload without assuming one byte per byte, and signs and verifies a real
message at every reported length (the size must equal the scanned size).

| Scheme | Digest form: max payload | Certificate form: max payload |
|---|---:|---:|
| ECDSA P-256 | 2386 B | 2248 B |
| Falcon-512 | 1782 B | **174 B** |
| ML-DSA-44 | **30 B** | none (6229 B at 0 B) |
| SLH-DSA-SHA2-128s | none (7887 B at 0 B) | none |

At 39 B the sizes are the published ones (Falcon 1615 / 736 / 2345 B,
ML-DSA digest 2490 B; regression-tested). Boundary tests check 174 B →
2481 B and 175 B → 2482 B for the Falcon certificate form, and that the length
determinant adds a byte at 128 and at 256. Only the payload varies; optional
certificate fields are not modelled.

## 7. Task 5 — documentation and tests

- README: legacy vs certificate-aware models; the 162 / 125 / 192 / 250
  vehicle figures and their assumptions; cold-start window (10 measured
  periods); the M/M/1 mention removed (the code has no Poisson model); the
  encoder described as a selected IEEE 1609.2-2022 subset with hypothetical
  PQ extensions; certificate validation described as a restricted one-level
  profile without `certIssuePermissions`, revocation, replay or misbehaviour
  checks; ratio portability qualified; library facts; new commands and
  flags; version-2 results.
- Source comments: the cold-start comment now says 10 measured periods; the
  `primitive` description in `cmd/bench` is corrected; the legacy `simulate`
  function documents what it does not model.
- Tests run in CI run 44822427947 (Linux, Go 1.25.6), all passed:
  - `go test ./...` (full, including SLH-DSA) in the `analysis` job;
  - `go test -short ./...` on both runner architectures;
  - the legacy simulations at 100,000 configurations reproduced
    `results/queuesim_fixed.md` and `results/queuesim_measured_run{1,2}.md`
    byte for byte (`simulations/legacy_regression.txt`).
- Run locally (Windows, Go 1.26.3): `gofmt`, `go vet ./...`,
  `go build ./...`, `go test` for `cmd/queuesim` and `internal/fleet`
  (passed). Before any change, the original simulator was confirmed to
  reproduce all three published simulation files at 100,000 configurations.

## 8. Unresolved issues and limitations

- **Local execution.** This Windows workstation's application-control policy
  blocks some newly built executables, including the `internal/receiver` test
  binary and Go 1.25.6 builds of `cmd/bench`. These were not worked around; the
  affected runs were done in CI instead. No new laptop measurement was taken.
  The workstation is not the i9-10980HK of the published runs.
- **No automotive hardware.** All new timings are from server cores in shared
  virtual machines. Values multiplied by 0.446 ms are estimates.
- **Hosted-runner noise and assignment.** Different x86 hosts in the two CI
  runs; shared VMs; tails are not representative of an embedded RTOS.
- **Laptop CPU features** were not recorded in the original runs; the
  SHA-hardware explanation for the laptop's higher certificate ratios is
  plausible but unverified for that machine.
- **Scheduler.** Idealised earliest-fit scheduling is an upper bound on what
  a distributed MAC could deliver. Collisions, half-duplex loss, SPS
  reselection and certificate fragmentation are not modelled. Under overload
  (N = 200) deferral saturates at the 100-ms window.
- **Jitter model.** "Measured jitter" draws independent samples from laptop
  timings, which loses any temporal correlation of interruptions.
- **Receiver profile.** One-level chains; no `certIssuePermissions`,
  revocation, replay, region or assurance checks.
- **Paths in comparison CSVs.** The `source_file` column refers to CI
  workspace paths (`results-analysis/…`, `artifacts/results-…`); the same
  files are under `simulations/` and `runs/` here.
- GitHub reported a Node.js 20 deprecation notice for the `checkout`,
  `setup-go` and `upload-artifact` actions; the runs succeeded.

## 9. How to reproduce

```sh
# Everything at once, natively on x86-64 and ARM64 (GitHub Actions):
gh workflow run bench.yml --ref bench-v2-methodology

# Or locally:
go test ./...
go run ./cmd/bench -tag <tag>                     # add -gogc -1 for the GC-off run
go run ./cmd/queuesim -mode fixed -configs 100000 -ext-out legacy_fixed_extended.md -csv legacy_fixed.csv
go run ./cmd/queuesim -arrival certaware -mode fixed -configs 100000 -csv certaware_fixed.csv
go run ./cmd/queuesim -arrival certaware -mode measured -tag intel-i9-10980hk-run1 -configs 100000 -senders 100,150,162,200 -csv certaware_measured_run1.csv
go run ./cmd/queuecompare -out queue_comparison legacy_fixed.csv certaware_fixed.csv ...
go run ./cmd/sizesweep -out payload
go run ./cmd/benchcompare -ref intel-i9-10980hk-run1 -out arm_benchmark_comparison results/bench_*.json <runner>/bench_*.json
```
