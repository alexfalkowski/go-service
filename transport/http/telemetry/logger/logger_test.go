package logger_test

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v3/id/uuid"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/meta"
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/net/http/client"
	"github.com/alexfalkowski/go-service/v3/net/http/content/stream"
	httpmeta "github.com/alexfalkowski/go-service/v3/net/http/meta"
	"github.com/alexfalkowski/go-service/v3/net/http/rpc"
	"github.com/alexfalkowski/go-service/v3/net/http/status"
	"github.com/alexfalkowski/go-service/v3/strings"
	"github.com/alexfalkowski/go-service/v3/telemetry/logger"
	httplogger "github.com/alexfalkowski/go-service/v3/transport/http/telemetry/logger"
	"github.com/alexfalkowski/go-service/v3/transport/http/token"
	"github.com/stretchr/testify/require"
)

func TestHandlerLogsResponse(t *testing.T) {
	tests := []struct {
		name  string
		level string
		code  int
	}{
		{name: "success", code: http.StatusOK, level: "INFO"},
		{name: "client error", code: http.StatusBadRequest, level: "WARN"},
		{name: "server error", code: http.StatusInternalServerError, level: "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
			handler := httplogger.NewHandler(http.NewRoutePolicy(), &logger.Logger{Logger: slogLogger})
			res := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/greeter/say-hello", http.NoBody)

			handler.ServeHTTP(res, req, func(res http.ResponseWriter, _ *http.Request) {
				res.WriteHeader(tt.code)
			})

			require.Contains(t, logs.String(), `"level":"`+tt.level+`"`)
			require.Contains(t, logs.String(), `"msg":"http: get /greeter/say-hello"`)
			require.Contains(t, logs.String(), `"system":"http"`)
			require.Contains(t, logs.String(), `"http.request.method":"get"`)
			require.Contains(t, logs.String(), `"url.path":"/greeter/say-hello"`)
			require.NotContains(t, logs.String(), `"http.route":`)
			require.Contains(t, logs.String(), `"code":`+strconv.Itoa(tt.code))
		})
	}
}

func TestHandlerLogsMatchedRoute(t *testing.T) {
	var logs bytes.Buffer
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	mux, policy := http.NewServeMux(), http.NewRoutePolicy()
	router := http.NewRouter(mux, policy)
	router.HandleRoute("GET example.com/repos/{owner}/{repo}", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusNoContent)
	}))
	metadata := httpmeta.NewHandler(test.UserAgent, test.Version, test.StaticIDGenerator("request-id"), mux, 1024)
	chain := http.NewChainedHandlers(metadata, httplogger.NewHandler(policy, &logger.Logger{Logger: slogLogger}))
	chain.UseHandler(mux)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodHead, "https://example.com/repos/alexfalkowski/abgleich?token=secret", http.NoBody)
	res := httptest.NewRecorder()
	chain.ServeHTTP(res, req)

	require.Equal(t, http.StatusNoContent, res.Code)
	require.Contains(t, logs.String(), `"msg":"http: head /repos/{owner}/{repo}"`)
	require.Contains(t, logs.String(), `"http.request.method":"head"`)
	require.Contains(t, logs.String(), `"http.route":"/repos/{owner}/{repo}"`)
	require.Contains(t, logs.String(), `"url.path":"/repos/alexfalkowski/abgleich"`)
	require.NotContains(t, logs.String(), "token=secret")
}

func TestHandlerLogsRPCOnlyForRegisteredRPCRoutes(t *testing.T) {
	const path = "/greet.v1.Greeter/SayHello"
	tests := []struct {
		name    string
		message string
		rpc     bool
	}{
		{name: "ordinary POST", message: strings.Concat("http: post ", path)},
		{name: "RPC POST", message: "http: greet.v1.Greeter/SayHello", rpc: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
			mux, policy := http.NewServeMux(), http.NewRoutePolicy()
			router := http.NewRouter(mux, policy)
			if tt.rpc {
				server := rpc.NewServer(router, test.UnaryContent, test.StreamContent, stream.Options{})
				server.Route(path, test.NoContent)
			} else {
				router.HandleRoute(strings.Concat("POST ", path), http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
					res.WriteHeader(http.StatusNoContent)
				}))
			}
			metadata := httpmeta.NewHandler(test.UserAgent, test.Version, test.StaticIDGenerator("request-id"), mux, 1024)
			chain := http.NewChainedHandlers(metadata, httplogger.NewHandler(policy, &logger.Logger{Logger: slogLogger}))
			chain.UseHandler(mux)
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, path, bytes.NewBufferString(`{}`))
			res := httptest.NewRecorder()
			chain.ServeHTTP(res, req)

			require.Less(t, res.Code, 400)
			require.Contains(t, logs.String(), strings.Concat(`"msg":"`, tt.message, `"`))
			require.Contains(t, logs.String(), `"http.request.method":"post"`)
			if tt.rpc {
				require.Contains(t, logs.String(), `"service":"greet.v1.Greeter"`)
				require.Contains(t, logs.String(), `"method":"SayHello"`)
			} else {
				require.NotContains(t, logs.String(), `"service":`)
				require.NotContains(t, logs.String(), `"method":`)
			}
		})
	}
}

func TestHandlerLogsStatusError(t *testing.T) {
	var logs bytes.Buffer
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	handler := httplogger.NewHandler(http.NewRoutePolicy(), &logger.Logger{Logger: slogLogger})
	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/greeter/say-hello", http.NoBody)

	handler.ServeHTTP(res, req, func(res http.ResponseWriter, req *http.Request) {
		require.NoError(t, status.WriteError(req.Context(), res, status.BadRequestError(test.ErrInvalid)))
	})

	require.Equal(t, http.StatusBadRequest, res.Code)
	test.RequireTrimmedResponseBody(t, res, "http: bad request")
	require.Contains(t, logs.String(), `"level":"WARN"`)
	require.Contains(t, logs.String(), `"error":"`+test.ErrInvalid.Error()+`"`)
}

func TestHandlerLogsRedactedPasetoTokenError(t *testing.T) {
	cfg := test.NewToken("paseto")
	verifier := token.NewToken(cfg, test.FS, uuid.NewGenerator())
	auth := token.NewHandler(http.NewRoutePolicy(), verifier)
	var logs bytes.Buffer
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	handler := httplogger.NewHandler(http.NewRoutePolicy(), &logger.Logger{Logger: slogLogger})
	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/greeter/say-hello", http.NoBody)
	req = req.WithContext(meta.WithAttributes(req.Context(), meta.WithAuthorization(meta.String("secret-header.secret-claims.signature"))))

	handler.ServeHTTP(res, req, func(res http.ResponseWriter, req *http.Request) {
		auth.ServeHTTP(res, req, func(http.ResponseWriter, *http.Request) {
			require.Fail(t, "unexpected authorized request")
		})
	})

	require.Equal(t, http.StatusUnauthorized, res.Code)
	require.Contains(t, logs.String(), `"error":"token: invalid match"`)
	require.NotContains(t, logs.String(), "secret-header")
	require.NotContains(t, logs.String(), "secret-claims")
}

func TestHandlerSkipsOperationPath(t *testing.T) {
	var logs bytes.Buffer
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	policy := http.NewRoutePolicy()
	http.NewRouter(http.NewServeMux(), policy).HandleRoute("GET /test/metrics", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), http.WithRouteOperation())
	handler := httplogger.NewHandler(policy, &logger.Logger{Logger: slogLogger})
	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test/metrics", http.NoBody)
	called := false

	handler.ServeHTTP(res, req, func(res http.ResponseWriter, _ *http.Request) {
		called = true
		res.WriteHeader(http.StatusNoContent)
	})

	require.True(t, called)
	require.Empty(t, logs.String())
}

func TestHandlerLogsRecoveredPanic(t *testing.T) {
	var logs bytes.Buffer
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	handler := httplogger.NewHandler(http.NewRoutePolicy(), &logger.Logger{Logger: slogLogger})
	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/greeter/say-hello", http.NoBody)

	require.PanicsWithValue(t, "boom", func() {
		handler.ServeHTTP(res, req, func(_ http.ResponseWriter, _ *http.Request) {
			panic("boom")
		})
	})

	require.Contains(t, logs.String(), `"msg":"http: panic"`)
	require.Contains(t, logs.String(), `"level":"ERROR"`)
}

func TestHandlerSuppressesAbortHandlerPanicLog(t *testing.T) {
	var logs bytes.Buffer
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	handler := httplogger.NewHandler(http.NewRoutePolicy(), &logger.Logger{Logger: slogLogger})
	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/greeter/say-hello", http.NoBody)
	streamErr := errors.New("stream aborted")

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req, func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			status.RecordError(req.Context(), streamErr)
			panic(http.ErrAbortHandler)
		})
	})

	require.NotContains(t, logs.String(), `"msg":"http: panic"`)
	require.Contains(t, logs.String(), `"msg":"http: get /greeter/say-hello"`)
	require.Contains(t, logs.String(), `"level":"INFO"`)
	require.Contains(t, logs.String(), `"code":200`)
	require.Contains(t, logs.String(), `"error":"stream aborted"`)
}

func TestRoundTripperLogsTransportError(t *testing.T) {
	var logs bytes.Buffer
	transportErr := errors.New("dial failed")
	transport := &test.ErrorRoundTripper{Err: transportErr}
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	rt := httplogger.NewRoundTripper(&logger.Logger{Logger: slogLogger}, transport)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/users", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.ErrorIs(t, err, transportErr)
	require.Nil(t, res)
	require.Contains(t, logs.String(), `"level":"ERROR"`)
	require.Contains(t, logs.String(), `"error":"dial failed"`)
}

func TestRoundTripperLogsStatusErrorCode(t *testing.T) {
	var logs bytes.Buffer
	err := status.Error(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
	transport := &test.ErrorRoundTripper{Err: err}
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	rt := httplogger.NewRoundTripper(&logger.Logger{Logger: slogLogger}, transport)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/users", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, logs.String(), `"level":"WARN"`)
	require.Contains(t, logs.String(), `"code":429`)
}

func TestRoundTripperLogsResponse(t *testing.T) {
	tests := []struct {
		name  string
		level string
		code  int
	}{
		{name: "success", code: http.StatusOK, level: "INFO"},
		{name: "client error", code: http.StatusBadRequest, level: "WARN"},
		{name: "server error", code: http.StatusInternalServerError, level: "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			transport := &test.StatusRoundTripper{Status: tt.code}
			slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
			rt := httplogger.NewRoundTripper(&logger.Logger{Logger: slogLogger}, transport)
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://api.github.com/repos/alexfalkowski/abgleich?token=secret", http.NoBody)
			require.NoError(t, err)

			res, err := rt.RoundTrip(req)
			require.NoError(t, err)
			require.Equal(t, tt.code, res.StatusCode)
			require.Contains(t, logs.String(), `"level":"`+tt.level+`"`)
			require.Contains(t, logs.String(), `"msg":"http: get api.github.com/repos/alexfalkowski/abgleich"`)
			require.Contains(t, logs.String(), `"system":"http"`)
			require.Contains(t, logs.String(), `"server.address":"api.github.com"`)
			require.Contains(t, logs.String(), `"http.request.method":"get"`)
			require.Contains(t, logs.String(), `"url.path":"/repos/alexfalkowski/abgleich"`)
			require.NotContains(t, logs.String(), `"service":`)
			require.NotContains(t, logs.String(), `"method":`)
			require.NotContains(t, logs.String(), "token=secret")
			require.Contains(t, logs.String(), `"code":`+strconv.Itoa(tt.code))
		})
	}
}

func TestRoundTripperLogsURLTemplate(t *testing.T) {
	var logs bytes.Buffer
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	redirectTransport := test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "api.github.com" {
			res := test.ResponseWithStatus(http.StatusFound)
			res.Header.Set(http.LocationKey, "https://example.com/repositories/123")
			return res, nil
		}
		return test.ResponseWithStatus(http.StatusNoContent), nil
	})
	rt := httplogger.NewRoundTripper(&logger.Logger{Logger: slogLogger}, redirectTransport)
	ctx := httpmeta.WithURLTemplate(t.Context(), "/repos/{owner}/{repo}")
	httpClient := test.NewContentClient(client.WithRoundTripper(rt), client.WithRedirect(client.RedirectFollow))

	require.NoError(t, httpClient.Get(ctx, "https://api.github.com/repos/alexfalkowski/abgleich", client.Options{}))
	require.Contains(t, logs.String(), `"msg":"http: get api.github.com/repos/{owner}/{repo}"`)
	require.Contains(t, logs.String(), `"url.template":"/repos/{owner}/{repo}"`)
	require.Contains(t, logs.String(), `"url.path":"/repos/alexfalkowski/abgleich"`)
	require.Contains(t, logs.String(), `"msg":"http: get example.com/repositories/123"`)
}

func TestRoundTripperLogsRPCClientMethod(t *testing.T) {
	var logs bytes.Buffer
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	rt := httplogger.NewRoundTripper(&logger.Logger{Logger: slogLogger}, &test.StatusRoundTripper{Status: http.StatusNoContent})
	httpClient := test.NewContentClient(client.WithRoundTripper(rt))
	rpcClient := rpc.NewClient("https://example.com/api/", httpClient)
	var res test.Response

	err := rpcClient.Post(t.Context(), "greet.v1.Greeter/SayHello?name=Bob", &test.Request{Name: "Bob"}, &res)

	require.NoError(t, err)
	require.Contains(t, logs.String(), `"msg":"http: greet.v1.Greeter/SayHello"`)
	require.Contains(t, logs.String(), `"service":"greet.v1.Greeter"`)
	require.Contains(t, logs.String(), `"method":"SayHello"`)
	require.Contains(t, logs.String(), `"http.request.method":"post"`)
	require.Contains(t, logs.String(), `"url.path":"/api/greet.v1.Greeter/SayHello"`)
}
