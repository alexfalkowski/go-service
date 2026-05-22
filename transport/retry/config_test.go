package retry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/retry"
	"github.com/stretchr/testify/require"
)

func TestConfigRejectsNegativeDurations(t *testing.T) {
	tests := []struct {
		cfg  *retry.Config
		name string
	}{
		{name: "timeout", cfg: &retry.Config{Timeout: -time.Second}},
		{name: "backoff", cfg: &retry.Config{Backoff: -time.Second}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, test.Validator.Struct(tt.cfg))
		})
	}
}

func TestConfigGetTimeout(t *testing.T) {
	tests := []struct {
		cfg  *retry.Config
		name string
		want time.Duration
	}{
		{name: "nil", want: time.DefaultTimeout},
		{name: "zero", cfg: &retry.Config{}, want: time.DefaultTimeout},
		{name: "negative", cfg: &retry.Config{Timeout: -time.Second}, want: time.DefaultTimeout},
		{name: "explicit", cfg: &retry.Config{Timeout: time.Second}, want: time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cfg.GetTimeout())
		})
	}
}

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
