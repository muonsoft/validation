package validation

type Constraint interface {
	Option
	GetCode() string
	GetMessageTemplate() string
	GetParameters() map[string]string
}

type NilConstraint interface {
	ValidateNil(options Options) error
}

type StringConstraint interface {
	ValidateString(value *string, options Options) error
}

type IntConstraint interface {
	ValidateInt(value *int64, options Options) error
}

type UintConstraint interface {
	ValidateUint(value *uint64, options Options) error
}

type FloatConstraint interface {
	ValidateFloat(value *float64, options Options) error
}
