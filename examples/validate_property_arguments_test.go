package examples

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func ValidatePropertyArguments(p Product) error {
	return validation.Validate(
		validation.StringProperty("name", &p.Name, it.IsNotBlank()),
		validation.IterableProperty("tags", p.Tags, it.HasMinCount(1)),
		validation.IterableProperty("components", p.Components, it.HasMinCount(1)),
	)
}

func ExampleValidate_propertyArguments() {
	p := Product{
		Name: "",
		Components: []Component{
			{
				ID:   1,
				Name: "",
			},
		},
	}

	err := ValidatePropertyArguments(p)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation at 'name': This value should not be blank.
	// violation at 'tags': This collection should contain 1 element or more.
	// violation at 'components[0].name': This value should not be blank.
	// violation at 'components[0].tags': This collection should contain 1 element or more.
}
