package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestNewTokenWithoutTokenConfig(t *testing.T) {
	tkn := http.NewToken(test.Name, test.NewHTTPTransportConfig(), nil, nil)
	require.Nil(t, tkn)
}

func TestNewTokenWithTokenConfig(t *testing.T) {
	cfg := test.NewHTTPTransportConfig()
	cfg.Token = test.NewToken("jwt")
	gen := uuid.NewGenerator()

	tkn := http.NewToken(test.Name, cfg, test.FS, gen)

	require.NotNil(t, tkn)
	require.NotNil(t, token.NewGenerator(tkn))
	require.NotNil(t, token.NewVerifier(tkn))
}

func TestNewServerLimiter(t *testing.T) {
	cfg := test.NewHTTPTransportConfig()
	cfg.Limiter = test.NewLimiterConfig("user-agent", "1s", 1)

	server, err := http.NewServerLimiter(fxtest.NewLifecycle(t), limiter.NewKeyMap(), cfg)

	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestNewServerLimiterDisabled(t *testing.T) {
	server, err := http.NewServerLimiter(fxtest.NewLifecycle(t), limiter.NewKeyMap(), nil)

	require.NoError(t, err)
	require.Nil(t, server)
}
