package bytes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/stretchr/testify/require"
)

// FuzzSizeTextRoundTrip explores human-readable size parsing and text marshaling used by config values.
func FuzzSizeTextRoundTrip(f *testing.F) {
	f.Add("0B")
	f.Add("64B")
	f.Add("4MB")
	f.Add("1048576B")
	f.Add("not-a-size")

	f.Fuzz(func(t *testing.T, text string) {
		if len(text) > 128 {
			t.Skip()
		}

		var size bytes.Size
		if err := size.UnmarshalText([]byte(text)); err != nil {
			return
		}
		if size < 0 || size > bytes.PB {
			t.Skip()
		}

		encoded, err := size.MarshalText()
		require.NoError(t, err)

		var decoded bytes.Size
		require.NoError(t, decoded.UnmarshalText(encoded))
		require.Equal(t, size, decoded)
	})
}

// FuzzSizeJSONRoundTrip explores JSON size parsing while preserving the text round-trip invariant.
func FuzzSizeJSONRoundTrip(f *testing.F) {
	f.Add([]byte(`"0B"`))
	f.Add([]byte(`"64B"`))
	f.Add([]byte(`"4MB"`))
	f.Add([]byte(`1048576`))
	f.Add([]byte(`"not-a-size"`))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1024 {
			t.Skip()
		}

		var size bytes.Size
		if err := json.Unmarshal(data, &size); err != nil {
			return
		}
		if size < 0 || size > bytes.PB {
			t.Skip()
		}

		encoded, err := json.Marshal(size)
		require.NoError(t, err)

		var decoded bytes.Size
		require.NoError(t, json.Unmarshal(encoded, &decoded))
		require.Equal(t, size, decoded)
	})
}
