package marshaller

// None for marshaller.
type None struct{}

// NewNone for marshaller.
func NewNone() *None {
	return &None{}
}

func (*None) Marshal(_ any) ([]byte, error) {
	return nil, nil
}

func (*None) Unmarshal(_ []byte, _ any) error {
	return nil
}
