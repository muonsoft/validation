# Agent skill evaluation: October 6, 2026

The final main-suite evaluation recorded **20/21 successful attempts with the
muonsoft-validation skill**, compared with **17/21 without it**. Total recorded
token volume fell by 56.3%, including cache reads, or 48.2% excluding cache reads
and writes. One attempt with the skill still missed an explicit requirement.

This report describes one measured run. It supports using the skill as practical
guidance; it does not establish a general success rate for AI agents.
For installation and usage, see the [agent skill guide](agent-skill.md).

## Setup

| Parameter | Value |
| --- | --- |
| Run identifier | `2026-10-06-113905-final/main` |
| Model | `opencode-go/deepseek-v4.1-flash` |
| Agent runtime | OpenCode v2.0.18 |
| Go toolchain | go1.27.1 linux/amd64 |
| Repository revision | `166eb62d49daf2817e3a71e79d55d129e0ac9ca4` |
| Suite | Main, revision 3, full profile |
| Repetitions | 3 per scenario and variant |
| Attempts | 7 scenarios × 2 variants × 3 repetitions = 42 |
| Timeout | 900 seconds per attempt |
| Compatibility checks | v0.19.0 and the current library snapshot at the recorded revision |

Each pair used the same task contract in fresh sessions. The treatment variant
had the public skill installed; the baseline did not. OpenCode's V2 workers used
isolated workspaces, profiles, session databases, and build caches. Hidden grader
tests and reference solutions were not mounted in worker filesystems. See the
[runbook](../evals/muonsoft-validation/README.md) and
[suite contract](../evals/muonsoft-validation/SUITE.md) for the evaluation design.

A successful attempt had to satisfy every declared requirement and pass the test
process on **both** library versions. Passing individual requirements does not
turn a failed scenario into a success. All 42 attempts completed and were graded;
there were no timeouts or ungraded attempts. This report contains test results
and targeted failure analysis, not a complete scored code review of all attempts.

## Correctness

| Scenario | Without skill | With skill |
| --- | ---: | ---: |
| `16-request`: composed request validation | 3/3 | 3/3 |
| `17-settings`: optional values and caller context | 0/3 | 3/3 |
| `18-tree`: recursive collections and paths | 2/3 | 2/3 |
| `19-update`: staged repository-backed updates | 3/3 | 3/3 |
| `20-reuse`: existing validation rules | 3/3 | 3/3 |
| `21-constraint`: custom constraint behavior | 3/3 | 3/3 |
| `22-display`: error grouping, negative control | 3/3 | 3/3 |
| **Total** | **17/21 (81.0%)** | **20/21 (95.2%)** |

The paired comparison contained four improvements, one regression, and sixteen
unchanged outcomes: a net increase of 14.29 percentage points. The tree failure
occurred in different repetitions in the two variants, producing both an
improvement and a regression despite the equal scenario totals.

Skill loading was observed in all 18 treatment attempts on applicable tasks.
No loading was observed in the three treatment attempts for the display-only
negative control. These observations come from recorded tool events; absence of
an observed load is not proof of non-use.

### Failure analysis

The three baseline settings attempts failed `TestSettingsContext`, reporting an
unexpected blank-value violation for `schedule.mode`; the third also failed
`TestSettingsOptional`. All corresponding treatment attempts passed.

In `18-tree-with-2`, the submitted list helper validated elements and sibling
uniqueness but omitted the required **1–3 elements** check for both root fields
and group children. The task states this requirement explicitly. As a result:

- `TestTreeAccumulation` received eight violations instead of nine, missing
  `ErrTooManyElements` at `form.fields[0].children`.
- `TestTreeDepthAndScope` received no violation for an empty root list, where
  `ErrTooFewElements` at `fields` was required.

The path test passed. The latter failing test also covers collection size; its
name should not be interpreted as evidence that depth handling itself failed.
The baseline attempt `18-tree-without-3` omitted the same size checks. These are
solution errors, not an ambiguity in the eval contract.

## Token usage and duration

All 21 graded pairs had complete usage records and are included, including
solutions that failed tests. Counters are normalized from the runtime's reported
usage; cache reads count repeated context across requests.

| Counter | Without skill | With skill | Change |
| --- | ---: | ---: | ---: |
| Uncached input | 2,885,553 | 1,408,141 | −51.2% |
| Output | 187,234 | 160,097 | −14.5% |
| Reasoning | 540,876 | 304,959 | −43.6% |
| Cache reads | 52,534,400 | 22,663,040 | −56.9% |
| Cache writes | 0 | 0 | — |
| Input + output + reasoning | 3,613,663 | 1,873,197 | −48.2% |
| All categories | 56,148,063 | 24,536,237 | −56.3% |
| Summed attempt duration | 8,758.344 s | 5,153.715 s | −41.2% |

These token volumes are **not monetary costs**. Summed attempt duration is not
the end-to-end wall-clock duration of the evaluation pipeline.

Sixteen of 21 pairs used fewer tokens with the skill. The median total per
attempt fell from 2,298,393 to 1,181,107 tokens. Savings were also present in
subsets that reduce the influence of failed solutions and the negative control:

| Subset | Pairs | Change including cache | Change excluding cache |
| --- | ---: | ---: | ---: |
| Applicable validation tasks | 18 | −56.5% | −48.8% |
| Both solutions passed | 16 | −42.2% | −45.9% |

Token usage was not uniformly lower. The display-only negative control used
98,038 tokens without the skill and 175,732 with it (+79.2%). No skill loading
was observed there, so this increase cannot be attributed to reading the skill.

## Interpretation and limits

The measured skill revision improved correctness on settings validation and
reduced aggregate token usage, including among pairs where both solutions passed.
It did not eliminate omissions: collection-size validation remains an observed
failure. Generated solutions still require application tests and code review.

The sample contains only three repetitions per scenario and one model/configuration.
The suite was used during skill development, so it is not an independent holdout.
No statistical significance or general improvement across models is claimed.
Test success also does not establish good rule ownership or reuse; the complete
manual rubric was not scored for this run. Earlier development and transfer runs
are not pooled into these numbers.

## Evidence and reproducibility

The source bundle is retained outside the repository under run identifier
`2026-10-06-113905-final/main`. This document is a curated summary; raw submissions,
logs, and provider configuration are not distributed with it. Within that bundle:

- `report.md` and `results.json` contain per-attempt outcomes.
- `tokens.json` contains paired counters and subset comparisons.
- `run.json` records the model, revisions, hashes, durations, and configuration digest.
- `materials/` and `libraries/` retain the exact evaluation inputs and library snapshots.
- `attempts/18-tree-with-2/workspace/validation.go`, lines 16–37, contains the
  list helper that omitted size validation. Its `prompt.md`, lines 3–4, states
  the requirement; `checks/current/requirements_test.go`, lines 13 and 43,
  checks the two observed failures. The corresponding `tests.jsonl.gz` records
  the assertions, with matching outcomes under `checks/v0.19.0/`.

The recorded skill SHA-256 is
`ab47878eef92ebe746c18a466e4a7fb2a84ca2097d716b658f856434e655634a`;
the suite SHA-256 is
`14257501967ad240e755ebbbf498f6bf9c6acf91afd6114726f4e8b50c32cc4c`.
The revision identifies the measured code, not subsequent documentation changes.
Use the [runbook](../evals/muonsoft-validation/README.md) to run a new comparison
in a fresh output directory. A rerun produces new samples and may differ from
the results above.
