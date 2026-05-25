package msgpack_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/encoding/msgpack"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder()
	msg := map[string]string{"test": "test"}

	require.NoError(t, encoder.Encode(buffer, msg))

	var actual map[string]string
	require.NoError(t, encoder.Decode(buffer, &actual))
	require.Equal(t, msg, actual)
}

func TestEncodeReturnsError(t *testing.T) {
	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder()

	require.Error(t, encoder.Encode(buffer, func() {}))
}

func TestDecodeReturnsError(t *testing.T) {
	encoder := msgpack.NewEncoder()

	var actual map[string]string
	require.Error(t, encoder.Decode(&test.ErrReaderCloser{}, &actual))
}

func TestMarshalUnmarshal(t *testing.T) {
	msg := map[string]string{"test": "test"}

	data, err := msgpack.Marshal(msg)
	require.NoError(t, err)

	var actual map[string]string
	require.NoError(t, msgpack.Unmarshal(data, &actual))
	require.Equal(t, msg, actual)
}

func TestDecodeRejectsTrailingValue(t *testing.T) {
	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder()
	require.NoError(t, encoder.Encode(buffer, map[string]string{"test": "test"}))
	require.NoError(t, encoder.Encode(buffer, map[string]string{"extra": "value"}))

	var actual map[string]string
	err := encoder.Decode(buffer, &actual)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}

func TestDecodeRejectsMalformedTrailingData(t *testing.T) {
	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder()
	require.NoError(t, encoder.Encode(buffer, map[string]string{"test": "test"}))
	_, err := buffer.WriteString("junk")
	require.NoError(t, err)

	var actual map[string]string
	err = encoder.Decode(buffer, &actual)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}

func TestUnmarshalRejectsTrailingValue(t *testing.T) {
	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := msgpack.NewEncoder()
	require.NoError(t, encoder.Encode(buffer, map[string]string{"test": "test"}))
	require.NoError(t, encoder.Encode(buffer, map[string]string{"extra": "value"}))

	var actual map[string]string
	err := msgpack.Unmarshal(test.Pool.Copy(buffer), &actual)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}

func TestUnmarshalRejectsLargeTrailingHeader(t *testing.T) {
	data, err := msgpack.Marshal(map[string]string{"test": "test"})
	require.NoError(t, err)
	data = append(data, 0xdd, 0xff, 0xff, 0xff, 0xff)

	var actual map[string]string
	err = msgpack.Unmarshal(data, &actual)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}

func TestUnmarshalReturnsDecodeError(t *testing.T) {
	var actual map[string]string

	require.Error(t, msgpack.Unmarshal([]byte("junk"), &actual))
}

func TestUnmarshalRejectsMalformedTrailingData(t *testing.T) {
	data, err := msgpack.Marshal(map[string]string{"test": "test"})
	require.NoError(t, err)
	data = append(data, []byte("junk")...)

	var actual map[string]string
	err = msgpack.Unmarshal(data, &actual)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}
