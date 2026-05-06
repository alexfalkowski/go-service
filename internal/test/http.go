package test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	transporthttp "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/body"
	"github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	httplimiter "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/klauspost/compress/gzhttp"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni/v3"
	"go.uber.org/fx/fxtest"
)

// HTTPTestOption configures an httptest-backed HTTP transport test server.
type HTTPTestOption func(*httpTestOptions)

type httpTestOptions struct {
	verifier      token.Verifier
	serverLimiter *httplimiter.Server
	hello         bool
}

// WithHTTPTestVerifier configures server-side token verification for an HTTP test server.
func WithHTTPTestVerifier(verifier token.Verifier) HTTPTestOption {
	return func(opts *httpTestOptions) {
		opts.verifier = verifier
	}
}

// WithHTTPTestServerLimiter configures server-side rate limiting for an HTTP test server.
func WithHTTPTestServerLimiter(limiter *httplimiter.Server) HTTPTestOption {
	return func(opts *httpTestOptions) {
		opts.serverLimiter = limiter
	}
}

// WithHTTPTestHello registers the shared hello route on an HTTP test server.
func WithHTTPTestHello() HTTPTestOption {
	return func(opts *httpTestOptions) {
		opts.hello = true
	}
}

// HTTPClientOption configures an httptest-backed HTTP transport test client.
type HTTPClientOption func(*httpClientOptions)

type httpClientOptions struct {
	generator      token.Generator
	clientLimiter  *httplimiter.Client
	breakerOptions []breaker.Option
	breaker        bool
}

// WithHTTPClientGenerator configures client-side token generation for an HTTP test client.
func WithHTTPClientGenerator(generator token.Generator) HTTPClientOption {
	return func(opts *httpClientOptions) {
		opts.generator = generator
	}
}

// WithHTTPClientLimiter configures client-side rate limiting for an HTTP test client.
func WithHTTPClientLimiter(limiter *httplimiter.Client) HTTPClientOption {
	return func(opts *httpClientOptions) {
		opts.clientLimiter = limiter
	}
}

// WithHTTPClientBreaker configures client-side circuit breaking for an HTTP test client.
func WithHTTPClientBreaker(opts ...breaker.Option) HTTPClientOption {
	return func(options *httpClientOptions) {
		options.breaker = true
		options.breakerOptions = opts
	}
}

// NewHTTPTestServer returns an httptest server wired with shared RPC, REST, content, and metadata helpers.
func NewHTTPTestServer(tb testing.TB, opts ...HTTPTestOption) *httptest.Server {
	tb.Helper()

	os := &httpTestOptions{}
	for _, opt := range opts {
		opt(os)
	}

	mux := http.NewServeMux()
	rpc.Register(rpc.RegisterParams{
		Mux:     mux,
		Pool:    Pool,
		Content: Content,
	})
	rest.Register(rest.RegisterParams{
		Mux:     mux,
		Pool:    Pool,
		Content: Content,
	})
	if os.hello {
		mux.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(strings.Bytes("hello!"))
		})
	}

	neg := negroni.New()
	neg.Use(meta.NewHandler(UserAgent, Version, uuid.NewGenerator()))
	neg.Use(body.NewHandler(server.DefaultMaxReceiveSize.Bytes()))
	if os.verifier != nil {
		neg.Use(token.NewHandler(UserID, os.verifier))
	}
	if os.serverLimiter != nil {
		neg.Use(httplimiter.NewHandler(os.serverLimiter))
	}
	neg.UseHandler(gzhttp.GzipHandler(mux))

	srv := httptest.NewServer(neg)
	tb.Cleanup(srv.Close)

	return srv
}

// NewHTTPTestClient returns an HTTP transport client that talks to srv.
func NewHTTPTestClient(tb testing.TB, srv *httptest.Server, opts ...HTTPClientOption) *http.Client {
	tb.Helper()

	os := &httpClientOptions{}
	for _, opt := range opts {
		opt(os)
	}

	clientOpts := []transporthttp.ClientOption{
		transporthttp.WithClientRoundTripper(srv.Client().Transport),
		transporthttp.WithClientUserAgent(UserAgent),
		transporthttp.WithClientID(uuid.NewGenerator()),
		transporthttp.WithClientLimiter(os.clientLimiter),
	}
	if os.generator != nil {
		clientOpts = append(clientOpts, transporthttp.WithClientTokenGenerator(UserID, os.generator))
	}
	if os.breaker {
		clientOpts = append(clientOpts, transporthttp.WithClientBreaker(os.breakerOptions...))
	}

	client, err := transporthttp.NewClient(clientOpts...)
	require.NoError(tb, err)
	tb.Cleanup(client.CloseIdleConnections)

	return client
}

// NewHTTPTestServerLimiter returns an HTTP server limiter bound to tb cleanup.
func NewHTTPTestServerLimiter(tb testing.TB, cfg *limiter.Config) *httplimiter.Server {
	tb.Helper()

	lc := fxtest.NewLifecycle(tb)
	limiter, err := NewHTTPServerLimiter(lc, LimiterKeyMap, cfg)
	require.NoError(tb, err)
	tb.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_ = limiter.Close(ctx)
	})

	return limiter
}

// NewHTTPTestClientLimiter returns an HTTP client limiter bound to tb cleanup.
func NewHTTPTestClientLimiter(tb testing.TB, cfg *limiter.Config) *httplimiter.Client {
	tb.Helper()

	lc := fxtest.NewLifecycle(tb)
	limiter, err := NewHTTPClientLimiter(lc, LimiterKeyMap, cfg)
	require.NoError(tb, err)
	tb.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_ = limiter.Close(ctx)
	})

	return limiter
}

// HTTPResponseWithBody issues a request to srv and returns the response with a trimmed body.
func HTTPResponseWithBody(tb testing.TB, srv *httptest.Server, method, url string, header http.Header, body io.Reader) (*http.Response, string) {
	tb.Helper()

	req, err := http.NewRequestWithContext(tb.Context(), method, url, body)
	require.NoError(tb, err)
	req.Header = header

	res, err := srv.Client().Do(req)
	require.NoError(tb, err)
	defer res.Body.Close()

	data, _, err := io.ReadAll(res.Body)
	require.NoError(tb, err)

	return res, bytes.String(bytes.TrimSpace(data))
}

// HTTPClientResponseWithBody issues a request with client and returns the response with a trimmed body.
func HTTPClientResponseWithBody(tb testing.TB, client *http.Client, method, url string, header http.Header, body io.Reader) (*http.Response, string, error) {
	tb.Helper()

	req, err := http.NewRequestWithContext(tb.Context(), method, url, body)
	require.NoError(tb, err)
	req.Header = header

	res, err := client.Do(req)
	if err != nil {
		return nil, strings.Empty, err
	}
	defer res.Body.Close()

	data, _, err := io.ReadAll(res.Body)
	require.NoError(tb, err)

	return res, bytes.String(bytes.TrimSpace(data)), nil
}

// ErrResponseWriter is an http.ResponseWriter test double whose writes fail with ErrFailed.
type ErrResponseWriter struct {
	Code int
}

// Header is always empty.
func (w *ErrResponseWriter) Header() http.Header {
	return http.Header{}
}

// Write returns ErrFailed.
func (w *ErrResponseWriter) Write([]byte) (int, error) {
	return 0, ErrFailed
}

// WriteHeader stores code in the Code field.
func (w *ErrResponseWriter) WriteHeader(code int) {
	w.Code = code
}
