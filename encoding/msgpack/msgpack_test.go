package msgpack_test

import (
	"testing"

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

func TestMarshalUnmarshal(t *testing.T) {
	msg := map[string]string{"test": "test"}

	data, err := msgpack.Marshal(msg)
	require.NoError(t, err)

	var actual map[string]string
	require.NoError(t, msgpack.Unmarshal(data, &actual))
	require.Equal(t, msg, actual)
}
