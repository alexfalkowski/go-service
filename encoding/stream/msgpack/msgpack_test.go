package msgpack_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/msgpack"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/stretchr/testify/require"
)

func TestEncoderRoundTripsValue(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder(buffer)
	require.NoError(t, encoder.Encode(map[string]string{"test": "one"}))
	require.NoError(t, encoder.Encode(map[string]string{"test": "two"}))
	require.NoError(t, encoder.Close())

	decoder := msgpack.NewDecoder(buffer)

	var first map[string]string
	require.NoError(t, decoder.Decode(&first))
	require.Equal(t, map[string]string{"test": "one"}, first)

	var second map[string]string
	require.NoError(t, decoder.Decode(&second))
	require.Equal(t, map[string]string{"test": "two"}, second)

	var third map[string]string
	require.ErrorIs(t, decoder.Decode(&third), io.EOF)
	require.NoError(t, decoder.Close())
}

func TestEncodeReturnsError(t *testing.T) {
	t.Parallel()

	encoder := msgpack.NewEncoder(test.ErrWriter{})

	require.Error(t, encoder.Encode(map[string]string{"test": "test"}))
}

func TestDecodeReturnsError(t *testing.T) {
	t.Parallel()

	decoder := msgpack.NewDecoder(&test.ErrReaderCloser{})

	var msg map[string]string
	require.Error(t, decoder.Decode(&msg))
}

func TestDecodeAllowsTrailingValues(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder(buffer)
	require.NoError(t, encoder.Encode(map[string]string{"test": "one"}))
	require.NoError(t, encoder.Encode(map[string]string{"test": "two"}))

	decoder := msgpack.NewDecoder(buffer)

	var first map[string]string
	require.NoError(t, decoder.Decode(&first))
	require.Equal(t, map[string]string{"test": "one"}, first)

	var second map[string]string
	require.NoError(t, decoder.Decode(&second))
	require.Equal(t, map[string]string{"test": "two"}, second)
}

func TestDecodeRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder(buffer)
	require.NoError(t, encoder.Encode(map[string]string{"test": "test", "extra": "ignored"}))

	decoder := msgpack.NewDecoder(buffer)
	msg := &message{}

	err := decoder.Decode(msg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "extra")
}

func TestDecodeDiscardsUnknownFields(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder(buffer)
	require.NoError(t, encoder.Encode(map[string]string{"test": "test", "extra": "ignored"}))

	decoder := msgpack.NewDecoder(buffer, codec.WithDiscardUnknown())
	msg := &message{}

	require.NoError(t, decoder.Decode(msg))
	require.Equal(t, "test", msg.Test)
}

type message struct {
	Test string `msgpack:"test"`
}

func TestEncoderCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder(buffer)
	require.NoError(t, encoder.Encode(map[string]string{"test": "test"}))
	require.NoError(t, encoder.Close())
	require.NoError(t, encoder.Close())
}

func TestDecoderCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	decoder := msgpack.NewDecoder(buffer)
	require.NoError(t, decoder.Close())
	require.NoError(t, decoder.Close())
}
