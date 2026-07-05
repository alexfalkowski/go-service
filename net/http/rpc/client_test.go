package rpc_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/stretchr/testify/require"
)

func TestPostRequiresRequest(t *testing.T) {
	client := rpc.NewClient("http://example.com")

	var res test.Response
	require.ErrorIs(t, client.Post(t.Context(), "/hello", nil, &res), rpc.ErrInvalidRequest)
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
	rpc.Register(rpc.RegisterParams{
		Mux:     mux,
		Content: test.Content,
		Pool:    test.Pool,
	})
	rpc.Route("/hello", test.SuccessSayHello)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := rpc.NewClient(server.URL,
		rpc.WithClientContentType(media.JSON),
		rpc.WithClientAccept(media.YAML),
	)
	var res test.Response

	err := client.Post(t.Context(), "/hello", &test.Request{Name: "Bob"}, &res)

	require.NoError(t, err)
	require.Equal(t, "Hello Bob", res.Greeting)
}
