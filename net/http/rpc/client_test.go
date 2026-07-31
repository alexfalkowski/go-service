package rpc_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/stretchr/testify/require"
)

func TestPostRequiresRequest(t *testing.T) {
	client := rpc.NewClient("http://example.com", rpc.WithClientTimeout("1s"))

	var res test.Response
	require.ErrorIs(t, client.Post(t.Context(), "/hello", nil, &res), rpc.ErrInvalidRequest)
}

func TestPostUsesConfiguredTimeout(t *testing.T) {
	rpc.Register(nil, test.Content, test.Pool, content.StreamOptions{})

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", media.JSON)
		res.WriteHeader(http.StatusOK)
		res.(http.Flusher).Flush()
		<-req.Context().Done()
	}))
	t.Cleanup(server.Close)

	client := rpc.NewClient(server.URL,
		rpc.WithClientContentType(media.JSON),
		rpc.WithClientTimeout("10ms"),
	)
	var res test.Response

	err := client.Post(t.Context(), "/hello", &test.Request{Name: "Bob"}, &res)

	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestPostRequiresNonNilTypedRequest(t *testing.T) {
	client := rpc.NewClient("http://example.com")

	var req *test.Request
	var res test.Response
	require.ErrorIs(t, client.Post(t.Context(), "/hello", req, &res), rpc.ErrInvalidRequest)
}

func TestPostRequiresResponse(t *testing.T) {
	client := rpc.NewClient("http://example.com")

	req := &test.Request{Name: "Bob"}
	require.ErrorIs(t, client.Post(t.Context(), "/hello", req, nil), rpc.ErrInvalidResponse)
}

func TestPostRequiresNonNilTypedResponse(t *testing.T) {
	client := rpc.NewClient("http://example.com")

	req := &test.Request{Name: "Bob"}
	var res *test.Response
	require.ErrorIs(t, client.Post(t.Context(), "/hello", req, res), rpc.ErrInvalidResponse)
}

func TestPostUsesAccept(t *testing.T) {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	rpc.Register(router, test.Content, test.Pool, content.StreamOptions{})
	rpc.Route("/hello", test.SuccessSayHello)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := rpc.NewClient(server.URL,
		rpc.WithClientContentType(media.JSON),
		rpc.WithClientAccept(media.YAML),
	)
	var res test.Response

	err := client.Post(t.Context(), "/hello", &test.Request{Name: "Bob"}, &res)

	require.NoError(t, err)
	require.Equal(t, "Hello Bob", res.Greeting)
}
