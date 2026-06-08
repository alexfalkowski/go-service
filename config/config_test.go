package config_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestValidFileConfig(t *testing.T) {
	files := []string{
		test.FilePath("configs/config.hjson"),
		test.FilePath("configs/config.toml"),
		test.FilePath("configs/config.yml"),
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			set := flag.NewFlagSet("test")
			set.AddConfig(file)

			decoder := test.NewDecoder(set)

			cfg, err := config.NewConfig[config.Config](decoder, test.Validator)
			require.NoError(t, err)
			verifyConfig(t, cfg)
		})
	}
}

func TestValidHomeFileConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	data, err := test.FS.ReadFile(test.Path("configs/config.yml"))
	require.NoError(t, err)
	require.NoError(t, test.FS.WriteFile(test.FS.Join(home, "config.yml"), data, 0o600))

	set := flag.NewFlagSet("test")
	set.AddConfig("file:~/config.yml")

	decoder := test.NewDecoder(set)

	cfg, err := config.NewConfig[config.Config](decoder, test.Validator)
	require.NoError(t, err)
	verifyConfig(t, cfg)
}

func TestInvalidFileConfig(t *testing.T) {
	files := []string{
		test.FilePath("configs/invalid.yml"),
		test.FilePath("configs/invalid_trace.yml"),
		test.FilePath("configs/missing.yml"),
		test.FilePath("configs/script.sh"),
		test.FilePath("config.go"),
		strings.Empty,
		"env:BOB",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			set := flag.NewFlagSet("test")
			set.AddConfig(file)

			decoder := test.NewDecoder(set)

			_, err := config.NewConfig[config.Config](decoder, test.Validator)
			require.Error(t, err)
		})
	}
}

func TestValidEnvConfig(t *testing.T) {
	tests := []struct {
		name string
		kind string
		path string
	}{
		{name: "yaml", kind: "yaml", path: "configs/config.yml"},
		{name: "hjson", kind: "hjson", path: "configs/config.hjson"},
		{name: "toml", kind: "toml", path: "configs/config.toml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := test.FS.ReadFile(test.Path(tt.path))
			require.NoError(t, err)

			t.Setenv("CONFIG", tt.kind+":"+base64.Encode(d))

			set := flag.NewFlagSet("test")
			set.AddConfig("env:CONFIG")

			decoder := test.NewDecoder(set)

			cfg, err := config.NewConfig[config.Config](decoder, test.Validator)
			require.NoError(t, err)
			verifyConfig(t, cfg)
		})
	}
}

func TestInvalidEnvMissingConfig(t *testing.T) {
	set := flag.NewFlagSet("test")
	set.AddConfig("env:CONFIG")

	decoder := test.NewDecoder(set)

	_, err := config.NewConfig[config.Config](decoder, test.Validator)
	require.ErrorIs(t, err, config.ErrEnvMissing)
}

func TestInvalidEnvKindConfig(t *testing.T) {
	d, err := test.FS.ReadFile(test.Path("configs/config.yml"))
	require.NoError(t, err)

	t.Setenv("CONFIG", "what:"+base64.Encode(d))

	set := flag.NewFlagSet("test")
	set.AddConfig("env:CONFIG")

	decoder := test.NewDecoder(set)

	_, err = config.NewConfig[config.Config](decoder, test.Validator)
	require.ErrorIs(t, err, config.ErrNoEncoder)
}

func TestInvalidEnvDataConfig(t *testing.T) {
	t.Setenv("CONFIG", "yaml:not_good")

	set := flag.NewFlagSet("test")
	set.AddConfig("env:CONFIG")

	decoder := test.NewDecoder(set)

	_, err := config.NewConfig[config.Config](decoder, test.Validator)
	require.Error(t, err)
}

func TestValidCommonConfig(t *testing.T) {
	tests := []struct {
		name string
		path string
		ext  string
	}{
		{name: "yaml", path: "configs/config.yml", ext: ".yml"},
		{name: "hjson", path: "configs/config.hjson", ext: ".hjson"},
		{name: "toml", path: "configs/config.toml", ext: ".toml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("XDG_CONFIG_HOME", test.FS.Join(home, ".config"))

			configDir := os.UserConfigDir()
			path := test.FS.Join(configDir, test.Name.String())

			require.NoError(t, test.FS.MkdirAll(path, 0o777))

			data, err := test.FS.ReadFile(test.Path(tt.path))
			require.NoError(t, err)

			require.NoError(t, test.FS.WriteFile(test.FS.Join(path, test.Name.String()+tt.ext), data, 0o600))

			set := flag.NewFlagSet("test")
			set.AddConfig(strings.Empty)

			decoder := test.NewDecoder(set)

			cfg, err := config.NewConfig[config.Config](decoder, test.Validator)
			require.NoError(t, err)
			verifyConfig(t, cfg)

			require.NoError(t, test.FS.RemoveAll(path))
		})
	}
}

func TestInvalidCommonConfig(t *testing.T) {
	set := flag.NewFlagSet("test")
	set.AddConfig(strings.Empty)

	decoder := test.NewDecoder(set)

	_, err := config.NewConfig[config.Config](decoder, test.Validator)
	require.Error(t, err)
}

func TestInvalidKindConfig(t *testing.T) {
	set := flag.NewFlagSet("test")
	set.AddConfig("test:test")

	decoder := test.NewDecoder(set)

	_, err := config.NewConfig[config.Config](decoder, test.Validator)
	require.Error(t, err)
}

func TestNewConfigRejectsEmptyDecodedConfig(t *testing.T) {
	_, err := config.NewConfig[config.Config](decoderFunc(func(any) error {
		return nil
	}), test.Validator)

	require.ErrorIs(t, err, config.ErrInvalidConfig)
}

func verifyConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	verifyDebugConfig(t, cfg)
	verifyCacheConfig(t, cfg)
	verifyCryptoConfig(t, cfg)
	verifyFeatureConfig(t, cfg)
	verifyIDConfig(t, cfg)
	verifyHooksConfig(t, cfg)
	verifySQLConfig(t, cfg)
	verifyTelemetryConfig(t, cfg)
	verifyTimeConfig(t, cfg)
	verifyGRPCConfig(t, cfg)
	verifyHTTPConfig(t, cfg)
}

func verifyDebugConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.True(t, cfg.Debug.IsEnabled())
	require.Equal(t, "tcp://localhost:6060", cfg.Debug.Address)
	require.False(t, cfg.Debug.TLS.IsEnabled())
}

func verifyCacheConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.Equal(t, "redis", cfg.Cache.Kind)
	require.Equal(t, "snappy", cfg.Cache.Compressor)
	require.Equal(t, "proto", cfg.Cache.Encoder)
	require.Equal(t, 4*bytes.MB, cfg.Cache.MaxSize)
	require.Equal(t, "file:../test/secrets/redis", cfg.Cache.Options["url"])
}

func verifyCryptoConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.True(t, cfg.Crypto.IsEnabled())
	require.True(t, cfg.Crypto.AES.IsEnabled())
	require.NotEmpty(t, cfg.Crypto.AES.Key)
	require.True(t, cfg.Crypto.Ed25519.IsEnabled())
	require.NotEmpty(t, cfg.Crypto.Ed25519.Public)
	require.NotEmpty(t, cfg.Crypto.Ed25519.Private)
	require.True(t, cfg.Crypto.HMAC.IsEnabled())
	require.NotEmpty(t, cfg.Crypto.HMAC.Key)
	require.True(t, cfg.Crypto.RSA.IsEnabled())
	require.NotEmpty(t, cfg.Crypto.RSA.Public)
	require.NotEmpty(t, cfg.Crypto.RSA.Private)
	require.True(t, cfg.Crypto.SSH.IsEnabled())
	require.NotEmpty(t, cfg.Crypto.SSH.Public)
	require.NotEmpty(t, cfg.Crypto.SSH.Private)
}

func verifyFeatureConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.Equal(t, "development", cfg.Environment.String())
	require.True(t, cfg.Feature.IsEnabled())
	require.Equal(t, "localhost:9000", cfg.Feature.Address)
	require.Equal(t, 10*time.Second, cfg.Feature.Timeout)
	require.Equal(t, 100*time.Millisecond, cfg.Feature.Retry.Backoff)
	require.Equal(t, time.Second, cfg.Feature.Retry.Timeout)
	require.Equal(t, uint64(3), cfg.Feature.Retry.Attempts)
}

func verifyIDConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.Equal(t, "uuid", cfg.ID.Kind)
}

func verifyHooksConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.Equal(t, "file:../test/secrets/hooks", cfg.Hooks.Secret)
}

func verifySQLConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.Len(t, cfg.SQL.PG.Masters, 1)
	require.Equal(t, "file:../test/secrets/pg", cfg.SQL.PG.Masters[0].URL)
	require.Len(t, cfg.SQL.PG.Slaves, 1)
	require.Equal(t, "file:../test/secrets/pg", cfg.SQL.PG.Slaves[0].URL)
	require.Equal(t, 5, cfg.SQL.PG.MaxIdleConns)
	require.Equal(t, 5, cfg.SQL.PG.MaxOpenConns)
	require.Equal(t, time.Hour, cfg.SQL.PG.ConnMaxLifetime)
}

func verifyTelemetryConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.Equal(t, "text", cfg.Telemetry.Logger.Kind)
	require.Equal(t, "info", cfg.Telemetry.Logger.Level)
	require.Equal(t, "prometheus", cfg.Telemetry.Metrics.Kind)
	require.Equal(t, "http://localhost:4318/v1/traces", cfg.Telemetry.Tracer.URL)
	require.Equal(t, "otlp", cfg.Telemetry.Tracer.Kind)
}

func verifyTimeConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.Equal(t, "nts", cfg.Time.Kind)
	require.Equal(t, "time.cloudflare.com", cfg.Time.Address)
}

func verifyGRPCConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.True(t, cfg.Transport.GRPC.Token.IsEnabled())
	require.Equal(t, "file:../test/configs/rbac.conf", cfg.Transport.GRPC.Token.Access.Model)
	require.Equal(t, "file:../test/configs/rbac.csv", cfg.Transport.GRPC.Token.Access.Policy)
	require.Equal(t, "jwt", cfg.Transport.GRPC.Token.Kind)
	require.Equal(t, time.Hour, cfg.Transport.GRPC.Token.JWT.Expiration)
	require.Equal(t, "iss", cfg.Transport.GRPC.Token.JWT.Issuer)
	require.Equal(t, "1234567890", cfg.Transport.GRPC.Token.JWT.Key)
	require.Equal(t, "file:../test/secrets/ed25519_public", cfg.Transport.GRPC.Token.JWT.Keys.Get("1234567890").Public)
	require.Equal(t, "file:../test/secrets/ed25519_private", cfg.Transport.GRPC.Token.JWT.Keys.Get("1234567890").Private)
	require.True(t, cfg.Transport.GRPC.Config.IsEnabled())
	require.Equal(t, "tcp://localhost:12000", cfg.Transport.GRPC.Address)
	require.Equal(t, 10*time.Second, cfg.Transport.GRPC.Timeout)
	require.Equal(t, 3*bytes.MB, cfg.Transport.GRPC.MaxReceiveSize)
	require.Equal(t, "user-agent", cfg.Transport.GRPC.Limiter.Kind)
	require.Equal(t, uint64(10), cfg.Transport.GRPC.Limiter.Tokens)
	require.Equal(t, time.Second, cfg.Transport.GRPC.Limiter.Interval)
	require.Equal(t,
		options.Map{
			"keepalive_enforcement_policy_ping_min_time": "10s",
			"keepalive_max_connection_idle":              "10s",
			"keepalive_max_connection_age":               "10s",
			"keepalive_max_connection_age_grace":         "10s",
			"keepalive_ping_time":                        "10s",
			"max_concurrent_streams":                     "64",
			"connection_timeout":                         "3s",
			"max_header_list_size":                       "16MB",
			"initial_window_size":                        "1MB",
			"initial_conn_window_size":                   "2MB",
			"max_send_msg_size":                          "8MB",
		},
		cfg.Transport.GRPC.Options,
	)
	require.False(t, cfg.Transport.GRPC.TLS.IsEnabled())
}

func verifyHTTPConfig(t *testing.T, cfg *config.Config) {
	t.Helper()

	require.True(t, cfg.Transport.HTTP.Token.IsEnabled())
	require.Equal(t, "file:../test/configs/rbac.conf", cfg.Transport.HTTP.Token.Access.Model)
	require.Equal(t, "file:../test/configs/rbac.csv", cfg.Transport.HTTP.Token.Access.Policy)
	require.Equal(t, "jwt", cfg.Transport.HTTP.Token.Kind)
	require.Equal(t, time.Hour, cfg.Transport.HTTP.Token.JWT.Expiration)
	require.Equal(t, "iss", cfg.Transport.HTTP.Token.JWT.Issuer)
	require.Equal(t, "1234567890", cfg.Transport.HTTP.Token.JWT.Key)
	require.Equal(t, "file:../test/secrets/ed25519_public", cfg.Transport.HTTP.Token.JWT.Keys.Get("1234567890").Public)
	require.Equal(t, "file:../test/secrets/ed25519_private", cfg.Transport.HTTP.Token.JWT.Keys.Get("1234567890").Private)
	require.True(t, cfg.Transport.HTTP.Config.IsEnabled())
	require.Equal(t, "tcp://localhost:11000", cfg.Transport.HTTP.Address)
	require.Equal(t, 10*time.Second, cfg.Transport.HTTP.Timeout)
	require.Equal(t, 2*bytes.MB, cfg.Transport.HTTP.MaxReceiveSize)
	require.Equal(t, "user-agent", cfg.Transport.HTTP.Limiter.Kind)
	require.Equal(t, uint64(10), cfg.Transport.HTTP.Limiter.Tokens)
	require.Equal(t, time.Second, cfg.Transport.HTTP.Limiter.Interval)
	require.Equal(t,
		options.Map{
			"read_timeout":        "10s",
			"write_timeout":       "10s",
			"idle_timeout":        "10s",
			"read_header_timeout": "10s",
			"max_header_bytes":    "32KB",
		},
		cfg.Transport.HTTP.Options,
	)
	require.False(t, cfg.Transport.HTTP.TLS.IsEnabled())
}

type decoderFunc func(any) error

func (f decoderFunc) Decode(v any) error {
	return f(v)
}
