package bytes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	encoding "github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	encoder := encoding.NewEncoder()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	require.NoError(t, encoder.Encode(buffer, bytes.NewBufferString("yes!")))
	require.Equal(t, "yes!", strings.TrimSpace(buffer.String()))

	var str string
	require.Error(t, encoder.Encode(buffer, &str))

	decode := test.Pool.Get()
	defer test.Pool.Put(decode)

	require.NoError(t, encoder.Decode(bytes.NewBufferString("test"), decode))
	require.Equal(t, "test", decode.String())

	require.Error(t, encoder.Decode(bytes.NewBufferString("test"), &str))
}

func TestDecodeResetsDestination(t *testing.T) {
	encoder := encoding.NewEncoder()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	buffer.WriteString("stale:")

	require.NoError(t, encoder.Decode(bytes.NewBufferString("fresh"), buffer))
	require.Equal(t, "fresh", buffer.String())
}

func TestInvalidTypedNilEncode(t *testing.T) {
	encoder := encoding.NewEncoder()

	err := encoder.Encode(io.Discard, (*bytes.Buffer)(nil))
	require.ErrorIs(t, err, errors.ErrInvalidType)
}

func TestEncodeReturnsWriteToError(t *testing.T) {
	encoder := encoding.NewEncoder()

	require.ErrorIs(t, encoder.Encode(io.Discard, test.ErrWriterTo{}), test.ErrFailed)
}

func TestInvalidTypedNilDecode(t *testing.T) {
	encoder := encoding.NewEncoder()

	err := encoder.Decode(bytes.NewBufferString("test"), (*bytes.Buffer)(nil))
	require.ErrorIs(t, err, errors.ErrInvalidType)
}
