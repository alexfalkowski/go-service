package xid

import "github.com/rs/xid"

// NewGenerator constructs an XID generator.
//
// The returned generator produces XID identifiers (globally unique identifiers that are compact and
// roughly sortable) via github.com/rs/xid.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator generates XID identifiers.
type Generator struct{}

// Generate returns a newly generated XID string.
//
// It calls xid.New and returns the canonical string representation of the identifier.
func (x *Generator) Generate() string {
	id := xid.New()
	return id.String()
}
