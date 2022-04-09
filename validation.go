// Copyright 2021 Igor Lazarev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package validation provides tools for data validation.
// It is designed to create complex validation rules with abilities to hook into the validation process.
package validation

import (
	"context"
)

// Validatable is interface for creating validatable types on the client side.
// By using it you can build complex validation rules on a set of objects used in other objects.
//
// Example
//  type Book struct {
//      Title    string
//      Author   string
//      Keywords []string
//  }
//
//  func (b Book) Validate(ctx context.Context, validator *validation.Validator) error {
//      return validator.Validate(
//          ctx,
//          validation.StringProperty("title", &b.Title, it.IsNotBlank()),
//          validation.StringProperty("author", &b.Author, it.IsNotBlank()),
//          validation.CountableProperty("keywords", len(b.Keywords), it.HasCountBetween(1, 10)),
//          validation.EachStringProperty("keywords", b.Keywords, it.IsNotBlank()),
//      )
//  }
type Validatable interface {
	Validate(ctx context.Context, validator *Validator) error
}

// Filter is used for processing the list of errors to return a single ViolationList.
// If there is at least one non-violation error it will return it instead.
func Filter(violations ...error) error {
	list := &ViolationList{}

	for _, violation := range violations {
		err := list.AppendFromError(violation)
		if err != nil {
			return err
		}
	}

	return list.AsError()
}

// ValidateByConstraintFunc is used for building validation functions for the values of specific types.
type ValidateByConstraintFunc func(constraint Constraint, scope Scope) error

type validateFunc func(scope Scope) (*ViolationList, error)

func newValidator(options []Option, validate ValidateByConstraintFunc) validateFunc {
	return func(scope Scope) (*ViolationList, error) {
		err := scope.applyOptions(options...)
		if err != nil {
			return nil, err
		}

		return validateOnScope(scope, options, validate)
	}
}

func validateOnScope(scope Scope, options []Option, validate ValidateByConstraintFunc) (*ViolationList, error) {
	violations := &ViolationList{}

	for _, option := range options {
		if constraint, ok := option.(controlConstraint); ok {
			vs, err := constraint.validate(scope, validate)
			if err != nil {
				return nil, err
			}

			violations.Join(vs)
		} else if constraint, ok := option.(Constraint); ok {
			err := violations.AppendFromError(validate(constraint, scope))
			if err != nil {
				return nil, err
			}
		}
	}

	return violations, nil
}
