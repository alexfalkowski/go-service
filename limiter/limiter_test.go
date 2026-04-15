package limiter_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestValidLimiter(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	m := limiter.KeyMap{"user-agent": meta.UserAgent}
	config := &limiter.Config{Kind: "user-agent", Tokens: 0, Interval: time.Second}

	limiter, err := limiter.NewLimiter(lc, m, config)
	require.NoError(t, err)
	require.NotNil(t, limiter)

	require.NoError(t, limiter.Close(t.Context()))
}

func TestMissingLimiter(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	m := limiter.KeyMap{}
	config := &limiter.Config{Kind: "user-agent", Tokens: 0, Interval: time.Second}

	_, err := limiter.NewLimiter(lc, m, config)
	require.Error(t, err)
}

func TestDisabledLimiter(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	m := limiter.KeyMap{"user-agent": meta.UserAgent}

	limiter, err := limiter.NewLimiter(lc, m, nil)
	require.NoError(t, err)
	require.Nil(t, limiter)
}

func TestClosedLimiter(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	m := limiter.KeyMap{"user-agent": meta.UserAgent}
	config := &limiter.Config{Kind: "user-agent", Tokens: 0, Interval: time.Second}

	limiter, err := limiter.NewLimiter(lc, m, config)
	require.NoError(t, err)
	require.NotNil(t, limiter)

	require.NoError(t, limiter.Close(t.Context()))

	_, _, err = limiter.Take(t.Context())
	require.Error(t, err)
}
