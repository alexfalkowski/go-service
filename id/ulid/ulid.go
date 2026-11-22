package ulid

import (
	"github.com/alexfalkowski/go-service/v3/crypto/rand"
	"github.com/alexfalkowski/go-service/v3/runtime"
	"github.com/alexfalkowski/go-service/v3/time"
	"github.com/oklog/ulid"
)

// ULID is an alias of [ulid.ULID].
type ULID = ulid.ULID

// Parse decodes a ULID string using [ulid.Parse].
// It assumes the input uses valid base32 characters.
func Parse(s string) (ULID, error) {
	return ulid.Parse(s)
}

// NewGenerator constructs a ULID generator.
//
// The returned generator produces ULIDs (Universally Unique Lexicographically Sortable Identifiers).
// ULIDs include a time component, making their string representation roughly sortable by creation time,
// and a randomness component sourced from the provided reader.
func NewGenerator(reader rand.Reader) *Generator {
	return &Generator{reader: reader}
}

// Generator generates ULID identifiers.
type Generator struct {
	// reader is the randomness source used for the ULID entropy component.
	reader rand.Reader
}

// Generate returns a newly generated ULID string.
//
// The ULID timestamp component is derived from [time.Now](), and the entropy component is read from the
// injected cryptographically secure reader.
//
// If ULID generation fails (for example due to reader errors), this method panics via [runtime.Must].
func (g *Generator) Generate() string {
	id, err := ulid.New(ulid.Timestamp(time.Now()), g.reader)
	runtime.Must(err)
	return id.String()
}
