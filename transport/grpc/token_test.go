package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/stretchr/testify/require"
)

func TestNewTokenWithoutTokenConfig(t *testing.T) {
	tkn := grpc.NewToken(test.Name, test.NewGRPCTransportConfig(), nil, nil)
	require.Nil(t, tkn)
}
