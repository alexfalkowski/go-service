package yaml_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/stretchr/testify/require"
)

// FuzzUnmarshal explores the strict YAML decoder surface and verifies accepted map payloads round-trip.
func FuzzUnmarshal(f *testing.F) {
	f.Add([]byte("test: test"))
	f.Add([]byte("test: ''"))
	f.Add([]byte("{}"))
	f.Add([]byte("test: test\n---\ntest: other"))
	f.Add([]byte(": invalid"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 16*1024 {
			t.Skip()
		}

		var msg map[string]string
		if err := yaml.Unmarshal(data, &msg); err != nil {
			return
		}

		encoded, err := yaml.Marshal(msg)
		require.NoError(t, err)

		var decoded map[string]string
		require.NoError(t, yaml.Unmarshal(encoded, &decoded))
		require.Len(t, decoded, len(msg))
		for key, value := range msg {
			require.Equal(t, value, decoded[key])
		}
	})
}
