package rand_test

import (
	"io"
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
	gen := rand.NewGenerator(staticReader{data: []byte{0x00, 0x01, 0x7f, 0x80, 0xff}})

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

type staticReader struct {
	data []byte
}

func (s staticReader) Read(p []byte) (int, error) {
	n := copy(p, s.data)
	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}
