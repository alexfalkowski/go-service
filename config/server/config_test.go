package server_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetMaxReceiveSize(t *testing.T) {
	tests := []struct {
		cfg  *server.Config
		name string
		want bytes.Size
	}{
		{name: "nil", want: server.DefaultMaxReceiveSize},
		{name: "zero", cfg: &server.Config{}, want: server.DefaultMaxReceiveSize},
		{name: "explicit", cfg: &server.Config{MaxReceiveSize: 64}, want: 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cfg.GetMaxReceiveSize())
		})
	}
}

func TestConfigRejectsNegativeMaxReceiveSize(t *testing.T) {
	cfg := &server.Config{MaxReceiveSize: -1}
	require.Error(t, test.Validator.Struct(cfg))
}
