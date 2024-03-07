package marshaller

// Marshaller allows to have different ways to marshal/unmarshal.
type Marshaller interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}
