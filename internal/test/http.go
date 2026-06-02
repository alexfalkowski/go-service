package test

import (
	"fmt"
	"io/fs"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
)

// MessageMediaType describes an HTTP content type and its corresponding encoder kind.
type MessageMediaType struct {
	Name        string
	ContentType string
	Kind        string
}

// MessageMediaTypes returns the message media types shared by HTTP transport tests.
func MessageMediaTypes() []MessageMediaType {
	return []MessageMediaType{
		{Name: "json", ContentType: media.JSON, Kind: "json"},
		{Name: "hjson", ContentType: media.HumanJSON, Kind: "hjson"},
		{Name: "yaml", ContentType: media.YAML, Kind: "yaml"},
		{Name: "yml", ContentType: "application/yml", Kind: "yml"},
		{Name: "toml", ContentType: media.TOML, Kind: "toml"},
	}
}

// ErrResponseWriter is an [http.ResponseWriter] test double whose writes fail with ErrFailed.
type ErrResponseWriter struct {
	Code int
}

// Header is always empty.
func (w *ErrResponseWriter) Header() http.Header {
	return http.Header{}
}

// Write returns ErrFailed.
func (w *ErrResponseWriter) Write([]byte) (int, error) {
	return 0, ErrFailed
}

// WriteHeader stores code in the Code field.
func (w *ErrResponseWriter) WriteHeader(code int) {
	w.Code = code
}

// RoundTripperFunc adapts a function to [http.RoundTripper].
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip calls f(req).
func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// ErrorRoundTripper is an [http.RoundTripper] test double that always returns Err.
type ErrorRoundTripper struct {
	Err   error
	Calls int
}

// RoundTrip records the call and returns Err.
func (r *ErrorRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	r.Calls++
	return nil, r.Err
}

// StatusRoundTripper is an [http.RoundTripper] test double that returns Status.
type StatusRoundTripper struct {
	Status int
	Calls  int
}

// RoundTrip records the call and returns a response with Status.
func (r *StatusRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	r.Calls++
	return ResponseWithStatus(r.Status), nil
}

// StatusSequenceRoundTripper returns one response status for each call.
type StatusSequenceRoundTripper struct {
	Codes []int
	Calls int
}

// RoundTrip records the call and returns the next status response.
func (r *StatusSequenceRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	code := r.Codes[r.Calls]
	r.Calls++

	return ResponseWithStatus(code), nil
}

// BodySequenceRoundTripper returns service unavailable responses with configured bodies.
type BodySequenceRoundTripper struct {
	Responses []string
	Calls     int
}

// RoundTrip records the call and returns the next response body.
func (r *BodySequenceRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	body := r.Responses[r.Calls]
	r.Calls++

	return &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Status:     StatusLine(http.StatusServiceUnavailable),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}, nil
}

// RequestBodyRecorderRoundTripper records request bodies and returns Status.
type RequestBodyRecorderRoundTripper struct {
	Bodies []string
	Status int
}

// RoundTrip records the request body and returns the configured status response.
func (r *RequestBodyRecorderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	body, _, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	r.Bodies = append(r.Bodies, string(body))

	return ResponseWithStatus(defaultStatus(r.Status, http.StatusServiceUnavailable)), nil
}

// TransportErrorThenSuccessRoundTripper fails once, then returns Status.
type TransportErrorThenSuccessRoundTripper struct {
	Err    error
	Bodies []string
	Status int
	Calls  int
}

// RoundTrip records the request body, returns Err on the first call, then succeeds.
func (r *TransportErrorThenSuccessRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	body, _, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	r.Calls++
	r.Bodies = append(r.Bodies, string(body))
	if r.Calls == 1 {
		return nil, defaultError(r.Err, io.ErrUnexpectedEOF)
	}

	return ResponseWithStatus(defaultStatus(r.Status, http.StatusOK)), nil
}

// OriginalBodyRoundTripper records whether the first request used Original.
type OriginalBodyRoundTripper struct {
	Original          io.ReadCloser
	Bodies            []string
	Calls             int
	FirstUsedOriginal bool
}

// RoundTrip records the request body and closes it.
func (r *OriginalBodyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.Calls == 0 {
		r.FirstUsedOriginal = req.Body == r.Original
	}

	body, _, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	if err := req.Body.Close(); err != nil {
		return nil, err
	}

	r.Calls++
	r.Bodies = append(r.Bodies, string(body))

	return ResponseWithStatus(http.StatusServiceUnavailable), nil
}

// AuthRoundTripper records Authorization headers and returns configured status codes.
type AuthRoundTripper struct {
	AuthValues []string
	AuthCounts []int
	Codes      []int
	Calls      int
}

// RoundTrip records the Authorization headers and returns the next status response.
func (r *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.AuthValues = append(r.AuthValues, req.Header.Get("Authorization"))
	r.AuthCounts = append(r.AuthCounts, len(req.Header.Values("Authorization")))
	code := r.Codes[r.Calls]
	r.Calls++

	return ResponseWithStatus(code), nil
}

// CauseRoundTripper records the request context error and cause.
type CauseRoundTripper struct {
	Cause error
	Err   error
	Wait  bool
}

// RoundTrip optionally waits for cancellation, then records context state.
func (r *CauseRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.Wait {
		<-req.Context().Done()
	}

	r.Cause = context.Cause(req.Context())
	r.Err = req.Context().Err()

	return ResponseWithStatus(http.StatusOK), nil
}

// NonReplayableReader is an [io.Reader] whose contents can only be read once.
type NonReplayableReader struct {
	Value string
	read  bool
}

// Read copies Value on the first read and EOF afterwards.
func (r *NonReplayableReader) Read(p []byte) (int, error) {
	if r.read {
		return 0, io.EOF
	}

	r.read = true
	copy(p, r.Value)
	return len(r.Value), io.EOF
}

// TrackedBody is an [io.ReadCloser] that records Close calls.
type TrackedBody struct {
	*strings.Reader
	Closed bool
}

// Close records that the body was closed.
func (b *TrackedBody) Close() error {
	b.Closed = true
	return nil
}

// UnknownLengthReader hides the concrete reader length from HTTP request construction.
type UnknownLengthReader struct {
	*strings.Reader
}

// HeaderDeletingRoundTripper deletes Header before calling RoundTripper.
type HeaderDeletingRoundTripper struct {
	RoundTripper http.RoundTripper
	Header       string
}

// RoundTrip deletes Header from req and calls the wrapped RoundTripper.
func (r *HeaderDeletingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Del(r.Header)

	return r.RoundTripper.RoundTrip(req)
}

// ErrFileSystem is an [fs.FS] test double whose file read returns a partial payload and ErrFailed.
type ErrFileSystem struct{}

// Open returns a failing file for asset.txt and [fs.ErrNotExist] for other names.
func (ErrFileSystem) Open(name string) (fs.File, error) {
	if name != "asset.txt" {
		return nil, fs.ErrNotExist
	}

	return &ErrFile{}, nil
}

// ErrFile is an [fs.File] test double that returns a partial payload and ErrFailed.
type ErrFile struct {
	read bool
}

// Stat returns ErrFileInfo.
func (f *ErrFile) Stat() (fs.FileInfo, error) {
	return ErrFileInfo{}, nil
}

// Read copies a partial payload and returns ErrFailed on the first read.
func (f *ErrFile) Read(p []byte) (int, error) {
	if f.read {
		return 0, io.EOF
	}

	f.read = true
	copy(p, "hello")

	return len("hello"), ErrFailed
}

// Close implements [fs.File] and always succeeds.
func (f *ErrFile) Close() error {
	return nil
}

// ErrFileInfo is an [fs.FileInfo] test double for ErrFile.
type ErrFileInfo struct{}

// Name returns asset.txt.
func (ErrFileInfo) Name() string {
	return "asset.txt"
}

// Size returns a size larger than the readable payload.
func (ErrFileInfo) Size() int64 {
	return 6
}

// Mode returns no file mode bits.
func (ErrFileInfo) Mode() fs.FileMode {
	return 0
}

// ModTime returns the zero time.
func (ErrFileInfo) ModTime() time.Time {
	return time.Time{}
}

// IsDir reports false.
func (ErrFileInfo) IsDir() bool {
	return false
}

// Sys returns nil.
func (ErrFileInfo) Sys() any {
	return nil
}

// ResponseWithStatus returns a response for code with an empty body.
func ResponseWithStatus(code int) *http.Response {
	return &http.Response{
		StatusCode: code,
		Status:     StatusLine(code),
		Body:       http.NoBody,
		Header:     make(http.Header),
	}
}

// StatusLine formats an HTTP status line fragment for code.
func StatusLine(code int) string {
	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func defaultStatus(status, fallback int) int {
	if status == 0 {
		return fallback
	}

	return status
}

func defaultError(err, fallback error) error {
	if err == nil {
		return fallback
	}

	return err
}
