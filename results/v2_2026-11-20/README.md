# Version-2 results (2026-11-20)

New measurements and simulations; the original results in `results/` are
unchanged and regression-tested.

| File | Content |
|---|---|
| `implementation_report.md` | what changed, how every number was produced, tests, limitations |
| `paper_update_values.md` | exact values and suggested sentences for the manuscript, by section and table |
| `arm_benchmark_comparison.{md,csv}` | laptop (published) vs native x86-64 and ARM64 Linux, GC on and off |
| `queue_comparison.{md,csv}` | legacy abstract C-V2X model vs certificate-aware scheduling, laptop ratios |
| `payload_sensitivity.{md,csv,json}` | exact COER size vs application payload, largest fitting payload |
| `simulations/` | every simulation report (standard and extended) and CSV; runner-ratio comparison; legacy regression note |
| `runs/ci-run<id>/<runner>/` | per-runner benchmark JSON/CSV/Markdown, raw samples (ns), environment, runner-ratio simulations |

CI runs: 44820736544 (commit a4a20d5) and 44822427947 (commit 0a262ef),
GitHub Actions, Go 1.25.6. The comparison tables use run 44822427947; run
44820736544 is kept as an independent repetition (different x86 host).
