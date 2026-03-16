package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/stretchr/testify/require"
)

func TestNewControllerWithoutTokenConfig(t *testing.T) {
	controller, err := grpc.NewController(test.NewGRPCTransportConfig())
	require.NoError(t, err)
	require.Nil(t, controller)
}

func TestNewTokenWithoutTokenConfig(t *testing.T) {
	tkn := grpc.NewToken(test.Name, test.NewGRPCTransportConfig(), nil, nil, nil, nil)
	require.Nil(t, tkn)
}
