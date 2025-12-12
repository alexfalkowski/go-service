package bytes_test

import (
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	eb "github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	encoder := eb.NewEncoder()

	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	require.NoError(t, encoder.Encode(buffer, bytes.NewBufferString("yes!")))
	require.Equal(t, "yes!", strings.TrimSpace(buffer.String()))

	var str string
	require.Error(t, encoder.Encode(buffer, &str))

	var buf bytes.Buffer
	require.NoError(t, encoder.Decode(bytes.NewBufferString("test"), &buf))
	require.Equal(t, "test", buf.String())

	require.Error(t, encoder.Decode(bytes.NewBufferString("test"), &str))
}
