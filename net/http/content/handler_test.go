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
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestNewRequestHandlerPrefersContentType(t *testing.T) {
	handler := content.NewRequestHandler(test.Content, func(_ context.Context, req *test.Request) (*test.Response, error) {
		return &test.Response{Greeting: "Hello " + req.Name}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(`{"name":"Bob"}`))
	req.Header.Set(content.TypeKey, media.JSON)
	req.Header.Set(content.AcceptKey, media.YAML)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, media.WithUTF8(media.JSON), res.Header().Get(content.TypeKey))
	require.JSONEq(t, `{"Greeting":"Hello Bob","Meta":null}`, res.Body.String())
}

func TestNewHandlerDoesNotLeakPartialBodyWhenEncodeFails(t *testing.T) {
	enc := encoding.NewMap(encoding.MapParams{})
	enc.Register("json", partialEncoder{})
	cont := content.NewContent(enc, test.Pool)

	handler := content.NewHandler(cont, func(_ context.Context) (*test.Response, error) {
		return &test.Response{Greeting: "Hello Bob"}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.TypeKey, media.JSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, media.WithUTF8(media.Error), res.Header().Get(content.TypeKey))
	require.Equal(t, test.ErrFailed.Error()+"\n", res.Body.String())
	require.NotContains(t, res.Body.String(), "partial")
}

func TestNotFoundHandlerWritesStatusError(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", http.NoBody)
	res := httptest.NewRecorder()

	require.True(t, content.NotFoundHandler()(res, req))
	require.Equal(t, http.StatusNotFound, res.Code)
	require.Equal(t, media.WithUTF8(media.Error), res.Header().Get(content.TypeKey))
	require.Equal(t, "Not Found\n", res.Body.String())
}

type partialEncoder struct{}

func (partialEncoder) Encode(w io.Writer, _ any) error {
	_, _ = io.WriteString(w, "partial")
	return test.ErrFailed
}

func (partialEncoder) Decode(io.Reader, any) error {
	return nil
}
