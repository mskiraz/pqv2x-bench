# pqv2x-bench

**What does it cost a vehicle to check a signed safety message when the
signature is post-quantum?** This repository measures it: how many bytes the
message and its certificate take, how long the receiving vehicle needs to
check them, and whether a receiver on a busy road can keep up.

It contains a Go implementation of the IEEE 1609.2 message and certificate
formats, a receiver that verifies them, a benchmark, a traffic simulation,
and the measured results.

This is the companion code for the paper *Post-Quantum V2X Broadcast
Authentication: A Feasibility Framework for Secure Migration*.
Every number in that paper's benchmark section comes from the programs here.

> **Not audited.** This is research code for measurement. It has had no
> security audit, is not constant-time, omits revocation and replay checks,
> and uses hypothetical post-quantum identifiers that no standard defines.
> Do not use it to protect real vehicles, messages, or keys. See
> [SECURITY.md](SECURITY.md).

> **Scope & Network Modeling Rationale.** This repository measures the
> **cryptographic wire format, receiver CPU verification overhead, and on-board
> queueing delay** of post-quantum broadcast authentication in V2X. It
> deliberately isolates the vehicle's onboard processing pipeline from wireless
> channel dynamics (such as RF path loss, fading, interference, and packet
> loss). This decoupling ensures that observed queue build-ups and deadline
> misses stem strictly from cryptographic verification and caching bottlenecks
> rather than radio-link noise. The empirical COER wire sizes and execution-time
> profiles produced here can directly serve as realistic inputs for full-stack
> vehicular network simulators (such as ns-3 coupled with SUMO).

## Contents

- [Background in one minute](#background-in-one-minute)
- [What it does](#what-it-does)
- [Quick start](#quick-start)
- [Signature schemes](#signature-schemes)
- [Repository layout](#repository-layout)
- [How a message is authenticated](#how-a-message-is-authenticated)
- [Encoding profile](#encoding-profile)
- [The benchmark: `cmd/bench`](#the-benchmark-cmdbench)
- [The queue simulation: `cmd/queuesim`](#the-queue-simulation-cmdqueuesim)
- [Results in this repository](#results-in-this-repository)
- [Results, version 2 (2026-11-20)](#results-version-2-2026-11-20)
- [Reproducing the paper's tables](#reproducing-the-papers-tables)
- [Running on ARM](#running-on-arm)
- [Getting stable measurements](#getting-stable-measurements)
- [Limitations](#limitations)
- [Third-party code](#third-party-code)
- [Citing](#citing)
- [License](#license)

## Background in one minute

You can skip this section if you know V2X security.

- **Vehicles broadcast safety messages.** In vehicle-to-everything (V2X)
  communication, each vehicle sends a short message with its position, speed
  and heading about ten times per second. In the US standards this is the
  Basic Safety Message (BSM). Nearby vehicles use it to warn the driver or
  brake.
- **Every message is digitally signed.** The signature lets a receiver check
  that the message comes from a legitimate vehicle and was not changed on the
  way. The format of the signed message and of the certificates is defined by
  the standard **IEEE 1609.2**.
- **Today's signature can be broken by a quantum computer.** Deployed systems
  use ECDSA P-256, with 64-byte signatures and 33-byte keys. A large quantum
  computer could forge such signatures. *Post-quantum* signature schemes
  (Falcon, ML-DSA, SLH-DSA) are designed to resist that, but their signatures
  and keys are 10 to 120 times larger.
- **That creates two practical problems.**
  - *Size.* A signed message has to fit into one radio transmission. In the
    C-V2X radio configuration the paper analyses, that limit is 2481 bytes.
  - *Time.* A receiver must check every message from every vehicle in range,
    which can be 2000 messages per second, and it must do so within
    milliseconds, because an old safety message is useless.

Four terms used throughout this README:

| Term | Meaning |
|---|---|
| **Certificate** | A small signed document that says "this public key belongs to an authorised vehicle". It is signed by an *issuing authority*. A receiver needs the sender's certificate to check the sender's messages. |
| **Certificate-carrying message** | A message that includes the sender's full certificate. To save space, a vehicle sends this form only once in every five messages. |
| **Digest message** | The other four messages. Instead of the certificate they carry an 8-byte fingerprint (hash) of it, called a *digest*. The receiver looks the certificate up by that fingerprint. |
| **Certificate cache** | The receiver's memory of certificates it has already checked, so that it does not check the same certificate again. |

## What it does

The repository does three things. Each one answers a question that earlier
work left open.

### 1. It builds real messages and measures their exact size

**Question:** how many bytes do IEEE 1609.2 certificates and signed messages
take when they carry post-quantum keys and signatures?

**Why it is open:** earlier studies *estimated* these sizes with a formula
(fixed overhead + key + signature). IEEE 1609.2 does not define post-quantum
algorithms yet, so nobody could simply look the sizes up. Estimates are risky
here, because a few bytes decide whether a message fits into one radio
transmission.

**What the code does:** it generates keys, creates certificates and signed
messages, and writes them out in the standard's wire format, the Canonical
Octet Encoding Rules (COER), field by field. The size reported is the length
of the actual byte string a vehicle would transmit. Nothing is estimated.

**Example result:**

| Scheme | Digest message | Certificate-carrying message | Fits in 2481 bytes? |
|---|---:|---:|---|
| ECDSA P-256 (today) | 132 B | 270 B | yes |
| Falcon-512 | 736 B | 2345 B | yes, with 136 B to spare |
| ML-DSA-44 | 2490 B | 6268 B | no: even the digest message is 9 B too large |

So with Falcon-512 a vehicle can still send a message and its certificate in
one transmission, and with ML-DSA-44 it cannot.

### 2. It times the complete check a receiver performs

**Question:** how much longer does a receiver need for the *whole* check of a
message than for the signature verification alone?

**Why it is open:** published benchmarks almost always time only the
mathematical signature verification, often called the *primitive*. A real
receiver has more to do for each message:

1. decode the received bytes into fields;
2. check that the message is recent;
3. find the sender's certificate in the cache, using the 8-byte digest;
4. if the message carries a certificate the receiver has not seen, verify
   that certificate's own signature against the issuing authority, and store
   it in the cache;
5. check that the certificate is still valid and allows this kind of message;
6. verify the message signature.

Step 4 is a second signature verification. Whether the extra steps cost 1% or
100% decides whether a receiver keeps up.

**What the code does:** it implements this receiver and times it, 20,000
times per case (2,000 for SLH-DSA, whose signing is slow), in four
situations:

| Case | Situation | What is timed |
|---|---|---|
| Primitive | reference | the signature verification alone |
| Digest, cached | the normal case, four messages out of five | full check; the certificate is found in the cache |
| Certificate, cached | one message out of five | full check; the message carries a certificate the receiver already knows |
| Certificate, new | first contact with a vehicle | full check, including verifying the new certificate |

**Example result for Falcon-512** (median, on the benchmark laptop):

| Case | Time | Compared with the primitive |
|---|---:|---:|
| Primitive | 63.1 µs | 1.000× |
| Digest, cached | 64.2 µs | 1.017× |
| Certificate, cached | 69.3 µs | 1.098× |
| Certificate, new | 140.9 µs | 2.233× |

So once a certificate is known, the complete check typically costs only 2% to
10% more than the bare signature verification. A new certificate more than
doubles the cost.

These are medians. Individual checks vary: on the benchmark laptop, 2% to 4%
of the digest checks and 8% to 13% of the cached-certificate checks took more
than 1.121 times the primitive, which is the per-message budget of the C-V2X
radio (see the simulation below). The benchmark therefore supports
"typically within budget", not a guaranteed worst case.

**Why the paper uses ratios, and how far they transfer.** The laptop used
here is much faster than a vehicle's on-board computer, so the microsecond
values do not describe a car. The paper therefore multiplies these ratios by
a published verification time for an ARMv8 processor (0.446 ms). The ratios
are **not** guaranteed to be portable: the non-signature work is mostly
SHA-256 hashing and decoding, whose relative cost depends on the processor's
hash instructions and on the signature library. Native runs on x86-64 and
ARM64 Linux (see [Results, version 2](#results-version-2-2026-11-20)) measure
how much they differ; results obtained by multiplying a ratio by 0.446 ms
are platform-scaled estimates, not measurements.

### 3. It simulates a receiver on a busy road

**Question:** given those times, can one receiver check all messages from 100
to 200 nearby vehicles quickly enough?

**The idea:** messages arrive at the receiver and wait in a queue; one
verifier checks them one at a time. If messages arrive faster than they can
be checked, the queue grows and messages wait. The simulation measures that
waiting time.

**What the code does:**

- It places 100, 150 or 200 vehicles around one receiver. Each sends 10
  messages per second, and one message in five carries the certificate.
- Messages arrive in 1 ms time slots (subframes). Two arrival models exist:
  - the **legacy abstract model** (used for the paper's tables so far) lets
    at most two Falcon-512 messages arrive per slot, whatever their form, and
    assumes every message is delivered;
  - the **certificate-aware model** charges a digest-form Falcon-512 message
    4 of the 10 subchannels of a slot and a certificate-carrying message all
    10, so a certificate-carrying message never shares its slot, and it
    counts the messages that do not fit (see
    [the queue simulation](#the-queue-simulation-cmdqueuesim)).
- Each message takes as long to check as measured in step 2, scaled to the
  automotive-class processor.
- It repeats this for 100,000 random arrangements of the vehicles' time slots
  and reports the average delay, the worst delay, and the share of messages
  that take longer than a 10 ms deadline.

It does this for three situations. The table below is from the **legacy
abstract model** with the original laptop ratios; the certificate-aware
results are in [Results, version 2](#results-version-2-2026-11-20).

| Situation | Meaning | Result for Falcon-512 |
|---|---|---|
| **Steady** | the receiver already knows every vehicle's certificate | with typical (median) check times the delay stays below 1 ms; with the measured timing variation the worst delay is about 9 ms at 162 vehicles |
| **No caching** | the receiver re-checks every certificate it is sent | works up to 162 vehicles (worst delay 12 ms); with 200 the queue grows without limit |
| **Cold start** | all vehicles are new at once, for example when a car enters a dense area | almost fully absorbed up to 162 vehicles; with 200, 11.9% of messages miss the 10 ms deadline in the first second |

Why 162, 125 and 200. The vehicle count a subchannel budget admits depends on
what is assumed about certificates:

| Assumption | Subchannels per message | Vehicles within 65% | Within 100% (physical) |
|---|---:|---:|---:|
| Digest-form messages only (certificates cost nothing on the radio) | 4 | 162 (0.65 x 10000 / 40) | 250 |
| One message in five carries the certificate in one slot (paper, Sec. V) | 0.8 x 4 + 0.2 x 10 = 5.2 | 125 (0.65 x 10000 / 52) | 192 |

The 162-vehicle case of the legacy simulation therefore matches only the
digest-only assumption; under the paper's certificate-aware average of 5.2
subchannels the illustrative 65% budget admits 125 vehicles. 200 vehicles
offer 104% of the physical subchannel capacity, so lossless single-slot
delivery is impossible; the certificate-aware model reports the messages that
cannot be scheduled. The 65% and 30% figures are illustrative policy budgets,
not physical limits and not measured channel busy ratios.

### What you get when you run it

- A table of encoded sizes and a table of timings for all four schemes
  (`results/bench_<tag>.md`), plus the same data as JSON and CSV, with the
  measurement environment (CPU and features, OS, Go version, build settings,
  GOMAXPROCS, garbage-collector setting, library versions).
- Every individual timing sample (`results/raw/`), so the distributions can
  be re-analysed.
- A table of queue delays for each situation and vehicle count, with full
  message accounting.
- An exact payload-size sweep of the encoded messages (`cmd/sizesweep`).

### What it is not

- **It is not a radio simulator.** There is no packet loss, interference, or
  signal propagation. The certificate-aware model accounts for subchannel
  capacity with an idealised scheduler; it is not a 3GPP Mode 4 or SAE
  J3161 MAC implementation. The goal is to isolate the receiver's computational and
  queueing bottleneck from channel noise; full-stack simulators (such as ns-3
  coupled with SUMO) can directly consume our measured wire sizes and timing
  profiles to evaluate channel-level packet delivery ratios.
- **It is not a security product.** The receiver leaves out checks a real system
  needs, and the code has not been audited (see [SECURITY.md](SECURITY.md)).

## Quick start

Requires Go 1.25 or later. No C compiler and no other tools are needed.

```sh
git clone https://github.com/mskiraz/pqv2x-bench.git
cd pqv2x-bench

# 1. Correctness and regression tests (about 1 minute; add -short to skip SLH-DSA)
go test ./...

# 2. Benchmark (about 2 minutes); writes results/bench_<tag>.{json,md,csv}
go run ./cmd/bench -tag my-machine

# 3. Queue simulation from the benchmark just produced
#    (legacy abstract model, then the certificate-aware model)
go run ./cmd/queuesim -mode measured -tag my-machine -configs 100000
go run ./cmd/queuesim -arrival certaware -mode measured -tag my-machine -configs 100000

# 4. Exact payload-size sweep of the encoded messages
go run ./cmd/sizesweep -out results/sizesweep
```

Step 2 prints two tables: encoded sizes and time per message. Step 3 prints
the queue results for a C-V2X receiver. Step 4 prints the largest payload
that fits the 2481-byte single-subframe limit for each scheme and form.

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
  queuesim/     periodic-arrival queue simulation (legacy and
                certificate-aware arrival models; + tests and golden files)
  sizesweep/    exact encoded size as a function of payload length
  benchcompare/ side-by-side report of several benchmark runs
  queuecompare/ side-by-side report of legacy and certificate-aware runs
internal/
  coer/         COER subset (ITU-T X.696): integers, lengths, choices, ...
  ieee1609/     IEEE 1609.2 SignedData and explicit certificate, encode/decode
  receiver/     certificate cache and the authentication path (+ tests)
  fleet/        test material: issuer, pseudonym certificates, signed messages
  sigs/         one interface over the four signature schemes
  hrtime/       high-resolution clock (performance counter on Windows)
  hostinfo/     measurement environment recorded with every benchmark run
third_party/
  dilithium-go/ patched copy of the ML-DSA library (see VENDORING.md there)
results/        outputs of the runs reported in the paper (never overwritten)
  raw/          one timing sample file per run, scheme and operation
  v2_2026-11-20/  version-2 results: native ARM64/x86-64 runs, certificate-
                aware simulation, payload sweep, report and paper values
.github/workflows/bench.yml   CI: native x86-64 and ARM64 benchmarks,
                analysis (tests, regression, simulations, sweep), report
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

**Scope of the validation.** This is a restricted, one-level experimental
profile, not full IEEE 1609.2 certificate processing. The issuer must already
be in the cache (provisioned out of band); its signature over the carried
certificate, its validity period and its algorithm are checked, but not its
`certIssuePermissions` or any chain-policy constraint. Not included:
revocation lists, misbehaviour detection, replay detection, region and
assurance-level checks, and chains longer than one issuer. Each omitted check
would add work to the measured paths.

The tests in `internal/receiver` run every scheme through the accepting
paths and through each rejection: truncated, over-long and empty input; a
certificate from an unknown issuer; a digest for an unknown certificate; a
PSID that the certificate does not permit; a tampered payload, signature or
carried certificate; a valid message relabelled with another sender's
digest; a stale message; and an expired certificate, cached or carried. A
cold-start test follows one receiver from an empty cache: digests are
rejected until their sender's certificate has been validated, an invalid
certificate is not cached and does not unlock its sender, and a certificate
with an unsupported key algorithm is refused.

## Encoding profile

The encoder implements a **selected subset** of the IEEE 1609.2-2022 ASN.1
types, with the options a SAE J2945/1 Basic Safety Message uses, plus
hypothetical post-quantum extension additions. It is not a full or
conformance-tested IEEE 1609.2 implementation: it encodes and accepts only
the structures listed below.

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
| `primitive` | the scheme's `Verify` call on the receiver's cached key value (see below for what that call contains) |
| `decode` | COER decoding of a digest-form message |
| `digest_hit` | full path; signer given as a digest; certificate cached |
| `cert_known` | full path; certificate carried and already cached |
| `cert_new` | full path; certificate carried, validated against its issuer, then cached |
| `digest_unknown` | digest-form message whose certificate is not cached: decoded, then discarded |

**What `primitive` contains differs by library** (checked against the
sources of the pinned versions):

- ECDSA P-256: the public key is parsed once into an `*ecdsa.PublicKey` when
  cached; `primitive` is `ecdsa.Verify`.
- Falcon-512 (`falcon-go` v0.1.1, `native` package): the cached key holds only
  the 897-byte encoded public key. Every `Verify` call decodes it, computes
  SHAKE256 over it, converts it to NTT form, decodes the signature, hashes the
  nonce and message to a point, and checks the norm. `primitive` therefore
  includes per-call key decoding; it is not "arithmetic on an already-parsed
  key". Verification is pure Go; the module's assembly files serve signing
  only.
- ML-DSA-44 (`dilithium-go` v0.3.0, vendored): no SIMD on any architecture.
  The upstream AVX2 assembly is missing, and in the vendored copy the function
  named `nttAVX2_8` is a scalar Go loop, so CPUs reporting AVX2 run scalar
  code as well. ML-DSA is used with an empty context over the 32-byte digest.
- SLH-DSA-SHA2-128s (`lattice-slh-dsa-go` v0.2.0): pure Go, no assembly; its
  SHA-2 calls go to the standard library, which uses SHA hardware
  instructions where the CPU provides them (`sha_ni` on x86-64, SHA2 on
  ARM64). That is why its times differ most between platforms.

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
- The program runs on one OS thread with the garbage collector enabled
  (the original methodology). `-gogc -1` disables the collector for a
  separate sensitivity run; `runtime.GC()` still runs once before each timed
  loop. GC-on and GC-off runs, and runs on different machines or operating
  systems, are different experimental conditions and are reported
  separately.
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
| `-gogc` | -2 (unchanged) | garbage collector for a sensitivity run: -1 disables it, a positive value sets GOGC |

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
  `ratio_to_primitive_mean`, `ratio_to_primitive_median`. Runs made with the
  current code also record, under `environment`, the CPU model and features,
  OS and kernel, Go version, compiler and build settings, GOMAXPROCS, the
  GOGC setting, and the module versions, and per scheme the iteration count,
  warm-up, fleet size and a description of the implementation. (The two
  original laptop runs predate these fields.)
- `results/bench_<tag>.csv`: one row per scheme and operation.
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
- Legacy arrival models (`-arrival legacy`, the default; these produced the
  paper's current tables and reproduce them byte for byte):
  - `unconstrained`: each sender has an independent random phase in the
    100 ms period, with no channel serialisation (a superposition of N
    periodic deterministic streams served by one deterministic server,
    often written N·D/D/1). There is no Poisson (M/M/1) model in the code.
  - `cv2x`: an **abstract** C-V2X receiver-arrival model. Each sender holds a
    distinct semi-persistent resource, drawn at random from 100 subframes
    with `m` message slots each. Every message, digest-form or
    certificate-carrying, uses its sender's slot, so up to `m` messages of
    any form arrive per subframe, and every message is delivered. All
    messages of a subframe reach the verifier together at the end of that
    subframe. `m = 2` corresponds to Falcon-512 *digest-form* messages
    occupying 4 of 10 subchannels; the radio cost of certificate-carrying
    messages (10 of 10 subchannels) is not modelled.
- Certificate-aware arrival model (`-arrival certaware`, model name
  `cv2x-cert`):
  - 1-ms subframes with 10 subchannels; a digest-form Falcon-512 message
    occupies 4, a certificate-carrying message all 10 (`-digest-subch`,
    `-cert-subch`, `-subchannels`). Two digest-form messages can share a
    subframe; a certificate-carrying message never shares one.
  - Senders, preferred subframes, certificate offsets and tie order are drawn
    exactly as in the legacy `cv2x` model from an identically seeded stream,
    so both models see the same random configurations.
  - Each message is generated at the start of its sender's preferred
    subframe and placed, in generation order, into the earliest subframe at
    or after that one with enough free subchannels, within a scheduling
    window (`-sched-window-ms`, default 100, at most one period). Messages
    that do not fit are counted as **not scheduled**; they are never
    delivered silently.
  - This is an idealised, deterministic, centralised scheduler: an upper
    bound on what a distributed MAC could deliver. It is not an
    implementation of 3GPP Mode 4 or SAE J3161 resource selection; there are
    no collisions, no half-duplex losses, no propagation loss and no
    retransmissions, and every sender is in range.
  - With the defaults, N senders offer 0.052 N of the physical subchannel
    capacity (0.52 at N = 100, 0.78 at 150, 0.842 at 162, 1.04 at 200). The
    report states whether the illustrative 30% and 65% budgets and the
    physical 100% limit are met.
- Each case runs many independent random configurations (phases or resource
  choices). A configuration is simulated for 20 periods; the first 10 are
  discarded as warm-up. Cold-start cases have no warm-up and measure the
  first 10 periods (1 s).
- Legacy model: a case with utilisation `rho >= 1` is reported as unstable
  and not simulated. Certificate-aware model: every case is simulated and
  marked unstable when the carried verification work reaches the
  verifier's capacity.

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
| `-m` | 2 | messages per C-V2X subframe (legacy) and digest-form SPS resources per subframe |
| `-seed` | 7 | random seed |
| `-senders` | `100,150,200` | comma-separated sender counts `N` |
| `-out` | none | also write the Markdown report to this file |
| `-arrival` | `legacy` | `legacy` or `certaware` |
| `-report` | `standard` | `standard` (original table) or `extended` (full accounting); `certaware` always uses `extended` |
| `-ext-out` | none | also write the extended report to this file |
| `-csv` | none | write the per-case accounting as CSV |
| `-subchannels`, `-digest-subch`, `-cert-subch` | 10, 4, 10 | certificate-aware radio model |
| `-sched-window-ms` | 100 | certificate-aware scheduling window (1 to 100) |

### Output columns

Standard table:

| Column | Meaning |
|---|---|
| `rho` | utilisation: arrival rate times mean service time in steady state |
| Mean T, p99.9 T, Max T | delay statistics over all authenticated messages |
| Pr[T > D] | deadline missed / authenticated |
| Configs with a miss | share of configurations with at least one late message |
| Discarded | unknown-signer discards / generated messages |

The extended table and the CSV add, for messages generated in the measured
periods: `generated`; `not_scheduled_digest` and `not_scheduled_cert`
(certificate-aware model only); `scheduled`; `unknown_signer_discarded`;
`authenticated`; `authentication_deadline_missed`;
`p_miss_given_authenticated` (missed / authenticated);
`on_time_fraction_of_generated` ((authenticated - missed) / generated, a
receiver-side figure that excludes all unmodelled network effects); p99 delay;
mean and maximum scheduling delay; verifier utilisation (busy time / measured
time); maximum backlog (messages waiting or in service, including the
arriving one); offered subchannel occupancy and the 30% / 65% / 100% budget
checks. Identities that the tests enforce: generated = scheduled + not
scheduled; scheduled = unknown-signer discarded + authenticated; missed <=
authenticated. Unknown-signer discards are never counted as authenticated and
never enter the delay statistics or their denominator.

The roughly 20% discard rate of the cold-start cases follows from the
transmission schedule: each sender's certificate offset is uniform over the
five messages of a 500 ms cycle, so on average 2 of its 10 messages in the
first second precede its first certificate. It is not a measured
over-the-air loss rate.

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
| Steady | 100 / 150 / 162 / 200 | 0.46 / 0.69 / 0.75 / 0.92 | 0.58 / 0.63 / 0.65 / 0.69 | 0.98 | 0 |
| No caching | 100 | 0.56 | 0.78 | 5.32 | 0 |
| No caching | 150 | 0.84 | 1.15 | 9.15 | 0 |
| No caching | 162 | 0.91 | 1.46 | 12.08 | 2.9e-6 |
| No caching | 200 (stress) | 1.12 | unstable | | |
| Cold start | 100 | | 0.66 | 5.22 | 0 |
| Cold start | 150 | | 0.78 | 8.62 | 0 |
| Cold start | 162 | | 0.83 | 11.24 | 3.9e-7 |
| Cold start | 200 (stress) | | 3.46 | 28.62 | 0.119 |

`N = 162` is the largest all-Falcon-512 load within a 65% subchannel
occupancy budget **if certificate-carrying messages cost no more radio
resources than digest-form ones** (4 subchannels); with the paper's
certificate-aware average of 5.2 subchannels the 65% budget admits 125
vehicles, and `N = 200` exceeds even the physical capacity (104%). This table
uses the legacy abstract arrival model; see version 2 for the
certificate-aware model. The table uses median service
times. With service times drawn from the measured samples ("measured
jitter" rows in the result files), the steady-state maximum delay is 9.0 ms
at `N = 162` and 12.3 ms at `N = 200`.

Files:

- `results/bench_intel-i9-10980hk-run1.{md,json}` and `-run2.{md,json}`
- `results/queuesim_fixed.md`
- `results/queuesim_measured_run1.md` and `_run2.md`
- `results/raw/`: 48 sample files, about 4 MB

## Results, version 2 (2026-11-20)

Directory [`results/v2_2026-11-20/`](results/v2_2026-11-20/); the original
results above are unchanged and regression-tested. Start with
[`implementation_report.md`](results/v2_2026-11-20/implementation_report.md);
[`paper_update_values.md`](results/v2_2026-11-20/paper_update_values.md)
lists the values and sentences for the manuscript. Labels: **[M]**
measurement, **[E]** ratio multiplied by 0.446 ms (platform-scaled estimate),
**[S]** simulation.

**Native receiver-path ratios, Falcon-512, median, GC enabled** [M]
([`arm_benchmark_comparison.md`](results/v2_2026-11-20/arm_benchmark_comparison.md)):

| Platform | Primitive (us) | Digest, cached | Certificate, cached | Certificate, new |
|---|---:|---:|---:|---:|
| Laptop i9-10980HK, Windows (published) | 63.10 | 1.017 | 1.098 | 2.233 |
| x86-64 Linux, Xeon 6975P-C (GitHub runner) | 38.46 | 1.011 | 1.041 | 2.061 |
| ARM64 Linux, Neoverse-N2 (GitHub runner) | 40.15 | 1.008 | 1.035 | 2.082 |

The cached-digest ratio transfers between platforms; the certificate paths
are lower on both Linux hosts than on the laptop, so the laptop ratios are
the most conservative. Scaled to 0.446 ms [E], the ARM64 ratios give 0.449,
0.462 and 0.929 ms instead of 0.454, 0.490 and 0.996 ms.

**Queue, Falcon-512, median work, laptop ratios** [S]
([`queue_comparison.md`](results/v2_2026-11-20/queue_comparison.md)):

| Scenario | N | Legacy max T | Cert-aware offered occupancy | Cert-aware not scheduled | Cert-aware max T | Deadline misses (cert-aware) |
|---|---:|---:|---:|---:|---:|---:|
| Steady | 100 / 162 | 0.98 ms | 0.52 / 0.84 | 0 | 0.91 ms | 0 |
| No caching | 100 / 162 | 5.32 / 12.07 ms | 0.52 / 0.84 | 0 | 1.00 ms | 0 |
| Cold start | 100 / 162 | 5.22 / 11.24 ms | 0.52 / 0.84 | 0 | 1.00 ms | 0 |
| Steady | 200 | 0.98 ms | 1.04 | 12.7% of generated | 0.91 ms | 0 |
| Cold start | 200 | 28.62 ms (11.9% misses) | 1.04 | 6.4% of generated | 1.00 ms | 0 |

When a certificate-carrying message occupies its subframe alone, each
subframe's verification work fits in 1 ms and the receiver never queues
across subframes; at 200 senders the radio, not the verifier, is the
binding constraint.

**Payload sensitivity** [M]
([`payload_sensitivity.md`](results/v2_2026-11-20/payload_sensitivity.md)):
the Falcon-512 certificate-carrying message fits the 2481-byte limit for
payloads up to 174 bytes (digest form: 1782 bytes); an ML-DSA-44 digest-form
message fits only up to 30 bytes.

## Reproducing the paper's tables

| Paper table | Command |
|---|---|
| Fixed-work receiver delay (legacy abstract model) | `go run ./cmd/queuesim -mode fixed -configs 100000` |
| Encoded IEEE 1609.2 sizes | `go run ./cmd/bench` (sizes table) |
| Receiver authentication time | `go run ./cmd/bench` (timing table) |
| C-V2X queue with measured ratios | `go run ./cmd/queuesim -mode measured -configs 100000 -senders 100,150,162,200 -tag <tag>` |

Sizes are deterministic. Timings depend on the machine; the ratios to the
primitive are what the paper uses. Simulation results with the default seed
are reproducible exactly for the `fixed` mode and, given the same raw sample
files, for the `measured` mode. To rerun the simulation on the published
samples:

```sh
go run ./cmd/queuesim -mode measured -configs 100000 -tag intel-i9-10980hk-run1 -senders 100,150,162,200
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

**Continuous integration.** `.github/workflows/bench.yml` has three jobs:

- `bench` (matrix `ubuntu-24.04` x86-64 and `ubuntu-24.04-arm` ARM64): short
  tests, the benchmark with the garbage collector enabled (original
  methodology) and, separately, disabled, the environment description, and
  both queue models from that runner's own Falcon-512 ratios;
- `analysis`: the full test suite, the legacy simulations compared byte for
  byte with `results/queuesim_*.md`, the certificate-aware simulations, the
  model comparison, and the payload-size sweep;
- `report`: the cross-platform benchmark comparison.

Every output, including raw samples, is uploaded as an artifact. Run it from
the Actions tab (bench, "Run workflow") or with
`gh workflow run bench.yml --ref <branch>`; it also runs on a push to `main`
that changes more than documentation. Hosted runners are shared virtual
machines with server CPUs (Neoverse-N2 on the ARM runner; the x86 host type
varies between runs), not automotive processors, so expect more timing noise
than on a dedicated board. The results of runs 44820736544 and 44822427947
are in [`results/v2_2026-11-20/`](results/v2_2026-11-20/).

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

- **Platform.** The published runs are on an x86 laptop; the version-2 runs
  are on x86-64 and ARM64 server cores. None is an automotive processor. The
  paper applies the path-to-primitive ratios to a published ARMv8 primitive
  time, which assumes the non-signature work (hashing, decoding, lookup)
  scales with the platform the way signature verification does. The native
  runs show that this holds closely for cached digests but not exactly for
  the certificate paths (see version 2 results).
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
  title  = {Post-Quantum {V2X} Broadcast Authentication: A Feasibility
            Framework for Secure Migration},
  year   = {2026},
  note   = {Manuscript under preparation}
}
```

## License

MIT License. See [LICENSE](LICENSE). The software is provided "as is",
without warranty of any kind, and **has not been audited**.
