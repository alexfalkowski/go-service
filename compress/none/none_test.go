package none_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestCompressorReturnsDataUnchanged(t *testing.T) {
	cmp := none.NewCompressor()
	data := strings.Bytes("hello")

	compressed, err := cmp.Compress(data, bytes.Size(len(data)))
	require.NoError(t, err)
	require.Equal(t, data, compressed)

	decompressed, err := cmp.Decompress(data, bytes.Size(len(data)))
	require.NoError(t, err)
	require.Equal(t, data, decompressed)
}

func TestCompressRejectsTooLarge(t *testing.T) {
	cmp := none.NewCompressor()

	_, err := cmp.Compress(strings.Bytes("hello"), 4)
	require.ErrorIs(t, err, errors.ErrTooLarge)
}

func TestDecompressRejectsTooLarge(t *testing.T) {
	cmp := none.NewCompressor()

	_, err := cmp.Decompress(strings.Bytes("hello"), 4)
	require.ErrorIs(t, err, errors.ErrTooLarge)
}
