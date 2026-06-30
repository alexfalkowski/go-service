package bytes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	encoding "github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/stretchr/testify/require"
)

// FuzzEncoder explores arbitrary byte streams through the repository's encode/decode byte-copy invariant.
func FuzzEncoder(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("test"))
	f.Add([]byte("line\nbreak"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 64*1024 {
			t.Skip()
		}

		encoder := encoding.NewEncoder()
		encoded := &bytes.Buffer{}
		require.NoError(t, encoder.Encode(encoded, bytes.NewBuffer(data)))

		decoded := &bytes.Buffer{}
		require.NoError(t, encoder.Decode(encoded, decoded))
		require.Equal(t, data, decoded.Bytes())
	})
}
