package compress_test

import (
	"math"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestMap(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)

			data := strings.Bytes("hello")
			d, err := cmp.Compress(data, bytes.KB)
			require.NoError(t, err)

			ns, err := cmp.Decompress(d, bytes.KB)
			require.NoError(t, err)
			require.Equal(t, data, ns)
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
	zstdCompressor := zstd.NewCompressor()
	s2Compressor := s2.NewCompressor()
	snappyCompressor := snappy.NewCompressor()
	noneCompressor := none.NewCompressor()

	compressors := compress.NewMap(compress.MapParams{
		Zstd:   zstdCompressor,
		S2:     s2Compressor,
		Snappy: snappyCompressor,
		None:   noneCompressor,
	})

	require.Same(t, zstdCompressor, compressors.Get("zstd"))
	require.Same(t, s2Compressor, compressors.Get("s2"))
	require.Same(t, snappyCompressor, compressors.Get("snappy"))
	require.Same(t, noneCompressor, compressors.Get("none"))
}

func TestMapRegister(t *testing.T) {
	compressors := compress.NewMap(compress.MapParams{})
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
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		t.Run(kind, func(t *testing.T) {
			require.NotNil(t, compressors.Get(kind))
		})
	}
}

func TestNoneCompressorReturnsDataUnchanged(t *testing.T) {
	cmp := none.NewCompressor()
	data := strings.Bytes("hello")

	compressed, err := cmp.Compress(data, bytes.Size(len(data)))
	require.NoError(t, err)
	require.Equal(t, data, compressed)

	decompressed, err := cmp.Decompress(data, bytes.Size(len(data)))
	require.NoError(t, err)
	require.Equal(t, data, decompressed)
}

func TestExactSizeLimits(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)
			data := strings.Bytes("hello")
			size := bytes.Size(len(data))

			compressed, err := cmp.Compress(data, size)
			require.NoError(t, err)

			decompressed, err := cmp.Decompress(compressed, size)
			require.NoError(t, err)
			require.Equal(t, data, decompressed)
		})
	}
}

func TestCompressRejectsTooLarge(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)

			_, err := cmp.Compress(strings.Bytes("hello"), 4)
			require.ErrorIs(t, err, errors.ErrTooLarge)
		})
	}
}

func TestDecompressRejectsTooLarge(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)

			data := strings.Bytes("hello")
			d, err := cmp.Compress(data, bytes.KB)
			require.NoError(t, err)

			_, err = cmp.Decompress(d, bytes.Size(len(data)-1))
			require.ErrorIs(t, err, errors.ErrTooLarge)
		})
	}
}

func TestDecompressRejectsInvalidData(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy"} {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)

			_, err := cmp.Decompress(strings.Bytes("invalid"), bytes.KB)
			require.Error(t, err)
			require.NotErrorIs(t, err, errors.ErrTooLarge)
		})
	}
}

func TestS2DecompressReturnsDecodedLenError(t *testing.T) {
	cmp := test.Compressor.Get("s2")
	data := []byte{0x80}

	_, err := cmp.Decompress(data, bytes.KB)
	require.Error(t, err)
	require.NotErrorIs(t, err, errors.ErrTooLarge)
}

func TestSnappyDecompressReturnsDecodedLenError(t *testing.T) {
	cmp := test.Compressor.Get("snappy")
	data := []byte{0x80}

	_, err := cmp.Decompress(data, bytes.KB)
	require.Error(t, err)
	require.NotErrorIs(t, err, errors.ErrTooLarge)
}

func TestZstdDecompressPreservesDecoderSizeError(t *testing.T) {
	cmp := zstd.NewCompressor()

	data := make([]byte, zstd.MinWindowSize+1)
	encoded, err := cmp.Compress(data, bytes.Size(len(data)))
	require.NoError(t, err)

	_, err = cmp.Decompress(encoded, zstd.MinWindowSize)
	require.ErrorIs(t, err, errors.ErrTooLarge)
	require.ErrorIs(t, err, zstd.ErrDecoderSizeExceeded)
}

func TestZstdDecompressRejectsInvalidLimits(t *testing.T) {
	cmp := zstd.NewCompressor()
	data := strings.Bytes("hello")
	encoded, err := cmp.Compress(data, bytes.KB)
	require.NoError(t, err)

	for _, tt := range []struct {
		name string
		size bytes.Size
	}{
		{name: "negative", size: -1},
		{name: "max int64", size: math.MaxInt64},
	} {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := cmp.Decompress(encoded, tt.size)
			require.ErrorIs(t, err, errors.ErrTooLarge)
			require.Nil(t, decoded)
		})
	}
}
