---
name: golang-code-review-comments
description: Applies Go community style from the official Code Review Comments wiki. Use when naming APIs, structuring code, or reviewing Go changes in this repository.
---

# Go Code Review Comments

Primary reference: **[Code Review Comments](https://go.dev/wiki/CodeReviewComments)** on the Go wiki.

When writing or reviewing Go code in this project, align with that document. Below are points that come up often for **public API naming** and day-to-day edits.

---

## Initialisms

Words that are acronyms or initialisms (URL, HTTP, ID, CIDR, JSON, etc.) use **consistent casing** in identifiers:

- Good: `ServeHTTP`, `parseURL`, `userID`, `IsCIDR`, `ErrInvalidCIDR`, `CIDROptions`
- Bad: `ServeHttp`, `parseUrl`, `userId`, `IsCidr`, `CidrOptions`

For multi-part names, either keep each initialism consistent (e.g. `xmlHTTPRequest`) or use the wiki’s patterns; **never** mix `Url` / `URL` or `Cidr` / `CIDR` in exported symbols.

When adding constraints or helpers named after a standard (CIDR, UUID, ISIN, …), check the acronym’s usual spelling in Go’s standard library and popular code (e.g. `net.ParseCIDR` — the concept is CIDR; exported wrappers here use `CIDR` in type and function names).

---

## Other high-signal items from the wiki

Skim the full wiki for the full list. In practice, also watch for:

- **Package comments** — `// Package foo ...` on the package clause file.
- **Doc comments** — complete sentences, start with the name being documented.
- **Error strings** — lowercase, no trailing punctuation (for wrapped errors / `errors.New` text where applicable).
- **Line length** — no hard limit in the wiki, but keep readable; wrap long signatures or struct literals.
- **Indent error paths** — handle errors first, reduce nesting (`if err != nil { return err }`).

---

## Relation to this repo

- **`golangci-lint`** (see `AGENTS.md`) enforces some overlapping rules (e.g. `predeclared`, `cyclop`). The wiki is the guide for **naming and idioms** that linters may not fully cover.
- When adding features documented in **`validation-add-constraint`**, apply **Initialisms** for any new exported names in `it`, `validate`, `is`, `message`, and root `validation` errors.

---

## Checklist (quick)

- [ ] Exported names use correct initialism casing (CIDR, URL, ID, …).
- [ ] Godoc on new exported symbols follows comment conventions.
- [ ] Error handling follows early-return style where it improves clarity.
- [ ] Full wiki reviewed for non-trivial or controversial changes.
