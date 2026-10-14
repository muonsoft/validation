# Application skill evaluation: October 13, 2026

DeepSeek V4.1 Flash passed **15/15 attempts with the skill and 15/15 without it**
on the application suite. With the skill, recorded token volume fell by **23.2%
excluding cache reads/writes**, or **41.5% including them**. Total attempt time
fell by 30.3%. Source review still found avoidable custom predicates and adapters.

These measurements describe the skill at the revision below, **before the
subsequent guidance changes for choices, patch sequencing, typed constraints and
numeric ranges**. Those changes have not been evaluated against a model. This is
a development suite, not an independent holdout or a general performance claim.
The [earlier main-suite report](agent-skill-evaluation.md) remains a separate result.

## Setup and validity

| Parameter | Value |
| --- | --- |
| Run identifier | `validation-application-isolated-full-v2` |
| Created at | 2026-10-13 19:20:41 UTC; reviewed October 14 |
| Model | `opencode-go/deepseek-v4.1-flash` |
| Runtime | OpenCode v2.0.18 |
| Go | go1.27.1 linux/amd64 |
| Repository revision | `b40bb9081ca341c847118d5d4dc3c2f437ba9ea0` |
| Suite | Application, revision 1, full profile |
| Attempts | 5 scenarios × 2 variants × 3 repetitions = 30 |
| Timeout | 900 seconds per attempt |
| Isolation | Separate filesystem, process and network namespaces; restricted model gateway |

Both variants received the same task contracts and isolated sessions. The skill
was observed loading in all 15 treatment attempts and none of the baselines.
Direct inventory reads were observed in 13/15 treatment attempts; the upload
reference was read in all three upload attempts. All attempts completed and were
graded, without provider errors or unexpected changes to protected files.

Each attempt passed all three behavioral requirements under both fixture labels,
`v0.20.0` and `current`. The two library snapshots have **identical source hashes**;
these checks do not independently establish compatibility with different APIs.
The reviewer examined prompts, submitted source, diffs, test logs and events for
all 30 attempts, and scored API use, rule ownership, execution flow and scope
separately from behavioral test results.

No access to hidden tests, reference solutions or the other variant's skill was
found in these logs. Baseline product attempts did inspect their own runtime
profiles and session databases. In `32-products-without-1`, event line 130, a
request to the worker's local `/api/skill` returned the built-in OpenCode guide
despite that skill being denied through the skill tool. This was not the
muonsoft-validation treatment skill or a hidden answer; no pair was excluded.
The network sandbox does not hide a worker's own runtime API/database from its
shell tools. The detours also mean token savings cannot be attributed solely to
faster discovery of validation APIs.

Earlier runs are not pooled into this comparison: one pre-isolation run exposed
references/tests through public GitHub, and an initial gateway run failed all
attempts because the required OpenCode session header was not forwarded.

## Behavior and token volume

| Scenario | Without skill | With skill | All tokens, without → with | Change |
| --- | ---: | ---: | ---: | ---: |
| Article refactoring | 3/3 | 3/3 | 3,392,530 → 2,913,110 | −14.1% |
| Profile PATCH | 3/3 | 3/3 | 2,608,353 → 2,013,561 | −22.8% |
| Product batches | 3/3 | 3/3 | 8,717,584 → 2,874,343 | −67.0% |
| Attachments | 3/3 | 3/3 | 5,218,067 → 3,780,025 | −27.6% |
| Tenant email policy | 3/3 | 3/3 | 3,934,144 → 2,386,588 | −39.3% |
| **Total** | **15/15** | **15/15** | **23,870,678 → 13,967,627** | **−41.5%** |

There were no paired correctness improvements or regressions: all 15 pairs passed.

| Counter | Without skill | With skill | Change |
| --- | ---: | ---: | ---: |
| Input | 1,145,178 | 898,906 | −21.5% |
| Output | 102,339 | 96,561 | −5.6% |
| Reasoning | 248,377 | 152,832 | −38.5% |
| Cache read | 22,374,784 | 12,819,328 | −42.7% |
| Cache write | 0 | 0 | — |
| Input + output + reasoning | 1,495,894 | 1,148,299 | −23.2% |
| All counters | 23,870,678 | 13,967,627 | −41.5% |
| Total attempt seconds | 3,897.801 | 2,717.957 | −30.3% |

Token volume is not monetary cost. All 15 pairs have complete usage. The skill
used fewer total tokens in 13/15 pairs; the median paired change was −31.9%.
The exceptions were article repeat 3 and attachments repeat 2. Excluding the
entire products scenario still gives −26.8% across all counters and −21.2%
excluding cache reads/writes. For attachments, uncached-plus-generated tokens
increased by 0.7% despite the reduction in the total including cache reads.

## Built-in and project-rule reuse

- **Articles:** both variants use requiredness, Unicode length, email, groups and
  custom error/message options in all three repeats. No reuse advantage observed.
- **PATCH:** none of the three treatment solutions uses `IsOneOf`; all implement
  membership through `OfStringBy` and equality/switch predicates, despite reading
  the inventory. One baseline uses `IsOneOf`, but checks format twice. Two treatment
  solutions use field-local `Sequentially`; the other repeats the format predicate.
- **Products:** both variants use `HasCountBetween`,
  `HasUniqueValuesBy(...).SkipEmptyKeys().At(...)`, and reusable `Product.Validate`
  rules in every repeat. Lower token usage does not mean a better final algorithm.
- **Attachments:** treatment uses `HasContentType` in 3/3 solutions, baseline in
  2/3; the other baseline manually adapts standalone validators and violations.
  Treatment repeat 3 still writes a manual size predicate; baseline repeat 3
  creates an unnecessary `imageSizeConstraint` instead of configuring `IsBetween`.
- **Tenant policy:** treatment uses `Sequentially` in 3/3, while every baseline
  repeats email syntax through `is.Email`. However, two treatment solutions wrap
  the existing policy in `ValidatableFunc` instead of passing its typed constraint
  directly. All solutions retain the supplied policy's observable behavior.

The follow-up skill revision adds examples for business allowlists with custom
errors, field-local patch prerequisites, direct use of an existing
`StringConstraint`, and built-in numeric ranges with custom errors. These are
corrections motivated by this run, not measured improvements. No scenarios or
scoring requirements were changed in response to these findings.

## Evidence and reproducibility

This is a curated report. The full run bundle is retained outside the repository
under the run identifier above; raw session logs and provider configuration are
not published here. Bundle-relative evidence includes:

- `run.json`, `results.json`, `tokens.json`: provenance, outcomes and paired counters.
- `review.json`, `review-audit.json`, `report.md`: all 30 source reviews with line
  references, tool-read evidence and grader results.
- `attempts/31-profile-with-1/workspace/validation.go`, lines 13–21 and 51–54:
  custom membership predicates and duplicate format evaluation.
- `attempts/31-profile-with-2/workspace/policy.go`, lines 8–38, and
  `attempts/31-profile-with-3/workspace/validation.go`, lines 10–16: custom choices.
- `attempts/33-attachments-with-3/workspace/validation.go`, lines 31–33: manual range.
- `attempts/34-domain-with-2/workspace/validation.go`, lines 15–17, and
  `attempts/34-domain-with-3/workspace/validation.go`, lines 18–20: extra adapters.
- `attempts/32-products-without-1/events.jsonl.gz`, lines 130, 139–145 and 180:
  local runtime-guide and session-data inspection, without hidden eval materials.

Recorded SHA-256 hashes:

| Material | SHA-256 |
| --- | --- |
| Measured skill | `3ac88e5662da6a31ea35086d2a723680971c82df61823cdc1dc98d3020d0eda1` |
| Fixture sources | `b64427cfbeb38459d8281bbc7e2579b94de1c1e252bb124c79df409257ebde79` |
| Application manifest | `6f2c11dafc9ec723267632fb46f1e1b62a53de3d8cdd1ed96b91261d23066588` |
| Runner | `b0acfdf81b776793fe523b13ae92a34db6988483b77c54c1659d745871912295` |
| Network gateway | `bed2bda7090f04270740b65dbf25fef6d3aed479eae6a5c4eb01f5c8afe5f827` |
| Each library snapshot | `bf046920bb7151a8cc8a125fb20252f20d74879705784b3da066ff2c07d6c11d` |

Use the [runbook](../evals/muonsoft-validation/README.md) for a new comparison in a
fresh output directory. Three repeats on a development set do not establish a
stable effect size across models or applications.
