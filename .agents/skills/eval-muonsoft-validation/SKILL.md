---
name: eval-muonsoft-validation
description: Run manual paired evaluations of the muonsoft-validation skill through OpenCode, review generated solutions, and produce a scored report. Requires this library checkout and an explicit eval request.
metadata:
  internal: true
---

# Evaluate the validation skill

Use only on an explicit eval request. Codex orchestrates OpenCode and reviews its
artifacts; do not implement the evaluated tasks yourself. Locate the checkout by
module `github.com/muonsoft/validation`, not a fixed machine path, and read the
[runbook](../../../evals/muonsoft-validation/README.md).

## Prepare and run

1. Require the exact OpenCode `provider/model` requested by the user, a
   provider-only JSON config using environment credentials, and a new output
   directory outside the checkout. The intended coding model is DeepSeek Flash
   4.1. Resolve its ID on the execution environment; never substitute a model.
2. Use `smoke` unless the user selects `full` or particular cases. Smoke runs four
   cases × two variants × one repeat; full runs twelve × two × three.
3. Run `doctor`: model availability and isolated skill discovery, no inference
   request. A failed preflight is an environment issue, not a model score.
4. On a new suite revision, run `check` and the runner unit tests. They check
   examples, reference solutions, failing starters, and runner behavior without
   model access. These evaluations are manual and are not part of CI.
5. Run the runner's `run` command. Poll with short tool waits to keep the user
   informed. Do not launch a second run while waiting for the first.

Use a dedicated eval environment. The runner installs only the public skill for
one variant and none for the baseline, using identical task prompts and fresh
sessions. Workers receive no independent tests or reference solutions. This is
workflow separation, not an operating-system sandbox; follow the runbook's clean
profile requirements and inspect logs for contamination.

Do not coach the worker, provide expected answers, edit its solution, or retry an
attempt invisibly. Cancellation stops further attempts and preserves artifacts.
A retry uses a new output directory. Never modify the public skill during a
measured run or mix revisions in one comparison.

## Review independently of test scores

After execution, read the task, submitted solution, diff, Go-test logs, and JSON
events for each graded attempt. A successful solution does not prove that it loaded
the skill. Only recorded skill/read tool events establish `observed`; otherwise
retain `not_observed` and inspect the raw log if needed.

Write `review.json` keyed by attempt ID. Each entry has these four criteria; every
criterion contains integer `score` (0–2) and `evidence` with a concrete explanation
and relative artifact link identifying lines or events:

| Key | 0 | 1 | 2 |
| --- | --- | --- | --- |
| `api` | Incorrect/invented API or discarded validator context | Working API with avoidable duplication/adaptation issues | Appropriate API; caller context and contract preserved |
| `rule_ownership` | Rules missing or inconsistently duplicated | Working but unnecessarily duplicated rules | Appropriate owner and reuse of rules |
| `execution_flow` | Lost violations, unsafe prerequisites, or masked technical errors | Correct result with unnecessary/unclear sequencing | Appropriate accumulation, dependencies, and error propagation |
| `scope` | Unrelated behavior/dependency changes or missed task | Solution with unnecessary abstractions | Focused, readable solution preserving unrelated behavior |

For the negative control, `api` means preserving the existing API without an
unrequested framework, and `execution_flow` means correct formatting without
altering validation. Do not penalize that task for not using this library.

Apply the same rubric to both variants. Judge equivalent implementations by
behavior and appropriateness, not resemblance to the reference. Do not score an
attempt with no gradable artifact. Keep subjective scores separate from test
percentages; neither score can conceal a critical failed requirement.

Run `report --review ...` to combine review JSON with captured results without
rerunning agents. Give the user links to `report.md`, `results.json`, and
`run.json`, and explain improvements, regressions, repeat variation, critical
failures, and missing evidence. A smoke run checks the workflow, not a statistical
claim of improvement. Recommend skill corrections from observed failures, leaving
the measured revision and original artifacts intact.
