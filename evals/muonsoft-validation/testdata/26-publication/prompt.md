# Validate publication settings

Implement Publication.Validate. All fields are optional; empty strings are valid.
Country must be a two-letter country/territory recognized by the installed x/text
classification (aliases are allowed); Language must be a recognized standalone
language code, not a locale. Locale may include script/region/extensions. Currency
must be a recognized ISO 4217 code. Week must be an exact YYYY-Www ISO calendar
week between 2020-W53 and 2021-W02 inclusive, including calendar week-53 validity.
Summary must contain 2–4 words when nonempty: Unicode letters/numbers with
combining marks form words, internal straight/curly apostrophes stay within a
word, and hyphens separate words; punctuation alone has zero words.
Accumulate independent failures at country, language, locale, currency, week,
and summary. Use the library's corresponding standard errors, including distinct
week format/number/min/max and word-count min/max errors. W00 is a format
error; a syntactically valid W53 in a year without week 53 is a week-number error.

Preserve contract.go and dependencies. Edit validation.go; private helpers in new
top-level non-test Go files are allowed. Remove temporary tests before submission.
Use existing library/project rules when their semantics match. Do not mutate input.
Preserve the supplied validator and its path prefix.
