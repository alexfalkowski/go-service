package zstd_test

import (
	"math"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
	"github.com/alexfalkowski/go-service/v2/strings"
	compress "github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"
)

func TestCompressor(t *testing.T) {
	t.Parallel()

	cmp := zstd.NewCompressor()
	data := strings.Bytes("hello")
	size := bytes.Size(len(data))

	compressed, err := cmp.Compress(data, size)
	require.NoError(t, err)

	decompressed, err := cmp.Decompress(compressed, size)
	require.NoError(t, err)
	require.Equal(t, data, decompressed)
}

func TestCompressRejectsTooLarge(t *testing.T) {
	t.Parallel()

	cmp := zstd.NewCompressor()

	_, err := cmp.Compress(strings.Bytes("hello"), 4)
	require.ErrorIs(t, err, errors.ErrTooLarge)
}

func TestDecompressRejectsTooLarge(t *testing.T) {
	t.Parallel()

	cmp := zstd.NewCompressor()
	data := strings.Bytes("hello")
	compressed, err := cmp.Compress(data, bytes.KB)
	require.NoError(t, err)

	_, err = cmp.Decompress(compressed, bytes.Size(len(data)-1))
	require.ErrorIs(t, err, errors.ErrTooLarge)
}

func TestDecompressRejectsInvalidData(t *testing.T) {
	t.Parallel()

	cmp := zstd.NewCompressor()

	_, err := cmp.Decompress(strings.Bytes("invalid"), bytes.KB)
	require.Error(t, err)
	require.NotErrorIs(t, err, errors.ErrTooLarge)
}

func TestDecompressPreservesDecoderSizeError(t *testing.T) {
	t.Parallel()

	cmp := zstd.NewCompressor()

	data := make([]byte, zstd.MinWindowSize+1)
	encoded, err := cmp.Compress(data, bytes.Size(len(data)))
	require.NoError(t, err)

	_, err = cmp.Decompress(encoded, zstd.MinWindowSize)
	require.ErrorIs(t, err, errors.ErrTooLarge)
	require.ErrorIs(t, err, zstd.ErrDecoderSizeExceeded)
}

func TestDecompressAllowsWindowLargerThanOutputLimit(t *testing.T) {
	t.Parallel()

	cmp := zstd.NewCompressor()
	data := make([]byte, 64<<10)
	encoder, err := compress.NewWriter(nil, compress.WithWindowSize(128<<10), compress.WithSingleSegment(false))
	require.NoError(t, err)
	defer encoder.Close()

	encoded := encoder.EncodeAll(data, nil)
	decoded, err := cmp.Decompress(encoded, bytes.Size(len(data)))
	require.NoError(t, err)
	require.Equal(t, data, decoded)
}

func TestDecompressRejectsInvalidLimits(t *testing.T) {
	t.Parallel()

	cmp := zstd.NewCompressor()
	data := strings.Bytes("hello")
	encoded, err := cmp.Compress(data, bytes.KB)
	require.NoError(t, err)

	for _, tt := range []struct {
		name string
		size bytes.Size
	}{
		{name: "negative", size: -1},
		{name: "max int64", size: math.MaxInt64},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := cmp.Decompress(encoded, tt.size)
			require.ErrorIs(t, err, errors.ErrTooLarge)
			require.Nil(t, decoded)
		})
	}
}
