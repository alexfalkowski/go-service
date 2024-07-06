package compress

// Compressor allows to have different ways to compress/decompress.
type Compressor interface {
	// Compress data.
	Compress(data []byte) []byte

	// Decompress data.
	Decompress(data []byte) ([]byte, error)
}
