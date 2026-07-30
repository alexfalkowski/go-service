package gob_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/gob"
	"github.com/stretchr/testify/require"
)

// FuzzUnmarshal explores the streaming gob decoder input space and verifies accepted map payloads
// round-trip through the stream encoder/decoder pair.
func FuzzUnmarshal(f *testing.F) {
	for _, msg := range []map[string]string{
		{},
		{"test": "test"},
		{"test": ""},
	} {
		buffer := &bytes.Buffer{}
		encoder := gob.NewEncoder(buffer)
		require.NoError(f, encoder.Encode(msg))
		f.Add(buffer.Bytes())
	}
	f.Add([]byte("junk"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 16*1024 {
			t.Skip()
		}

		decoder := gob.NewDecoder(bytes.NewReader(data))

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
		encoder := gob.NewEncoder(buffer)
		for _, msg := range msgs {
			require.NoError(t, encoder.Encode(msg))
		}
		require.NoError(t, encoder.Close())

		decoded := gob.NewDecoder(buffer)
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
