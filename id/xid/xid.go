package xid

import "github.com/rs/xid"

// ID is an alias of [xid.ID].
type ID = xid.ID

// Parse decodes an XID string using [xid.FromString].
func Parse(s string) (ID, error) {
	return xid.FromString(s)
}

// NewGenerator constructs an XID generator.
//
// The returned generator produces XID identifiers (globally unique identifiers that are compact and
// roughly sortable) via [github.com/rs/xid].
//
// XIDs intentionally expose ordering characteristics and are not designed to be opaque or
// unpredictable. Use this generator for compact sortable IDs, not for secrets or bearer values.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator generates XID identifiers.
//
// XID values are compact and roughly sortable, but they are not opaque security tokens.
type Generator struct{}

// Generate returns a newly generated XID string.
//
// It calls [xid.New] and returns the canonical string representation of the identifier.
func (x *Generator) Generate() string {
	id := xid.New()
	return id.String()
}
