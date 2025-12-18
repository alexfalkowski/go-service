package yaml_test

import (
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	encoder := yaml.NewEncoder()
	msg := map[string]string{"test": "test"}

	require.NoError(t, encoder.Encode(bytes, msg))
	require.Equal(t, "test: test", strings.TrimSpace(bytes.String()))
}

func TestDecode(t *testing.T) {
	encoder := yaml.NewEncoder()
	var msg map[string]string

	require.NoError(t, encoder.Decode(bytes.NewBufferString("test: test"), &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}
