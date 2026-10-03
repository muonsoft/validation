# Manual validation skill evaluations

Codex orchestrates the run and reviews the code. OpenCode is the coding agent;
the intended model is **DeepSeek V4.1 Flash**. Supply its exact model ID; V2 can
resolve its provider from the existing local connection. No model is substituted.
These evals are not part of CI or the library's normal test script.

Invoke `$eval-muonsoft-validation` from this repository. The
[orchestrator skill](../../.agents/skills/eval-muonsoft-validation/SKILL.md)
contains the rubric. It is repository-local; public skill users do not need it.

## Prerequisites and isolation

Use a dedicated Linux environment with Python 3.10+, Git, Go 1.24+ with race
support, and OpenCode. The adapters target OpenCode 1.14.19 and 2.0.18;
preflight rejects unsupported skill-discovery output. Keep the existing Git tag
`v0.19.0` available; fetch it if missing, never create a release tag for evals.
Go dependencies must be cached or downloadable.
V2 additionally requires Linux bubblewrap (`bwrap`) and enabled user namespaces.
Cache the Go module dependencies before starting: the worker's module cache is
mounted read-only.
For a controlled comparison, use a dedicated dependency-only module cache, not a
shared cache containing copies of this repository and its maintainer skills:
`GOMODCACHE=/tmp/validation-eval-modcache go mod download`. Pass the same
`GOMODCACHE` to `doctor` and `run`.

On V2, omit `--config` to reuse the running local OpenCode Console connection
(`opencode` or `opencode-go`). The runner reads the local model catalog and the
active credential from SQLite without modifying the user's database. A bare model
ID must resolve uniquely or match that exact model in recent-model preferences;
an explicit `provider/model` avoids ambiguity. Select reasoning variants with
`#variant`; UI variant preferences are not inherited.

Only the selected model definition is copied into the isolated config. Its access
token and custom headers are passed through environment variables, never saved in
the config or copied as login storage. An expired OAuth token requires refreshing
the connection in OpenCode; the runner does not rotate the user's credentials.
Other providers require `--config`.
Workers receive an allowlisted runtime environment plus only the credentials
referenced by that provider configuration. Unrelated API keys and shell startup
hooks are excluded. The selected provider credentials remain available to the
OpenCode process; treat raw worker logs as private because shell output can expose
its environment. Do not publish raw logs without checking for secrets.

Alternatively, supply a provider-only JSON config containing only `$schema` and
`provider` (V1) or `providers` (native V2).
Configure the exact model and use `{env:VARIABLE}` for credentials, including
custom authorization headers. See the [OpenCode configuration guide](https://opencode.ai/docs/config/).
The runner uses [`opencode run --model ... --format json`](https://opencode.ai/docs/cli/#run)
in fresh sessions. Authentication must work from the explicit config and
environment; existing OpenCode login storage is never copied into eval profiles.

V1 uses `--pure`. V2 starts a private, temporary loopback server for each operation,
waits for the requested model's asynchronous catalog registration, and connects
with `--server`. The server is stopped even on timeout or cancellation. This avoids
both touching the user's background service and mistaking an initially empty V2
catalog for an unavailable model. Skills are inspected through V2's `/api/skill`.
The runner sets each child's `PWD` to its actual workspace: V2's run command
prefers that variable over the process working directory.
Each V2 operation gets a fresh session database and Go build cache. Its server
runs inside bubblewrap: only the task workspace and private profile are writable;
system tools, Go module dependencies, and library snapshots are read-only. The
checkout, previous attempts, reference solutions, and the real home directory
are absent from that filesystem. There is no fallback to an unsandboxed V2 worker.
See the [V2 CLI documentation](https://opencode.ai/v2/docs/cli/commands/).

The runner creates isolated XDG directories, disables external/Claude skills and
plugins, and installs the treatment skill under the workspace's native
`.opencode/skills/` directory. Discovery must return exactly the expected user-skill
set. V2's built-in `opencode` and `report` skills are recorded and accepted only
at their built-in paths; tool permissions deny loading them in both variants.
`HOME` and user settings are not modified. V2's `OPENCODE_TEST_HOME` points at an
empty profile directory to isolate its separate Claude/Agents skill discovery.
Do not use provider setups that
inject remote organizational instructions. Preflight cannot prove absence of
all managed settings: the dedicated environment must have none.

Workers can inspect library Go source through their module replacement. They are
not given eval tests, reference solutions, or orchestrator instructions. The V1
adapter provides workflow isolation only: use a dedicated clean environment for
it. V2 also isolates the worker filesystem and PID namespace, while retaining
network access for the provider. Inspect logs for contamination in either mode.
Web/MCP tools and additional agent delegation are disabled; local shell access
remains necessary for editing and Go checks.

## Commands

Run from the repository root. Output paths must be new directories outside the
checkout. Replace placeholders with paths and the exact model ID on your machine.

```bash
# V2: use the existing OpenCode Console login, without a separate config.
python3 evals/muonsoft-validation/run.py doctor \
  --model opencode-go/deepseek-v4.1-flash --output /tmp/validation-eval-v2-preflight
python3 evals/muonsoft-validation/run.py run \
  --model opencode-go/deepseek-v4.1-flash \
  --profile smoke --output /tmp/validation-eval-v2-smoke

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

# Ninety attempts: fifteen cases, two variants, three repeats.
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

To continue only attempts that have never started, pass `--resume /path/to/old-run`
to `run` along with a new `--output` and the same model, config, profile, selected
cases, and timeout. The runner requires matching skill, suite, case manifest,
runner, library snapshots, Go/OpenCode versions, and provider configuration hashes.
It rejects active runs and concurrent continuations using a process lock. Completed,
failed, cancelled, and interrupted attempts are inherited without retries; stale
`running`/`preparing` statuses become `interrupted` only in the new manifest.
Previous artifact directories are linked read-only by convention and must be kept
at their original locations; the continuation never writes through those links.
An existing `review.json` is carried forward. Further continuations can use the
latest output. A changed revision requires a separate measured run.

Continuation is supported only for runs recorded by this version of the runner
(with `runner.lock`, configuration hash, and selected case IDs). Historical runs
without those fields cannot be automatically resumed. To retry a failed attempt,
start a separately recorded run; do not overwrite or merge it as a replacement.

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

The full suite now includes 15 cases: the original 12, an already-scoped batch
collection (`13-scoped`), and two further negative controls (`14-errors`, `15-copy`).
The new cases are failure-informed follow-ups, not an independent holdout. Smoke
remains four cases; use `--case 07-batch --case 13-scoped` for focused path checks.

## Reports and scoring

The [2026-10-03 report](reports/2026-10-03-deepseek-v4.1-flash.md) and its sanitized
JSON preserve the original 12-case run. They do not measure the subsequent skill
correction or suite expansion.

Every attempt preserves its prompt, solution, diff, JSON events, separate stderr,
and independent test logs. `run.json` records source hashes, configuration, tool
versions, status, duration, skill-loading evidence, and provider-reported tokens
and cost when available. Absent usage is unknown, not zero; cost is not independently
verified.
V2 additionally saves `completion.json`: the session database must report success
in the expected workspace, and the completed final assistant text must match the
CLI log. This handles V2's final-text reconciliation without treating a missing
`step_finish` event alone as failure or trusting a zero exit code alone.

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

`tokens.json` and the report compare complete graded pairs with complete normalized
usage, retaining test failures but excluding execution failures symmetrically.
They show input, output, reasoning, cache reads/writes, totals, medians, per-case
comparisons, and successful-pair/applicable-case subsets. Absent counters remain
unknown, including for interrupted attempts. Summed cache reads count repeated
contexts on every request; token volume is not a monetary charge. No price table
or billing verification is implied. Saved case requirements prevent new suite
revisions from changing the interpretation of earlier results.

Keep original artifacts when improving the skill: each revision gets a separate
run. Do not automatically commit reports, provider configs, workspaces, or keys.
A smoke run validates the workflow, not statistical improvement. No real model
results are implied by deterministic harness tests passing.
