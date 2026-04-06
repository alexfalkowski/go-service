package retry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/stretchr/testify/require"
)

func TestConfigMaxAttempts(t *testing.T) {
	tests := []struct {
		name      string
		attempts  uint64
		nilConfig bool
		want      uint64
	}{
		{name: "nil", nilConfig: true, want: 0},
		{name: "zero", want: 0},
		{name: "one", attempts: 1, want: 1},
		{name: "three", attempts: 3, want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg *retry.Config
			if !tt.nilConfig {
				cfg = &retry.Config{Attempts: tt.attempts}
			}

			require.Equal(t, tt.want, cfg.MaxAttempts())
		})
	}
}

func TestConfigMaxRetries(t *testing.T) {
	tests := []struct {
		name      string
		attempts  uint64
		nilConfig bool
		want      uint64
	}{
		{name: "nil", nilConfig: true, want: 0},
		{name: "zero", want: 0},
		{name: "one", attempts: 1, want: 0},
		{name: "two", attempts: 2, want: 1},
		{name: "three", attempts: 3, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg *retry.Config
			if !tt.nilConfig {
				cfg = &retry.Config{Attempts: tt.attempts}
			}

			require.Equal(t, tt.want, cfg.MaxRetries())
		})
	}
}
