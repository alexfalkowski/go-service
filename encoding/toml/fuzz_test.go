package toml_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/stretchr/testify/require"
)

// FuzzUnmarshal explores TOML decoder input space and verifies accepted map payloads round-trip.
func FuzzUnmarshal(f *testing.F) {
	f.Add([]byte(`test = "test"`))
	f.Add([]byte(`test = ""`))
	f.Add([]byte(""))
	f.Add([]byte("test = \"test\"\nextra = \"value\""))
	f.Add([]byte("="))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 16*1024 {
			t.Skip()
		}

		var msg map[string]string
		if err := toml.Unmarshal(data, &msg); err != nil {
			return
		}

		encoded, err := toml.Marshal(msg)
		require.NoError(t, err)

		var decoded map[string]string
		require.NoError(t, toml.Unmarshal(encoded, &decoded))
		require.Len(t, decoded, len(msg))
		for key, value := range msg {
			require.Equal(t, value, decoded[key])
		}
	})
}
