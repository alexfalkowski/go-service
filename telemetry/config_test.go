package telemetry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/bytes"
	"github.com/alexfalkowski/go-service/v3/encoding/hjson"
	"github.com/alexfalkowski/go-service/v3/encoding/json"
	"github.com/alexfalkowski/go-service/v3/encoding/toml"
	"github.com/alexfalkowski/go-service/v3/encoding/yaml"
	"github.com/alexfalkowski/go-service/v3/strings"
	"github.com/alexfalkowski/go-service/v3/telemetry"
	"github.com/alexfalkowski/go-service/v3/telemetry/logger"
	"github.com/alexfalkowski/go-service/v3/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v3/telemetry/tracer"
	"github.com/stretchr/testify/require"
)

func TestConfigDecodesExporterEndpoints(t *testing.T) {
	t.Parallel()

	const object = `{
		"logger": {"endpoint": "localhost:4317"},
		"metrics": {"endpoint": "http://localhost:4318/v1/metrics"},
		"tracer": {"endpoint": "https://collector.example.com/v1/traces"}
	}`
	const tables = `
[logger]
endpoint = "localhost:4317"
[metrics]
endpoint = "http://localhost:4318/v1/metrics"
[tracer]
endpoint = "https://collector.example.com/v1/traces"
`
	want := telemetry.Config{
		Logger:  &logger.Config{Endpoint: "localhost:4317"},
		Metrics: &metrics.Config{Endpoint: "http://localhost:4318/v1/metrics"},
		Tracer:  &tracer.Config{Endpoint: "https://collector.example.com/v1/traces"},
	}
	tests := []struct {
		name      string
		data      string
		unmarshal func([]byte, any) error
	}{
		{name: "JSON", data: object, unmarshal: json.Unmarshal},
		{name: "HJSON", data: object, unmarshal: hjson.Unmarshal},
		{name: "YAML", data: object, unmarshal: yaml.Unmarshal},
		{name: "TOML", data: tables, unmarshal: toml.Unmarshal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg telemetry.Config
			require.NoError(t, tt.unmarshal([]byte(tt.data), &cfg))
			require.Equal(t, want, cfg)

			legacy := strings.ReplaceAll(tt.data, "endpoint", "url")
			require.Error(t, tt.unmarshal([]byte(legacy), &cfg))
		})
	}
}

func TestMetadataConfigGetMaxValueSize(t *testing.T) {
	for _, test := range []struct {
		name   string
		config *telemetry.MetadataConfig
		want   bytes.Size
	}{
		{name: "nil", want: 1024},
		{name: "zero", config: &telemetry.MetadataConfig{}, want: 1024},
		{name: "configured", config: &telemetry.MetadataConfig{MaxValueSize: 4 * bytes.KB}, want: 4 * bytes.KB},
	} {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, bytes.Size(test.config.GetMaxValueSize()))
		})
	}
}
