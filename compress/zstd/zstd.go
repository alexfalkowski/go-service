package zstd

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/bytes"
	compress "github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/klauspost/compress/zstd"
)

// MinWindowSize is the minimum decoder window size used by the underlying Zstandard implementation.
const MinWindowSize = zstd.MinWindowSize

// ErrDecoderSizeExceeded is returned by the underlying Zstandard decoder when decoded size exceeds its limit.
var ErrDecoderSizeExceeded = zstd.ErrDecoderSizeExceeded

// NewCompressor constructs a Zstandard (zstd) compressor implementation.
//
// The returned value implements `github.com/alexfalkowski/go-service/v2/compress.Compressor`.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements Zstandard (zstd) compression.
//
// It satisfies the `github.com/alexfalkowski/go-service/v2/compress.Compressor` interface.
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

	e, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}
	defer e.Close()

	return e.EncodeAll(data, nil), nil
}

// Decompress returns the decompressed representation of data.
//
// An error is returned if data is not valid zstd-encoded content or the decompressed data exceeds size.
func (c *Compressor) Decompress(data []byte, size bytes.Size) ([]byte, error) {
	limit := size.Bytes()
	maxMemory := uint64(max(limit, int64(MinWindowSize)))
	d, err := zstd.NewReader(
		bytes.NewReader(data),
		zstd.WithDecoderMaxMemory(maxMemory),
	)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	decoded, _, err := io.ReadAll(io.LimitReader(d, limit+1))
	if err != nil {
		if errors.Is(err, ErrDecoderSizeExceeded) {
			return nil, fmt.Errorf("%w: %w", compress.ErrTooLarge, err)
		}

		return nil, err
	}
	if int64(len(decoded)) > limit {
		return nil, compress.ErrTooLarge
	}

	return decoded, nil
}
