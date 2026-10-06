package scenario

import "github.com/muonsoft/validation"

type Form struct{ Fields []Field }
type Field struct {
	Name, Kind, Value string
	Children          []Field
	Attributes        map[string]string
}

var ErrDepth = validation.NewError("depth", "Maximum field depth exceeded.")
