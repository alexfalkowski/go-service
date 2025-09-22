package none

// NewCompressor for none.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor for none.
type Compressor struct{}

func (c *Compressor) Compress(data []byte) []byte {
	return data
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}
