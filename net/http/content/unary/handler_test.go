package unary_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestNewRequestHandlerUsesAcceptForResponse(t *testing.T) {
	for _, tc := range []struct {
		mediaType string
		kind      string
	}{
		{mediaType: media.TOML, kind: "toml"},
		{mediaType: media.YAML, kind: "yaml"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			handler := unary.NewRequestHandler(test.UnaryContent, func(_ context.Context, req *test.Request) (*test.Response, error) {
				return &test.Response{Greeting: "Hello " + req.Name}, nil
			})

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(`{"name":"Bob"}`))
			req.Header.Set(http.ContentTypeKey, media.JSON)
			req.Header.Set(http.AcceptKey, tc.mediaType)
			res := httptest.NewRecorder()
			res.Header().Set(http.VaryKey, "Accept-Encoding, Accept")

			handler.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)
			require.Equal(t, tc.mediaType, res.Header().Get(http.ContentTypeKey))
			require.Equal(t, []string{"Accept-Encoding, Accept", http.ContentTypeKey}, res.Header().Values(http.VaryKey))
			var response test.Response
			require.NoError(t, test.Encoder.Get(tc.kind).Decode(res.Body, &response))
			require.Equal(t, "Hello Bob", response.Greeting)
		})
	}
}

func TestNewRequestHandlerRejectsUnsafeRequestBody(t *testing.T) {
	for _, tc := range []struct {
		mediaType string
		kind      string
	}{
		{mediaType: media.TOML, kind: "toml"},
		{mediaType: media.TOML + "; profile=test", kind: "toml"},
		{mediaType: "application/gob", kind: "gob"},
		{mediaType: media.MessagePack + "; profile=test", kind: "msgpack"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			called := false
			handler := unary.NewRequestHandler(test.UnaryContent, func(_ context.Context, req *test.Request) (*test.Response, error) {
				called = true
				return &test.Response{Greeting: "Hello " + req.Name}, nil
			})
			body := bytes.NewBuffer(nil)
			require.NoError(t, test.Encoder.Get(tc.kind).Encode(body, &test.Request{Name: "Bob"}))
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", body)
			req.Header.Set(http.ContentTypeKey, tc.mediaType)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.False(t, called)
			require.Equal(t, http.StatusUnsupportedMediaType, res.Code)
			require.Equal(t, "text/error; charset=utf-8", res.Header().Get(http.ContentTypeKey))
			require.Equal(t, []string{http.AcceptKey, http.ContentTypeKey}, res.Header().Values(http.VaryKey))
			test.RequireTrimmedResponseBody(t, res, "http: unsupported media type")
		})
	}
}

func TestNewRequestHandlerTreatsInternalErrorContentTypeAsText(t *testing.T) {
	handler := unary.NewRequestHandler(test.UnaryContent, func(_ context.Context, req *bytes.Buffer) (*bytes.Buffer, error) {
		return bytes.NewBufferString("Hello " + req.String()), nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader("Bob"))
	req.Header.Set(http.ContentTypeKey, media.Error)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "text/plain; charset=utf-8", res.Header().Get(http.ContentTypeKey))
	test.RequireResponseBody(t, res, "Hello Bob")
}

func TestNewHandlerTreatsInternalErrorAcceptAsText(t *testing.T) {
	handler := unary.NewHandler(test.UnaryContent, func(_ context.Context) (*bytes.Buffer, error) {
		return bytes.NewBufferString("Hello Bob"), nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.Error)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "text/plain; charset=utf-8", res.Header().Get(http.ContentTypeKey))
	test.RequireResponseBody(t, res, "Hello Bob")
}

func TestNewHandlerReplacesExistingContentType(t *testing.T) {
	handler := unary.NewHandler(test.UnaryContent, func(_ context.Context) (*bytes.Buffer, error) {
		return bytes.NewBufferString("Hello Bob"), nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.Text)
	res := httptest.NewRecorder()
	res.Header().Set(http.ContentTypeKey, media.HTML)

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, []string{"text/plain; charset=utf-8"}, res.Header().Values(http.ContentTypeKey))
	test.RequireResponseBody(t, res, "Hello Bob")
}

func TestNewHandlerCanonicalizesResponseContentType(t *testing.T) {
	handler := unary.NewHandler(test.UnaryContent, func(_ context.Context) (*test.Response, error) {
		return &test.Response{Greeting: "Hello Héllo"}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, "application/json; charset=iso-8859-1; q=0.5")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, media.JSON, res.Header().Get(http.ContentTypeKey))
	test.RequireResponseBodyContains(t, res, "Héllo")
}

func TestNewHandlerDoesNotLeakPartialBodyWhenEncodeFails(t *testing.T) {
	enc := encoding.NewMap()
	enc.Register("json", test.PartialEncoder{})
	cont := unary.NewContent(enc, test.Pool)

	handler := unary.NewHandler(cont, func(_ context.Context) (*test.Response, error) {
		return &test.Response{Greeting: "Hello Bob"}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.ContentTypeKey, media.JSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, "text/error; charset=utf-8", res.Header().Get(http.ContentTypeKey))
	test.RequireTrimmedResponseBody(t, res, "http: internal server error")
	test.RequireResponseBodyNotContains(t, res, "partial")
}

func TestNewHandlerRejectsUnencodableAcceptAsNotAcceptable(t *testing.T) {
	for _, mediaType := range []string{
		media.Text, "application/octet-stream", media.Protobuf, media.ProtobufText, media.ProtobufJSON,
	} {
		t.Run(mediaType, func(t *testing.T) {
			called := false
			handler := unary.NewHandler(test.UnaryContent, func(_ context.Context) (*test.Response, error) {
				called = true

				return &test.Response{Greeting: "Hello Bob"}, nil
			})

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
			req.Header.Set(http.AcceptKey, mediaType)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.True(t, called)
			require.Equal(t, http.StatusNotAcceptable, res.Code)
			require.Equal(t, "text/error; charset=utf-8", res.Header().Get(http.ContentTypeKey))
			test.RequireTrimmedResponseBody(t, res, "http: not acceptable")
		})
	}
}

func TestNewHandlerRejectsUnavailableResponseCodec(t *testing.T) {
	enc := encoding.NewMap()
	enc.Register("json", nil)
	cont := unary.NewContent(enc, test.Pool)

	called := false
	handler := unary.NewHandler(cont, func(_ context.Context) (*test.Response, error) {
		called = true

		return &test.Response{Greeting: "Hello Bob"}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.JSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.False(t, called)
	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, "text/error; charset=utf-8", res.Header().Get(http.ContentTypeKey))
	test.RequireTrimmedResponseBody(t, res, "http: internal server error")
}
