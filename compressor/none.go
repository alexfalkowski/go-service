package compressor

// None for compressor.
type None struct{}

// NewNone compressor.
func NewNone() *None {
	return &None{}
}

func (c *None) Compress(src []byte) []byte {
	return src
}

func (c *None) Decompress(src []byte) ([]byte, error) {
	return src, nil
}
