package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestMaxBytesHandler(t *testing.T) {
	handler := http.MaxBytesHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, _, err := io.ReadAll(req.Body)
		var maxBytesError *http.MaxBytesError
		require.ErrorAs(t, err, &maxBytesError)

		_, _ = res.Write([]byte("ok"))
	}), 1)

	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", bytes.NewBufferString("too large"))

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestNewServerRejectsNegativeTimeoutOption(t *testing.T) {
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

	for _, key := range []string{"read_timeout", "write_timeout", "idle_timeout", "read_header_timeout"} {
		t.Run(key, func(t *testing.T) {
			require.Panics(t, func() {
				http.NewServer(options.Map{key: "-1s"}, time.Second, handler)
			})
		})
	}
}

func TestProtocols(t *testing.T) {
	protocols := http.Protocols()

	require.True(t, protocols.HTTP1())
	require.True(t, protocols.HTTP2())
	require.True(t, protocols.UnencryptedHTTP2())
}

func TestParseTime(t *testing.T) {
	now := time.Now().UTC()
	value := now.Format(http.TimeFormat)

	parsed, err := http.ParseTime(value)

	require.NoError(t, err)
	require.Equal(t, now.Truncate(time.Second.Duration()), parsed)
}

func TestTransport(t *testing.T) {
	cfg := &tls.Config{}

	transport := http.Transport(cfg)

	require.NotNil(t, transport.Proxy)
	require.NotNil(t, transport.DialContext)
	require.True(t, transport.ForceAttemptHTTP2)
	require.Equal(t, 100, transport.MaxIdleConns)
	require.Equal(t, 100, transport.MaxIdleConnsPerHost)
	require.Equal(t, 100, transport.MaxConnsPerHost)
	require.Equal(t, (90 * time.Second).Duration(), transport.IdleConnTimeout)
	require.Equal(t, (10 * time.Second).Duration(), transport.TLSHandshakeTimeout)
	require.Equal(t, time.Second.Duration(), transport.ExpectContinueTimeout)
	require.Same(t, cfg, transport.TLSClientConfig)
	require.NotNil(t, transport.Protocols)
	require.True(t, transport.Protocols.HTTP1())
	require.True(t, transport.Protocols.HTTP2())
	require.True(t, transport.Protocols.UnencryptedHTTP2())
}

func TestNewServerSetsProtocols(t *testing.T) {
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

	server := http.NewServer(options.Map{}, time.Second, handler)

	require.NotNil(t, server.Protocols)
	require.True(t, server.Protocols.HTTP1())
	require.True(t, server.Protocols.HTTP2())
	require.True(t, server.Protocols.UnencryptedHTTP2())
}

func TestSameOriginRedirect(t *testing.T) {
	tests := []struct {
		want error
		name string
		next string
	}{
		{name: "same origin", next: "https://example.com/next", want: nil},
		{name: "different host", next: "https://other.example.com/next", want: http.ErrUseLastResponse},
		{name: "different scheme", next: "http://example.com/next", want: http.ErrUseLastResponse},
	}

	prev, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/start", http.NoBody)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next, err := http.NewRequestWithContext(t.Context(), http.MethodGet, tt.next, http.NoBody)
			require.NoError(t, err)
			err = http.SameOriginRedirect(next, []*http.Request{prev})
			if tt.want == nil {
				require.NoError(t, err)
				return
			}

			require.ErrorIs(t, err, tt.want)
		})
	}
}

func TestSameOrigin(t *testing.T) {
	prev, err := url.Parse("https://example.com/start")
	require.NoError(t, err)

	same, err := url.Parse("https://example.com/next")
	require.NoError(t, err)

	different, err := url.Parse("https://other.example.com/next")
	require.NoError(t, err)

	require.True(t, http.SameOrigin(prev, same))
	require.False(t, http.SameOrigin(prev, different))
	require.False(t, http.SameOrigin(nil, same))
	require.False(t, http.SameOrigin(prev, nil))
}

func TestIsCrossOriginRedirect(t *testing.T) {
	prev, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/start", http.NoBody)
	require.NoError(t, err)

	same, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/next", http.NoBody)
	require.NoError(t, err)
	same.Response = &http.Response{Request: prev}

	different, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://other.example.com/next", http.NoBody)
	require.NoError(t, err)
	different.Response = &http.Response{Request: prev}

	require.False(t, http.IsCrossOriginRedirect(nil))
	require.False(t, http.IsCrossOriginRedirect(prev))
	require.False(t, http.IsCrossOriginRedirect(same))
	require.True(t, http.IsCrossOriginRedirect(different))
}

func TestClosingRoundTripperClosesBodyWhenRequested(t *testing.T) {
	rt := http.ClosingRoundTripper(func(*http.Request) (*http.Response, error, bool) {
		return nil, io.ErrUnexpectedEOF, true
	})
	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", body)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.True(t, body.Closed)
}

func TestClosingRoundTripperLeavesDelegatedBodyOpen(t *testing.T) {
	rt := http.ClosingRoundTripper(func(*http.Request) (*http.Response, error, bool) {
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil, false
	})
	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", body)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.False(t, body.Closed)
}

func TestIgnoreRedirect(t *testing.T) {
	err := http.IgnoreRedirect(nil, nil)
	require.ErrorIs(t, err, http.ErrUseLastResponse)
}

func TestParseServiceMethod(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		method  string
		service string
		action  string
	}{
		{name: "service route", method: http.MethodGet, url: "/test/hello", service: "test", action: "hello"},
		{name: "deep service route", method: http.MethodPost, url: "/test/users/123", service: "test", action: "users/123"},
		{name: "root", method: http.MethodGet, url: "/", service: "root", action: "get"},
		{name: "single segment", method: http.MethodPost, url: "/health", service: "health", action: "post"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), test.method, test.url, http.NoBody)
			service, action := http.ParseServiceMethod(req)
			require.Equal(t, test.service, service)
			require.Equal(t, test.action, action)
		})
	}
}

func TestHandleWhenTelemetryDisabled(t *testing.T) {
	require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})

	mux := http.NewServeMux()
	called := false

	http.Handle(mux, "/hello", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		called = true
		_, _ = res.Write([]byte("hello"))
	}))

	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)

	mux.ServeHTTP(res, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, res.Code)
	test.RequireResponseBody(t, res, "hello")
}

func TestHandleWhenMetricsEnabled(t *testing.T) {
	t.Cleanup(func() {
		require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Config:    &metrics.Config{},
		Reader:    metrics.NewManualReader(),
	})

	mux := http.NewServeMux()

	http.Handle(mux, "/hello", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		_, _ = res.Write([]byte("hello"))
	}))

	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	test.RequireResponseBody(t, res, "hello")
}
