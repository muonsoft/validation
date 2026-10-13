package scenario_test

import (
	"reflect"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

func goodProduct(sku string) s.Product { return s.Product{SKU: sku, Name: "Product", Quantity: 1} }
func TestProductsAccumulation(t *testing.T) {
	x := s.ProductBatch{Items: []s.Product{{SKU: "A", Price: -1}, goodProduct(""), {SKU: "A", Name: "OK", Quantity: -1}, goodProduct("B")}}
	before := append([]s.Product(nil), x.Items...)
	err := x.Validate(t.Context(), newValidator(t).AtProperty("payload"))
	want(t, err, expected{"payload.items", validation.ErrTooManyElements}, expected{"payload.items[0].sku", validation.ErrNotUnique}, expected{"payload.items[2].sku", validation.ErrNotUnique}, expected{"payload.items[0].name", validation.ErrIsBlank}, expected{"payload.items[0].quantity", s.ErrQuantity}, expected{"payload.items[0].price", s.ErrPrice}, expected{"payload.items[2].quantity", s.ErrQuantity})
	wantPath(t, err, "payload.items[2].sku", validation.PropertyName("payload"), validation.PropertyName("items"), validation.ArrayIndex(2), validation.PropertyName("sku"))
	if !reflect.DeepEqual(before, x.Items) {
		t.Fatal("changed batch")
	}
}
func TestProductsOptionalKeysAndBounds(t *testing.T) {
	want(t, (s.ProductBatch{}).Validate(t.Context(), newValidator(t)), expected{"items", validation.ErrTooFewElements})
	for _, items := range [][]s.Product{{goodProduct("")}, {goodProduct(""), goodProduct("")}, {goodProduct("A"), goodProduct("a"), goodProduct("")}} {
		want(t, (s.ProductBatch{Items: items}).Validate(t.Context(), newValidator(t)))
	}
	want(t, (s.ProductBatch{Items: []s.Product{goodProduct("A"), goodProduct("A"), goodProduct("A")}}).Validate(t.Context(), newValidator(t)), expected{"items[0].sku", validation.ErrNotUnique}, expected{"items[1].sku", validation.ErrNotUnique}, expected{"items[2].sku", validation.ErrNotUnique})
}
func TestProductsOwnRules(t *testing.T) {
	want(t, (s.Product{Price: -0.01}).Validate(t.Context(), newValidator(t).AtProperty("product")), expected{"product.name", validation.ErrIsBlank}, expected{"product.quantity", s.ErrQuantity}, expected{"product.price", s.ErrPrice})
	want(t, goodProduct("").Validate(t.Context(), newValidator(t)))
}
