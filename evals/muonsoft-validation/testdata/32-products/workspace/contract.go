package scenario

import "github.com/muonsoft/validation"

var ErrQuantity = validation.NewError("quantity", "Quantity must be positive.")
var ErrPrice = validation.NewError("price", "Price cannot be negative.")

type Product struct {
	SKU, Name string
	Quantity  int
	Price     float64
}
type ProductBatch struct{ Items []Product }
