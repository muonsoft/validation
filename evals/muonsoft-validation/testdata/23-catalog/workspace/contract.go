package scenario

import "github.com/muonsoft/validation"

var ErrDepth = validation.NewError("catalog_depth", "Catalog is too deep.")
var ErrMissing = validation.NewError("catalog_missing", "Node is missing.")

type Node struct {
	Label    string
	Children map[string]*Node
}
type Catalog struct{ Roots map[string]*Node }
