package nsq

import (
	"github.com/alexfalkowski/go-service/marshaller"
)

// NewMsgPackMarshaller for NSQ.
func NewMsgPackMarshaller() Marshaller {
	return marshaller.NewMsgPack()
}
