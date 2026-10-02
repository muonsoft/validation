package validate_test

import (
	"bytes"
	"errors"
	"image/png"
	"strings"
	"testing"

	"github.com/muonsoft/validation/validate"
)

func TestFileName(t *testing.T) {
	cases := []struct {
		name           string
		basic, windows bool
	}{
		{"", true, true},
		{"photo.jpg", true, true},
		{"отчёт 日本語.pdf", true, true},
		{".env", true, true},
		{"a b.txt", true, true},
		{"report.final.pdf", true, true},
		{"..file", true, true},
		{".", false, false},
		{"..", false, false},
		{"../a", false, false},
		{`..\a`, false, false},
		{"/a", false, false},
		{`C:\a`, false, false},
		{"C:a", false, false},
		{`\\server\a`, false, false},
		{"a\x00b", false, false},
		{"a\nb", false, false},
		{"a\u0085b", false, false},
		{"\xff", false, false},
		{"CON", true, false},
		{"con.txt", true, false},
		{"NuL.tar.gz", true, false},
		{"NUL .txt", true, false},
		{"PRN", true, false},
		{"AUX", true, false},
		{"COM1", true, false},
		{"LPT9", true, false},
		{"COM¹.txt", true, false},
		{"COM²", true, false},
		{"LPT³", true, false},
		{"CONIN$", true, false},
		{"CONOUT$.txt", true, false},
		{"COM0", true, true},
		{"COM10", true, true},
		{"a.", true, false},
		{"a ", true, false},
		{"file:stream", true, false},
		{"a<b", true, false},
		{"a>b", true, false},
		{`a"b`, true, false},
		{"a|b", true, false},
		{"a?b", true, false},
		{"a*b", true, false},
		{"%2e%2e", true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := validate.FileName(tc.name); (err == nil) != tc.basic {
				t.Errorf("basic: %v", err)
			}
			if err := validate.FileName(tc.name, validate.WithWindowsFileNameRestrictions()); (err == nil) != tc.windows {
				t.Errorf("windows: %v", err)
			}
		})
	}
}

func TestFileExtension(t *testing.T) {
	cases := []struct {
		name       string
		extensions []string
		want       error
	}{
		{"", []string{"jpg"}, nil},
		{"photo.JPG", []string{".jpg"}, nil},
		{"photo.jpg", []string{"JPG"}, nil},
		{"a.tar.GZ", []string{"tar.gz"}, nil},
		{"a.tar.gz", []string{"gz"}, nil},
		{"report.final.pdf", []string{"pdf"}, nil},
		{".env.txt", []string{"txt"}, nil},
		{"photo.jpg.exe", []string{"jpg"}, validate.ErrInvalidFileExtension},
		{"photo.jpg ", []string{"jpg"}, validate.ErrInvalidFileExtension},
		{".jpg", []string{"jpg"}, validate.ErrInvalidFileExtension},
		{".tar.gz", []string{"tar.gz"}, validate.ErrInvalidFileExtension},
		{"photo", []string{"jpg"}, validate.ErrInvalidFileExtension},
		{"a.", []string{"jpg"}, validate.ErrInvalidFileExtension},
		{"../a.jpg", []string{"jpg"}, validate.ErrInvalidFileExtension},
		{`a\b.jpg`, []string{"jpg"}, validate.ErrInvalidFileExtension},
		{"a\x00.jpg", []string{"jpg"}, validate.ErrInvalidFileExtension},
		{"a.K", []string{"k"}, validate.ErrInvalidFileExtension},
		{"", nil, validate.ErrInvalidFileExtensions},
		{"a.jpg", []string{"jpg", ""}, validate.ErrInvalidFileExtensions},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := validate.FileExtension(tc.name, tc.extensions...); !errors.Is(err, tc.want) {
				t.Errorf("got %v, want %v", err, tc.want)
			}
		})
	}
	for _, extension := range []string{"", ".", "..jpg", "tar..gz", "jpg.", "*.jpg", "jpg/png", `a\b`, "jpg ", " jpg", "a:b", "é"} {
		if err := validate.FileExtension("", extension); !errors.Is(err, validate.ErrInvalidFileExtensions) {
			t.Errorf("extension %q: %v", extension, err)
		}
	}
}

func TestMIMEType(t *testing.T) {
	cases := []struct {
		value string
		types []string
		want  error
	}{
		{"", nil, nil},
		{"image/png", nil, nil},
		{" Image/PNG ", []string{"IMAGE/png"}, nil},
		{"text/plain; charset=utf-8", []string{"text/plain"}, nil},
		{`text/plain; title="a;b"`, []string{"text/plain"}, nil},
		{"application/vnd.example+json", nil, nil},
		{"text/plain", []string{"image/png"}, validate.ErrInvalidMIMEType},
		{"attachment", nil, validate.ErrInvalidMIMEType},
		{"image/*", nil, validate.ErrInvalidMIMEType},
		{"image/", nil, validate.ErrInvalidMIMEType},
		{"/png", nil, validate.ErrInvalidMIMEType},
		{"text/plain; charset", nil, validate.ErrInvalidMIMEType},
		{"text/plain; charset=utf-8; charset=ascii", nil, validate.ErrInvalidMIMEType},
		{"image/png, image/jpeg", nil, validate.ErrInvalidMIMEType},
		{"", []string{""}, validate.ErrInvalidMIMETypes},
		{"", []string{"image/*"}, validate.ErrInvalidMIMETypes},
		{"", []string{"text/plain; charset=utf-8"}, validate.ErrInvalidMIMETypes},
		{"", []string{"text/plain;"}, validate.ErrInvalidMIMETypes},
	}
	for _, tc := range cases {
		t.Run(tc.value, func(t *testing.T) {
			if err := validate.MIMEType(tc.value, tc.types...); !errors.Is(err, tc.want) {
				t.Errorf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestContentType(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		mime string
	}{
		{"png signature only", []byte("\x89PNG\r\n\x1a\n"), "image/png"},
		{"jpeg signature only", []byte{0xff, 0xd8, 0xff}, "image/jpeg"},
		{"gif87", []byte("GIF87a"), "image/gif"},
		{"gif89", []byte("GIF89a"), "image/gif"},
		{"pdf", []byte("%PDF-1.7\n"), "application/pdf"},
		{"zip", []byte("PK\x03\x04"), "application/zip"},
		{"text", []byte("hello"), "text/plain"},
		{"unknown", []byte{0, 1, 2}, "application/octet-stream"},
		{"truncated png", []byte("\x89PNG"), "text/plain"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			original := bytes.Clone(tc.data)
			if err := validate.ContentType(tc.data, []string{tc.mime}); err != nil {
				t.Fatal(err)
			}
			if err := validate.ContentType(tc.data, []string{"application/x-not-allowed"}); !errors.Is(err, validate.ErrInvalidContentType) {
				t.Fatal(err)
			}
			if !bytes.Equal(original, tc.data) {
				t.Fatal("input changed")
			}
		})
	}
	for _, data := range [][]byte{nil, {}} {
		if err := validate.ContentType(data, []string{"image/png"}, validate.WithContentTypeDetector(func([]byte) string { t.Fatal("called for empty input"); return "" })); err != nil {
			t.Fatal(err)
		}
		if err := validate.ContentType(data, nil); !errors.Is(err, validate.ErrInvalidMIMETypes) {
			t.Fatal(err)
		}
	}
	data := []byte(strings.Repeat("a", 1024))
	var calls int
	detector := func(got []byte) string {
		calls++
		if len(got) != len(data) {
			t.Fatal("truncated data")
		}
		return "Image/PNG; x=y"
	}
	if err := validate.ContentType(data, []string{"image/png"}, validate.WithContentTypeDetector(detector)); err != nil || calls != 1 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
	if err := validate.ContentType(data, []string{"text/plain"}, validate.WithContentTypeDetector(detector), validate.WithContentTypeDetector(nil)); err != nil {
		t.Fatal(err)
	}
	for _, result := range []string{"", "attachment", "image/*", "bad/type; broken"} {
		if err := validate.ContentType(data, []string{"image/png"}, validate.WithContentTypeDetector(func([]byte) string { return result })); !errors.Is(err, validate.ErrInvalidDetectedContentType) {
			t.Fatalf("%q: %v", result, err)
		}
	}
	// A signature after byte 512 must not affect the standard detector.
	data = append(bytes.Repeat([]byte("a"), 512), []byte("\x89PNG\r\n\x1a\n")...)
	if err := validate.ContentType(data, []string{"text/plain"}); err != nil {
		t.Fatal(err)
	}
}

func FuzzFileName(f *testing.F) {
	for _, name := range []string{"file.jpg", "CON.txt", "../a", "\xff"} {
		f.Add(name)
	}
	f.Fuzz(func(t *testing.T, name string) {
		basic := validate.FileName(name)
		windows := validate.FileName(name, validate.WithWindowsFileNameRestrictions())
		if windows == nil && basic != nil {
			t.Fatal("Windows restrictions weakened basic validation")
		}
	})
}

func FuzzMIMEType(f *testing.F) {
	for _, value := range []string{"text/plain", "attachment", "text/plain; charset=utf-8", "\x00"} {
		f.Add(value)
	}
	f.Fuzz(func(t *testing.T, value string) {
		if validate.MIMEType(value) == nil && value != "" && !strings.Contains(value, "/") {
			t.Fatal("accepted non-media type")
		}
	})
}

func TestContentTypeDoesNotValidateImageIntegrity(t *testing.T) {
	signature := []byte("\x89PNG\r\n\x1a\n")
	if err := validate.ContentType(signature, []string{"image/png"}); err != nil {
		t.Fatal(err)
	}
	if _, err := png.Decode(bytes.NewReader(signature)); err == nil {
		t.Fatal("signature alone must not be a complete image")
	}
}
