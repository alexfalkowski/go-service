package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/id/uuid"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/time"
	"github.com/alexfalkowski/go-service/v3/transport/http"
	"github.com/alexfalkowski/go-service/v3/transport/http/token"
	"github.com/alexfalkowski/go-service/v3/transport/limiter"
	"github.com/alexfalkowski/go-signal/v2/signaltest"
	"github.com/stretchr/testify/require"
)

func TestNewTokenUsesDefaultConfiguration(t *testing.T) {
	tkn := http.NewToken(test.NewHTTPTransportConfig(), nil, nil)
	require.Nil(t, tkn)
}

func TestNewTokenWithoutSharedServerConfig(t *testing.T) {
	tkn := http.NewToken(&http.Config{Timeout: time.Second}, nil, nil)
	require.Nil(t, tkn)
}

func TestNewTokenUsesConfiguredTokenSettings(t *testing.T) {
	cfg := test.NewHTTPTransportConfig()
	cfg.Token = test.NewToken("jwt")
	gen := uuid.NewGenerator()

	tkn := http.NewToken(cfg, test.FS, gen)

	require.NotNil(t, tkn)
	require.NotNil(t, token.NewGenerator(tkn))
	require.NotNil(t, token.NewVerifier(tkn))
}

func TestNewServerLimiter(t *testing.T) {
	cfg := test.NewHTTPTransportConfig()
	cfg.Limiter = test.NewLimiterConfig("user-agent", "1s", 1)

	server, err := http.NewServerLimiter(signaltest.NewLifecycle(t), limiter.NewKeyMap(), cfg)

	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestNewServerLimiterWithoutSharedServerConfig(t *testing.T) {
	server, err := http.NewServerLimiter(signaltest.NewLifecycle(t), limiter.NewKeyMap(), &http.Config{Timeout: time.Second})

	require.NoError(t, err)
	require.Nil(t, server)
}

func TestNewServerLimiterDisabled(t *testing.T) {
	server, err := http.NewServerLimiter(signaltest.NewLifecycle(t), limiter.NewKeyMap(), nil)

	require.NoError(t, err)
	require.Nil(t, server)
}
