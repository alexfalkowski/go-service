package marshaller

import (
	"github.com/alexfalkowski/go-service/marshaller"
)

// NewMsgPack for NSQ.
func NewMsgPack() Marshaller {
	return marshaller.NewMsgPack()
}
