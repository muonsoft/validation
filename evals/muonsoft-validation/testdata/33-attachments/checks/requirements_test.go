package scenario_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

var png = []byte{0x89, 'P', 'N', 'G', 13, 10, 26, 10}
var jpeg = []byte{0xff, 0xd8, 0xff, 0xe0}

type trackedReader struct {
	source       io.Reader
	read, closed int
}

func (r *trackedReader) Read(p []byte) (int, error) {
	n, e := r.source.Read(p)
	r.read += n
	return n, e
}
func (r *trackedReader) Close() error { r.closed++; return nil }

type partialFailure struct{ err error }

func (r partialFailure) Read(p []byte) (int, error) { return copy(p, png), r.err }
func TestAttachmentsMetadata(t *testing.T) {
	for _, tc := range []struct {
		x    s.Attachment
		want []expected
	}{
		{s.Attachment{Name: "CON.png", MIME: "text/plain", Size: 0}, []expected{{"name", validation.ErrInvalidFileName}, {"mime", validation.ErrInvalidMIMEType}, {"size", s.ErrSize}}},
		{s.Attachment{Name: "x.png.exe", MIME: "image/png", Size: 1048577}, []expected{{"name", validation.ErrInvalidFileExtension}, {"size", s.ErrSize}}},
		{s.Attachment{Name: "x.png", MIME: "image/jpeg", Size: 1}, []expected{{"mime", validation.ErrInvalidMIMEType}}},
		{s.Attachment{}, []expected{{"name", validation.ErrIsBlank}, {"mime", validation.ErrIsBlank}, {"size", s.ErrSize}}},
	} {
		body := &trackedReader{source: partialFailure{errors.New("must not read")}}
		want(t, tc.x.Validate(t.Context(), newValidator(t), body), tc.want...)
		if body.read != 0 || body.closed != 0 {
			t.Fatal("touched body before metadata passed")
		}
	}
}
func TestAttachmentsPairs(t *testing.T) {
	for _, tc := range []struct {
		name, mime string
		good, bad  []byte
	}{{"x.PNG", "IMAGE/PNG; charset=binary", png, jpeg}, {"x.jpg", "image/jpeg", jpeg, png}, {"x.JPEG", "image/jpeg", jpeg, png}} {
		for _, size := range []int64{1, 1048576} {
			x := s.Attachment{Name: tc.name, MIME: tc.mime, Size: size}
			want(t, x.Validate(t.Context(), newValidator(t), bytes.NewReader(tc.good)))
			want(t, x.Validate(t.Context(), newValidator(t).AtProperty("attachment"), bytes.NewReader(tc.bad)), expected{"attachment.content", validation.ErrInvalidContentType})
		}
	}
}
func TestAttachmentsReadContract(t *testing.T) {
	x := s.Attachment{Name: "x.png", MIME: "image/png", Size: 1000}
	data := append(append([]byte(nil), png...), make([]byte, 2048)...)
	body := &trackedReader{source: bytes.NewReader(data)}
	want(t, x.Validate(t.Context(), newValidator(t), body))
	if body.read != 512 || body.closed != 0 {
		t.Fatalf("read=%d closed=%d", body.read, body.closed)
	}
	want(t, x.Validate(t.Context(), newValidator(t), bytes.NewReader(nil)), expected{"content", validation.ErrIsBlank})
	failure := errors.New("read failed")
	err := x.Validate(t.Context(), newValidator(t), partialFailure{failure})
	if !errors.Is(err, failure) {
		t.Fatalf("lost error: %v", err)
	}
	if _, ok := validation.UnwrapViolations(err); ok {
		t.Fatal("technical failure became violation")
	}
}
