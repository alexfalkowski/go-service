package id

import (
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/oklog/ulid"
)

// NewULID creates a new ULID generator.
func NewULID(reader rand.Reader) *ULID {
	return &ULID{
		reader: reader,
	}
}

// ULID generator.
type ULID struct {
	reader rand.Reader
}

// Generate a ULID.
func (k *ULID) Generate() string {
	id, err := ulid.New(ulid.Timestamp(time.Now()), k.reader)
	runtime.Must(err)

	return id.String()
}
