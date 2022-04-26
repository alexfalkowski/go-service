package marshaller

import (
	"github.com/alexfalkowski/go-service/marshaller"
)

// NewMsgPack for NSQ.
// nolint:ireturn
func NewMsgPack() Marshaller {
	return marshaller.NewMsgPack()
}
