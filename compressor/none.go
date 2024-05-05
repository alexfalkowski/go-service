package compressor

// None for compressor.
type None struct{}

// NewNone compressor.
func NewNone() *None {
	return &None{}
}

func (c *None) Compress(data []byte) []byte {
	return data
}

func (c *None) Decompress(data []byte) ([]byte, error) {
	return data, nil
}
