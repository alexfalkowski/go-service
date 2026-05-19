package client_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
)

func TestDoAllowsResponseAtMaxResponseSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.Text)
		_, _ = io.WriteString(res, "hello")
	}))
	defer server.Close()

	c := client.NewClient(test.Content, test.Pool, client.WithMaxResponseSize(5))

	err := c.Get(t.Context(), server.URL, client.Options{})
	require.NoError(t, err)
}

func TestDoRejectsResponseOverMaxResponseSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.Text)
		_, _ = io.WriteString(res, "hello!")
	}))
	defer server.Close()

	c := client.NewClient(test.Content, test.Pool, client.WithMaxResponseSize(5))

	err := c.Get(t.Context(), server.URL, client.Options{})
	require.Error(t, err)
	require.EqualError(t, err, "http: request entity too large")
	require.Equal(t, http.StatusRequestEntityTooLarge, status.Code(err))
}

func TestDoRejectsErrorResponseOverMaxResponseSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.Error)
		res.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(res, "too large")
	}))
	defer server.Close()

	c := client.NewClient(test.Content, test.Pool, client.WithMaxResponseSize(5))

	err := c.Get(t.Context(), server.URL, client.Options{})
	require.Error(t, err)
	require.EqualError(t, err, "http: request entity too large")
	require.Equal(t, http.StatusRequestEntityTooLarge, status.Code(err))
}

func TestDoUsesDefaultMaxResponseSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.Text)
		_, _ = io.WriteString(res, "hello")
	}))
	defer server.Close()

	c := client.NewClient(test.Content, test.Pool, client.WithMaxResponseSize(0))

	err := c.Get(t.Context(), server.URL, client.Options{})
	require.NoError(t, err)
	require.Equal(t, 4*bytes.MB, client.DefaultMaxResponseSize)
}
