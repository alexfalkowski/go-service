package marshaller

// Marshaller allows to have different ways to marshal/unmarshal.
type Marshaller interface {
	// Marshal value.
	Marshal(v any) ([]byte, error)

	// Unmarshal data to value.
	Unmarshal(data []byte, v any) error
}
