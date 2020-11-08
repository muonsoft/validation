package validation

import (
	"fmt"
	"strconv"
	"strings"
)

type PropertyPathElement interface {
	IsIndex() bool
	fmt.Stringer
}

type PropertyNameElement struct {
	Value string
}

func (p PropertyNameElement) IsIndex() bool {
	return false
}

func (p PropertyNameElement) String() string {
	return p.Value
}

type ArrayIndexElement struct {
	Value int
}

func (a ArrayIndexElement) IsIndex() bool {
	return true
}

func (a ArrayIndexElement) String() string {
	return strconv.Itoa(a.Value)
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
