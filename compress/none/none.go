package none

// Compressor for none.
type Compressor struct{}

// NewNone for none.
func NewCompressor() *Compressor {
	return &Compressor{}
}

func (c *Compressor) Compress(data []byte) []byte {
	return data
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}
