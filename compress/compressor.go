package compress

import "github.com/alexfalkowski/go-service/v2/bytes"

// Compressor provides a common interface for compressing and decompressing byte slices.
//
// Implementations may use different algorithms and may return errors during
// decompression if the input is invalid or corrupted.
//
// Size limits are applied to uncompressed data: [Compressor.Compress] rejects
// input larger than size, and [Compressor.Decompress] rejects decoded output
// larger than size. Implementations return
// [github.com/alexfalkowski/go-service/v2/compress/errors.ErrTooLarge] when a
// size limit is exceeded.
type Compressor interface {
	// Compress returns a compressed representation of data.
	//
	// The input data must not exceed size. Implementations return
	// [github.com/alexfalkowski/go-service/v2/compress/errors.ErrTooLarge]
	// when data is larger than size.
	Compress(data []byte, size bytes.Size) ([]byte, error)

	// Decompress returns the decompressed representation of data.
	//
	// The decoded output must not exceed size. Implementations return
	// [github.com/alexfalkowski/go-service/v2/compress/errors.ErrTooLarge]
	// when decoded data is larger than size.
	Decompress(data []byte, size bytes.Size) ([]byte, error)
}
