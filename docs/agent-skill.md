# Use validation with AI agents

The library includes an optional `muonsoft-validation` skill: instructions and
task-specific references that a coding agent can read when implementing Go
validation. It helps the agent choose the library's APIs, reuse existing rules,
and preserve application behavior while changing validation code.

## Installation

Install the Go library in your application as described in [Installation](installation.md).
Install the skill separately from your project directory:

```bash
npx skills add muonsoft/validation --skill muonsoft-validation
```

This command requires Node.js and `npx`. Select your coding agent and installation
scope in the installer. The skill is optional; `go get` does not install it, and
the Go library has no runtime dependency on an AI service.

For manual installation, copy the entire
[`skills/muonsoft-validation`](../skills/muonsoft-validation/SKILL.md) directory
into your agent's supported skills location. Keep `SKILL.md` and `references/`
together so that the relative links remain usable. Follow your agent's skill
discovery instructions to make the installation available in a session.

## What it helps with

The skill guides agents through:

- Defining reusable entity and DTO rules with `Validatable` and existing constraints.
- Accumulating independent violations while preserving paths, groups, translations,
  and the caller's validator.
- Validating nested objects, slices, maps, recursive trees, and unique values.
- Handling optional fields, enums, patch semantics, and conditional checks.
- Sequencing database or HTTP checks after their prerequisites, preserving
  technical errors, and preventing invalid data from being saved.
- Testing valid input, simultaneous violations, boundary cases, and dependency calls.

The [skill entrypoint](../skills/muonsoft-validation/SKILL.md) links to focused
references for each task. Its examples target v0.19.0 and are checked against that
version and the library checkout by the manual evaluation workflow. The agent
is instructed to inspect your installed version and existing application rules
before making changes.

## Example task

Give the agent the application requirements and the code it should modify. For example:

> Use the muonsoft-validation skill to implement validation for CreateOrder.
> Require an email address and 1–10 items. Each item needs a product ID and a
> positive quantity. Collect independent errors with item indices in their paths.
> Check product availability only after local validation succeeds, preserve
> repository errors, and test that invalid orders are never saved.

Agents that support automatic skill discovery can select it for validation work;
you can also name it explicitly using your agent's supported invocation mechanism.
The skill is intended for changing validation behavior. A task that only formats
or groups existing errors does not need it.

## Evaluation and expectations

In the recorded October 6, 2026 evaluation, DeepSeek V4.1 Flash passed 20 of 21
attempts with the skill and 17 of 21 without it, using fewer tokens overall.
One solution with the skill still omitted a required collection-size check.
These are results for one model and a development suite, not a guarantee for
other models or applications. Review generated code and run application tests.

See the [evaluation report](agent-skill-evaluation.md) for the setup, full
breakdown, remaining failure, and interpretation. Maintainers can use the
[evaluation runbook](../evals/muonsoft-validation/README.md) to run new comparisons.
