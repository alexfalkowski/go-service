package compress_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		cmp := test.Compressor.Get(kind)

		data := strings.Bytes("hello")
		d := cmp.Compress(data)

		ns, err := cmp.Decompress(d)
		require.NoError(t, err)
		require.Equal(t, data, ns)
	}

	for _, key := range []string{"test", "bob"} {
		cmp := test.Compressor.Get(key)
		require.Nil(t, cmp)
	}
}
