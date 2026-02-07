# Creating Custom Constraints

Everything you need to create a custom constraint is to implement one of the interfaces:

* `BoolConstraint` - for validating boolean values;
* `NumberConstraint` - for validating numeric values;
* `StringConstraint` - for validating string values;
* `ComparableConstraint` - for validating generic comparable values;
* `ComparablesConstraint` - for validating slice of generic comparable values;
* `CountableConstraint` - for validating iterable values based only on the count of elements;
* `TimeConstraint` - for validating date/time values.

Also, you can combine several types of constraints. See examples for more details:

* [custom static constraint](https://pkg.go.dev/github.com/muonsoft/validation#example-Validator.Validate-CustomConstraint);
* [custom constraint as a service](https://pkg.go.dev/github.com/muonsoft/validation#example-Validator.GetConstraint-CustomServiceConstraint);
* [custom constraint with custom argument for domain type](https://pkg.go.dev/github.com/muonsoft/validation#example-NewArgument-CustomArgumentConstraintValidator);
* [Func as Constraint and Each/EachProperty for slices of any type](https://pkg.go.dev/github.com/muonsoft/validation#example-Func).
