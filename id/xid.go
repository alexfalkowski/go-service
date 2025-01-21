package id

import (
	"github.com/rs/xid"
)

// XID generator.
type XID struct{}

// Generate an XID.
func (x *XID) Generate() string {
	id := xid.New()

	return id.String()
}
