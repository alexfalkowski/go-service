package json_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/json"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/stretchr/testify/require"
)

func TestEncoderRoundTripsValue(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := json.NewEncoder(buffer)
	require.NoError(t, encoder.Encode(map[string]string{"test": "one"}))
	require.NoError(t, encoder.Encode(map[string]string{"test": "two"}))
	require.NoError(t, encoder.Close())

	require.NotContains(t, buffer.String(), "  ")

	decoder := json.NewDecoder(buffer)

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

	encoder := json.NewEncoder(test.ErrWriter{})

	require.Error(t, encoder.Encode(map[string]string{"test": "test"}))
}

func TestDecodeReturnsError(t *testing.T) {
	t.Parallel()

	decoder := json.NewDecoder(&test.ErrReaderCloser{})

	var msg map[string]string
	require.Error(t, decoder.Decode(&msg))
}

func TestDecodeRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	_, err := buffer.WriteString(`{"test":"test","extra":"ignored"}`)
	require.NoError(t, err)

	decoder := json.NewDecoder(buffer)
	msg := &message{}

	err = decoder.Decode(msg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "extra")
}

func TestDecodeDiscardsUnknownFields(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	_, err := buffer.WriteString(`{"test":"test","extra":"ignored"}`)
	require.NoError(t, err)

	decoder := json.NewDecoder(buffer, codec.WithDiscardUnknown())
	msg := &message{}

	require.NoError(t, decoder.Decode(msg))
	require.Equal(t, "test", msg.Test)
}

func TestEncoderCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := json.NewEncoder(buffer)
	require.NoError(t, encoder.Encode(map[string]string{"test": "test"}))
	require.NoError(t, encoder.Close())
	require.NoError(t, encoder.Close())
}

func TestDecoderCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	decoder := json.NewDecoder(buffer)
	require.NoError(t, decoder.Close())
	require.NoError(t, decoder.Close())
}

type message struct {
	Test string `json:"test"`
}
