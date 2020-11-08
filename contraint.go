package validation

type Constraint interface {
	Option
	Code() string
}

type StringConstraint interface {
	ValidateString(value *string, options Options) error
}

type IntConstraint interface {
	ValidateInt(value *int, options Options) error
}

type FloatConstraint interface {
	ValidateFloat(value *float64, options Options) error
}
