---
name: eval-muonsoft-validation
description: Analyze saved muonsoft-validation evaluation results, review generated solutions, and optionally run the standalone OpenCode evaluator when requested.
metadata:
  internal: true
---

# Evaluate the validation skill

Read the [runbook](../../../evals/muonsoft-validation/README.md). The default workflow
is a standalone CLI run followed by analysis in Codex. An analysis request does
not authorize model inference. For saved results, start with `report.md`,
`results.json`, `tokens.json` and `run.json`; use the exact materials saved with
that run rather than interpreting it through a newer suite. Gzip logs are lossless
and remain available for evidence inspection. No original machine or session
database is needed for ordinary review.

When explicitly asked to run evaluations, locate the checkout by module
`github.com/muonsoft/validation`, not a fixed machine path. OpenCode implements
the measured tasks; do not coach it or implement its submissions yourself.

## Prepare and run

1. Require the exact model requested by the user and a new output directory
   outside the checkout. The intended coding model is DeepSeek V4.1 Flash.
   On OpenCode V2, omit `--config` to reuse the existing local Console connection;
   the runner resolves an exact bare model ID through the catalog and recent-model
   preferences. Do not ask for a separate provider config before checking this.
   V1 and other providers require a provider-only JSON config with environment
   credentials. Never substitute a model; select reasoning variants only when requested.
2. Use `smoke` unless the user selects `full` or particular cases. Main smoke runs
   three cases × two variants × one repeat; full runs seven × two × three.
   The earlier 15 cases are available with `--suite legacy`; keep their results separate.
3. Run `doctor`: model availability and isolated skill discovery, no inference
   request. A failed preflight is an environment issue, not a model score.
4. On a new suite revision, run `check` and the runner unit tests. They check
   examples, reference solutions, failing starters, and runner behavior without
   model access. These evaluations are manual and are not part of CI.
5. Run the runner's `run` command. It checks materials automatically, generates
   reports, and defaults to compact retention. Use `--skip-check` only after
   verifying the exact revision separately. Poll with short tool waits to keep
   the user informed. Do not launch a second run while waiting for the first.

Use a dedicated eval environment. The runner installs only the public skill for
one variant and none for the baseline, using identical task prompts and fresh
sessions. Workers receive no independent tests or reference solutions. V2 requires
bubblewrap and isolates each worker's filesystem, session database, and build
cache. V1 provides workflow separation only; follow the runbook's clean-profile
requirements. Inspect logs for contamination in either mode and exclude any run
where a worker read hidden tests, references, or another variant's skill.

Do not coach the worker, provide expected answers, edit its solution, or retry an
attempt invisibly. Cancellation stops further attempts and preserves artifacts.
A retry uses a new output directory. To continue unstarted attempts of an unchanged
run, use `--resume` with a new output as described in the runbook; recorded attempts
are preserved and never retried by continuation. Never modify the public skill during a
measured run or mix revisions in one comparison.

## Review independently of test scores

After execution, read the task, submitted source files, diff, Go-test logs, and JSON
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

In the main suite, report built-in/project-rule reuse explicitly under
`rule_ownership`, especially for `20-reuse`. Tests cannot establish that a manual
check reused an existing rule. Cite concrete code; do not treat loops or helpers
as intrinsically wrong. `21-constraint` intentionally requires a custom rule.
Paths are JavaScript-style `PropertyPath.String()` and typed segments, not JSON
Pointer. Multiple independent violations may share one path. See the
[suite contract](../../../evals/muonsoft-validation/SUITE.md).

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
failures, missing evidence, and paired token volume from `tokens.json`. Distinguish
cache reads from other counters and token volume from monetary cost. A smoke run checks the workflow, not a statistical
claim of improvement. Recommend skill corrections from observed failures, leaving
the measured revision and original artifacts intact.
