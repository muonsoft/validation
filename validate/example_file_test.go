package validate_test

import (
	"fmt"

	"github.com/muonsoft/validation/validate"
)

func ExampleFileName() {
	fmt.Println(validate.FileName("отчёт.pdf"))
	fmt.Println(validate.FileName("NUL.txt", validate.WithWindowsFileNameRestrictions()))
	// Output:
	// <nil>
	// invalid file name
}

func ExampleFileExtension() {
	fmt.Println(validate.FileExtension("backup.tar.gz", "tar.gz"))
	fmt.Println(validate.FileExtension("photo.jpg.exe", "jpg"))
	// Output:
	// <nil>
	// invalid file extension
}

func ExampleMIMEType() {
	fmt.Println(validate.MIMEType("text/plain; charset=utf-8", "text/plain"))
	fmt.Println(validate.MIMEType("attachment"))
	// Output:
	// <nil>
	// invalid MIME type
}

func ExampleContentType() {
	fmt.Println(validate.ContentType([]byte("hello"), []string{"text/plain"}))
	fmt.Println(validate.ContentType([]byte{0, 1, 2}, []string{"text/plain"}))
	// Output:
	// <nil>
	// invalid content type
}
