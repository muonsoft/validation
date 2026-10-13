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

`check` executes public skill examples on the current checkout and each
reference/starter on both fixture library versions with the race detector. It then applies declared
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


## Built-in constraints suite (revision 1)

Select `--suite builtin`; `builtin-cases.json` owns four scenarios targeting
v0.20.0 and current, with 8 smoke or 24 full attempts. Main/transfer fixtures and
the historical evaluation report are unchanged. This is a development suite,
authored alongside the skill; it is not an independent holdout. No model run or
improvement is claimed by adding these fixtures.

Prompts state domain semantics and required error contracts without naming the
constructors or the skill. Both variants can inspect the installed library.
References are hidden examples of correct implementations. Tests accept
behaviorally equivalent solutions; they do not grep source for API names.

| Case | Behavioral focus | Built-in reuse to inspect during review |
| --- | --- | --- |
| `26-publication` | Country/language/locale/currency distinctions, calendar week 53 and inclusive bounds, Unicode word boundaries, simultaneous failures | Existing format constraints, week options, word-count constraint instead of lists, regexes, calendar parsing, or `strings.Fields` |
| `27-payment` | Selected schemes and independent checksum, requiredness, optional IBAN/BIC, unmodified input | Compose card scheme and Luhn rules; reuse IBAN/BIC validation instead of implementing checksums or formats |
| `28-password` | Nil versus empty, inclusive score threshold, custom estimator calls, invalid scores, secret handling | Strength constraint with minimum/estimator options, preserving built-in configuration-error behavior |
| `29-upload` | Accumulated metadata before I/O, Windows names, suffix/MIME/content separation, bounded reading, technical errors | Compose name/extension/MIME/content constraints; keep I/O and requiredness in application code |

The four mutations respectively remove the lower week bound, remove Luhn,
lower the password threshold, and allow PNG content for a text upload. Each must
compile and fail its named behavioral test on both supported versions. Nil-return
starters must compile and fail tests. `check` also validates all portable skill
examples on the current checkout, independently of fixture release baselines.

Review built-in reuse explicitly under `rule_ownership`, separately from tests:

- **2:** Appropriate built-ins express matching rules, with options/composition
  preserving the full contract; custom glue is limited to actual application needs.
- **1:** Correct behavior but avoidable reimplementation of an available rule,
  or low-level wrappers that unnecessarily reproduce constraint behavior.
- **0:** Duplicated/replaced rules omit required semantics or diverge from the contract.

Cite the submitted code and name reused or unnecessarily reimplemented rules.
Report behavior pass rates and reuse scores per case for both variants. Also
record whether the built-in reference was read when the logs establish it;
entrypoint loading alone does not establish use of that reference. Do not treat
manual loops, safe adapters, or absence of a particular constructor as automatic
failure. Inspect secret handling in password messages/parameters/logs and total
read behavior in uploads. Generated code may pass tests without proving all of
these properties. Keep measurements separate from main/transfer and freeze this
revision before the manual model run.


## Application suite (revision 1)

Select `--suite application`; `application-cases.json` owns five additional cases,
all included in smoke (10 attempts); full repeats them three times (30 attempts).
Fixtures target v0.20.0 and current. Main, transfer, and builtin fixtures are
unchanged. The public skill is unchanged for this suite so a manual run can assess
its existing guidance on more composed tasks.

| Case | Application requirement | Reuse and behavior to review |
| --- | --- | --- |
| `30-editor` | Refactor existing article validation, correcting incomplete checks without changing its public error contract | Unicode length, empty-value semantics, custom errors/messages, editor and publish-only groups, independent failures |
| `31-profile` | Three-state country/locale patch with business allowlists | Omission ignores stale values; null clears; requiredness, format, and policy are sequenced per field without a global barrier |
| `32-products` | Product batch with optional stock codes | Inclusive collection bounds; validate each product's own rules; skip empty uniqueness keys; report every duplicate with typed indices even in an oversized batch |
| `33-attachments` | Accept PNG/JPEG attachments with metadata and content agreement | Reuse file rules; compose an explicit extension/MIME/content policy; reject metadata before reading; bound reads and preserve partial-read errors |
| `34-domain` | Combine email syntax with an injected tenant policy | Preserve the protected policy; run it once only after email prerequisites, including when name fails; retain context, groups, language, factory and technical errors |

Prompts describe product contracts and observable errors without naming the
constructors to choose. Protected contracts expose error identities and existing
project APIs; they are not recipes for the implementation. Equivalent code passes
behavioral tests. Source review separately decides whether built-in and project
rules were reused appropriately. A handwritten reimplementation passing tests is
not evidence of good rule reuse; conversely, small field-sequencing helpers and
extension-to-content policy maps are appropriate application logic.

Use the existing four review criteria with these anchors for `rule_ownership`:
2 means appropriate built-in/project reuse, 1 means correct but avoidably duplicated
logic, and 0 means required rules are absent or semantically altered. Cite concrete
code per case. Test scores do not prove that a refactoring used built-ins. In
`34-domain`, tests also observe calls to a changing policy, so copying one tenant's
policy cannot satisfy the contract. In `32-products`, calling Product.Validate
alone verifies that item rules have a reusable owner.

Reference solutions must pass both fixture versions. Starters must compile and
fail behavioral tests. Mutation checks remove publish-group format validation,
inspect an omitted patch field, count empty uniqueness keys, skip element checks,
accept both content types independently of extension, or omit the tenant policy.
Each mutation must compile and fail its named requirement, never merely fail to
build. `check` also executes the current skill's portable examples.

This suite was authored after reviewing builtin development results. It is a
separate development set, not an independently authored holdout. Keep its scores
separate from builtin; do not retroactively change that run's denominator. No new
model results are claimed. Freeze these materials before inference and use a new
output directory for each measured revision or retry.
