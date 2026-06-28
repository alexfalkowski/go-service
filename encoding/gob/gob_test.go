package gob_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	encoder := gob.NewEncoder()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	require.NoError(t, encoder.Encode(buffer, map[string]string{"test": "test"}))

	var msg map[string]string
	require.NoError(t, encoder.Decode(buffer, &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}

func TestMarshalUnmarshal(t *testing.T) {
	t.Parallel()

	msg := map[string]string{"test": "test"}

	data, err := gob.Marshal(msg)
	require.NoError(t, err)

	var actual map[string]string
	require.NoError(t, gob.Unmarshal(data, &actual))
	require.Equal(t, msg, actual)
}

func TestMarshalReturnsError(t *testing.T) {
	t.Parallel()

	_, err := gob.Marshal(func() {})

	require.Error(t, err)
}

func TestEncodeReturnsError(t *testing.T) {
	t.Parallel()

	encoder := gob.NewEncoder()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	require.Error(t, encoder.Encode(buffer, func() {}))
}

func TestDecodeReturnsError(t *testing.T) {
	t.Parallel()

	encoder := gob.NewEncoder()

	var msg map[string]string
	require.Error(t, encoder.Decode(&test.ErrReaderCloser{}, &msg))
}

func TestDecodeRejectsTrailingValue(t *testing.T) {
	t.Parallel()

	encoder := gob.NewEncoder()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	require.NoError(t, encoder.Encode(buffer, map[string]string{"test": "test"}))
	require.NoError(t, encoder.Encode(buffer, map[string]string{"extra": "value"}))

	var msg map[string]string
	err := encoder.Decode(buffer, &msg)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}

func TestUnmarshalRejectsTrailingValue(t *testing.T) {
	t.Parallel()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	require.NoError(t, gob.NewEncoder().Encode(buffer, map[string]string{"test": "test"}))
	require.NoError(t, gob.NewEncoder().Encode(buffer, map[string]string{"extra": "value"}))

	var msg map[string]string
	err := gob.Unmarshal(test.Pool.Copy(buffer), &msg)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}

func TestDecodeRejectsMalformedTrailingData(t *testing.T) {
	t.Parallel()

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
