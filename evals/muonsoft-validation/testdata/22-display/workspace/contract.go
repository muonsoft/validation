package scenario

type Issue struct{ Path, Message string }
type Group struct {
	Path     string
	Messages []string
}

func ExistingValidate(name string) []Issue {
	if name == "" {
		return []Issue{{Path: "name", Message: "Required"}}
	}
	return nil
}
