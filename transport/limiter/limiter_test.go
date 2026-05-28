package limiter_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestConfigRejectsNegativeInterval(t *testing.T) {
	cfg := &limiter.Config{Interval: -time.Second}
	require.Error(t, test.Validator.Struct(cfg))
}

func TestValidLimiter(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	m := limiter.KeyMap{"user-agent": meta.UserAgent}
	config := &limiter.Config{Kind: "user-agent", Tokens: 0, Interval: time.Second}

	limiter, err := limiter.NewLimiter(lc, m, config)
	require.NoError(t, err)
	require.NotNil(t, limiter)

	require.NoError(t, limiter.Close(t.Context()))
}

func TestTake(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	m := limiter.KeyMap{"user-agent": meta.UserAgent}
	config := &limiter.Config{Kind: "user-agent", Tokens: 2, Interval: time.Second}

	limiter, err := limiter.NewLimiter(lc, m, config)
	require.NoError(t, err)
	require.NotNil(t, limiter)
	defer func() {
		require.NoError(t, limiter.Close(t.Context()))
	}()

	ctx := meta.WithAttributes(t.Context(), meta.WithUserAgent(meta.String("test-agent")))

	ok, header, err := limiter.Take(ctx)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "limit=2, remaining=1", header)
}

func TestNewKeyMapIncludesServiceMethod(t *testing.T) {
	key := limiter.NewKeyMap()["service-method"]
	ctx := meta.WithAttributes(t.Context(), meta.WithServiceMethod(meta.String("GET /users/{id}")))

	require.Equal(t, meta.String("GET /users/{id}"), key(ctx))
}

func TestMissingLimiter(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	m := limiter.KeyMap{}
	config := &limiter.Config{Kind: "user-agent", Tokens: 0, Interval: time.Second}

	_, err := limiter.NewLimiter(lc, m, config)
	require.ErrorIs(t, err, limiter.ErrMissingKey)
}

func TestNilKeyFuncLimiter(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	m := limiter.KeyMap{"user-agent": nil}
	config := &limiter.Config{Kind: "user-agent", Tokens: 0, Interval: time.Second}

	_, err := limiter.NewLimiter(lc, m, config)
	require.ErrorIs(t, err, limiter.ErrMissingKey)
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
