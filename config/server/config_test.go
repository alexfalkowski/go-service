package server_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetMaxReceiveBytes(t *testing.T) {
	tests := []struct {
		cfg  *server.Config
		name string
		want int64
	}{
		{name: "nil", want: server.DefaultMaxReceiveBytes},
		{name: "zero", cfg: &server.Config{}, want: server.DefaultMaxReceiveBytes},
		{name: "explicit", cfg: &server.Config{MaxReceiveBytes: 64}, want: 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cfg.GetMaxReceiveBytes())
		})
	}
}

func TestConfigRejectsNegativeMaxReceiveBytes(t *testing.T) {
	cfg := &server.Config{MaxReceiveBytes: -1}
	require.Error(t, test.Validator.Struct(cfg))
}
