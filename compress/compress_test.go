package compress_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)

			data := strings.Bytes("hello")
			d, err := cmp.Compress(data, bytes.KB)
			require.NoError(t, err)

			ns, err := cmp.Decompress(d, bytes.KB)
			require.NoError(t, err)
			require.Equal(t, data, ns)
		})
	}

	for _, key := range []string{"test", "bob"} {
		t.Run(key, func(t *testing.T) {
			cmp := test.Compressor.Get(key)
			require.Nil(t, cmp)
		})
	}
}

func TestCompressRejectsTooLarge(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)

			_, err := cmp.Compress(strings.Bytes("hello"), 4)
			require.ErrorIs(t, err, errors.ErrTooLarge)
		})
	}
}

func TestDecompressRejectsTooLarge(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)

			data := strings.Bytes("hello")
			d, err := cmp.Compress(data, bytes.KB)
			require.NoError(t, err)

			_, err = cmp.Decompress(d, bytes.Size(len(data)-1))
			require.ErrorIs(t, err, errors.ErrTooLarge)
		})
	}
}

func TestDecompressRejectsInvalidData(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy"} {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)

			_, err := cmp.Decompress(strings.Bytes("invalid"), bytes.KB)
			require.Error(t, err)
			require.NotErrorIs(t, err, errors.ErrTooLarge)
		})
	}
}

func TestS2DecompressReturnsDecodedLenError(t *testing.T) {
	cmp := test.Compressor.Get("s2")
	data := []byte{0x80}

	_, err := cmp.Decompress(data, bytes.KB)
	require.Error(t, err)
	require.NotErrorIs(t, err, errors.ErrTooLarge)
}

func TestSnappyDecompressReturnsDecodedLenError(t *testing.T) {
	cmp := test.Compressor.Get("snappy")
	data := []byte{0x80}

	_, err := cmp.Decompress(data, bytes.KB)
	require.Error(t, err)
	require.NotErrorIs(t, err, errors.ErrTooLarge)
}

func TestZstdDecompressPreservesDecoderSizeError(t *testing.T) {
	cmp := zstd.NewCompressor()

	data := make([]byte, zstd.MinWindowSize+1)
	encoded, err := cmp.Compress(data, bytes.Size(len(data)))
	require.NoError(t, err)

	_, err = cmp.Decompress(encoded, zstd.MinWindowSize)
	require.ErrorIs(t, err, errors.ErrTooLarge)
	require.ErrorIs(t, err, zstd.ErrDecoderSizeExceeded)
}
