package test

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
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
		{Name: "gob", ContentType: "application/gob", Kind: "gob"},
		{Name: "msgpack", ContentType: media.MessagePack, Kind: "msgpack"},
	}
}

// ErrResponseWriter is an http.ResponseWriter test double whose writes fail with ErrFailed.
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
