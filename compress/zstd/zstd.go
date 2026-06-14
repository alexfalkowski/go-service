package zstd

import (
	"fmt"
	"math"

	"github.com/alexfalkowski/go-service/v2/bytes"
	compress "github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/klauspost/compress/zstd"
)

// MinWindowSize is the minimum decoder window size used by the underlying Zstandard implementation.
const MinWindowSize = zstd.MinWindowSize

const decoderMaxMemory = bytes.MaxConfigSize

// ErrDecoderSizeExceeded is returned by the underlying Zstandard decoder when decoded size exceeds its limit.
var ErrDecoderSizeExceeded = zstd.ErrDecoderSizeExceeded

// NewCompressor constructs a Zstandard (zstd) compressor implementation.
//
// The returned value implements [github.com/alexfalkowski/go-service/v2/compress.Compressor].
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements Zstandard (zstd) compression.
//
// It satisfies the [github.com/alexfalkowski/go-service/v2/compress.Compressor] interface.
type Compressor struct{}

// Compress returns the zstd-compressed representation of data.
//
// This method uses the klauspost/compress zstd encoder.
//
// An error is returned if data exceeds size.
func (c *Compressor) Compress(data []byte, size bytes.Size) ([]byte, error) {
	if int64(len(data)) > size.Bytes() {
		return nil, compress.ErrTooLarge
	}

	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}
	defer encoder.Close()

	return encoder.EncodeAll(data, nil), nil
}

// Decompress returns the decompressed representation of data.
//
// An error is returned if data is not valid zstd-encoded content or the decompressed data exceeds size.
func (c *Compressor) Decompress(data []byte, size bytes.Size) ([]byte, error) {
	limit := size.Bytes()
	if limit < 0 || limit == math.MaxInt64 {
		return nil, compress.ErrTooLarge
	}

	decoder, err := zstd.NewReader(
		bytes.NewReader(data),
		zstd.WithDecoderMaxMemory(uint64(decoderMaxMemory)),
	)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	decoded, _, err := io.ReadAll(io.LimitReader(decoder, limit+1))
	if err != nil {
		if errors.Is(err, ErrDecoderSizeExceeded) {
			return nil, fmt.Errorf("%w: %w", compress.ErrTooLarge, err)
		}

		return nil, err
	}
	if int64(len(decoded)) > limit {
		return nil, fmt.Errorf("%w: %w", compress.ErrTooLarge, ErrDecoderSizeExceeded)
	}

	return decoded, nil
}
