package rpc_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestPostRequiresRequest(t *testing.T) {
	rpcClient := rpc.NewClient("http://example.com", test.NewContentClient(client.WithTimeout(time.MustParseDuration("1s"))))

	var res test.Response
	require.ErrorIs(t, rpcClient.Post(t.Context(), "/hello", nil, &res), rpc.ErrInvalidRequest)
}

func TestPostUsesConfiguredTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", media.JSON)
		res.WriteHeader(http.StatusOK)
		res.(http.Flusher).Flush()
		<-req.Context().Done()
	}))
	t.Cleanup(server.Close)

	rpcClient := rpc.NewClient(server.URL, test.NewContentClient(client.WithTimeout(time.MustParseDuration("10ms"))),
		rpc.WithClientContentType(media.JSON),
	)
	var res test.Response

	err := rpcClient.Post(t.Context(), "/hello", &test.Request{Name: "Bob"}, &res)

	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestPostRequiresNonNilTypedRequest(t *testing.T) {
	client := rpc.NewClient("http://example.com", test.NewContentClient())

	var req *test.Request
	var res test.Response
	require.ErrorIs(t, client.Post(t.Context(), "/hello", req, &res), rpc.ErrInvalidRequest)
}

func TestPostRequiresResponse(t *testing.T) {
	client := rpc.NewClient("http://example.com", test.NewContentClient())

	req := &test.Request{Name: "Bob"}
	require.ErrorIs(t, client.Post(t.Context(), "/hello", req, nil), rpc.ErrInvalidResponse)
}

func TestPostRequiresNonNilTypedResponse(t *testing.T) {
	client := rpc.NewClient("http://example.com", test.NewContentClient())

	req := &test.Request{Name: "Bob"}
	var res *test.Response
	require.ErrorIs(t, client.Post(t.Context(), "/hello", req, res), rpc.ErrInvalidResponse)
}

func TestPostUsesAccept(t *testing.T) {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	server := rpc.NewServer(router, test.UnaryContent, test.StreamContent, stream.Options{})
	server.Route("/hello", test.SuccessSayHello)

	httpServer := httptest.NewServer(mux)
	t.Cleanup(httpServer.Close)

	client := rpc.NewClient(httpServer.URL, test.NewContentClient(),
		rpc.WithClientContentType(media.JSON),
		rpc.WithClientAccept(media.YAML),
	)
	var res test.Response

	err := client.Post(t.Context(), "/hello", &test.Request{Name: "Bob"}, &res)

	require.NoError(t, err)
	require.Equal(t, "Hello Bob", res.Greeting)
}
