package bytes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	hello := strings.Bytes("hello")
	helloCopy := bytes.Copy(hello)

	require.Equal(t, hello, helloCopy)
}
