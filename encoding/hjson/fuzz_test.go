package hjson_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/hjson"
	"github.com/stretchr/testify/require"
)

func FuzzUnmarshal(f *testing.F) {
	f.Add([]byte("{test: test}"))
	f.Add([]byte("{test: ''}"))
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

		var decoded map[string]string
		require.NoError(t, hjson.Unmarshal(encoded, &decoded))
		require.Len(t, decoded, len(msg))
		for key, value := range msg {
			require.Equal(t, value, decoded[key])
		}
	})
}
