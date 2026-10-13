# Evaluation suite contract

This suite uses fictional applications and independent fixtures. No private
project source, imports, services, identifiers or documentation are required.
It is a deliberate development set, not an independently authored holdout. Freeze
its revision before measuring a model and keep each revision's results separate.
Main and transfer results are reported separately. Keep historical runs with their
saved materials; do not resume or regrade them against a changed suite.

Main revision 3 explicitly defines the repository call in `19-update` as
`NameTaken(ctx, normalizedName, loaded.ID)`. Earlier revisions did not specify the
order of the two string arguments, so their argument-order failures are not clean
evidence of a skill failure.

## Scenarios

| ID | Change requested | Critical observations | Smoke |
| --- | --- | --- | --- |
| `16-request` | Implement a request with contact, visits and nested tasks | Parent/collection/child accumulation, compound sibling uniqueness, indexed paths, standalone children, bounds | Yes |
| `17-settings` | Extend an optional scheduling block | Nil vs explicit empty, independent and dependent bounds, multiple violations at one path, publish groups, richer-method context and language | No |
| `18-tree` | Repair a recursive form | Sibling-local uniqueness, depth barrier with sibling independence, literal map keys and typed path elements | No |
| `19-update` | Normalize and update a resource bundle | Full-state replacement, immutable loaded state, local barrier, one batch lookup, every missing occurrence, original indices, technical errors, save only after success | Yes |
| `20-reuse` | Extend and simplify existing profile validation | Built-in rules, existing child and injected project constraint, nullable fields, eager accumulation; reuse reviewed separately | No |
| `21-constraint` | Create a category-aware string constraint and use it twice | Requiredness composition, dependency errors, codes, parameters, EN/RU, caller context and repeat/concurrent reuse | No |
| `22-display` | Group existing errors for display | Exact paths, message order and multiplicity, independent results, unchanged validation; negative skill-selection control | Yes |

Each prompt specifies the observable contract, protected exports, error codes,
path placement and permitted edits. Implementations span existing files and may
add private top-level Go helpers. Supplied contracts/dependencies stay protected.
Cases require no actual database, HTTP server or private service. Dependencies
are fakes exercised by hidden tests. Workers never see checks or references.

Main full: 7 cases × 2 variants × 3 repeats = 42 attempts. Main smoke: cases 16,
19 and 22 × 2 variants × 1 repeat = 6 attempts. More complex attempts may use
more tokens, so fewer
attempts do not by themselves establish a cost reduction.

## Paths and error sets

The primary path contract is `PropertyPath.String()`, for example
`request.visits[1].tasks[1].duration`. There is no JSON Pointer requirement.
Selected tests also compare each element's value and `IsIndex()` type. Literal
keys such as `"0"` must be properties, not array indices. The form tests cover
empty keys, dots, slashes, tildes, quotes and backslashes using the library's
bracket notation and escaping.

Expected violations are multisets of (path, error identity). Multiple distinct
errors on the same property are valid. Extra, missing or repeated unexpected
violations fail. Error-list order is not prescribed unless the task explicitly
requires it; map iteration must not cause flaky tests. Presentation grouping has
an explicit first-seen order contract.

Eager validation means accumulating independent violations. It does not mean
performing a dependent lookup on invalid prerequisites. Technical errors preserve
identity and stop later work. Tests observe dependency calls and saved state as
well as returned errors.

## Behavioral grading and review

A scenario passes only if every named requirement and the test process pass on
both library versions. A category rate is diagnostic, not partial scenario
success. Model/provider failures and incomplete attempts are reported separately.
Only complete matched pairs enter the main token comparison. Repeated attempts
of the same case are not independent task samples.

Use the same review rubric for both variants:

- `api`: correct installed API; supplied validator, context and dependency reuse.
- `rule_ownership`: use available matching built-ins and project rules; reuse child
  validation rather than duplicate it at the parent/command boundary.
- `execution_flow`: accumulate independent errors, honor prerequisites, preserve
  operational failures and avoid mutation during validation.
- `scope`: focused changes, no unnecessary framework or fluent builder.

In `20-reuse`, explicitly identify each handwritten replacement of a matching
existing rule and explain the mismatch or avoidable duplication. Do not infer
reuse from passing tests or penalize manual loops needed for actual domain work.
`21-constraint` deliberately asks for a new domain rule; a custom implementation
there is expected. Do not score code by similarity to the reference or count of
`if` statements. Prefer reviewing code without variant labels when practicable;
record whether review was blinded. Report missing reviews as unreviewed.

## Verifying the evaluator

`check` executes public skill examples and each reference/starter on both library
versions with the race detector for fixture tests. It then applies declared
mutations to reference code and requires a specific behavioral test to fail:

- lost parent prefix;
- early return that hides independent violations;
- incorrect nested property segment;
- shifted batch indices;
- technical errors converted to violations;
- bypassed existing rule;
- lost duplicate display messages.

A compilation error is not an acceptable mutation result. These checks validate
the evaluator, not model quality. The immutable mutation recipes live in
`cases.json`; modifying a reference requires reviewing its mutation too.

The runner tests simulate OpenCode without model inference. They verify isolation,
completion/usage accounting, multi-file grading, compaction, lock handling,
continuation, and detached report regeneration. See the [runbook](README.md) for
commands and retention policy.

## Transfer suite (revision 1)

Select `--suite transfer`; `transfer-cases.json` owns its manifest. Main fixtures
and score denominators remain unchanged. All three cases are in smoke
(6 attempts); full repeats each three times (18 attempts). They target v0.19.0 and
the current snapshot with the same paired isolation, race checks, and scoring.

| ID | Different contract | Observable requirements |
| --- | --- | --- |
| `23-catalog` | Map-rooted tree, depth starts at zero, nil nodes, repeated labels allowed | Independent branches, literal typed keys, inclusive depth barrier, no borrowed sibling-uniqueness rule |
| `24-patch` | Explicit presence flag plus nullable integer enum | Omitted vs clear vs set, custom error identities, activate without default group, Unicode boundaries and translations |
| `25-branches` | Lookup prerequisites are local to each item | Invalid siblings do not create a global barrier, original order and repeated IDs, technical error stops later calls |

This set was authored with knowledge of development evals and has already been
measured; it is **not an independently authored holdout**. Do not use its own results to
revise the skill and still call subsequent runs unseen validation. Freeze skill,
runner, manifest, and fixtures before inference; their saved hashes identify the
measured bytes. Keep every revision and retry in a separate output directory.
No new model results are claimed by adding or checking these fixtures.

Reference solutions must pass both versions. Nil-return starters must compile and
fail behavioral tests. The mutations respectively change root depth, permit zero
enum values, and replace a local branch continuation with a global return; each
must fail its named test without compilation errors. Tests assess the contract,
not resemblance to the reference solution.

For an independent holdout, have an author who has not seen development failures
supply additional tasks and freeze them before showing results to skill authors.
Measure existing main scenarios separately; transfer scores do not establish
a token improvement on a different task.
