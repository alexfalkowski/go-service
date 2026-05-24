package rand_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestValidRand(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	c, err := gen.GenerateBytes(5)
	require.NoError(t, err)
	require.Len(t, c, 5)

	s, err := gen.GenerateText(32)
	require.NoError(t, err)
	require.Len(t, s, 32)
}

func TestGenerateBytesUsesRawBytes(t *testing.T) {
	gen := rand.NewGenerator(test.StaticReader{Data: []byte{0x00, 0x01, 0x7f, 0x80, 0xff}})

	data, err := gen.GenerateBytes(5)
	require.NoError(t, err)
	require.Equal(t, []byte{0x00, 0x01, 0x7f, 0x80, 0xff}, data)
}

func TestInvalidRand(t *testing.T) {
	gen := rand.NewGenerator(&test.ErrReaderCloser{})

	_, err := gen.GenerateBytes(5)
	require.Error(t, err)

	gen = rand.NewGenerator(&test.ErrReaderCloser{})
	_, err = gen.GenerateText(5)
	require.Error(t, err)
}

func TestInvalidSize(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	var data []byte
	var err error
	require.NotPanics(t, func() {
		data, err = gen.GenerateBytes(-1)
	})
	require.Nil(t, data)
	require.ErrorIs(t, err, rand.ErrInvalidSize)

	var text string
	require.NotPanics(t, func() {
		text, err = gen.GenerateText(-1)
	})
	require.Empty(t, text)
	require.ErrorIs(t, err, rand.ErrInvalidSize)
}
