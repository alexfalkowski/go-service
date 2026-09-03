package compress_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

var compressorKinds = []string{"zstd", "s2", "snappy", "none"}

func TestMapReturnsRegisteredCompressors(t *testing.T) {
	for _, kind := range compressorKinds {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)

			data := strings.Bytes("hello")
			compressed, err := cmp.Compress(data, bytes.KB)
			require.NoError(t, err)

			decompressed, err := cmp.Decompress(compressed, bytes.KB)
			require.NoError(t, err)
			require.Equal(t, data, decompressed)
		})
	}

	for _, key := range []string{"test", "bob"} {
		t.Run("returns nil for "+key, func(t *testing.T) {
			cmp := test.Compressor.Get(key)
			require.Nil(t, cmp)
		})
	}
}

func TestMapRegisterReplacesExistingCompressor(t *testing.T) {
	compressors := compress.NewMap()
	custom := test.NewCompressor(test.ErrFailed)
	replacement := none.NewCompressor()

	compressors.Register("custom", custom)
	require.Same(t, custom, compressors.Get("custom"))

	compressors.Register("custom", replacement)
	require.Same(t, replacement, compressors.Get("custom"))
}
