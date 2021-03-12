// Copyright 2021 Igor Lazarev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package validation provides tools for data validation.
// It is designed to create complex validation rules with abilities to hook into the validation process.
package validation

import (
	"reflect"
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
//  func (b Book) Validate(validator *validation.Validator) error {
//      return validator.Validate(
//          validation.StringProperty("title", &b.Title, it.IsNotBlank()),
//          validation.StringProperty("author", &b.Author, it.IsNotBlank()),
//          validation.CountableProperty("keywords", len(b.Keywords), it.HasCountBetween(1, 10)),
//          validation.EachStringProperty("keywords", b.Keywords, it.IsNotBlank()),
//      )
//  }
type Validatable interface {
	Validate(validator *Validator) error
}

// Filter is used for processing the list of errors to return a single ViolationList.
// If there is at least one non-violation error it will return it instead.
func Filter(violations ...error) error {
	filteredViolations := make(ViolationList, 0, len(violations))

	for _, err := range violations {
		addErr := filteredViolations.AppendFromError(err)
		if addErr != nil {
			return addErr
		}
	}

	return filteredViolations.AsError()
}

var validatableType = reflect.TypeOf((*Validatable)(nil)).Elem()
