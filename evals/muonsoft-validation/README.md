# Manual validation skill evaluations

Codex orchestrates the run and reviews the code. OpenCode is the coding agent;
the intended model is **DeepSeek Flash 4.1**, identified by an explicit
`provider/model` ID on the execution environment. No model is chosen implicitly.
These evals are not part of CI or the library's normal test script.

Invoke `$eval-muonsoft-validation` from this repository. The
[orchestrator skill](../../.agents/skills/eval-muonsoft-validation/SKILL.md)
contains the rubric. It is repository-local; public skill users do not need it.

## Prerequisites and isolation

Use a dedicated Linux environment with Python 3.10+, Git, Go 1.24+ with race
support, and OpenCode. The adapter targets the inspected OpenCode 1.14.19 CLI;
preflight rejects unsupported skill-discovery output. Keep the existing Git tag
`v0.19.0` available; fetch it if missing, never create a release tag for evals.
Go dependencies must be cached or downloadable.

Supply a provider-only JSON config containing only `$schema` and `provider`.
Configure the exact model and use `{env:VARIABLE}` for credentials, including
custom authorization headers. See the [OpenCode configuration guide](https://opencode.ai/docs/config/).
The runner uses [`opencode run --model ... --format json`](https://opencode.ai/docs/cli/#run)
in fresh sessions. Authentication must work from the explicit config and
environment; existing OpenCode login storage is deliberately not inherited.

The runner creates isolated XDG directories, disables external/Claude skills and
plugins, and installs the treatment skill under the workspace's native
`.opencode/skills/` directory. Discovery must return exactly the expected skill
set. `HOME` and user settings are not modified. Do not use provider setups that
inject remote organizational instructions. Preflight cannot prove absence of
all managed settings: the dedicated environment must have none.

Workers can inspect library Go source through their module replacement. They are
not given eval tests, reference solutions, or orchestrator instructions. This is
workflow isolation, **not a hostile-code sandbox**: a shell process under the same
OS user can access other files. Use the dedicated environment and inspect logs
for contamination. Web/MCP tools and additional agent delegation are disabled;
local shell access remains necessary for editing and Go checks.

## Commands

Run from the repository root. Output paths must be new directories outside the
checkout. Replace placeholders with paths and the exact model ID on your machine.

```bash
# No inference: inspect availability and skill discovery.
python3 evals/muonsoft-validation/run.py doctor \
  --model '<provider/model>' --config /path/to/provider.json \
  --output /tmp/validation-eval-preflight

# No inference: check the runner, examples, references, and failing starters.
python3 -B evals/muonsoft-validation/run_test.py
python3 evals/muonsoft-validation/run.py check \
  --output /tmp/validation-eval-materials

# Eight attempts: four cases, two variants, one repeat.
python3 evals/muonsoft-validation/run.py run \
  --model '<provider/model>' --config /path/to/provider.json \
  --profile smoke --output /tmp/validation-eval-smoke

# Seventy-two attempts: twelve cases, two variants, three repeats.
python3 evals/muonsoft-validation/run.py run \
  --model '<provider/model>' --config /path/to/provider.json \
  --profile full --output /tmp/validation-eval-full

# Rebuild with Codex review; no agent reruns.
python3 evals/muonsoft-validation/run.py report \
  --run /tmp/validation-eval-smoke --review /tmp/validation-eval-smoke/review.json
```

`--case 05-stages` selects a case regardless of the profile's list; repeat this flag
for several cases. The profile still determines repetitions. `--timeout` defaults
to 900 seconds per coding attempt. Independent Go checks have a 180-second limit
per version. Failures are not silently retried. Restart into a new output directory.

## Cases and compatibility

`cases.json` records schema version, ID, title, smoke membership, expected skill
applicability, editable files, and named tests representing observable requirements.
Each `testdata/<case>/` contains:

- `prompt.md`: identical task for both variants, without naming the skill.
- `workspace/`: immutable exported contract and compilable incomplete/buggy solution.
- `checks/`: black-box tests added only after the coding process exits.
- `reference/`: one correct solution for validating tests, never shown to workers.

Fixtures become independent temporary Go modules and are not library packages.
Only `solution.go` is editable. Changes to the contract, module files, installed
skill, or additional files invalidate the submission. Temporary scratch tests must
be removed before completion. The grader compiles the submitted solution with the
original contract and independent tests, not worker-authored tests.

Tasks target v0.19.0. Each solution is checked with `go test -race -json` on both
v0.19.0 and a snapshot of the current implementation. Hashes record actual bytes,
including uncommitted library changes. Snapshots contain Go implementation/module
files, not repository instructions or eval answers. `check` also installs a copy
of the public skill by itself to check portable links and executes every complete
Go example on both versions. Every reference must pass and every deliberately
incomplete/buggy starter must fail.

## Reports and scoring

Every attempt preserves its prompt, solution, diff, JSON events, separate stderr,
and independent test logs. `run.json` records source hashes, configuration, tool
versions, status, duration, skill-loading evidence, and provider-reported tokens
and cost when available. Absent usage is unknown, not zero; cost is not independently
verified.

`report.md` and `results.json` keep deterministic outcomes separate from Codex
review. A requirement passes only when its named test passes. A scenario passes
only when all requirements and the test process pass on both versions. Compilation
failures, failing tests, timeout, unavailable environment, incomplete event logs,
and cancellation remain distinguishable. Ungraded attempts never count as success.
The report shows each repeat, version compatibility, paired improvements/regressions,
and variation across repeats.

A completed skill call or read of its entrypoint establishes `observed` loading.
Other access, including shell reads, remains `not_observed`; this is not proof of
non-use. Inspect transcripts. In the negative control, observed loading is a
potential discovery false positive, separate from correctness of the code.

Codex writes a JSON object keyed by attempt ID. Each reviewed attempt has `api`,
`rule_ownership`, `execution_flow`, and `scope`. Each criterion contains `score`
(integer 0–2) and `evidence` (explanation with a relative artifact link and line
number). The [orchestrator skill](../../.agents/skills/eval-muonsoft-validation/SKILL.md)
defines score anchors. Missing reviews remain explicitly unreviewed. Scores never
change objective test results, and equivalent implementations are accepted.

Keep original artifacts when improving the skill: each revision gets a separate
run. Do not automatically commit reports, provider configs, workspaces, or keys.
A smoke run validates the workflow, not statistical improvement. No real model
results are implied by deterministic harness tests passing.
