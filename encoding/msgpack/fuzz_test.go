package msgpack_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/msgpack"
	"github.com/stretchr/testify/require"
)

func FuzzUnmarshal(f *testing.F) {
	for _, msg := range []map[string]string{
		{},
		{"test": "test"},
		{"test": ""},
	} {
		data, err := msgpack.Marshal(msg)
		require.NoError(f, err)
		f.Add(data)
	}
	f.Add([]byte("junk"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 16*1024 {
			t.Skip()
		}

		var msg map[string]string
		if err := msgpack.Unmarshal(data, &msg); err != nil {
			return
		}

		encoded, err := msgpack.Marshal(msg)
		require.NoError(t, err)

		var decoded map[string]string
		require.NoError(t, msgpack.Unmarshal(encoded, &decoded))
		require.Len(t, decoded, len(msg))
		for key, value := range msg {
			require.Equal(t, value, decoded[key])
		}
	})
}
