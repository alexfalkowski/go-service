package gob_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	encoder := gob.NewEncoder()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, map[string]string{"test": "test"}))

	var msg map[string]string
	require.NoError(t, encoder.Decode(bytes, &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}

func TestEncodeReturnsError(t *testing.T) {
	encoder := gob.NewEncoder()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.Error(t, encoder.Encode(bytes, func() {}))
}

func TestDecodeReturnsError(t *testing.T) {
	encoder := gob.NewEncoder()

	var msg map[string]string
	require.Error(t, encoder.Decode(&test.ErrReaderCloser{}, &msg))
}

func TestDecodeRejectsTrailingValue(t *testing.T) {
	encoder := gob.NewEncoder()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, map[string]string{"test": "test"}))
	require.NoError(t, encoder.Encode(bytes, map[string]string{"extra": "value"}))

	var msg map[string]string
	err := encoder.Decode(bytes, &msg)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}

func TestDecodeRejectsMalformedTrailingData(t *testing.T) {
	encoder := gob.NewEncoder()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	require.NoError(t, encoder.Encode(buffer, map[string]string{"test": "test"}))
	_, err := buffer.WriteString("junk")
	require.NoError(t, err)

	var msg map[string]string
	err = encoder.Decode(buffer, &msg)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}
