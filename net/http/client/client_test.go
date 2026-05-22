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
	"github.com/alexfalkowski/go-service/v2/strings"
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

func TestDoPreservesErrorMediaStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.Error)
		res.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(res, "bad request")
	}))
	defer server.Close()

	c := client.NewClient(test.Content, test.Pool)

	err := c.Get(t.Context(), server.URL, client.Options{})
	require.Error(t, err)
	require.EqualError(t, err, "bad request")
	require.Equal(t, http.StatusBadRequest, status.Code(err))
}

func TestDoNormalizesErrorMediaSuccessStatusCode(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{name: "ok", code: http.StatusOK},
		{name: "redirect", code: 302},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
				res.Header().Set(content.TypeKey, media.Error)
				res.WriteHeader(tt.code)
				_, _ = io.WriteString(res, "upstream error")
			}))
			defer server.Close()

			c := client.NewClient(test.Content, test.Pool)

			err := c.Get(t.Context(), server.URL, client.Options{Response: &struct{}{}})
			require.Error(t, err)
			require.EqualError(t, err, "upstream error")
			require.Equal(t, http.StatusInternalServerError, status.Code(err))
		})
	}
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
	require.Equal(t, 4*bytes.MB, bytes.DefaultSize)
}

func TestDoUsesMsgPack(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var request test.Request
		require.Equal(t, media.MessagePack, req.Header.Get(content.TypeKey))
		require.NoError(t, test.Encoder.Get("msgpack").Decode(req.Body, &request))
		res.Header().Set(content.TypeKey, media.MessagePack)
		require.NoError(t, test.Encoder.Get("msgpack").Encode(res, &test.Response{Greeting: "Hello " + request.Name}))
	}))
	defer server.Close()

	var response test.Response
	c := client.NewClient(test.Content, test.Pool)

	err := c.Post(t.Context(), server.URL, client.Options{
		ContentType: media.MessagePack,
		Request:     &test.Request{Name: "Bob"},
		Response:    &response,
	})
	require.NoError(t, err)
	require.Equal(t, "Hello Bob", response.Greeting)
}

func TestDoDetachesRequestBodyFromResponseBuffer(t *testing.T) {
	var body io.ReadCloser
	c := client.NewClient(test.Content, test.Pool, client.WithRoundTripper(test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		body = req.Body
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{content.TypeKey: []string{media.Text}},
			Body:       io.NopCloser(strings.NewReader("response")),
		}, nil
	})))

	err := c.Post(t.Context(), "http://example.com", client.Options{
		ContentType: media.JSON,
		Request:     &test.Request{Name: "Bob"},
	})
	require.NoError(t, err)
	require.NotNil(t, body)

	data, _, err := io.ReadAll(body)
	require.NoError(t, err)
	require.JSONEq(t, `{"Name":"Bob"}`, string(data))
}
