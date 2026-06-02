package yaml_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := yaml.NewEncoder()
	msg := map[string]string{"test": "test"}

	require.NoError(t, encoder.Encode(buffer, msg))
	require.Equal(t, "test: test", strings.TrimSpace(buffer.String()))
}

func TestMarshalUnmarshal(t *testing.T) {
	msg := map[string]string{"test": "test"}

	data, err := yaml.Marshal(msg)
	require.NoError(t, err)
	require.Equal(t, "test: test", strings.TrimSpace(string(data)))

	var actual map[string]string
	require.NoError(t, yaml.Unmarshal(data, &actual))
	require.Equal(t, msg, actual)
}

func TestMarshalReturnsError(t *testing.T) {
	_, err := yaml.Marshal(marshalError{})

	require.ErrorIs(t, err, test.ErrFailed)
}

func TestEncodeReturnsError(t *testing.T) {
	encoder := yaml.NewEncoder()

	require.Error(t, encoder.Encode(test.ErrWriter{}, map[string]string{"test": "test"}))
}

func TestDecode(t *testing.T) {
	encoder := yaml.NewEncoder()
	var msg map[string]string

	require.NoError(t, encoder.Decode(bytes.NewBufferString("test: test"), &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}

func TestDecodeRejectsTrailingDocument(t *testing.T) {
	for _, input := range []string{
		"test: test\n---\ntest: other",
		"test: test\n---\n: invalid",
	} {
		t.Run(input, func(t *testing.T) {
			encoder := yaml.NewEncoder()
			var msg map[string]string

			err := encoder.Decode(bytes.NewBufferString(input), &msg)

			require.ErrorIs(t, err, errors.ErrTrailingData)
		})
	}
}

func TestUnmarshalRejectsTrailingDocument(t *testing.T) {
	var msg map[string]string

	err := yaml.Unmarshal([]byte("test: test\n---\ntest: other"), &msg)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}

type marshalError struct{}

func (m marshalError) MarshalYAML() (any, error) {
	return nil, test.ErrFailed
}
