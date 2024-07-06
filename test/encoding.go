package test

import (
	"github.com/alexfalkowski/go-service/encoding"
)

// Marshaller for tests.
var Marshaller = encoding.NewMap()

// NewMarshaller for test.
func NewMarshaller(err error) encoding.Marshaller {
	return &mar{err: err}
}

type mar struct {
	err error
}

func (m *mar) Marshal(_ any) ([]byte, error) {
	return nil, m.err
}

func (m *mar) Unmarshal(_ []byte, _ any) error {
	return m.err
}
