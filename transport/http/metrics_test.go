package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/stretchr/testify/require"
)

func TestPrometheusHTTP(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("prometheus"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
	world.Register()

	_, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()

	ctx, cancel := test.Timeout()
	defer cancel()

	header := http.Header{}
	url := world.NamedServerURL("http", "metrics")

	res, body, err := world.ResponseWithBody(ctx, url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Contains(t, body, "go_info")
	require.Contains(t, body, `db_system="redis"`)
	require.Contains(t, body, `db_system_name="postgresql"`)
	require.Contains(t, body, "system")
	require.Contains(t, body, "process")
	require.Contains(t, body, "runtime")

	world.RequireStop()
}

func TestPrometheusAuthHTTP(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

	world := test.NewWorld(t,
		test.WithWorldTelemetry("prometheus"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
		test.WithWorldToken(tkn, tkn),
		test.WithWorldHTTP(),
	)
	world.Register()

	_, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()

	header := http.Header{}
	url := world.NamedServerURL("http", "metrics")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Contains(t, body, "go_info")
	require.Contains(t, body, `db_system="redis"`)
	require.Contains(t, body, `db_system_name="pg"`)
	require.Contains(t, body, "system")
	require.Contains(t, body, "process")
	require.Contains(t, body, "runtime")

	world.RequireStop()
}
