package rand_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/stretchr/testify/require"
)

func TestValidRand(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	c, err := gen.GenerateBytes(5)
	require.NoError(t, err)
	require.Len(t, c, 5)

	text, err := gen.GenerateText(32)
	require.NoError(t, err)
	require.Len(t, text, 32)
}

func TestGenerateTextUsesAlphanumericAlphabet(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	text, err := gen.GenerateText(256)
	require.NoError(t, err)

	for _, r := range text {
		require.Contains(t, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", string(r))
	}
}

func TestGenerateBytesUsesRawBytes(t *testing.T) {
	gen := rand.NewGenerator(test.StaticReader{Data: []byte{0x00, 0x01, 0x7f, 0x80, 0xff}})

	data, err := gen.GenerateBytes(5)
	require.NoError(t, err)
	require.Equal(t, []byte{0x00, 0x01, 0x7f, 0x80, 0xff}, data)
}

func TestGenerateBytesReadsShortChunksFully(t *testing.T) {
	gen := rand.NewGenerator(&shortReader{data: []byte{0x00, 0x01, 0x7f, 0x80, 0xff}, size: 2})

	data, err := gen.GenerateBytes(5)
	require.NoError(t, err)
	require.Equal(t, []byte{0x00, 0x01, 0x7f, 0x80, 0xff}, data)
}

func TestInvalidRand(t *testing.T) {
	t.Run("bytes", func(t *testing.T) {
		gen := rand.NewGenerator(&test.ErrReaderCloser{})

		_, err := gen.GenerateBytes(5)
		require.Error(t, err)
	})

	t.Run("text", func(t *testing.T) {
		gen := rand.NewGenerator(&test.ErrReaderCloser{})

		_, err := gen.GenerateText(5)
		require.Error(t, err)
	})
}

func TestInvalidSize(t *testing.T) {
	t.Run("bytes", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		var data []byte
		var err error
		require.NotPanics(t, func() {
			data, err = gen.GenerateBytes(-1)
		})
		require.Nil(t, data)
		require.ErrorIs(t, err, rand.ErrInvalidSize)
	})

	t.Run("text", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		var text string
		var err error
		require.NotPanics(t, func() {
			text, err = gen.GenerateText(-1)
		})
		require.Empty(t, text)
		require.ErrorIs(t, err, rand.ErrInvalidSize)
	})
}

type shortReader struct {
	data []byte
	size int
}

func (r *shortReader) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}

	n := min(len(p), min(r.size, len(r.data)))
	copy(p, r.data[:n])
	r.data = r.data[n:]

	return n, nil
}
