package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

type Property struct {
	Name       string
	Properties Properties
}

type Properties []Property

func (p Property) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		ctx,
		validation.StringProperty("name", p.Name, it.IsNotBlank()),
		validation.ValidProperty("properties", p.Properties),
	)
}

func (properties Properties) Validate(ctx context.Context, validator *validation.Validator) error {
	violations := validation.ViolationList{}

	for i := range properties {
		err := validator.AtIndex(i).Validate(ctx, validation.Valid(properties[i]))
		err = violations.AppendFromError(err)
		if err != nil {
			return err
		}
	}

	return violations.AsError()
}

func BenchmarkViolationsGeneration(b *testing.B) {
	properties := makeProperties(1000)

	validator, err := validation.NewValidator()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(context.Background(), validation.Valid(properties))
	}
}

func makeProperties(n int) Properties {
	if n <= 0 {
		return nil
	}
	properties := make(Properties, n)
	for i := range properties {
		properties[i].Properties = makeProperties(n / 10)
	}
	return properties
}
