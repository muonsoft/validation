package scenario_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

type observedReader struct {
	source               io.Reader
	bytes, calls, closes int
}

func (r *observedReader) Read(p []byte) (int, error) {
	r.calls++
	n, e := r.source.Read(p)
	r.bytes += n
	return n, e
}
func (r *observedReader) Close() error { r.closes++; return nil }

type brokenReader struct{ err error }

func (r brokenReader) Read([]byte) (int, error) { return 0, r.err }
func TestUploadMetadata(t *testing.T) {
	for _, tc := range []struct {
		x        s.Upload
		failures []expected
	}{
		{s.Upload{}, []expected{{"name", validation.ErrIsBlank}, {"mime", validation.ErrIsBlank}}},
		{s.Upload{Name: "CON.txt", MIME: "image/png"}, []expected{{"name", validation.ErrInvalidFileName}, {"mime", validation.ErrInvalidMIMEType}}},
		{s.Upload{Name: "../notes.txt", MIME: "text/plain"}, []expected{{"name", validation.ErrInvalidFileName}, {"name", validation.ErrInvalidFileExtension}}},
		{s.Upload{Name: "notes.txt.exe", MIME: "text/plain; broken"}, []expected{{"name", validation.ErrInvalidFileExtension}, {"mime", validation.ErrInvalidMIMEType}}},
	} {
		body := &observedReader{source: strings.NewReader("hello")}
		want(t, tc.x.Validate(t.Context(), newValidator(t), body), tc.failures...)
		if body.calls != 0 || body.closes != 0 {
			t.Fatal("read or closed body before metadata passed")
		}
	}
}
func TestUploadContent(t *testing.T) {
	x := s.Upload{Name: "notes.TXT", MIME: "TEXT/PLAIN; charset=utf-8"}
	want(t, x.Validate(t.Context(), newValidator(t), strings.NewReader("Hello world")))
	want(t, x.Validate(t.Context(), newValidator(t).AtProperty("upload"), strings.NewReader("")), expected{"upload.content", validation.ErrIsBlank})
	want(t, x.Validate(t.Context(), newValidator(t).AtProperty("upload"), bytes.NewReader([]byte{0x89, 'P', 'N', 'G', 13, 10, 26, 10})), expected{"upload.content", validation.ErrInvalidContentType})
}
func TestUploadReadBoundary(t *testing.T) {
	x := s.Upload{Name: "notes.txt", MIME: "text/plain"}
	body := &observedReader{source: strings.NewReader(strings.Repeat("a", 2048))}
	want(t, x.Validate(t.Context(), newValidator(t), body))
	if body.bytes != 512 || body.closes != 0 {
		t.Fatalf("read %d bytes; closes %d", body.bytes, body.closes)
	}
	failure := errors.New("read unavailable")
	err := x.Validate(t.Context(), newValidator(t), brokenReader{failure})
	if !errors.Is(err, failure) {
		t.Fatalf("lost read error: %v", err)
	}
	if _, ok := validation.UnwrapViolations(err); ok {
		t.Fatal("I/O error became input violation")
	}
}
