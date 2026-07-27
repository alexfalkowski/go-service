package compress_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

var compressorKinds = []string{"zstd", "s2", "snappy", "none"}

func TestMap(t *testing.T) {
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
		t.Run(key, func(t *testing.T) {
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

func TestMapRegister(t *testing.T) {
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

	app := fx.New(
		compress.Module,
		fx.Populate(&compressors),
		fx.NopLogger,
	)

	require.NoError(t, app.Err())
	for _, kind := range compressorKinds {
		t.Run(kind, func(t *testing.T) {
			require.NotNil(t, compressors.Get(kind))
		})
	}
}

func TestModuleDoesNotProvideZstdCompressor(t *testing.T) {
	var compressor *zstd.Compressor

	app := fx.New(
		compress.Module,
		fx.Populate(&compressor),
		fx.NopLogger,
	)

	require.Error(t, app.Err())
}
