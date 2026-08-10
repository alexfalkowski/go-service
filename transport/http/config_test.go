package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/stretchr/testify/require"
)

func TestConfigResolvesUnaryAndStreamingTimeoutsIndependently(t *testing.T) {
	cfg := &http.Config{
		Config:  &server.Config{Options: options.Map{"read_timeout": "10s", "write_timeout": "20s"}},
		Timeout: 5 * time.Second,
	}

	require.Equal(t, 5*time.Second, cfg.GetTimeout())
	require.Equal(t, 10*time.Second, cfg.GetReadTimeout())
	require.Equal(t, 20*time.Second, cfg.GetWriteTimeout())
}

func TestConfigDefaultsTimeoutsIndependently(t *testing.T) {
	var cfg *http.Config

	require.Equal(t, time.DefaultTimeout, cfg.GetTimeout())
	require.Equal(t, time.DefaultTimeout, cfg.GetReadTimeout())
	require.Equal(t, time.DefaultTimeout, cfg.GetWriteTimeout())
}

func TestConfigRejectsNegativeUnaryTimeout(t *testing.T) {
	cfg := &http.Config{Config: &server.Config{}, Timeout: -time.Second}

	require.Error(t, test.Validator.Struct(cfg))
}
