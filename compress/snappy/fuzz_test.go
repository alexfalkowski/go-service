package snappy_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
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

		cmp := snappy.NewCompressor()
		size := bytes.Size(limit)

		compressed, err := cmp.Compress(data, size)
		if int64(len(data)) > size.Bytes() {
			require.ErrorIs(t, err, errors.ErrTooLarge)
			return
		}
		require.NoError(t, err)

		decompressed, err := cmp.Decompress(compressed, size)
		require.NoError(t, err)
		require.Equal(t, string(data), string(decompressed))
	})
}

func FuzzDecompress(f *testing.F) {
	cmp := snappy.NewCompressor()
	encoded, err := cmp.Compress([]byte("hello"), bytes.KB)
	require.NoError(f, err)

	f.Add(encoded, uint16(5))
	f.Add([]byte("invalid"), uint16(1024))
	f.Add([]byte{0x80}, uint16(1024))

	f.Fuzz(func(t *testing.T, data []byte, limit uint16) {
		if len(data) > 64*1024 {
			t.Skip()
		}

		decoded, err := snappy.NewCompressor().Decompress(data, bytes.Size(limit))
		if err != nil {
			return
		}

		require.LessOrEqual(t, int64(len(decoded)), int64(limit))
	})
}
