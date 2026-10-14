# Validate a product batch with optional stock codes

Implement Product.Validate and ProductBatch.Validate for an import endpoint.
A batch must contain 1–3 products inclusive. Every product requires a nonempty
Name, Quantity greater than zero, and Price at least zero. Missing Name uses the
standard blank-value error; quantity and price use ErrQuantity and ErrPrice with
their messages. These rules must also work when Product.Validate is called alone.
SKU is optional. All nonempty SKUs must be unique within the batch, using exact
case-sensitive comparison; repeated empty SKUs are allowed. A duplicate must
produce the standard uniqueness error at the sku field of EVERY occurrence,
including the first. Collection count uses standard count-bound errors at items.
Accumulate count, duplicates, and element failures independently, including when
the batch is too large. Paths use items, actual array indices, name, quantity,
price, sku under the caller's prefix. Never reorder items or change their values.

Preserve contract.go, existing protected files, and dependencies. Edit validation.go, product.go.
Private helpers in new top-level non-test Go files are allowed. Remove temporary
tests before submission. Use existing library/project rules when their semantics
match. Do not change input, normalize values, or create a replacement validator.
