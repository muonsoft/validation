package it_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func ExampleIsFileName() {
	v, _ := validation.NewValidator()
	for _, name := range []string{"report.pdf", "CON.txt"} {
		fmt.Println(v.Validate(context.Background(), validation.String(name, it.IsFileName().WithWindowsRestrictions())))
	}
	// Output:
	// <nil>
	// violation: "This value is not a valid file name."
}

func ExampleHasFileExtension() {
	v, _ := validation.NewValidator()
	for _, name := range []string{"backup.TAR.GZ", "backup.tar.gz.exe"} {
		fmt.Println(v.Validate(context.Background(), validation.String(name, it.HasFileExtension("tar.gz"))))
	}
	// Output:
	// <nil>
	// violation: "This file extension is not allowed."
}

func ExampleIsMIMEType() {
	v, _ := validation.NewValidator()
	for _, metadata := range []string{"text/plain; charset=utf-8", "image/png"} {
		fmt.Println(v.Validate(context.Background(), validation.String(metadata, it.IsMIMEType("text/plain"))))
	}
	// Output:
	// <nil>
	// violation: "This MIME type is invalid or not allowed."
}

func ExampleHasContentType() {
	v, _ := validation.NewValidator()
	for _, data := range [][]byte{[]byte("hello"), {0, 1, 2}} {
		fmt.Println(v.Validate(context.Background(), validation.This(data, it.HasContentType("text/plain"))))
	}
	// Output:
	// <nil>
	// violation: "This content type is not allowed."
}

func ExampleContentTypeConstraint_WithDetector() {
	v, _ := validation.NewValidator()
	rule := it.HasContentType("application/x-example").WithDetector(func(data []byte) string {
		if len(data) >= 4 && string(data[:4]) == "EXMP" {
			return "application/x-example"
		}
		return "application/octet-stream"
	})
	fmt.Println(v.Validate(context.Background(), validation.This([]byte("EXMP payload"), rule)))
	// Output:
	// <nil>
}
