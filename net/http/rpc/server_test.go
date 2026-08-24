package rpc_test

import (
	"bufio"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestStreamRouteRejectsHTTP1(t *testing.T) {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	server := rpc.NewServer(router, test.UnaryContent, test.StreamContent, stream.Options{})

	server.StreamRoute("/hello", func(_ context.Context, _ *stream.RequestStream[test.Request, test.Response]) error {
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

func TestRouteAppliesRouteOptions(t *testing.T) {
	mux := http.NewServeMux()
	policy := http.NewRoutePolicy()
	router := http.NewRouter(mux, policy)
	server := rpc.NewServer(router, test.UnaryContent, test.StreamContent, stream.Options{})

	server.Route("/hello", test.SuccessSayHello, http.WithRouteOperation(), http.WithRouteUnauthenticated())

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	require.True(t, policy.IsOperation(req))
	require.True(t, policy.IsUnauthenticated(req))
	require.False(t, policy.IsRequestStreaming(req))
	require.False(t, policy.IsResponseStreaming(req))
}

func TestStreamRouteRecvAndSendOverHTTP2(t *testing.T) {
	mux := http.NewServeMux()
	policy := http.NewRoutePolicy()
	router := http.NewRouter(mux, policy)
	server := rpc.NewServer(router, test.UnaryContent, test.StreamContent, stream.Options{})

	server.StreamRoute("/hello", func(_ context.Context, stream *stream.RequestStream[test.Request, test.Response]) error {
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
	}, http.WithRouteOperation(), http.WithRouteUnauthenticated())

	body := "{\"Name\":\"Bob\"}\n{\"Name\":\"Alice\"}\n"
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(body))
	req.ProtoMajor = 2
	req.Header.Set(http.ContentTypeKey, media.NDJSON)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.True(t, policy.IsOperation(req))
	require.True(t, policy.IsUnauthenticated(req))
	require.True(t, policy.IsRequestStreaming(req))
	require.True(t, policy.IsResponseStreaming(req))

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeGreeting(t, scanner))
	require.Equal(t, "Hello Alice", decodeGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func decodeGreeting(t *testing.T, scanner *bufio.Scanner) string {
	t.Helper()

	require.True(t, scanner.Scan())

	var res map[string]any
	require.NoError(t, json.Unmarshal(scanner.Bytes(), &res))

	greeting, _ := res["Greeting"].(string)
	return greeting
}
