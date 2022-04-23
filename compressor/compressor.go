package compressor

// Compressor allows to have different ways to compress/decompress.
type Compressor interface {
	Compress(src []byte) []byte
	Decompress(src []byte) ([]byte, error)
}
