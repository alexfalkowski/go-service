package cache_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	drivererrors "github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	compresserrors "github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/ptr"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestValidCache(t *testing.T) {
	for _, tt := range cacheRoundTripCases() {
		t.Run(tt.name, func(t *testing.T) {
			test.RequireCacheRoundTrip(t, tt.config, tt.persist(), tt.get())
		})
	}
}

func TestGenericValidCache(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	require.NoError(t, cache.Persist(t.Context(), "test", new("hello?"), time.Minute))

	value, ok, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "hello?", *value)

	require.NoError(t, world.Remove(t.Context(), "test"))
}

func TestGetOrPersist(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	calls := 0
	fn := func() (string, error) {
		calls++

		return "hello?", nil
	}

	value, err := cache.GetOrPersist(t.Context(), "test", time.Minute, fn)
	require.NoError(t, err)
	require.Equal(t, "hello?", *value)
	require.Equal(t, 1, calls)

	value, err = cache.GetOrPersist(t.Context(), "test", time.Minute, fn)
	require.NoError(t, err)
	require.Equal(t, "hello?", *value)
	require.Equal(t, 1, calls)

	require.NoError(t, world.Remove(t.Context(), "test"))
}

func TestGetOrPersistRedis(t *testing.T) {
	cfg := test.NewCacheConfig("redis", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())
	ctx := t.Context()
	require.NoError(t, world.Remove(ctx, "test"))

	calls := 0
	fn := func() (string, error) {
		calls++

		return "hello?", nil
	}

	value, err := cache.GetOrPersist(ctx, "test", time.Minute, fn)
	require.NoError(t, err)
	require.Equal(t, "hello?", *value)
	require.Equal(t, 1, calls)

	value, err = cache.GetOrPersist(ctx, "test", time.Minute, fn)
	require.NoError(t, err)
	require.Equal(t, "hello?", *value)
	require.Equal(t, 1, calls)

	require.NoError(t, world.Remove(ctx, "test"))
}

func TestGetOrPersistSingleWinner(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())
	ctx := t.Context()

	const goroutines = 16
	var group sync.WaitGroup
	values := make([]string, goroutines)
	errs := make([]error, goroutines)

	group.Add(goroutines)
	for i := range goroutines {
		go func(i int) {
			defer group.Done()

			value, err := cache.GetOrPersist(ctx, "test", time.Minute, func() (string, error) {
				return fmt.Sprintf("value-%d", i), nil
			})
			errs[i] = err
			if err == nil {
				values[i] = *value
			}
		}(i)
	}
	group.Wait()

	for i := range goroutines {
		require.NoError(t, errs[i])
		require.Equal(t, values[0], values[i])
	}
	require.NotEmpty(t, values[0])

	require.NoError(t, world.Remove(ctx, "test"))
}

func TestGetOrPersistReleasesWaiterOnItsOwnCancellation(t *testing.T) {
	drv := &blockingGetDriver{gets: make(chan string, 16)}
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldCacheDriver(drv))

	leaderStarted := make(chan struct{})
	release := make(chan struct{})
	leaderErr := make(chan error, 1)
	leaderValue := "unchanged"
	go func() {
		leaderErr <- world.GetOrPersist(context.Background(), "test", &leaderValue, time.Minute, func() error {
			close(leaderStarted)
			<-release
			leaderValue = "leader"

			return nil
		})
	}()

	<-leaderStarted
	<-drv.gets // leader's initial Get
	<-drv.gets // leader's single-flight recheck Get

	dupCtx, dupCancel := context.WithCancel(context.Background())
	dupErr := make(chan error, 1)
	dupValue := "unchanged"
	go func() {
		dupErr <- world.GetOrPersist(dupCtx, "test", &dupValue, time.Minute, func() error {
			dupValue = "duplicate"

			return nil
		})
	}()

	<-drv.gets // duplicate missed and joined the in-flight fill
	dupCancel()

	select {
	case err := <-dupErr:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("duplicate blocked on the shared fill instead of returning on its own cancellation")
	}
	require.Equal(t, "unchanged", dupValue)

	close(release)
	require.NoError(t, <-leaderErr)
	require.Equal(t, "leader", leaderValue)
}

func TestGetOrPersistLeaderCancellationFailsWaiters(t *testing.T) {
	drv := &blockingGetDriver{gets: make(chan string, 16)}
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldCacheDriver(drv))

	leaderStarted := make(chan struct{})
	release := make(chan struct{})
	leaderCtx, leaderCancel := context.WithCancel(context.Background())
	leaderErr := make(chan error, 1)
	leaderValue := "unchanged"
	go func() {
		leaderErr <- world.GetOrPersist(leaderCtx, "test", &leaderValue, time.Minute, func() error {
			close(leaderStarted)
			<-release
			leaderValue = "leader"

			return nil
		})
	}()

	<-leaderStarted
	<-drv.gets
	<-drv.gets

	followerCalled := make(chan struct{}, 1)
	followerErr := make(chan error, 1)
	followerValue := "unchanged"
	go func() {
		followerErr <- world.GetOrPersist(context.Background(), "test", &followerValue, time.Minute, func() error {
			followerCalled <- struct{}{}

			return nil
		})
	}()

	<-drv.gets // follower missed and is joining the shared fill
	// Give the follower time to park on the shared fill before it publishes.
	time.Sleep(100 * time.Millisecond)

	// The shared fill runs under the leader context, so canceling the leader and
	// then letting the fill publish fails the still-live follower too.
	leaderCancel()
	close(release)

	require.ErrorIs(t, <-followerErr, context.Canceled)
	require.ErrorIs(t, <-leaderErr, context.Canceled)

	select {
	case <-followerCalled:
		t.Fatal("follower ran its own loader instead of sharing the leader's fill")
	default:
	}
}

func TestGenericGetOrPersistDisabledCache(t *testing.T) {
	cache.Register(nil)
	t.Cleanup(func() {
		cache.Register(nil)
	})

	called := false
	value, err := cache.GetOrPersist(t.Context(), "test", time.Minute, func() (string, error) {
		called = true

		return "hello?", nil
	})
	require.NoError(t, err)
	require.Nil(t, value)
	require.False(t, called)
}

func TestGetOrPersistLoaderError(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())
	ctx := t.Context()

	_, err := cache.GetOrPersist(ctx, "test", time.Minute, func() (string, error) {
		return strings.Empty, test.ErrFailed
	})
	require.ErrorIs(t, err, test.ErrFailed)

	value, ok, err := cache.Get[string](ctx, "test")
	require.NoError(t, err)
	require.False(t, ok)
	require.Empty(t, *value)

	require.NoError(t, world.Remove(ctx, "test"))
}

func TestGetOrPersistPublish(t *testing.T) {
	winner := base64.Encode(strings.Bytes(`"winner"`))
	tests := []struct {
		driver  driver.Driver
		name    string
		want    string
		wantErr bool
	}{
		{name: "returns existing winner", driver: getOrSaveCacheDriver{existing: winner, loaded: true}, want: "winner"},
		{name: "propagates save error", driver: getOrSaveCacheDriver{err: test.ErrFailed}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
			world := test.NewStartedWorld(t,
				test.WithWorldCacheConfig(cfg),
				test.WithWorldCacheDriver(tt.driver),
			)

			value := "unchanged"
			called := false
			err := world.GetOrPersist(t.Context(), "test", &value, time.Minute, func() error {
				called = true
				value = "produced"

				return nil
			})

			require.True(t, called)
			if tt.wantErr {
				require.Error(t, err)

				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, value)
		})
	}
}

func TestGetOrPersistEncodeError(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "error", "redis")
	test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	_, err := cache.GetOrPersist(t.Context(), "test", time.Minute, func() (string, error) {
		return "hello?", nil
	})
	require.Error(t, err)
}

func TestGetOrPersistDriverGetError(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	world := test.NewStartedWorld(t,
		test.WithWorldCacheConfig(cfg),
		test.WithWorldCacheDriver(&test.ErrCache{}),
	)

	value := "unchanged"
	called := false
	err := world.GetOrPersist(t.Context(), "test", &value, time.Minute, func() error {
		called = true

		return nil
	})
	require.ErrorIs(t, err, test.ErrFailed)
	require.False(t, called)
}

func TestGetOrPersistRecheck(t *testing.T) {
	winner := base64.Encode(strings.Bytes(`"winner"`))
	tests := []struct {
		driver  *recheckCacheDriver
		name    string
		wantErr bool
	}{
		{name: "recheck finds concurrently stored value", driver: &recheckCacheDriver{value: winner}},
		{name: "recheck surfaces backend error", driver: &recheckCacheDriver{recheck: test.ErrFailed}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
			world := test.NewStartedWorld(t,
				test.WithWorldCacheConfig(cfg),
				test.WithWorldCacheDriver(tt.driver),
			)

			value := "unchanged"
			called := false
			err := world.GetOrPersist(t.Context(), "test", &value, time.Minute, func() error {
				called = true

				return nil
			})

			require.False(t, called)
			if tt.wantErr {
				require.Error(t, err)

				return
			}
			require.NoError(t, err)
			require.Equal(t, "winner", value)
		})
	}
}

func TestModuleRegistersGenericCache(t *testing.T) {
	cache.Register(nil)
	t.Cleanup(func() {
		cache.Register(nil)
	})

	app := di.New(
		di.NoLogger,
		cache.Module,
		di.Constructor(func() *logger.Logger { return nil }),
		fx.Supply(
			test.NewCacheConfig("ttlcache", "none", "json", "redis"),
			test.FS,
			test.Encoder,
			test.Pool,
			test.Compressor,
		),
	)
	require.NoError(t, app.Err())

	require.NoError(t, app.Start(t.Context()))
	t.Cleanup(func() {
		require.NoError(t, app.Stop(context.Background()))
	})

	require.NoError(t, cache.Persist(t.Context(), "test", new("hello?"), time.Minute))

	value, ok, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "hello?", *value)
}

func TestGenericDisabledCache(t *testing.T) {
	cache.Register(nil)

	require.NoError(t, cache.Persist(t.Context(), "test", new("hello?"), time.Minute))

	value, ok, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, value)
}

func TestGenericGetValidEmptyCache(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	require.NoError(t, cache.Persist(t.Context(), "test", new(""), time.Minute))

	value, ok, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, strings.Empty, *value)

	require.NoError(t, world.Remove(t.Context(), "test"))
}

func TestGenericGetMissingCache(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	value, ok, err := cache.Get[string](t.Context(), "missing")
	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, strings.Empty, *value)

	require.NoError(t, world.Remove(t.Context(), "missing"))
}

func TestGenericGetDisabledCache(t *testing.T) {
	cache.Register(nil)
	t.Cleanup(func() {
		cache.Register(nil)
	})

	require.NoError(t, cache.Persist(t.Context(), "test", new("hello?"), time.Minute))

	value, ok, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, value)
}

func TestMaxSizeOnPersist(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	err := world.Persist(t.Context(), "test", new("hello?"), time.Minute)
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
}

func TestMaxSizeOnPersistCompressedValue(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	err := world.Persist(t.Context(), "test", new("a"), time.Minute)
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
}

func TestMaxSizeOnPersistStopsOversizedEncodedValue(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())
	value := &oversizedWriterTo{size: cfg.MaxSize.Bytes() + 1}

	err := world.Persist(t.Context(), "test", value, time.Minute)
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
	require.Equal(t, cfg.MaxSize.Bytes(), value.written)
}

func TestMaxSizeOnGet(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	require.NoError(t, world.Persist(t.Context(), "test", new("hello?"), time.Minute))

	cfg.MaxSize = 4

	_, err := world.Get(t.Context(), "test", ptr.Zero[string]())
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
}

func TestMaxSizeOnGetAllowsEncodedValueLargerThanLimit(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t,
		test.WithWorldCacheConfig(cfg),
		test.WithWorldCacheDriver(&test.Cache{Value: base64.Encode(strings.Bytes(`"ok"`))}),
	)

	value := ptr.Zero[string]()
	ok, err := world.Get(t.Context(), "test", value)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "ok", *value)
}

func TestMaxSizeOnGetRejectsEncodedValueTooLarge(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t,
		test.WithWorldCacheConfig(cfg),
		test.WithWorldCacheDriver(&test.Cache{Value: base64.Encode(make([]byte, 7))}),
	)

	_, err := world.Get(t.Context(), "test", ptr.Zero[string]())
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
}

func TestMaxSizeOnGetAllowsHugeConfiguredLimit(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	cfg.MaxSize = bytes.Size(math.MaxInt64)
	world := test.NewStartedWorld(t,
		test.WithWorldCacheConfig(cfg),
		test.WithWorldCacheDriver(&test.Cache{Value: base64.Encode(strings.Bytes(`"ok"`))}),
	)

	value := ptr.Zero[string]()
	ok, err := world.Get(t.Context(), "test", value)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "ok", *value)
}

func TestGetMissesAfterCacheFormatChange(t *testing.T) {
	tests := []struct {
		persist   *config.Config
		retrieval *config.Config
		name      string
	}{
		{
			persist:   test.NewCacheConfig("ttlcache", "none", "json", "redis"),
			retrieval: test.NewCacheConfig("ttlcache", "snappy", "json", "redis"),
			name:      "compressor",
		},
		{
			persist:   test.NewCacheConfig("ttlcache", "none", "json", "redis"),
			retrieval: test.NewCacheConfig("ttlcache", "none", "gob", "redis"),
			name:      "encoder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drv, err := driver.NewDriver(driver.DriverParams{Config: tt.persist})
			require.NoError(t, err)

			writer := test.NewStartedWorld(t,
				test.WithWorldCacheConfig(tt.persist),
				test.WithWorldCacheDriver(drv),
			)
			require.NoError(t, writer.Persist(t.Context(), "test", new("hello?"), time.Minute))

			reader := test.NewStartedWorld(t,
				test.WithWorldCacheConfig(tt.retrieval),
				test.WithWorldCacheDriver(drv),
			)
			value := "unchanged"
			ok, err := reader.Get(t.Context(), "test", &value)
			require.NoError(t, err)
			require.False(t, ok)
			require.Equal(t, "unchanged", value)
		})
	}
}

func TestExpiredCache(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())
	require.NoError(t, cache.Persist(t.Context(), "test", new("hello?"), time.Nanosecond))

	// Simulate expiry.
	time.Sleep(time.Second)

	_, ok, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.False(t, ok)

	require.NoError(t, world.Remove(t.Context(), "test"))
}

func TestErroneousCache(t *testing.T) {
	tests := []struct {
		config *config.Config
		name   string
	}{
		{name: "missing driver secret", config: test.NewCacheConfig("redis", "snappy", "json", "none")},
		{name: "invalid driver secret", config: test.NewCacheConfig("redis", "snappy", "json", "hooks")},
		{name: "unsupported cache kind", config: test.NewCacheConfig("test", "snappy", "json", "hooks")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := driver.NewDriver(driver.DriverParams{
				Lifecycle: fxtest.NewLifecycle(t),
				FS:        test.FS,
				Config:    tt.config,
			})
			require.Error(t, err)
		})
	}
}

func TestDisabledCache(t *testing.T) {
	t.Run("driver", func(t *testing.T) {
		_, err := driver.NewDriver(driver.DriverParams{
			Lifecycle: fxtest.NewLifecycle(t),
			FS:        test.FS,
		})
		require.NoError(t, err)
	})

	t.Run("cache", func(t *testing.T) {
		world := test.NewStartedWorld(t, test.WithWorldCacheConfig(nil), test.WithWorldRegisterCache())

		require.Nil(t, world.Cache)
	})
}

func TestErroneousSave(t *testing.T) {
	t.Run("invalid encoder", func(t *testing.T) {
		cfg := test.NewCacheConfig("ttlcache", "snappy", "error", "redis")
		test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())
		require.Error(t, cache.Persist(t.Context(), "test", new("test"), time.Minute))
	})

	t.Run("read from only falls back to configured encoder", func(t *testing.T) {
		cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
		world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg))
		require.NoError(t, world.Persist(t.Context(), "test", &test.ReadFromOnly{Name: "hello?"}, time.Minute))

		var get test.Request
		ok, err := world.Get(t.Context(), "test", &get)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "hello?", get.Name)
	})
}

func TestErroneousGet(t *testing.T) {
	tests := []struct {
		driver driver.Driver
		config *config.Config
		name   string
	}{
		{name: "decode error", config: test.NewCacheConfig("ttlcache", "snappy", "error", "redis"), driver: &test.Cache{Value: "d2hhdD8="}},
		{name: "decompress error", config: test.NewCacheConfig("ttlcache", "error", "json", "redis"), driver: &test.Cache{Value: "d2hhdD8="}},
		{name: "unmarshal error", config: test.NewCacheConfig("ttlcache", "snappy", "json", "redis"), driver: &test.Cache{Value: "what?"}},
		{name: "driver error", config: test.NewCacheConfig("ttlcache", "snappy", "json", "redis"), driver: &test.ErrCache{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldCacheConfig(tt.config),
				test.WithWorldCacheDriver(tt.driver),
			)

			_, err := world.Get(t.Context(), "test", ptr.Zero[string]())
			require.Error(t, err)
		})
	}
}

func TestWriteToOnlyUsesConfiguredEncoderOnGet(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg))
	require.NoError(t, world.Persist(t.Context(), "test", &test.Request{Name: "hello?"}, time.Minute))

	get := &test.WriteToOnly{}
	ok, err := world.Get(t.Context(), "test", get)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "hello?", get.Name)
}

func TestGetValidEmptyStringCache(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg))

	require.NoError(t, world.Persist(t.Context(), "test", new(""), time.Minute))

	value := "existing"
	ok, err := world.Get(t.Context(), "test", &value)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, strings.Empty, value)
}

func TestGetValidEmptyBufferCache(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg))

	require.NoError(t, world.Persist(t.Context(), "test", bytes.NewBufferString(strings.Empty), time.Minute))

	value := bytes.NewBufferString("existing")
	ok, err := world.Get(t.Context(), "test", value)
	require.NoError(t, err)
	require.True(t, ok)
	require.Empty(t, value.Bytes())
}

func TestGetMissingCache(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &config.Config{Kind: "ttlcache", MaxEntries: config.DefaultMaxEntries}

	d, err := driver.NewDriver(driver.DriverParams{
		Lifecycle: lc,
		Config:    cfg,
	})
	require.NoError(t, err)

	params := cache.CacheParams{
		Config:     cfg,
		Compressor: test.Compressor,
		Encoder:    test.Encoder,
		Pool:       test.Pool,
		Driver:     d,
	}

	c := cache.NewCache(params)
	t.Cleanup(func() {
		require.NoError(t, c.Flush(context.Background()))
	})
	value := "existing"

	ok, err := c.Get(t.Context(), "missing", &value)
	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, "existing", value)
}

func TestGetExpiredCacheLeavesDestinationUnchanged(t *testing.T) {
	cfg := &config.Config{Kind: "ttlcache", MaxEntries: config.DefaultMaxEntries}
	c := cache.NewCache(cache.CacheParams{
		Config:     cfg,
		Compressor: test.Compressor,
		Encoder:    test.Encoder,
		Pool:       test.Pool,
		Driver:     expiredCacheDriver{},
	})
	value := "existing"

	ok, err := c.Get(t.Context(), "expired", &value)
	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, "existing", value)
}

func cacheRoundTripCases() []cacheRoundTripCase {
	compressors := []string{"none", "snappy", "s2", "zstd"}
	encoders := []string{"json", "hjson", "yaml", "yml", "toml", "gob", "msgpack"}
	cases := []cacheRoundTripCase{
		{
			name:    "redis/default/request",
			config:  test.NewCacheConfig("redis", "snappy", strings.Empty, "redis"),
			persist: func() any { return &test.Request{Name: "hello?"} },
			get:     func() any { return &test.Request{} },
		},
		{
			name:    "ttlcache/default/string",
			config:  test.NewCacheConfig("ttlcache", strings.Empty, strings.Empty, "redis"),
			persist: func() any { return new("hello?") },
			get:     func() any { return ptr.Zero[string]() },
		},
		{
			name:    "ttlcache/default/buffer",
			config:  test.NewCacheConfig("ttlcache", strings.Empty, strings.Empty, "redis"),
			persist: func() any { return bytes.NewBufferString("hello?") },
			get:     func() any { return &bytes.Buffer{} },
		},
		{
			name:    "ttlcache/default/protobuf",
			config:  test.NewCacheConfig("ttlcache", strings.Empty, strings.Empty, "redis"),
			persist: func() any { return &v1.SayHelloRequest{Name: "hello?"} },
			get:     func() any { return &v1.SayHelloRequest{} },
		},
		{
			name:    "ttlcache/unknown/unknown/request",
			config:  test.NewCacheConfig("ttlcache", "unknown", "unknown", "redis"),
			persist: func() any { return &test.Request{Name: "hello?"} },
			get:     func() any { return &test.Request{} },
		},
	}

	tests := make([]cacheRoundTripCase, 0, len(cases)+len(compressors)*len(encoders))
	tests = append(tests, cases...)

	for _, compressor := range compressors {
		for _, encoder := range encoders {
			tests = append(tests, cacheRoundTripCase{
				name:    "ttlcache/" + compressor + "/" + encoder + "/request",
				config:  test.NewCacheConfig("ttlcache", compressor, encoder, "redis"),
				persist: func() any { return &test.Request{Name: "hello?"} },
				get:     func() any { return &test.Request{} },
			})
		}
	}

	return tests
}

type cacheRoundTripCase struct {
	persist func() any
	get     func() any
	config  *config.Config
	name    string
}

type oversizedWriterTo struct {
	size    int64
	written int64
}

func (w *oversizedWriterTo) WriteTo(dst io.Writer) (int64, error) {
	n, err := dst.Write(make([]byte, w.size))
	w.written = int64(n)

	return w.written, err
}

type expiredCacheDriver struct{}

func (expiredCacheDriver) Delete(context.Context, string) error {
	return nil
}

func (expiredCacheDriver) Get(context.Context, string) (string, error) {
	return strings.Empty, drivererrors.ErrExpired
}

func (expiredCacheDriver) Flush(context.Context) error {
	return nil
}

func (expiredCacheDriver) Save(context.Context, string, string, time.Duration) error {
	return nil
}

func (expiredCacheDriver) GetOrSave(context.Context, string, string, time.Duration) (string, bool, error) {
	return strings.Empty, false, nil
}

// getOrSaveCacheDriver misses on Get and returns a configured result from GetOrSave, exercising the publish
// path of Cache.GetOrPersist without a live backend.
type getOrSaveCacheDriver struct {
	err      error
	existing string
	loaded   bool
}

func (getOrSaveCacheDriver) Delete(context.Context, string) error {
	return nil
}

func (getOrSaveCacheDriver) Get(context.Context, string) (string, error) {
	return strings.Empty, drivererrors.ErrMissing
}

func (getOrSaveCacheDriver) Flush(context.Context) error {
	return nil
}

func (getOrSaveCacheDriver) Save(context.Context, string, string, time.Duration) error {
	return nil
}

func (d getOrSaveCacheDriver) GetOrSave(context.Context, string, string, time.Duration) (string, bool, error) {
	return d.existing, d.loaded, d.err
}

// recheckCacheDriver misses the first Get, then hits or errors on later Gets, simulating a value stored or a
// backend failure between Cache.GetOrPersist's initial read and its single-flight re-check.
type recheckCacheDriver struct {
	recheck error
	value   string
	gets    int
}

func (*recheckCacheDriver) Delete(context.Context, string) error {
	return nil
}

func (d *recheckCacheDriver) Get(context.Context, string) (string, error) {
	d.gets++
	if d.gets == 1 {
		return strings.Empty, drivererrors.ErrMissing
	}
	if d.recheck != nil {
		return strings.Empty, d.recheck
	}

	return d.value, nil
}

func (*recheckCacheDriver) Flush(context.Context) error {
	return nil
}

func (*recheckCacheDriver) Save(context.Context, string, string, time.Duration) error {
	return nil
}

func (*recheckCacheDriver) GetOrSave(context.Context, string, string, time.Duration) (string, bool, error) {
	return strings.Empty, false, nil
}

// blockingGetDriver reports every Get on gets and always misses, letting a test observe when callers reach
// Cache.GetOrPersist's single-flight path while a loader is held blocked.
type blockingGetDriver struct {
	gets chan string
}

func (*blockingGetDriver) Delete(context.Context, string) error {
	return nil
}

func (d *blockingGetDriver) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return strings.Empty, err
	}
	d.gets <- key

	return strings.Empty, drivererrors.ErrMissing
}

func (*blockingGetDriver) Flush(context.Context) error {
	return nil
}

func (*blockingGetDriver) Save(context.Context, string, string, time.Duration) error {
	return nil
}

func (*blockingGetDriver) GetOrSave(ctx context.Context, _, _ string, _ time.Duration) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return strings.Empty, false, err
	}

	return strings.Empty, false, nil
}
