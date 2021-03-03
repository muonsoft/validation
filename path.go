package validation

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type PropertyPathElement interface {
	IsIndex() bool
	fmt.Stringer
}

type PropertyNameElement string

func (p PropertyNameElement) IsIndex() bool {
	return false
}

func (p PropertyNameElement) String() string {
	return string(p)
}

type ArrayIndexElement int

func (a ArrayIndexElement) IsIndex() bool {
	return true
}

func (a ArrayIndexElement) String() string {
	return strconv.Itoa(int(a))
}

type PropertyPath []PropertyPathElement

func (path PropertyPath) Format() string {
	var s strings.Builder

	for i, element := range path {
		if i > 0 && !element.IsIndex() {
			s.WriteString(".")
		}
		if element.IsIndex() {
			s.WriteString("[" + element.String() + "]")
		} else {
			s.WriteString(element.String())
		}
	}

	return s.String()
}

func (path PropertyPath) MarshalJSON() ([]byte, error) {
	return json.Marshal(path.Format())
}
