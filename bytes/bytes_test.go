package bytes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestCloneReturnsEqualBytes(t *testing.T) {
	hello := strings.Bytes("hello")
	require.Equal(t, hello, bytes.Clone(hello))
}

func TestContainsRuneRecognizesRune(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		b    []byte
		r    rune
		want bool
	}{
		{name: "NUL present", b: []byte("service\x00admin"), r: '\x00', want: true},
		{name: "NUL absent", b: []byte("service/admin"), r: '\x00', want: false},
		{name: "UTF-8 rune", b: []byte("café"), r: 'é', want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.want, bytes.ContainsRune(test.b, test.r))
		})
	}
}
