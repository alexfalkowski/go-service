package compress

// Compressor provides a common interface for compressing and decompressing byte slices.
//
// Implementations may use different algorithms and may return errors during decompression if the input
// is invalid or corrupted.
type Compressor interface {
	// Compress returns a compressed representation of data.
	Compress(data []byte) []byte

	// Decompress returns the decompressed representation of data.
	Decompress(data []byte) ([]byte, error)
}
