package id

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	nanoid "github.com/matoous/go-nanoid"
)

// NanoID generator.
type NanoID struct{}

// Generate a NanoID.
func (n *NanoID) Generate() string {
	id, err := nanoid.Nanoid()
	runtime.Must(err)
	return id
}
