# pqv2x-bench

Receiver-side authentication benchmark and queue simulation for post-quantum
V2X broadcast authentication.

This is the companion code for the paper *Post-Quantum Migration for V2X: A
Framework for Broadcast Authentication in Mixed Classical/PQ Deployments*. It
produces the exact encoded sizes, the receiver-path timings, and the queue
simulations reported there.

> **Not audited.** This is research code for measurement. It has had no
> security audit, is not constant-time, omits revocation and replay checks,
> and uses hypothetical post-quantum identifiers that no standard defines.
> Do not use it to protect real vehicles, messages, or keys. See
> [SECURITY.md](SECURITY.md).

## Contents

- [What it does](#what-it-does)
- [Quick start](#quick-start)
- [Signature schemes](#signature-schemes)
- [Repository layout](#repository-layout)
- [How a message is authenticated](#how-a-message-is-authenticated)
- [Encoding profile](#encoding-profile)
- [The benchmark: `cmd/bench`](#the-benchmark-cmdbench)
- [The queue simulation: `cmd/queuesim`](#the-queue-simulation-cmdqueuesim)
- [Results in this repository](#results-in-this-repository)
- [Reproducing the paper's tables](#reproducing-the-papers-tables)
- [Running on ARM](#running-on-arm)
- [Getting stable measurements](#getting-stable-measurements)
- [Limitations](#limitations)
- [Third-party code](#third-party-code)
- [Citing](#citing)
- [License](#license)

## What it does

The code answers two questions.

1. **How large are IEEE 1609.2 certificates and signed messages when they
   carry post-quantum keys and signatures?** The structures are encoded with
   the Canonical Octet Encoding Rules (COER), byte for byte. Nothing is
   estimated.
2. **How much longer does the complete receiver authentication path take than
   the bare signature verification?** Published benchmarks usually report the
   signature primitive only. A receiver also decodes the message, resolves the
   signer through a certificate cache, validates a new certificate against its
   issuer, and checks validity and permissions.

A third program takes the measured per-path times and simulates the
authentication queue of one receiver that hears 100 to 200 vehicles, each
sending ten messages per second.

## Quick start

Requires Go 1.25 or later. No C compiler and no other tools are needed.

```sh
git clone https://github.com/mskiraz/pqv2x-bench.git
cd pqv2x-bench

# 1. Correctness tests (about 15 s; add -short to skip SLH-DSA)
go test ./internal/...

# 2. Benchmark (about 2 minutes); writes results/bench_<tag>.{json,md}
go run ./cmd/bench -tag my-machine

# 3. Queue simulation from the benchmark just produced
go run ./cmd/queuesim -mode measured -tag my-machine -configs 100000
```

Step 2 prints two tables: encoded sizes and time per message. Step 3 prints
the queue results for a C-V2X receiver.

## Signature schemes

| Scheme | Public key | Signature | Implementation |
|---|---:|---:|---|
| ECDSA P-256 | 33 B | 64 B | Go standard library (`crypto/ecdsa`) |
| Falcon-512, padded format | 897 B | 666 B | [`lattice-safe/falcon-go`](https://github.com/lattice-safe/falcon-go) v0.1.1, `native` package |
| ML-DSA-44 | 1312 B | 2420 B | [`lattice-safe/dilithium-go`](https://github.com/lattice-safe/dilithium-go) v0.3.0, vendored |
| SLH-DSA-SHA2-128s | 32 B | 7856 B | [`lattice-safe/lattice-slh-dsa-go`](https://github.com/lattice-safe/lattice-slh-dsa-go) v0.2.0 |

Notes:

- Falcon-512 uses the fixed-size 666-byte padded signature of the Falcon
  specification, through the library's `native` package. The library's
  high-level API produces 809-byte constant-time-format signatures instead.
- ECDSA public keys are compressed points, as IEEE 1609.2 carries them.
- Every scheme signs the same 32-byte value (see
  [Encoding profile](#encoding-profile)), so the schemes are compared on
  identical inputs.

## Repository layout

```
cmd/
  bench/        benchmark: encoded sizes and per-path timing
  queuesim/     periodic-arrival queue simulation
internal/
  coer/         COER subset (ITU-T X.696): integers, lengths, choices, ...
  ieee1609/     IEEE 1609.2 SignedData and explicit certificate, encode/decode
  receiver/     certificate cache and the authentication path (+ tests)
  fleet/        test material: issuer, pseudonym certificates, signed messages
  sigs/         one interface over the four signature schemes
  hrtime/       high-resolution clock (performance counter on Windows)
third_party/
  dilithium-go/ patched copy of the ML-DSA library (see VENDORING.md there)
results/        outputs of the runs reported in the paper
  raw/          one timing sample file per run, scheme and operation
.github/workflows/bench.yml   CI job for x86-64 and 64-bit ARM runners
```

How the pieces fit together:

```
 fleet ──builds──> issuer cert, 50 pseudonym certs, signed messages
   │                      (encoded by ieee1609 using coer, signed via sigs)
   ▼
 bench ──feeds messages to──> receiver.Process ──> decode, cache, verify
   │                                   timed with hrtime
   ▼
 results/bench_<tag>.json, .md, raw/*.txt
   │
   ▼
 queuesim ──replays the measured path times──> queue delay statistics
```

## How a message is authenticated

`receiver.Process` handles one encoded message in six steps.

1. **Decode** the COER-encoded secured message.
2. **Relevance:** reject it if its generation time is outside the accepted
   window (5 s in the test material).
3. **Resolve the signer.**
   - Digest form: look up the 8-byte `HashedId8` in the certificate cache.
     If it is not there, the message cannot be verified and is discarded.
   - Certificate form: hash the carried certificate and look it up.
4. **Validate a new certificate** (cache miss on a carried certificate): find
   the issuer in the cache, verify the certificate's signature with the
   issuer's key, and store a private copy in the cache.
5. **Consistency:** check the certificate's validity period, that its
   application permission matches the message's PSID, and that the signature
   algorithm matches the certified key.
6. **Verify the message signature.**

The path returns an outcome (`Accepted`, `Malformed`, `Stale`,
`UnknownSigner`, `UnknownIssuer`, `BadCertificate`, `BadSignature`) and the
path taken (`PathDigestHit`, `PathCertKnown`, `PathCertNew`).

Not included: revocation lists, misbehaviour detection, replay detection, and
chains longer than one issuer.

## Encoding profile

The structures follow the IEEE 1609.2-2022 ASN.1 definitions with the options
a SAE J2945/1 Basic Safety Message uses.

**Secured message** (`Ieee1609Dot2Data` containing `SignedData`):

- hash algorithm SHA-256;
- payload: an unsecured 39-byte BSM;
- header: `psid` (0x20) and `generationTime`;
- signer: a `HashedId8` digest, or a sequence of one certificate;
- signature.

**Certificate** (explicit):

- issuer identified by `sha256AndDigest`;
- identified by linkage data (`iCert` and a 9-byte linkage value);
- `cracaId`, `crlSeries`, validity as start time plus a duration in hours;
- one application permission (PSID 0x20);
- a verification key and the issuer's signature;
- no region and no assurance level.

**Post-quantum keys and signatures.** IEEE 1609.2 has no post-quantum
algorithm identifiers. This code adds them as *extension additions* to the
`PublicVerificationKey` and `Signature` CHOICE types: index 5 for
Falcon-512, 6 for ML-DSA-44, 7 for SLH-DSA-SHA2-128s. COER encodes an
extension addition as an open type, that is, with a length determinant in
front of the value. The reported sizes include that length (3 bytes for all
three schemes' signatures and for the Falcon and ML-DSA keys).

**Signature input.** As IEEE 1609.2 specifies for explicit certificates:

```
SHA-256( SHA-256(to-be-signed data) || SHA-256(signer certificate) )
```

For a self-signed issuer certificate the second hash is the hash of the empty
string. The issuing authority uses the same scheme as the pseudonym
certificates.

**Worked example: ECDSA digest-form message, 132 bytes.**

| Field | Bytes |
|---|---:|
| protocol version, content choice, hash algorithm | 3 |
| payload preamble | 1 |
| inner data: version, choice, length | 3 |
| BSM payload | 39 |
| header: preamble, PSID (length + value), generation time | 11 |
| signer: choice + HashedId8 | 9 |
| signature: choice, r as x-only point (1 + 32), s (32) | 66 |
| **Total** | **132** |

**Worked example: ECDSA certificate, 144 bytes.**

| Field | Bytes |
|---|---:|
| preamble, version, type | 3 |
| issuer: choice + HashedId8 | 9 |
| to-be-signed preamble | 1 |
| identifier: choice, preamble, iCert, linkage value | 13 |
| cracaId, crlSeries | 5 |
| validity: start, duration choice, hours | 7 |
| application permissions: count, preamble, PSID | 5 |
| verification key: two choices, compressed point (1 + 32) | 35 |
| issuer signature | 66 |
| **Total** | **144** |

## The benchmark: `cmd/bench`

### What is timed

| Operation | Meaning |
|---|---|
| `primitive` | signature verification only, with the public key already parsed |
| `decode` | COER decoding of a digest-form message |
| `digest_hit` | full path; signer given as a digest; certificate cached |
| `cert_known` | full path; certificate carried and already cached |
| `cert_new` | full path; certificate carried, validated against its issuer, then cached |
| `digest_unknown` | digest-form message whose certificate is not cached: decoded, then discarded |

### Method

- For each scheme the program generates one issuer, 50 pseudonym
  certificates, and 4 messages of each form per sender (4 senders and 2
  messages for SLH-DSA, whose signing is slow).
- Each operation runs 1,000 untimed warm-up iterations, then 20,000 timed
  iterations (50 and 2,000 for SLH-DSA). Every iteration is timed
  individually, so the output is a distribution, not only a mean.
- Timing uses a high-resolution clock. On Windows the Go clock advances in
  0.5 ms steps, so `internal/hrtime` reads the performance counter instead
  (100 ns resolution). On other systems `time.Now` is used.
- `decode` and `digest_unknown` take about 0.1 microseconds. They are timed
  in batches of 64 and the batch time is divided by 64.
- For `cert_new`, the certificate is removed from the cache before each
  iteration, outside the timed region.
- The program runs on one OS thread with the garbage collector enabled.
- Every timed call's outcome is checked. An unexpected outcome aborts the
  run, so a result table always comes from correct verifications.

### Flags

| Flag | Default | Meaning |
|---|---|---|
| `-iters` | 20000 | timed iterations per operation |
| `-slh-iters` | 2000 | timed iterations per operation for SLH-DSA |
| `-warmup` | 1000 | untimed warm-up iterations per operation |
| `-senders` | 50 | pseudonym certificates per scheme |
| `-per` | 4 | messages of each form per sender |
| `-out` | `results` | output directory |
| `-tag` | `<os>-<arch>` | label for this run, used in file names |
| `-schemes` | all | comma-separated subset: `ecdsa,falcon,mldsa,slhdsa` |

Example, Falcon only with fewer iterations:

```sh
go run ./cmd/bench -schemes falcon -iters 5000 -tag quick
```

### Output

- `results/bench_<tag>.md`: the two tables, ready to read.
- `results/bench_<tag>.json`: the same data, machine-readable. Per scheme:
  `sizes` (`public_key`, `signature`, `certificate`, `digest_spdu`,
  `certificate_spdu`, in bytes) and, per operation, `n`, `mean_us`, `std_us`,
  `min_us`, `p50_us`, `p90_us`, `p99_us`, `p999_us`, `max_us`,
  `ratio_to_primitive_mean`, `ratio_to_primitive_median`. The header records
  the CPU, operating system, architecture, and Go version.
- `results/raw/<tag>_<scheme>_<operation>.txt`: every timed sample in
  nanoseconds, one per line. Scheme names in file names are `ecdsap256`,
  `falcon512`, `mldsa44`, `slhdsasha2128s`.

## The queue simulation: `cmd/queuesim`

### Model

- `N` senders each transmit every 100 ms (10 Hz). One message in five
  carries the sender's certificate (SAE J2945/1), at a random offset per
  sender.
- One receiver with a single serial, first-come-first-served verifier.
- The delay `T` of a message is its waiting time plus its service time.
- Two arrival models:
  - `unconstrained`: each sender has an independent random phase in the
    100 ms period, with no channel serialisation. This is the N·D/D/1 queue.
  - `cv2x`: SAE J3161 C-V2X Mode 4. Each sender holds a distinct
    semi-persistent resource, drawn at random from 100 subframes with `m`
    message slots each. All messages of a subframe reach the verifier
    together at the end of that subframe. `m = 2` corresponds to Falcon-512
    digest messages occupying 4 of 10 subchannels.
- Each case runs many independent random configurations (phases or resource
  choices). A configuration is simulated for 20 periods; the first 10 are
  discarded as warm-up.
- A case with utilisation `rho >= 1` is reported as unstable and not
  simulated.

### Mode `fixed`

Service time is a fixed multiple of the published ARMv8 Falcon-512 primitive
time `S` (0.446 ms by default):

| Scenario | Digest message | Certificate-carrying message |
|---|---|---|
| `S` | S | S |
| `S/2S` | S | 2S |
| `1.5S` | 1.5S | 1.5S |
| `2S` | 2S | 2S |

### Mode `measured`

Service times come from a `cmd/bench` run for Falcon-512. The simulator
computes, for each receiver path, the ratio of its median to the primitive
median, and multiplies that ratio by the ARMv8 primitive time. This transfers
the *ratio* of the paths to the target platform, not the absolute laptop
times.

| Scenario | Meaning |
|---|---|
| `steady` | every certificate is cached: digest messages take the `digest_hit` path and certificate-carrying messages the `cert_known` path |
| `no caching` | every certificate-carrying message takes the `cert_new` path |
| `cold start` | all senders are unknown at time 0 and the queue is empty. Digest messages from a still-unknown sender are discarded; a sender's first certificate-carrying message takes the `cert_new` path. Statistics cover the first second |

Each scenario is run twice: with *median work* (deterministic service time
per path) and with *measured jitter* (service times drawn from the raw
samples, rescaled so that the primitive mean equals the ARMv8 primitive
time). The second variant carries the benchmark machine's execution-time
jitter into the simulation and is a sensitivity check, not a prediction.

### Flags

| Flag | Default | Meaning |
|---|---|---|
| `-mode` | `measured` | `fixed` or `measured` |
| `-results` | `results` | directory written by `cmd/bench` |
| `-tag` | the only run found | benchmark run to use in `measured` mode |
| `-arm-primitive-ms` | 0.446 | target-platform primitive verification time |
| `-deadline-ms` | 10 | authentication deadline `D` |
| `-configs` | 10000 | random configurations per case |
| `-m` | 2 | messages per C-V2X subframe |
| `-seed` | 7 | random seed |
| `-out` | none | also write the Markdown report to this file |

### Output columns

| Column | Meaning |
|---|---|
| `rho` | utilisation: arrival rate times mean service time in steady state |
| Mean T, p99.9 T, Max T | delay statistics over all authenticated messages |
| Pr[T > D] | share of authenticated messages later than the deadline |
| Configs with a miss | share of configurations with at least one late message |
| Discarded | share of messages discarded because the signer was unknown |

## Results in this repository

Platform: Intel Core i9-10980HK laptop, Windows 11, Go 1.25.6. Two
independent runs (`run1`, `run2`) agree on every median to within 0.3% for
ECDSA, Falcon-512 and ML-DSA-44, and to within 2% for SLH-DSA.

**Encoded sizes (bytes), issuer using the same scheme**

| Scheme | Certificate | Digest message | Certificate-carrying message |
|---|---:|---:|---:|
| ECDSA P-256 | 144 | 132 | 270 |
| Falcon-512 | 1615 | 736 | 2345 |
| ML-DSA-44 | 3784 | 2490 | 6268 |
| SLH-DSA-SHA2-128s | 7938 | 7926 | 15858 |

For reference, the largest single-subframe C-V2X transport block in the
configuration the paper analyses is 2481 bytes.

**Median time per message (microseconds) and ratio to the primitive**

| Scheme | Primitive | Digest, cached | Certificate, cached | Certificate, new |
|---|---:|---:|---:|---:|
| ECDSA P-256 | 91.0 | 92.1 (1.012) | 93.1 (1.023) | 196.6 (2.160) |
| Falcon-512 | 63.1 | 64.2 (1.017) | 69.3 (1.098) | 140.9 (2.233) |
| ML-DSA-44 | 124.2 | 124.3 (1.001) | 136.2 (1.097) | 277.0 (2.230) |
| SLH-DSA-SHA2-128s | 1184.9 | 1187.0 (1.002) | 1233.4 (1.041) | 2461.9 (2.078) |

Decoding alone takes 0.12 to 0.13 microseconds.

**C-V2X receiver queue, Falcon-512, measured ratios applied to 0.446 ms**
(`m = 2`, deadline 10 ms, 100,000 configurations, median work)

| Scenario | N | rho | Mean T (ms) | Max T (ms) | Pr[T > 10 ms] |
|---|---:|---:|---:|---:|---:|
| Steady | 100 / 150 / 200 | 0.46 / 0.69 / 0.92 | 0.58 / 0.63 / 0.69 | 0.98 | 0 |
| No caching | 100 | 0.56 | 0.78 | 5.32 | 0 |
| No caching | 150 | 0.84 | 1.15 | 9.15 | 0 |
| No caching | 200 | 1.12 | unstable | | |
| Cold start | 100 | | 0.66 | 5.22 | 0 |
| Cold start | 150 | | 0.78 | 8.62 | 0 |
| Cold start | 200 | | 3.46 | 28.62 | 0.119 |

Files:

- `results/bench_intel-i9-10980hk-run1.{md,json}` and `-run2.{md,json}`
- `results/queuesim_fixed.md`
- `results/queuesim_measured_run1.md` and `_run2.md`
- `results/raw/`: 48 sample files, about 4 MB

## Reproducing the paper's tables

| Paper table | Command |
|---|---|
| Fixed-work receiver delay (M/M/1 rows are analytic) | `go run ./cmd/queuesim -mode fixed -configs 100000` |
| Encoded IEEE 1609.2 sizes | `go run ./cmd/bench` (sizes table) |
| Receiver authentication time | `go run ./cmd/bench` (timing table) |
| C-V2X queue with measured ratios | `go run ./cmd/queuesim -mode measured -configs 100000 -tag <tag>` |

Sizes are deterministic. Timings depend on the machine; the ratios to the
primitive are what the paper uses. Simulation results with the default seed
are reproducible exactly for the `fixed` mode and, given the same raw sample
files, for the `measured` mode. To rerun the simulation on the published
samples:

```sh
go run ./cmd/queuesim -mode measured -configs 100000 -tag intel-i9-10980hk-run1
```

## Running on ARM

The code is pure Go and cross-compiles without changes.

Build for a 64-bit ARM Linux board, for example a Raspberry Pi 4 or 5:

```sh
GOOS=linux GOARCH=arm64 go build -o bench-arm64 ./cmd/bench
GOOS=linux GOARCH=arm64 go build -o queuesim-arm64 ./cmd/queuesim
```

In PowerShell:

```powershell
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o bench-arm64 ./cmd/bench
```

Copy the binary to the board and run `./bench-arm64 -tag rpi4`.

**Continuous integration.** `.github/workflows/bench.yml` runs the tests, the
benchmark, and the measured-mode simulation on GitHub's `ubuntu-24.04`
(x86-64) and `ubuntu-24.04-arm` (64-bit ARM) runners, and uploads the result
files as artifacts. It runs on a push to `main` or manually from the Actions
tab. Hosted runners are shared machines, so expect more timing noise than on
a dedicated board. This workflow has not been run yet.

## Getting stable measurements

- Close other applications and plug the laptop in; power-saving modes change
  CPU frequency.
- Run the benchmark twice with different `-tag` values and compare medians.
  The runs in `results/` differ by under 0.3%.
- Prefer medians and the median ratios over means: means include occasional
  operating-system interruptions (visible in the p99.9 and max columns).
- On Linux, pinning the process to one core (`taskset -c 2 ./bench ...`)
  reduces jitter further.

## Limitations

- **Platform.** The published runs are on an x86 laptop, not an automotive
  processor. The paper applies the path-to-primitive ratios to a published
  ARMv8 primitive time. That assumes the non-signature work (hashing,
  decoding, lookup) scales with the platform the way signature verification
  does.
- **Implementations.** All schemes are pure Go. Optimised C or assembly will
  have different absolute times and possibly different ratios.
- **Receiver path.** No revocation check, misbehaviour check, or replay
  cache. Each adds work.
- **Encoding.** One profile, with hypothetical post-quantum identifiers and
  an issuer using the same scheme as the pseudonym certificate. A region or
  assurance field would add bytes. Pseudonym certificates in deployed
  systems are often implicit (ECQV) certificates, which this code does not
  model.
- **Simulation.** One serial verifier; random, collision-free resource
  selection; no radio effects (no packet loss, half-duplex, or propagation).
  Every sender is assumed to be in range of the receiver.
- **Security.** Not audited. See [SECURITY.md](SECURITY.md).

## Third-party code

- `third_party/dilithium-go` is a copy of
  `github.com/lattice-safe/dilithium-go` v0.3.0 (MIT License, upstream
  `LICENSE` retained). Version 0.3.0 does not build on amd64: a generated
  assembly file is missing upstream and its generator fails. The copy has two
  build-file changes that select the portable Go number-theoretic transform
  on every architecture. No cryptographic code is changed, and the library's
  own tests, including the NIST known-answer tests, pass on the copy. See
  `third_party/dilithium-go/VENDORING.md`. `go.mod` points to the copy with a
  `replace` directive.
- `github.com/lattice-safe/falcon-go` and
  `github.com/lattice-safe/lattice-slh-dsa-go` are used unmodified as Go
  modules (MIT License).
- `golang.org/x/crypto` and `golang.org/x/sys` are indirect dependencies
  (BSD-3-Clause).

## Citing

If you use this code or its results, please cite the paper:

```bibtex
@article{kiraz_pqv2x,
  author = {Kiraz, Mehmet Sabir and Kardas, Suleyman and Uzunkol, Osmanbey
            and AlMarshoud, Mishri and Al-Bayatti, Ali},
  title  = {Post-Quantum Migration for {V2X}: A Framework for Broadcast
            Authentication in Mixed Classical/{PQ} Deployments},
  year   = {2026},
  note   = {Manuscript under preparation}
}
```

## License

MIT License. See [LICENSE](LICENSE). The software is provided "as is",
without warranty of any kind, and **has not been audited**.
