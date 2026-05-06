package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
)

func TestUnlimited(t *testing.T) {
	cfg := test.NewLimiterConfig("user-agent", "1s", 100)
	server := test.NewHTTPTestServer(t,
		test.WithHTTPTestHello(),
		test.WithHTTPTestServerLimiter(test.NewHTTPTestServerLimiter(t, cfg)),
	)
	client := test.NewHTTPTestClient(t, server, test.WithHTTPClientLimiter(test.NewHTTPTestClientLimiter(t, cfg)))

	url := server.URL + "/hello"

	_, _, err := test.HTTPClientResponseWithBody(t, client, http.MethodGet, url, http.Header{}, http.NoBody)
	require.NoError(t, err)

	res, body, err := test.HTTPClientResponseWithBody(t, client, http.MethodGet, url, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "hello!", body)
}

func TestServerLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
		t.Run(f, func(t *testing.T) {
			server := test.NewHTTPTestServer(t,
				test.WithHTTPTestHello(),
				test.WithHTTPTestServerLimiter(test.NewHTTPTestServerLimiter(t, test.NewLimiterConfig(f, "1s", 0))),
			)

			url := server.URL + "/hello"

			_, _ = test.HTTPResponseWithBody(t, server, http.MethodGet, url, http.Header{}, http.NoBody)

			res, _ := test.HTTPResponseWithBody(t, server, http.MethodGet, url, http.Header{}, http.NoBody)
			require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
			require.NotEmpty(t, res.Header.Get("Ratelimit"))
		})
	}
}

func TestClientLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
		t.Run(f, func(t *testing.T) {
			server := test.NewHTTPTestServer(t, test.WithHTTPTestHello())
			client := test.NewHTTPTestClient(t, server,
				test.WithHTTPClientLimiter(test.NewHTTPTestClientLimiter(t, test.NewLimiterConfig(f, "1s", 0))),
			)

			url := server.URL + "/hello"

			_, _, err := test.HTTPClientResponseWithBody(t, client, http.MethodGet, url, http.Header{}, http.NoBody)
			require.NoError(t, err)

			_, _, err = test.HTTPClientResponseWithBody(t, client, http.MethodGet, url, http.Header{}, http.NoBody)
			require.Error(t, err)
			require.Equal(t, http.StatusTooManyRequests, status.Code(err))
		})
	}
}

func TestServerClosedLimiter(t *testing.T) {
	limiter := test.NewHTTPTestServerLimiter(t, test.NewLimiterConfig("user-agent", "1s", 100))
	server := test.NewHTTPTestServer(t,
		test.WithHTTPTestHello(),
		test.WithHTTPTestServerLimiter(limiter),
	)

	err := limiter.Close(t.Context())
	require.NoError(t, err)

	url := server.URL + "/hello"

	res, _ := test.HTTPResponseWithBody(t, server, http.MethodGet, url, http.Header{}, http.NoBody)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestClientClosedLimiter(t *testing.T) {
	server := test.NewHTTPTestServer(t, test.WithHTTPTestHello())
	limiter := test.NewHTTPTestClientLimiter(t, test.NewLimiterConfig("user-agent", "1s", 100))
	client := test.NewHTTPTestClient(t, server,
		test.WithHTTPClientLimiter(limiter),
	)

	url := server.URL + "/hello"

	err := limiter.Close(t.Context())
	require.NoError(t, err)

	_, _, err = test.HTTPClientResponseWithBody(t, client, http.MethodGet, url, http.Header{}, http.NoBody)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, status.Code(err))
}
