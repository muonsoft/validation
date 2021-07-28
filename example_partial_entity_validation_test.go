package validation_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type File struct {
	Name string
	Data []byte
}

// This validation will always check that file is valid.
// Partial validation will be applied by AllowedFileExtensionConstraint
// and AllowedFileSizeConstraint.
func (f File) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		ctx,
		validation.StringProperty("name", &f.Name, it.HasLengthBetween(5, 50)),
	)
}

type FileUploadRequest struct {
	Section string
	File    *File
}

type FileConstraint interface {
	validation.Constraint
	ValidateFile(file *File, scope validation.Scope) error
}

func FileArgument(file *File, options ...validation.Option) validation.Argument {
	return validation.NewArgument(options, func(constraint validation.Constraint, scope validation.Scope) error {
		if c, ok := constraint.(FileConstraint); ok {
			return c.ValidateFile(file, scope)
		}
		// If you want to use built-in constraints for checking for nil or empty values
		// such as it.IsNil() or it.IsBlank().
		if c, ok := constraint.(validation.NilConstraint); ok {
			if file == nil {
				return c.ValidateNil(scope)
			}
			return nil
		}

		return validation.NewInapplicableConstraintError(constraint, "File")
	})
}

// AllowedFileExtensionConstraint used to check that file has one of allowed extensions.
// This constraint can be used for partial validation.
type AllowedFileExtensionConstraint struct {
	extensions []string
}

func FileHasAllowedExtension(extensions ...string) AllowedFileExtensionConstraint {
	return AllowedFileExtensionConstraint{extensions: extensions}
}

func (c AllowedFileExtensionConstraint) SetUp() error {
	return nil
}

func (c AllowedFileExtensionConstraint) Name() string {
	return "AllowedFileExtensionConstraint"
}

func (c AllowedFileExtensionConstraint) ValidateFile(file *File, scope validation.Scope) error {
	if file == nil {
		return nil
	}

	extension := strings.ReplaceAll(filepath.Ext(file.Name), ".", "")

	return scope.Validator().AtProperty("name").ValidateString(
		scope.Context(),
		&extension,
		it.IsOneOfStrings(c.extensions...).Message("Not allowed extension. Must be one of: {{ choices }}."),
	)
}

// AllowedFileSizeConstraint used to check that file has limited size.
// This constraint can be used for partial validation.
type AllowedFileSizeConstraint struct {
	minSize int64
	maxSize int64
}

func FileHasAllowedSize(min, max int64) AllowedFileSizeConstraint {
	return AllowedFileSizeConstraint{minSize: min, maxSize: max}
}

func (c AllowedFileSizeConstraint) SetUp() error {
	return nil
}

func (c AllowedFileSizeConstraint) Name() string {
	return "AllowedFileSizeConstraint"
}

func (c AllowedFileSizeConstraint) ValidateFile(file *File, scope validation.Scope) error {
	if file == nil {
		return nil
	}

	size := len(file.Data)

	return scope.Validator().ValidateNumber(
		scope.Context(),
		size,
		it.IsGreaterThanInteger(c.minSize).Message("File size is too small."),
		it.IsLessThanInteger(c.maxSize).Message("File size is too large."),
	)
}

func ExampleScope_Validator() {
	// this constraints will be applied to all files uploaded as avatars
	avatarConstraints := []validation.Option{
		FileHasAllowedExtension("jpeg", "jpg", "gif"),
		FileHasAllowedSize(100, 1000),
	}
	// this constraints will be applied to all files uploaded as documents
	documentConstraints := []validation.Option{
		FileHasAllowedExtension("doc", "pdf", "txt"),
		FileHasAllowedSize(1000, 100000),
	}

	requests := []FileUploadRequest{
		{
			Section: "avatars",
			File:    &File{Name: "avatar.png", Data: bytes.Repeat([]byte{0}, 99)},
		},
		{
			Section: "documents",
			File:    &File{Name: "sheet.xls", Data: bytes.Repeat([]byte{0}, 100001)},
		},
	}

	for _, request := range requests {
		switch request.Section {
		case "avatars":
			err := validator.Validate(
				context.Background(),
				// common validation of validatable
				validation.Valid(request.File),
				// specific validation for file storage section
				FileArgument(request.File, avatarConstraints...),
			)
			fmt.Println(err)
		case "documents":
			err := validator.Validate(
				context.Background(),
				// common validation of validatable
				validation.Valid(request.File),
				// specific validation for file storage section
				FileArgument(request.File, documentConstraints...),
			)
			fmt.Println(err)
		}
	}

	// Output:
	// violation at 'name': Not allowed extension. Must be one of: jpeg, jpg, gif.; violation: File size is too small.
	// violation at 'name': Not allowed extension. Must be one of: doc, pdf, txt.; violation: File size is too large.
}
