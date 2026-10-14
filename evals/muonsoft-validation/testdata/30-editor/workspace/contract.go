package scenario

import "github.com/muonsoft/validation"

var ErrTitleRequired = validation.NewError("title_required", "Title is required.")
var ErrTitleLength = validation.NewError("title_length", "Title must have 3 to 12 characters.")
var ErrContactRequired = validation.NewError("contact_required", "Contact is required for publishing.")
var ErrContactFormat = validation.NewError("contact_format", "Contact must be an email address.")

type Article struct{ Title, Contact string }
