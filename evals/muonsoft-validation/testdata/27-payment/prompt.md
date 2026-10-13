# Validate payment details

Implement Payment.Validate. Card is required and must match Visa or Mastercard
prefix/length rules AND have a valid Luhn checksum. Do not remove spaces or
hyphens. IBAN and BIC are optional; validate a supplied IBAN's country-specific
format/checksum and a supplied BIC's default strict uppercase format. No BIC/IBAN
country relationship is required. Report independent failures together at card,
iban, bic with the corresponding library errors. Scheme and checksum are separate
rules and can both fail for the same card; empty card yields only ErrIsBlank.
Do not make network requests or infer that a matching card/account exists.

Preserve contract.go and dependencies. Edit validation.go; private helpers in new
top-level non-test Go files are allowed. Remove temporary tests before submission.
Use existing library/project rules when their semantics match. Do not mutate input.
Preserve the supplied validator and its path prefix.
