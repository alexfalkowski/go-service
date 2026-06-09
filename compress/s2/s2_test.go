package s2_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestCompressor(t *testing.T) {
	cmp := s2.NewCompressor()
	data := strings.Bytes("hello")
	size := bytes.Size(len(data))

	compressed, err := cmp.Compress(data, size)
	require.NoError(t, err)

	decompressed, err := cmp.Decompress(compressed, size)
	require.NoError(t, err)
	require.Equal(t, data, decompressed)
}

func TestCompressRejectsTooLarge(t *testing.T) {
	cmp := s2.NewCompressor()

	_, err := cmp.Compress(strings.Bytes("hello"), 4)
	require.ErrorIs(t, err, errors.ErrTooLarge)
}

func TestDecompressRejectsTooLarge(t *testing.T) {
	cmp := s2.NewCompressor()
	data := strings.Bytes("hello")
	compressed, err := cmp.Compress(data, bytes.KB)
	require.NoError(t, err)

	_, err = cmp.Decompress(compressed, bytes.Size(len(data)-1))
	require.ErrorIs(t, err, errors.ErrTooLarge)
}

func TestDecompressRejectsInvalidData(t *testing.T) {
	cmp := s2.NewCompressor()

	_, err := cmp.Decompress(strings.Bytes("invalid"), bytes.KB)
	require.Error(t, err)
	require.NotErrorIs(t, err, errors.ErrTooLarge)
}

func TestDecompressReturnsDecodedLenError(t *testing.T) {
	cmp := s2.NewCompressor()

	_, err := cmp.Decompress([]byte{0x80}, bytes.KB)
	require.Error(t, err)
	require.NotErrorIs(t, err, errors.ErrTooLarge)
}
