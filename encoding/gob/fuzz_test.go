package gob_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/stretchr/testify/require"
)

// FuzzUnmarshal explores gob decoder input space and verifies accepted map payloads round-trip.
func FuzzUnmarshal(f *testing.F) {
	for _, msg := range []map[string]string{
		{},
		{"test": "test"},
		{"test": ""},
	} {
		data, err := gob.Marshal(msg)
		require.NoError(f, err)
		f.Add(data)
	}
	f.Add([]byte("junk"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 16*1024 {
			t.Skip()
		}

		var msg map[string]string
		if err := gob.Unmarshal(data, &msg); err != nil {
			return
		}

		encoded, err := gob.Marshal(msg)
		require.NoError(t, err)

		var decoded map[string]string
		require.NoError(t, gob.Unmarshal(encoded, &decoded))
		require.Len(t, decoded, len(msg))
		for key, value := range msg {
			require.Equal(t, value, decoded[key])
		}
	})
}
