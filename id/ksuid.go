package id

import (
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/segmentio/ksuid"
)

// KSUID generator.
type KSUID struct{}

// Generate a KSUID.
func (k *KSUID) Generate() string {
	id, err := ksuid.NewRandom()
	runtime.Must(err)

	return id.String()
}
