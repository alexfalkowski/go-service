package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestNewTokenWithoutSharedServerConfig(t *testing.T) {
	tkn := grpc.NewToken(&grpc.Config{Timeout: time.Second}, nil, nil)
	require.Nil(t, tkn)
}

func TestNewServerLimiterWithoutSharedServerConfig(t *testing.T) {
	server, err := grpc.NewServerLimiter(fxtest.NewLifecycle(t), limiter.NewKeyMap(), &grpc.Config{Timeout: time.Second}, nil)

	require.NoError(t, err)
	require.Nil(t, server)
}

func TestNewTokenUsesDefaultConfiguration(t *testing.T) {
	tkn := grpc.NewToken(test.NewGRPCTransportConfig(), nil, nil)
	require.Nil(t, tkn)
}
