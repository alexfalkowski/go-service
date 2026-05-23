package rpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
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
