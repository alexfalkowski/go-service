package xid

import "github.com/rs/xid"

// NewGenerator creates a new XID generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator for XIDs.
type Generator struct{}

// Generate an XID.
func (x *Generator) Generate() string {
	id := xid.New()
	return id.String()
}
