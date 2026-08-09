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
		{name: "max backoff", cfg: &retry.Config{MaxBackoff: -time.Second}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, test.Validator.Struct(tt.cfg))
		})
	}
}

func TestConfigRejectsAttemptsAboveMax(t *testing.T) {
	require.NoError(t, test.Validator.Struct(&retry.Config{Attempts: retry.MaxAttempts}))
	require.Error(t, test.Validator.Struct(&retry.Config{Attempts: retry.MaxAttempts + 1}))
}

func TestConfigGetDurations(t *testing.T) {
	tests := []struct {
		cfg            *retry.Config
		name           string
		wantTimeout    time.Duration
		wantBackoff    time.Duration
		wantMaxBackoff time.Duration
	}{
		{name: "nil", wantTimeout: time.DefaultTimeout, wantBackoff: retry.DefaultBackoff},
		{name: "zero", cfg: &retry.Config{}, wantTimeout: time.DefaultTimeout, wantBackoff: retry.DefaultBackoff},
		{
			name:        "negative timeout",
			cfg:         &retry.Config{Timeout: -time.Second},
			wantTimeout: time.DefaultTimeout,
			wantBackoff: retry.DefaultBackoff,
		},
		{
			name:           "negative backoff",
			cfg:            &retry.Config{Backoff: -time.Second},
			wantTimeout:    time.DefaultTimeout,
			wantBackoff:    retry.DefaultBackoff,
			wantMaxBackoff: 0,
		},
		{
			name:           "negative max backoff",
			cfg:            &retry.Config{MaxBackoff: -time.Second},
			wantTimeout:    time.DefaultTimeout,
			wantBackoff:    retry.DefaultBackoff,
			wantMaxBackoff: 0,
		},
		{
			name:           "explicit timeout",
			cfg:            &retry.Config{Timeout: time.Second},
			wantTimeout:    time.Second,
			wantBackoff:    retry.DefaultBackoff,
			wantMaxBackoff: 0,
		},
		{
			name:           "explicit backoff",
			cfg:            &retry.Config{Backoff: time.Second},
			wantTimeout:    time.DefaultTimeout,
			wantBackoff:    time.Second,
			wantMaxBackoff: 0,
		},
		{
			name:           "explicit max backoff",
			cfg:            &retry.Config{MaxBackoff: time.Second},
			wantTimeout:    time.DefaultTimeout,
			wantBackoff:    retry.DefaultBackoff,
			wantMaxBackoff: time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantTimeout, tt.cfg.GetTimeout())
			require.Equal(t, tt.wantBackoff, tt.cfg.GetBackoff())
			require.Equal(t, tt.wantMaxBackoff, tt.cfg.GetMaxBackoff())
		})
	}
}

func TestConfigRejectsInvalidStrategy(t *testing.T) {
	for _, strategy := range []string{"", "constant", "exponential", "fibonacci"} {
		require.NoError(t, test.Validator.Struct(&retry.Config{Strategy: strategy}))
	}

	require.Error(t, test.Validator.Struct(&retry.Config{Strategy: "bogus"}))
}

func TestConfigGetStrategy(t *testing.T) {
	tests := []struct {
		cfg  *retry.Config
		name string
		want string
	}{
		{name: "nil", want: "constant"},
		{name: "empty", cfg: &retry.Config{}, want: "constant"},
		{name: "explicit", cfg: &retry.Config{Strategy: "exponential"}, want: "exponential"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cfg.GetStrategy())
		})
	}
}

func TestConfigAttemptsReturnsConfiguredOrDefaultValue(t *testing.T) {
	tests := []struct {
		cfg          *retry.Config
		name         string
		wantAttempts uint64
		wantRetries  uint64
	}{
		{name: "nil"},
		{name: "zero", cfg: &retry.Config{}},
		{name: "one", cfg: &retry.Config{Attempts: 1}, wantAttempts: 1},
		{name: "two", cfg: &retry.Config{Attempts: 2}, wantAttempts: 2, wantRetries: 1},
		{name: "three", cfg: &retry.Config{Attempts: 3}, wantAttempts: 3, wantRetries: 2},
		{name: "above max", cfg: &retry.Config{Attempts: retry.MaxAttempts + 1}, wantAttempts: retry.MaxAttempts, wantRetries: retry.MaxAttempts - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantAttempts, tt.cfg.MaxAttempts())
			require.Equal(t, tt.wantRetries, tt.cfg.MaxRetries())
		})
	}
}
