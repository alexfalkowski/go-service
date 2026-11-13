package compress_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/bytes"
	"github.com/alexfalkowski/go-service/v3/compress"
	"github.com/alexfalkowski/go-service/v3/compress/none"
	"github.com/alexfalkowski/go-service/v3/compress/s2"
	"github.com/alexfalkowski/go-service/v3/compress/snappy"
	"github.com/alexfalkowski/go-service/v3/compress/zstd"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/strings"
	"github.com/alexfalkowski/go-signal/v2"
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

func TestNewMapRegistersDefaultCompressors(t *testing.T) {
	compressors := compress.NewMap()

	expected := map[string]compress.Compressor{
		"zstd":   zstd.NewCompressor(),
		"s2":     s2.NewCompressor(),
		"snappy": snappy.NewCompressor(),
		"none":   none.NewCompressor(),
	}

	for kind, expectedCompressor := range expected {
		t.Run(kind, func(t *testing.T) {
			require.IsType(t, expectedCompressor, compressors.Get(kind))
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

func TestModuleProvidesDefaultCompressors(t *testing.T) {
	var compressors *compress.Map

	_, err := signal.NewApp(
		compress.Module,
		signal.Register(func(value *compress.Map) { compressors = value }),
	)

	require.NoError(t, err)
	for _, kind := range compressorKinds {
		t.Run(kind, func(t *testing.T) {
			require.NotNil(t, compressors.Get(kind))
		})
	}
}

func TestModuleDoesNotProvideZstdCompressor(t *testing.T) {
	_, err := signal.NewApp(
		compress.Module,
		signal.Register(func(*zstd.Compressor) {}),
	)

	require.Error(t, err)
}
