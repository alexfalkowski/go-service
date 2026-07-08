package hjson_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/hjson"
	"github.com/stretchr/testify/require"
)

// FuzzUnmarshal explores the HJSON decoder surface, including duplicate-key handling, for accepted map payloads.
func FuzzUnmarshal(f *testing.F) {
	f.Add([]byte("{test: test}"))
	f.Add([]byte("{test: ''}"))
	f.Add([]byte("{test: ' <\"'}"))
	f.Add([]byte("{}"))
	f.Add([]byte("{test: first, test: second}"))
	f.Add([]byte("{"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 16*1024 {
			t.Skip()
		}

		var msg map[string]string
		if err := hjson.Unmarshal(data, &msg); err != nil {
			return
		}

		encoded, err := hjson.Marshal(msg)
		require.NoError(t, err)

		// hjson-go can normalize some accepted decoded strings while formatting, so
		// assert that emitted HJSON reaches a stable representation.
		var normalized map[string]string
		require.NoError(t, hjson.Unmarshal(encoded, &normalized))

		reencoded, err := hjson.Marshal(normalized)
		require.NoError(t, err)

		var decoded map[string]string
		require.NoError(t, hjson.Unmarshal(reencoded, &decoded))
		require.Equal(t, normalized, decoded)
	})
}
