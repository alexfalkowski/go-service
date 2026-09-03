package yaml_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/yaml"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/stretchr/testify/require"
)

func TestEncoderRoundTripsValue(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := yaml.NewEncoder(buffer)
	require.NoError(t, encoder.Encode(map[string]string{"test": "one"}))
	require.NoError(t, encoder.Encode(map[string]string{"test": "two"}))
	require.NoError(t, encoder.Close())

	require.Contains(t, buffer.String(), "---")

	decoder := yaml.NewDecoder(buffer)

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

func TestEncodeOrCloseReturnsError(t *testing.T) {
	t.Parallel()

	encoder := yaml.NewEncoder(test.ErrWriter{})

	require.Error(t, errors.Join(
		encoder.Encode(map[string]string{"test": "test"}),
		encoder.Close(),
	))
}

func TestDecodeReturnsError(t *testing.T) {
	t.Parallel()

	decoder := yaml.NewDecoder(&test.ErrReaderCloser{})

	var msg map[string]string
	require.Error(t, decoder.Decode(&msg))
}

func TestDecodeRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	_, err := buffer.WriteString("test: test\nextra: ignored")
	require.NoError(t, err)

	decoder := yaml.NewDecoder(buffer)
	msg := &message{}

	err = decoder.Decode(msg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "extra")
}

func TestDecodeDiscardsUnknownFields(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	_, err := buffer.WriteString("test: test\nextra: ignored")
	require.NoError(t, err)

	decoder := yaml.NewDecoder(buffer, codec.WithDiscardUnknown())
	msg := &message{}

	require.NoError(t, decoder.Decode(msg))
	require.Equal(t, "test", msg.Test)
}

func TestEncoderCloseIsNotIdempotent(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := yaml.NewEncoder(buffer)
	require.NoError(t, encoder.Encode(map[string]string{"test": "test"}))
	require.NoError(t, encoder.Close())
	require.Error(t, encoder.Close())
}

func TestDecoderCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	decoder := yaml.NewDecoder(buffer)
	require.NoError(t, decoder.Close())
	require.NoError(t, decoder.Close())
}

type message struct {
	Test string `yaml:"test"`
}
