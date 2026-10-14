package scenario

import "github.com/muonsoft/validation"

var ErrSize = validation.NewError("upload_size", "Image must be between 1 byte and 1 MiB.")

type Attachment struct {
	Name, MIME string
	Size       int64
}
