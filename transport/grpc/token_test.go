package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/time"
	"github.com/alexfalkowski/go-service/v3/transport/grpc"
	"github.com/alexfalkowski/go-service/v3/transport/limiter"
	"github.com/alexfalkowski/go-signal/v2/signaltest"
	"github.com/stretchr/testify/require"
)

func TestNewTokenWithoutSharedServerConfig(t *testing.T) {
	tkn := grpc.NewToken(&grpc.Config{Timeout: time.Second}, nil, nil)
	require.Nil(t, tkn)
}

func TestNewServerLimiterWithoutSharedServerConfig(t *testing.T) {
	server, err := grpc.NewServerLimiter(signaltest.NewLifecycle(t), limiter.NewKeyMap(), &grpc.Config{Timeout: time.Second}, nil)

	require.NoError(t, err)
	require.Nil(t, server)
}

func TestNewTokenUsesDefaultConfiguration(t *testing.T) {
	tkn := grpc.NewToken(test.NewGRPCTransportConfig(), nil, nil)
	require.Nil(t, tkn)
}
