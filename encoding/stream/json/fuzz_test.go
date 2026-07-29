package json_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/json"
	"github.com/stretchr/testify/require"
)

// FuzzUnmarshal explores the strict streaming JSON decoder surface and verifies accepted map
// payloads round-trip through the stream encoder/decoder pair.
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

		decoder := json.NewDecoder(bytes.NewReader(data))

		var msgs []map[string]string
		for {
			var msg map[string]string
			if err := decoder.Decode(&msg); err != nil {
				break
			}
			msgs = append(msgs, msg)
		}
		require.NoError(t, decoder.Close())

		if len(msgs) == 0 {
			return
		}

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		for _, msg := range msgs {
			require.NoError(t, encoder.Encode(msg))
		}
		require.NoError(t, encoder.Close())

		decoded := json.NewDecoder(buffer)
		for _, msg := range msgs {
			var actual map[string]string
			require.NoError(t, decoded.Decode(&actual))
			require.Len(t, actual, len(msg))
			for key, value := range msg {
				require.Equal(t, value, actual[key])
			}
		}
		require.NoError(t, decoded.Close())
	})
}
