package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/stretchr/testify/require"
)

func TestConfigResolvesUnaryTimeout(t *testing.T) {
	cfg := &grpc.Config{Config: &server.Config{}, Timeout: 5 * time.Second}

	require.Equal(t, 5*time.Second, cfg.GetTimeout())
}

func TestConfigDefaultsUnaryTimeout(t *testing.T) {
	var cfg *grpc.Config

	require.Equal(t, time.DefaultTimeout, cfg.GetTimeout())
}

func TestConfigRejectsNegativeUnaryTimeout(t *testing.T) {
	cfg := &grpc.Config{Config: &server.Config{}, Timeout: -time.Second}

	require.Error(t, test.Validator.Struct(cfg))
}
