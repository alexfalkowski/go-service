package ulid

import (
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/oklog/ulid"
)

// NewGenerator creates a new ULID generator.
func NewGenerator(reader rand.Reader) *Generator {
	return &Generator{reader: reader}
}

// Generator for ULID.
type Generator struct {
	reader rand.Reader
}

// Generate a ULID.
func (k *Generator) Generate() string {
	id, err := ulid.New(ulid.Timestamp(time.Now()), k.reader)
	runtime.Must(err)
	return id.String()
}
