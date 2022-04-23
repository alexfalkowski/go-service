package marshaller

// Compressor allows to have different ways to marshal/unmarshal.
type Marshaller interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}
