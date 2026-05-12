package content_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/mime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestNewRequestHandlerRejectsErrorContentType(t *testing.T) {
	called := false

	handler := content.NewRequestHandler(test.Content, func(_ context.Context, _ *test.Request) (*test.Response, error) {
		called = true
		return &test.Response{}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(`{"name":"Bob"}`))
	req.Header.Set(content.TypeKey, mime.ErrorMediaType)

	res := httptest.NewRecorder()

	require.NotPanics(t, func() {
		handler.ServeHTTP(res, req)
	})

	require.False(t, called)
	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, mime.ErrorMediaType, res.Header().Get(content.TypeKey))
	require.Contains(t, res.Body.String(), `content: invalid request media type "text/error;charset=utf-8"`)
}

func TestNewHandlerRejectsErrorContentType(t *testing.T) {
	called := false

	handler := content.NewHandler(test.Content, func(_ context.Context) (*test.Response, error) {
		called = true
		return &test.Response{}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.TypeKey, mime.ErrorMediaType)

	res := httptest.NewRecorder()

	require.NotPanics(t, func() {
		handler.ServeHTTP(res, req)
	})
	require.False(t, called)
	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, mime.ErrorMediaType, res.Header().Get(content.TypeKey))
	require.Contains(t, res.Body.String(), `content: invalid request media type "text/error;charset=utf-8"`)
}

func TestNewHandlerDoesNotLeakPartialBodyWhenEncodeFails(t *testing.T) {
	enc := encoding.NewMap(encoding.MapParams{})
	enc.Register("json", partialEncoder{})
	cont := content.NewContent(enc, test.Pool)

	handler := content.NewHandler(cont, func(_ context.Context) (*test.Response, error) {
		return &test.Response{Greeting: "Hello Bob"}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.TypeKey, mime.JSONMediaType)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, mime.ErrorMediaType, res.Header().Get(content.TypeKey))
	require.Equal(t, test.ErrFailed.Error()+"\n", res.Body.String())
	require.NotContains(t, res.Body.String(), "partial")
}

type partialEncoder struct{}

func (partialEncoder) Encode(w io.Writer, _ any) error {
	_, _ = io.WriteString(w, "partial")
	return test.ErrFailed
}

func (partialEncoder) Decode(io.Reader, any) error {
	return nil
}
