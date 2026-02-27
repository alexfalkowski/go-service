package none

// NewCompressor constructs a no-op compressor implementation.
//
// The returned value implements `github.com/alexfalkowski/go-service/v2/compress.Compressor` and is
// typically registered under kind "none" to represent "no compression".
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements a no-op compression codec.
//
// Compress returns the input unchanged, and Decompress returns the input unchanged with a nil error.
// This is useful when you want to disable compression while still satisfying the common compression
// interface.
type Compressor struct{}

// Compress returns data unchanged.
func (c *Compressor) Compress(data []byte) []byte {
	return data
}

// Decompress returns data unchanged.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}
