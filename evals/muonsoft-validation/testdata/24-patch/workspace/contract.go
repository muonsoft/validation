package scenario

import "github.com/muonsoft/validation"

var ErrMode = validation.NewError("patch_mode", "Unsupported mode.")
var ErrClear = validation.NewError("patch_clear", "Mode cannot be cleared.")

type ModePatch struct {
	Present bool
	Value   *int
}
type Patch struct {
	Mode ModePatch
	Note *string
}
