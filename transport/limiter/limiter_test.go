package limiter_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
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
	config := &limiter.Config{Kind: "user-agent", Tokens: 1, Interval: time.Second}

	limiter, err := limiter.NewLimiter(lc, m, config)
	require.NoError(t, err)
	require.NotNil(t, limiter)
	defer func() {
		require.NoError(t, limiter.Close(t.Context()))
	}()

	first := meta.WithAttributes(t.Context(), meta.WithUserAgent(meta.String("first-agent")))
	second := meta.WithAttributes(t.Context(), meta.WithUserAgent(meta.String("second-agent")))

	ok, header, err := limiter.Take(first)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "limit=1, remaining=0", header)

	ok, header, err = limiter.Take(first)
	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, "limit=1, remaining=0", header)

	ok, header, err = limiter.Take(second)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "limit=1, remaining=0", header)
}

func TestNewKeyMap(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want meta.Value
	}{
		{
			name: "user-agent",
			ctx:  meta.WithAttributes(t.Context(), meta.WithUserAgent(meta.String("test-agent"))),
			want: meta.String("test-agent"),
		},
		{
			name: "ip",
			ctx:  meta.WithAttributes(t.Context(), meta.WithIPAddr(meta.String("192.0.2.1"))),
			want: meta.String("192.0.2.1"),
		},
		{
			name: "user-id",
			ctx:  meta.WithAttributes(t.Context(), meta.WithUserID(meta.String("test-user"))),
			want: meta.String("test-user"),
		},
		{
			name: "service-method",
			ctx:  meta.WithAttributes(t.Context(), meta.WithServiceMethod(meta.String("GET /users/{id}"))),
			want: meta.String("GET /users/{id}"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := limiter.NewKeyMap()[tt.name]
			require.NotNil(t, key)
			require.Equal(t, tt.want, key(tt.ctx))
		})
	}
}

func TestLifecycleStopsLimiter(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	m := limiter.KeyMap{"user-agent": meta.UserAgent}
	config := &limiter.Config{Kind: "user-agent", Tokens: 1, Interval: time.Second}

	limiter, err := limiter.NewLimiter(lc, m, config)
	require.NoError(t, err)
	require.NotNil(t, limiter)

	lc.RequireStart()
	lc.RequireStop()

	_, _, err = limiter.Take(t.Context())
	require.Error(t, err)
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
