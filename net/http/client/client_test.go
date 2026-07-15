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
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.Pool, client.WithMaxResponseSize(5))

	err := c.Get(t.Context(), server.URL, client.Options{})
	require.NoError(t, err)
}

func TestDoRejectsResponseOverMaxResponseSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.Text)
		_, _ = io.WriteString(res, "hello!")
	}))
	t.Cleanup(server.Close)

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
	t.Cleanup(server.Close)

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
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.Pool)

	err := c.Get(t.Context(), server.URL, client.Options{})
	require.Error(t, err)
	require.EqualError(t, err, "bad request")
	require.Equal(t, http.StatusBadRequest, status.Code(err))
}

func TestDoUsesDefaultMessageForNonstandardErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.Text)
		res.WriteHeader(http.StatusClientClosedRequest)
	}))
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.Pool)

	err := c.Get(t.Context(), server.URL, client.Options{})
	require.Error(t, err)
	require.EqualError(t, err, "http: client closed request")
	require.Equal(t, http.StatusClientClosedRequest, status.Code(err))
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
			t.Cleanup(server.Close)

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
	t.Cleanup(server.Close)

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
		res.Header().Set(content.TypeKey, media.MessagePack+"; profile=test")
		require.NoError(t, test.Encoder.Get("msgpack").Encode(res, &test.Response{Greeting: "Hello " + request.Name}))
	}))
	t.Cleanup(server.Close)

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

func TestDoSendsAccept(t *testing.T) {
	c := client.NewClient(test.Content, test.Pool, client.WithRoundTripper(test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		var request test.Request
		require.Equal(t, media.JSON, req.Header.Get(content.TypeKey))
		require.Equal(t, media.YAML, req.Header.Get(content.AcceptKey))
		require.NoError(t, test.Encoder.Get("json").Decode(req.Body, &request))
		body := bytes.NewBuffer(nil)
		require.NoError(t, test.Encoder.Get("yaml").Encode(body, &test.Response{Greeting: "Hello " + request.Name}))

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{content.TypeKey: []string{media.YAML}},
			Body:       io.NopCloser(body),
		}, nil
	})))

	var response test.Response
	err := c.Post(t.Context(), "http://example.com", client.Options{
		ContentType: media.JSON,
		Accept:      media.YAML,
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

func TestDoDefaultRedirectRejectsCrossOriginCredentialRoundTrip(t *testing.T) {
	var urls []string
	rt := test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		urls = append(urls, req.URL.String())
		return &http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Header:     http.Header{"Location": []string{"https://other.example.com/target"}},
			Body:       http.NoBody,
			Request:    req,
		}, nil
	})
	c := client.NewClient(test.Content, test.Pool, client.WithRoundTripper(authRoundTripper{RoundTripper: rt}))

	err := c.Get(t.Context(), "https://example.com/start", client.Options{})

	require.NoError(t, err)
	require.Equal(t, []string{"https://example.com/start"}, urls)
}

func TestDoRedirectFollowAllowsExplicitCrossOriginCredentialRoundTrip(t *testing.T) {
	var urls []string
	rt := redirectingAuthRoundTripper(t, &urls, "https://other.example.com/target")
	c := client.NewClient(
		test.Content,
		test.Pool,
		client.WithRoundTripper(authRoundTripper{RoundTripper: rt}),
		client.WithRedirect(client.RedirectFollow),
	)

	err := c.Get(t.Context(), "https://example.com/start", client.Options{})

	require.NoError(t, err)
	require.Equal(t, []string{"https://example.com/start", "https://other.example.com/target"}, urls)
}

func TestDoRedirectSameOriginRejectsCrossOriginCredentialRoundTrip(t *testing.T) {
	var urls []string
	rt := test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		urls = append(urls, req.URL.String())
		return &http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Header:     http.Header{"Location": []string{"https://other.example.com/target"}},
			Body:       http.NoBody,
			Request:    req,
		}, nil
	})
	c := client.NewClient(
		test.Content,
		test.Pool,
		client.WithRoundTripper(authRoundTripper{RoundTripper: rt}),
		client.WithRedirect(client.RedirectSameOrigin),
	)

	err := c.Get(t.Context(), "https://example.com/start", client.Options{})

	require.NoError(t, err)
	require.Equal(t, []string{"https://example.com/start"}, urls)
}

func TestDoRedirectSameOriginAllowsSameOriginCredentialRoundTrip(t *testing.T) {
	var urls []string
	rt := redirectingAuthRoundTripper(t, &urls, "https://example.com/target")
	c := client.NewClient(
		test.Content,
		test.Pool,
		client.WithRoundTripper(authRoundTripper{RoundTripper: rt}),
		client.WithRedirect(client.RedirectSameOrigin),
	)

	err := c.Get(t.Context(), "https://example.com/start", client.Options{})

	require.NoError(t, err)
	require.Equal(t, []string{"https://example.com/start", "https://example.com/target"}, urls)
}

func TestDoRedirectIgnoreRejectsRedirectRoundTrip(t *testing.T) {
	var urls []string
	rt := test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		urls = append(urls, req.URL.String())
		return &http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Header:     http.Header{"Location": []string{"https://example.com/target"}},
			Body:       http.NoBody,
			Request:    req,
		}, nil
	})
	c := client.NewClient(
		test.Content,
		test.Pool,
		client.WithRoundTripper(rt),
		client.WithRedirect(client.RedirectIgnore),
	)

	err := c.Get(t.Context(), "https://example.com/start", client.Options{})

	require.NoError(t, err)
	require.Equal(t, []string{"https://example.com/start"}, urls)
}

type authRoundTripper struct {
	http.RoundTripper
}

func (r authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	cloned.Header.Set("Authorization", "Bearer secret")
	return r.RoundTripper.RoundTrip(cloned)
}

func redirectingAuthRoundTripper(t *testing.T, urls *[]string, target string) http.RoundTripper {
	t.Helper()

	return test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		*urls = append(*urls, req.URL.String())
		if len(*urls) == 1 {
			return &http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header:     http.Header{"Location": []string{target}},
				Body:       http.NoBody,
				Request:    req,
			}, nil
		}

		require.Equal(t, "Bearer secret", req.Header.Get("Authorization"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{content.TypeKey: []string{media.Text}},
			Body:       io.NopCloser(strings.NewReader("ok")),
			Request:    req,
		}, nil
	})
}
