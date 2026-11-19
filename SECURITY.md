# Security notice

**This software has not been audited.** It is research code written to
measure sizes and timings for a paper. It is not a V2X security stack and
must not be used to protect real vehicles, real messages, or real keys.

Specifically:

- No independent security review or audit has been performed on this code or
  on the way it uses its cryptographic dependencies.
- The post-quantum implementations come from third-party libraries
  (`github.com/lattice-safe/*`). Their own security status is described in
  their repositories; this project makes no claim about it.
- The code is not written to be constant-time and has not been assessed for
  side channels or fault attacks.
- The receiver path is deliberately incomplete: it has no revocation check,
  no misbehaviour detection, no replay cache, and no certificate-chain
  processing beyond one issuer.
- IEEE 1609.2 defines no post-quantum algorithm identifiers. The identifiers
  and encodings used here for post-quantum keys and signatures are
  hypothetical and will not interoperate with any standard implementation.
- Private keys are generated in memory for test purposes and are not
  protected or erased.
- The software is provided "as is", without warranty of any kind, under the
  MIT License (see `LICENSE`).

## Reporting a problem

If you find a bug that affects the measurements or the encodings reported in
the paper, please open an issue in the repository.
