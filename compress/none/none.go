package none

// NewCompressor constructs a no-op compressor.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements a no-op compressor that returns input unchanged.
type Compressor struct{}

// Compress returns the input unchanged.
func (c *Compressor) Compress(data []byte) []byte {
	return data
}

// Decompress returns the input unchanged.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}
