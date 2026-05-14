package compress

import "github.com/alexfalkowski/go-service/v2/bytes"

// Compressor provides a common interface for compressing and decompressing byte slices.
//
// Implementations may use different algorithms and may return errors during decompression if the input
// is invalid or corrupted.
type Compressor interface {
	// Compress returns a compressed representation of data.
	//
	// The input data must not exceed size.
	Compress(data []byte, size bytes.Size) ([]byte, error)

	// Decompress returns the decompressed representation of data.
	//
	// The returned data must not exceed size.
	Decompress(data []byte, size bytes.Size) ([]byte, error)
}
