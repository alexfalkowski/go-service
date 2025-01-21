package id

import (
	"crypto/rand"
	"time"

	"github.com/alexfalkowski/go-service/runtime"
	"github.com/oklog/ulid"
)

// ULID generator.
type ULID struct{}

// Generate a ULID.
func (k *ULID) Generate() string {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	runtime.Must(err)

	return id.String()
}
