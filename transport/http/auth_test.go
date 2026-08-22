package http_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/stretchr/testify/require"
)

func TestUnaryHandlerAuthenticatesSupportedToken(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			gen := uuid.NewGenerator()
			tkn := token.NewToken(cfg, test.FS, gen)

			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(tkn, tkn), test.WithWorldHTTP())

			world.RPCServer.Route("/hello", test.SuccessSayHello)

			header := http.Header{}
			header.Set(http.ContentTypeKey, media.JSON)
			header.Set("Request-Id", "test")
			header.Set("X-Forwarded-For", "127.0.0.1")
			header.Set("Geolocation", "geo:47,11")

			url := world.PathServerURL("http", "hello")

			res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.NotEmpty(t, body)
		})
	}
}

func TestUnaryHandlerRejectsUnknownTokenKind(t *testing.T) {
	cfg := test.NewToken("none")
	tkn := token.NewToken(cfg, test.FS, nil)

	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("test", nil), tkn), test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Equal(t, "http: unauthorized", body)
}

func TestUnaryHandlerAcceptsValidAuthentication(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")), test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set("Request-Id", "test")
	header.Set("X-Forwarded-For", "127.0.0.1")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NotEmpty(t, body)
}

func TestUnaryHandlerRejectsDeniedAccess(t *testing.T) {
	controller, err := access.NewController(&access.Config{
		Model:  test.FilePath("configs/rbac.conf"),
		Policy: "p, " + test.UserID.String() + ", http:POST /other, invoke",
	}, test.FS)
	require.NoError(t, err)

	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")),
		test.WithWorldAccessController(controller),
		test.WithWorldHTTP(),
	)

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, res.StatusCode)
	require.Equal(t, "http: forbidden", body)
}

func TestAuthDoesNotBypassApplicationMetricsPath(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP())

	world.HandleRoute("GET /admin/metrics", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		_, _ = res.Write([]byte("secret"))
	}))
	world.Start()

	res, body, err := world.ResponseWithBody(t.Context(), world.PathServerURL("http", "admin/metrics"), http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Equal(t, "http: unauthorized", body)
}

func TestRejectsInvalidAuthUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")), test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Equal(t, "http: unauthorized", body)
}

func TestAuthUnaryWithAppend(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set("Request-Id", "test")
	header.Set("Authorization", "What Invalid")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.NotEmpty(t, body)
}

func TestUnaryHandlerAcceptsLowercaseBearerScheme(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set("Request-Id", "test")
	header.Set("Authorization", "bearer test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NotEmpty(t, body)
}

func TestUnaryHandlerRejectsMissingAuthentication(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Equal(t, "http: unauthorized", body)
}

func TestUnaryHandlerRejectsEmptyAuthentication(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator(strings.Empty, nil), test.NewVerifier("test")), test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "authorization is invalid")
}

func TestUnaryHandlerRejectsMissingClientAuthentication(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Equal(t, "http: unauthorized", body)
}

func TestUnaryHandlerPropagatesTokenError(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, test.ErrGenerate), test.NewVerifier("test")),
		test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "token: generation issue")
}

func TestUnaryHandlerReturnsBreakerError(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, test.ErrGenerate), test.NewVerifier("test")),
		test.WithWorldHTTP(),
	)

	var err error
	url := world.PathServerURL("http", "hello")

	for i := range 10 {
		t.Run("attempt-"+strconv.Itoa(i+1), func(t *testing.T) {
			header := http.Header{}
			header.Set(http.ContentTypeKey, media.JSON)
			header.Set("Request-Id", "test")

			_, _, err = world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
		})
	}
	require.Error(t, err)
}
