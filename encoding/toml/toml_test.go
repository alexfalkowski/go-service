package toml_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	encoder := toml.NewEncoder()
	msg := map[string]string{"test": "test"}

	require.NoError(t, encoder.Encode(bytes, msg))
	require.Equal(t, `test = "test"`, strings.TrimSpace(bytes.String()))
}

func TestMarshalUnmarshal(t *testing.T) {
	msg := map[string]string{"test": "test"}

	data, err := toml.Marshal(msg)
	require.NoError(t, err)
	require.Equal(t, `test = "test"`, strings.TrimSpace(string(data)))

	var actual map[string]string
	require.NoError(t, toml.Unmarshal(data, &actual))
	require.Equal(t, msg, actual)
}

func TestMarshalReturnsError(t *testing.T) {
	_, err := toml.Marshal(func() {})

	require.Error(t, err)
}

func TestEncodeReturnsError(t *testing.T) {
	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	encoder := toml.NewEncoder()

	require.Error(t, encoder.Encode(bytes, func() {}))
}

func TestDecode(t *testing.T) {
	encoder := toml.NewEncoder()
	var msg map[string]string

	require.NoError(t, encoder.Decode(bytes.NewBufferString(`test = "test"`), &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}

func TestDecodeIgnoresUndecodedMetadata(t *testing.T) {
	encoder := toml.NewEncoder()
	msg := &message{}

	require.NoError(t, encoder.Decode(bytes.NewBufferString("test = \"test\"\nextra = \"ignored\""), msg))
	require.Equal(t, "test", msg.Test)
}

type message struct {
	Test string `toml:"test"`
}
