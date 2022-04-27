package test

import (
	"github.com/alexfalkowski/go-service/marshaller"
)

// NewMarshaller for test.
func NewMarshaller(err error) marshaller.Marshaller {
	return &mar{err: err}
}

type mar struct {
	err error
}

func (m *mar) Marshal(v any) ([]byte, error) {
	return nil, m.err
}

func (m *mar) Unmarshal(data []byte, v any) error {
	return m.err
}
