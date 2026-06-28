package json_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/stretchr/testify/require"
)

func FuzzUnmarshal(f *testing.F) {
	f.Add([]byte(`{"test":"test"}`))
	f.Add([]byte(`{"test":""}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"test":"test"}{"extra":"value"}`))
	f.Add([]byte(`{`))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 16*1024 {
			t.Skip()
		}

		var msg map[string]string
		if err := json.Unmarshal(data, &msg); err != nil {
			return
		}

		encoded, err := json.Marshal(msg)
		require.NoError(t, err)

		var decoded map[string]string
		require.NoError(t, json.Unmarshal(encoded, &decoded))
		require.Len(t, decoded, len(msg))
		for key, value := range msg {
			require.Equal(t, value, decoded[key])
		}
	})
}
