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
			set.AddInput(file)

			decoder := test.NewDecoder(set)

			config, err := config.NewConfig[config.Config](decoder, test.Validator)
			require.NoError(t, err)
			verifyConfig(t, config)
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
	set.AddInput("file:~/config.yml")

	decoder := test.NewDecoder(set)

	config, err := config.NewConfig[config.Config](decoder, test.Validator)
	require.NoError(t, err)
	verifyConfig(t, config)
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
			set.AddInput(file)

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
			set.AddInput("env:CONFIG")

			decoder := test.NewDecoder(set)

			config, err := config.NewConfig[config.Config](decoder, test.Validator)
			require.NoError(t, err)
			verifyConfig(t, config)
		})
	}
}

func TestInvalidEnvMissingConfig(t *testing.T) {
	set := flag.NewFlagSet("test")
	set.AddInput("env:CONFIG")

	decoder := test.NewDecoder(set)

	_, err := config.NewConfig[config.Config](decoder, test.Validator)
	require.ErrorIs(t, err, config.ErrEnvMissing)
}

func TestInvalidEnvKindConfig(t *testing.T) {
	d, err := test.FS.ReadFile(test.Path("configs/config.yml"))
	require.NoError(t, err)

	t.Setenv("CONFIG", "what:"+base64.Encode(d))

	set := flag.NewFlagSet("test")
	set.AddInput("env:CONFIG")

	decoder := test.NewDecoder(set)

	_, err = config.NewConfig[config.Config](decoder, test.Validator)
	require.ErrorIs(t, err, config.ErrNoEncoder)
}

func TestInvalidEnvDataConfig(t *testing.T) {
	t.Setenv("CONFIG", "yaml:not_good")

	set := flag.NewFlagSet("test")
	set.AddInput("env:CONFIG")

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
			set.AddInput(strings.Empty)

			decoder := test.NewDecoder(set)

			config, err := config.NewConfig[config.Config](decoder, test.Validator)
			require.NoError(t, err)
			verifyConfig(t, config)

			require.NoError(t, test.FS.RemoveAll(path))
		})
	}
}

func TestInvalidCommonConfig(t *testing.T) {
	set := flag.NewFlagSet("test")
	set.AddInput(strings.Empty)

	decoder := test.NewDecoder(set)

	_, err := config.NewConfig[config.Config](decoder, test.Validator)
	require.Error(t, err)
}

func TestInvalidKindConfig(t *testing.T) {
	set := flag.NewFlagSet("test")
	set.AddInput("test:test")

	decoder := test.NewDecoder(set)

	_, err := config.NewConfig[config.Config](decoder, test.Validator)
	require.Error(t, err)
}

//nolint:funlen,gosec
func verifyConfig(t *testing.T, config *config.Config) {
	t.Helper()

	require.Equal(t, "tcp://localhost:6060", config.Debug.Address)
	require.False(t, config.Debug.TLS.IsEnabled())
	require.True(t, config.Crypto.RSA.IsEnabled())
	require.Equal(t, "redis", config.Cache.Kind)
	require.Equal(t, "snappy", config.Cache.Compressor)
	require.Equal(t, "proto", config.Cache.Encoder)
	require.Equal(t, 4*bytes.MB, config.Cache.MaxSize)
	require.Equal(t, "file:../test/secrets/redis", config.Cache.Options["url"])
	require.True(t, config.Crypto.IsEnabled())
	require.True(t, config.Crypto.AES.IsEnabled())
	require.NotEmpty(t, config.Crypto.AES.Key)
	require.True(t, config.Crypto.Ed25519.IsEnabled())
	require.NotEmpty(t, config.Crypto.Ed25519.Public)
	require.NotEmpty(t, config.Crypto.Ed25519.Private)
	require.True(t, config.Crypto.HMAC.IsEnabled())
	require.NotEmpty(t, config.Crypto.HMAC.Key)
	require.True(t, config.Crypto.RSA.IsEnabled())
	require.NotEmpty(t, config.Crypto.RSA.Public)
	require.NotEmpty(t, config.Crypto.RSA.Private)
	require.True(t, config.Crypto.SSH.IsEnabled())
	require.NotEmpty(t, config.Crypto.SSH.Public)
	require.NotEmpty(t, config.Crypto.SSH.Private)
	require.True(t, config.Debug.IsEnabled())
	require.Equal(t, "development", config.Environment.String())
	require.True(t, config.Feature.IsEnabled())
	require.Equal(t, "localhost:9000", config.Feature.Address)
	require.Equal(t, "uuid", config.ID.Kind)
	require.Equal(t, "file:../test/secrets/hooks", config.Hooks.Secret)
	require.Len(t, config.SQL.PG.Masters, 1)
	require.Equal(t, "file:../test/secrets/pg", config.SQL.PG.Masters[0].URL)
	require.Len(t, config.SQL.PG.Slaves, 1)
	require.Equal(t, "file:../test/secrets/pg", config.SQL.PG.Slaves[0].URL)
	require.Equal(t, 5, config.SQL.PG.MaxIdleConns)
	require.Equal(t, 5, config.SQL.PG.MaxOpenConns)
	require.Equal(t, time.Hour, config.SQL.PG.ConnMaxLifetime)
	require.Equal(t, "text", config.Telemetry.Logger.Kind)
	require.Equal(t, "info", config.Telemetry.Logger.Level)
	require.Equal(t, "prometheus", config.Telemetry.Metrics.Kind)
	require.Equal(t, "nts", config.Time.Kind)
	require.Equal(t, "time.cloudflare.com", config.Time.Address)
	require.Equal(t, "http://localhost:4318/v1/traces", config.Telemetry.Tracer.URL)
	require.Equal(t, "otlp", config.Telemetry.Tracer.Kind)
	require.True(t, config.Transport.GRPC.Token.IsEnabled())
	require.Equal(t, "file:../test/configs/rbac.conf", config.Transport.GRPC.Token.Access.Model)
	require.Equal(t, "file:../test/configs/rbac.csv", config.Transport.GRPC.Token.Access.Policy)
	require.Equal(t, "jwt", config.Transport.GRPC.Token.Kind)
	require.Equal(t, time.Hour, config.Transport.GRPC.Token.JWT.Expiration)
	require.Equal(t, "iss", config.Transport.GRPC.Token.JWT.Issuer)
	require.Equal(t, "1234567890", config.Transport.GRPC.Token.JWT.KeyID)
	require.True(t, config.Transport.GRPC.Config.IsEnabled())
	require.Equal(t, "tcp://localhost:12000", config.Transport.GRPC.Address)
	require.Equal(t, 3*bytes.MB, config.Transport.GRPC.MaxReceiveSize)
	require.Equal(t, "user-agent", config.Transport.GRPC.Limiter.Kind)
	require.Equal(t, 10, int(config.Transport.GRPC.Limiter.Tokens))
	require.Equal(t, time.Second, config.Transport.GRPC.Limiter.Interval)
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
		config.Transport.GRPC.Options,
	)
	require.False(t, config.Transport.GRPC.TLS.IsEnabled())
	require.True(t, config.Transport.HTTP.Token.IsEnabled())
	require.Equal(t, "file:../test/configs/rbac.conf", config.Transport.HTTP.Token.Access.Model)
	require.Equal(t, "file:../test/configs/rbac.csv", config.Transport.HTTP.Token.Access.Policy)
	require.Equal(t, "jwt", config.Transport.HTTP.Token.Kind)
	require.Equal(t, time.Hour, config.Transport.HTTP.Token.JWT.Expiration)
	require.Equal(t, "iss", config.Transport.HTTP.Token.JWT.Issuer)
	require.Equal(t, "1234567890", config.Transport.HTTP.Token.JWT.KeyID)
	require.True(t, config.Transport.HTTP.Config.IsEnabled())
	require.Equal(t, "tcp://localhost:11000", config.Transport.HTTP.Address)
	require.Equal(t, 2*bytes.MB, config.Transport.HTTP.MaxReceiveSize)
	require.Equal(t, "user-agent", config.Transport.HTTP.Limiter.Kind)
	require.Equal(t, 10, int(config.Transport.HTTP.Limiter.Tokens))
	require.Equal(t, time.Second, config.Transport.HTTP.Limiter.Interval)
	require.Equal(t,
		options.Map{
			"read_timeout":        "10s",
			"write_timeout":       "10s",
			"idle_timeout":        "10s",
			"read_header_timeout": "10s",
			"max_header_bytes":    "32KB",
		},
		config.Transport.HTTP.Options,
	)
	require.False(t, config.Transport.HTTP.TLS.IsEnabled())
}
