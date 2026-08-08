package telemetry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/stretchr/testify/require"
)

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
