package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type Vehicle struct {
	Model    string
	MaxSpeed int
}

func (v Vehicle) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(ctx,
		validation.StringProperty("model", v.Model, it.IsNotBlank(), it.HasMaxLength(100)),
		validation.NumberProperty[int]("maxSpeed", v.MaxSpeed, it.IsBetween(50, 200)),
	)
}

type Car struct {
	Vehicle
	PassengerSeats int
}

func (c Car) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(ctx,
		validation.CheckNoViolations(c.Vehicle.Validate(ctx, validator)),
		validation.NumberProperty[int]("passengerSeats", c.PassengerSeats, it.IsBetween(2, 6)),
	)
}

type Truck struct {
	Vehicle
	LoadCapacity float64
}

func (t Truck) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(ctx,
		validation.CheckNoViolations(t.Vehicle.Validate(ctx, validator)),
		validation.NumberProperty[float64]("loadCapacity", t.LoadCapacity, it.IsBetween(10.0, 200.0)),
	)
}

func ExampleCheckNoViolations() {
	vehicles := []validation.Validatable{
		Car{
			Vehicle: Vehicle{
				Model:    "Audi",
				MaxSpeed: 10,
			},
			PassengerSeats: 1,
		},
		Truck{
			Vehicle: Vehicle{
				Model:    "Benz",
				MaxSpeed: 20,
			},
			LoadCapacity: 5,
		},
	}

	for i, vehicle := range vehicles {
		err := validator.ValidateIt(context.Background(), vehicle)
		if violations, ok := validation.UnwrapViolationList(err); ok {
			fmt.Println("vehicle", i, "is not valid:")
			for violation := violations.First(); violation != nil; violation = violation.Next() {
				fmt.Println(violation)
			}
		}
	}

	// Output:
	// vehicle 0 is not valid:
	// violation at "maxSpeed": "This value should be between 50 and 200."
	// violation at "passengerSeats": "This value should be between 2 and 6."
	// vehicle 1 is not valid:
	// violation at "maxSpeed": "This value should be between 50 and 200."
	// violation at "loadCapacity": "This value should be between 10 and 200."
}
