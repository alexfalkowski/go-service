package content_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestNewRequestHandlerUsesAcceptForResponse(t *testing.T) {
	handler := content.NewRequestHandler(test.Content, func(_ context.Context, req *test.Request) (*test.Response, error) {
		return &test.Response{Greeting: "Hello " + req.Name}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(`{"name":"Bob"}`))
	req.Header.Set(content.TypeKey, media.JSON)
	req.Header.Set(content.AcceptKey, media.YAML)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, media.YAML, res.Header().Get(content.TypeKey))
	var response test.Response
	require.NoError(t, test.Encoder.Get("yaml").Decode(res.Body, &response))
	require.Equal(t, "Hello Bob", response.Greeting)
}

func TestNewRequestHandlerRejectsUnsafeBinaryRequestBody(t *testing.T) {
	for _, tc := range []struct {
		mediaType string
		kind      string
	}{
		{mediaType: "application/gob", kind: "gob"},
		{mediaType: media.MessagePack + "; profile=test", kind: "msgpack"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			called := false
			handler := content.NewRequestHandler(test.Content, func(_ context.Context, req *test.Request) (*test.Response, error) {
				called = true
				return &test.Response{Greeting: "Hello " + req.Name}, nil
			})
			body := bytes.NewBuffer(nil)
			require.NoError(t, test.Encoder.Get(tc.kind).Encode(body, &test.Request{Name: "Bob"}))
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", body)
			req.Header.Set(content.TypeKey, tc.mediaType)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.False(t, called)
			require.Equal(t, http.StatusUnsupportedMediaType, res.Code)
			require.Equal(t, "text/error; charset=utf-8", res.Header().Get(content.TypeKey))
			test.RequireTrimmedResponseBody(t, res, "http: unsupported media type")
		})
	}
}

func TestNewRequestHandlerTreatsInternalErrorContentTypeAsText(t *testing.T) {
	handler := content.NewRequestHandler(test.Content, func(_ context.Context, req *bytes.Buffer) (*bytes.Buffer, error) {
		return bytes.NewBufferString("Hello " + req.String()), nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader("Bob"))
	req.Header.Set(content.TypeKey, media.Error)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "text/plain; charset=utf-8", res.Header().Get(content.TypeKey))
	test.RequireResponseBody(t, res, "Hello Bob")
}

func TestNewHandlerTreatsInternalErrorAcceptAsText(t *testing.T) {
	handler := content.NewHandler(test.Content, func(_ context.Context) (*bytes.Buffer, error) {
		return bytes.NewBufferString("Hello Bob"), nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.Error)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "text/plain; charset=utf-8", res.Header().Get(content.TypeKey))
	test.RequireResponseBody(t, res, "Hello Bob")
}

func TestNewHandlerReplacesExistingContentType(t *testing.T) {
	handler := content.NewHandler(test.Content, func(_ context.Context) (*bytes.Buffer, error) {
		return bytes.NewBufferString("Hello Bob"), nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.Text)
	res := httptest.NewRecorder()
	res.Header().Set(content.TypeKey, media.HTML)

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, []string{"text/plain; charset=utf-8"}, res.Header().Values(content.TypeKey))
	test.RequireResponseBody(t, res, "Hello Bob")
}

func TestNewHandlerDoesNotLeakPartialBodyWhenEncodeFails(t *testing.T) {
	enc := encoding.NewMap(encoding.MapParams{})
	enc.Register("json", test.PartialEncoder{})
	cont := content.NewContent(enc, test.Pool)

	handler := content.NewHandler(cont, func(_ context.Context) (*test.Response, error) {
		return &test.Response{Greeting: "Hello Bob"}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.TypeKey, media.JSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, "text/error; charset=utf-8", res.Header().Get(content.TypeKey))
	test.RequireTrimmedResponseBody(t, res, "http: internal server error")
	test.RequireResponseBodyNotContains(t, res, "partial")
}

func TestNotFoundHandlerWritesStatusError(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", http.NoBody)
	res := httptest.NewRecorder()

	require.True(t, content.NotFoundHandler()(res, req))
	require.Equal(t, http.StatusNotFound, res.Code)
	require.Equal(t, "text/error; charset=utf-8", res.Header().Get(content.TypeKey))
	test.RequireTrimmedResponseBody(t, res, "http: not found")
}
