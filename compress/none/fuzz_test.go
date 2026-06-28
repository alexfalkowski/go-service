package none_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/stretchr/testify/require"
)

func FuzzCompressor(f *testing.F) {
	f.Add([]byte(""), uint16(0))
	f.Add([]byte("hello"), uint16(5))
	f.Add([]byte("hello"), uint16(4))

	f.Fuzz(func(t *testing.T, data []byte, limit uint16) {
		if len(data) > 64*1024 {
			t.Skip()
		}

		cmp := none.NewCompressor()
		size := bytes.Size(limit)

		compressed, err := cmp.Compress(data, size)
		if int64(len(data)) > size.Bytes() {
			require.ErrorIs(t, err, errors.ErrTooLarge)
			return
		}
		require.NoError(t, err)
		require.Equal(t, data, compressed)

		decompressed, err := cmp.Decompress(compressed, size)
		require.NoError(t, err)
		require.Equal(t, data, decompressed)
	})
}
