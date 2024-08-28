package test

import (
	"io"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/sync"
)

// Pool for tests.
var Pool = sync.NewBufferPool()

// Marshaller for tests.
var Marshaller = encoding.NewMarshallerMap()

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

// Encoder for tests.
var Encoder = encoding.NewEncoderMap()

// NewEncoder for test.
func NewEncoder(err error) encoding.Encoder {
	return &enc{err: err}
}

type enc struct {
	err error
}

func (e *enc) Encode(_ io.Writer, _ any) error {
	return e.err
}

func (e *enc) Decode(_ io.Reader, _ any) error {
	return e.err
}
