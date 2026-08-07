package rest_test

import (
	"bufio"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/transport/http/body"
	"github.com/stretchr/testify/require"
)

func TestGetAcceptsUnaryHandler(t *testing.T) {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	rest.Register(router, test.UnaryContent, test.StreamContent, test.Pool, stream.Options{})

	var handler unary.Handler[test.Response] = func(context.Context) (*test.Response, error) {
		return &test.Response{Greeting: "Hello Bob"}, nil
	}
	rest.Get("/hello", handler)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestStreamGetSendsValuesAndMarksRouteStreaming(t *testing.T) {
	mux := http.NewServeMux()
	policy := http.NewRoutePolicy()
	router := http.NewRouter(mux, policy)
	rest.Register(router, test.UnaryContent, test.StreamContent, test.Pool, stream.Options{})

	rest.StreamGet("/hello", func(_ context.Context, stream *stream.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return stream.Send(&test.Response{Greeting: "Hello Alice"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, media.NDJSON, res.Header().Get(http.ContentTypeKey))
	require.True(t, policy.IsStreaming(req))
	require.False(t, policy.IsRequestStreaming(req))

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeGreeting(t, scanner))
	require.Equal(t, "Hello Alice", decodeGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestStreamRouteLimitsRequestBody(t *testing.T) {
	mux := http.NewServeMux()
	policy := http.NewRoutePolicy()
	router := http.NewRouter(mux, policy)
	rest.Register(router, test.UnaryContent, test.StreamContent, test.Pool, stream.Options{})

	called := false
	rest.StreamRoute("POST /hello", func(context.Context, *stream.Stream[test.Response]) error {
		called = true
		return nil
	})

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", &test.UnknownLengthReader{Reader: strings.NewReader(strings.Repeat("a", 100))})
	require.NoError(t, err)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	body.NewHandler(policy, 64).ServeHTTP(res, req, mux.ServeHTTP)

	require.Equal(t, http.StatusRequestEntityTooLarge, res.Code)
	require.False(t, called)
}

func TestStreamPostRejectsHTTP1(t *testing.T) {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	rest.Register(router, test.UnaryContent, test.StreamContent, test.Pool, stream.Options{})

	rest.StreamPost("/hello", func(_ context.Context, _ *stream.RequestStream[test.Request, test.Response]) error {
		return nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 1
	req.Header.Set(http.ContentTypeKey, media.NDJSON)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusHTTPVersionNotSupported, res.Code)
}

func TestStreamPostPutPatchRecvAndSendOverHTTP2(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(string, stream.RequestHandler[test.Request, test.Response])
	}{
		{name: "post", method: http.MethodPost, call: rest.StreamPost[test.Request, test.Response]},
		{name: "put", method: http.MethodPut, call: rest.StreamPut[test.Request, test.Response]},
		{name: "patch", method: http.MethodPatch, call: rest.StreamPatch[test.Request, test.Response]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			policy := http.NewRoutePolicy()
			router := http.NewRouter(mux, policy)
			rest.Register(router, test.UnaryContent, test.StreamContent, test.Pool, stream.Options{})

			tt.call("/hello", func(_ context.Context, stream *stream.RequestStream[test.Request, test.Response]) error {
				for {
					req, err := stream.Recv()
					if err != nil {
						if stream.IsFinished(err) {
							return nil
						}

						return err
					}

					if err := stream.Send(&test.Response{Greeting: "Hello " + req.Name}); err != nil {
						return err
					}
				}
			})

			body := "{\"Name\":\"Bob\"}\n{\"Name\":\"Alice\"}\n"
			req := httptest.NewRequestWithContext(t.Context(), tt.method, "/hello", strings.NewReader(body))
			req.ProtoMajor = 2
			req.Header.Set(http.ContentTypeKey, media.NDJSON)
			req.Header.Set(http.AcceptKey, media.NDJSON)
			res := httptest.NewRecorder()

			mux.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)
			require.True(t, policy.IsStreaming(req))
			require.True(t, policy.IsRequestStreaming(req))

			scanner := bufio.NewScanner(res.Body)
			require.Equal(t, "Hello Bob", decodeGreeting(t, scanner))
			require.Equal(t, "Hello Alice", decodeGreeting(t, scanner))
			require.False(t, scanner.Scan())
		})
	}
}

func decodeGreeting(t *testing.T, scanner *bufio.Scanner) string {
	t.Helper()

	require.True(t, scanner.Scan())

	var res map[string]any
	require.NoError(t, json.Unmarshal(scanner.Bytes(), &res))

	greeting, _ := res["Greeting"].(string)
	return greeting
}
