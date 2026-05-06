package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/stretchr/testify/require"
)

func TestRestNoContent(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		t.Run(method, func(t *testing.T) {
			server := test.NewHTTPTestServer(t)

			test.RegisterHandlers("/hello", test.RestNoContent)

			url := server.URL + http.Pattern(test.Name, "/hello")
			client := rest.NewClient(rest.WithClientRoundTripper(server.Client().Transport), rest.WithClientTimeout("10s"))
			err := client.Do(t.Context(), method, url, rest.NoOptions)
			require.NoError(t, err)
		})
	}
}

func TestRestRequestNoContent(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			server := test.NewHTTPTestServer(t)

			test.RegisterRequestHandlers("/hello", test.RestRequestNoContent)

			url := server.URL + http.Pattern(test.Name, "/hello")
			req := &test.Request{Name: "test"}
			opts := &rest.Options{
				ContentType: mime.JSONMediaType,
				Request:     req,
			}
			client := rest.NewClient(rest.WithClientRoundTripper(server.Client().Transport), rest.WithClientTimeout("10s"))
			err := client.Do(t.Context(), method, url, opts)
			require.NoError(t, err)
		})
	}
}

func TestRestError(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		t.Run(method, func(t *testing.T) {
			server := test.NewHTTPTestServer(t)

			test.RegisterHandlers("/hello", test.RestError)

			url := server.URL + http.Pattern(test.Name, "/hello")
			client := rest.NewClient(rest.WithClientRoundTripper(server.Client().Transport), rest.WithClientTimeout("10s"))
			err := client.Do(t.Context(), method, url, rest.NoOptions)
			require.Error(t, err)
		})
	}
}

func TestRestRequestError(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			server := test.NewHTTPTestServer(t)

			test.RegisterRequestHandlers("/hello", test.RestRequestError)

			url := server.URL + http.Pattern(test.Name, "/hello")
			req := &test.Request{Name: "test"}
			opts := &rest.Options{
				ContentType: mime.JSONMediaType,
				Request:     req,
			}
			client := rest.NewClient(rest.WithClientRoundTripper(server.Client().Transport), rest.WithClientTimeout("10s"))
			err := client.Do(t.Context(), method, url, opts)
			require.Error(t, err)
		})
	}
}

func TestRestWithContent(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		t.Run(method, func(t *testing.T) {
			server := test.NewHTTPTestServer(t)

			test.RegisterHandlers("/hello", test.RestContent)

			url := server.URL + http.Pattern(test.Name, "/hello")
			resp := &test.Response{}
			opts := &rest.Options{
				Response: resp,
			}
			client := rest.NewClient(rest.WithClientRoundTripper(server.Client().Transport), rest.WithClientTimeout("10s"))
			err := client.Do(t.Context(), method, url, opts)
			require.NoError(t, err)
			require.Equal(t, "Hello Bob", resp.Greeting)
		})
	}
}

func TestRestRequestWithContent(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			server := test.NewHTTPTestServer(t)

			test.RegisterRequestHandlers("/hello", test.RestRequestContent)

			url := server.URL + http.Pattern(test.Name, "/hello")
			req := &test.Request{Name: "test"}
			resp := &test.Response{}
			opts := &rest.Options{
				ContentType: mime.JSONMediaType,
				Request:     req,
				Response:    resp,
			}
			client := rest.NewClient(rest.WithClientRoundTripper(server.Client().Transport), rest.WithClientTimeout("10s"))
			err := client.Do(t.Context(), method, url, opts)
			require.NoError(t, err)
			require.Equal(t, "Hello test", resp.Greeting)
		})
	}
}

func TestRestInvalidStatusCode(t *testing.T) {
	server := test.NewHTTPTestServer(t)
	client := rest.NewClient(rest.WithClientRoundTripper(server.Client().Transport), rest.WithClientTimeout("10s"))

	test.RegisterHandlers("/hello", test.RestInvalidStatusCode)

	url := server.URL + http.Pattern(test.Name, "/hello")
	err := client.Get(t.Context(), url, rest.NoOptions)
	require.Error(t, err)

	err = client.Delete(t.Context(), url, rest.NoOptions)
	require.Error(t, err)

	test.RegisterRequestHandlers("/hello", test.RestRequestInvalidStatusCode)

	url = server.URL + http.Pattern(test.Name, "/hello")
	req := &test.Request{}
	opts := &rest.Options{Request: req}

	err = client.Post(t.Context(), url, opts)
	require.Error(t, err)

	err = client.Put(t.Context(), url, opts)
	require.Error(t, err)

	err = client.Patch(t.Context(), url, opts)
	require.Error(t, err)
}
