package encoder

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Message to be encoded.
func Message(message interface{}) string {
	p, ok := message.(proto.Message)
	if !ok {
		return ""
	}

	m, err := protojson.Marshal(p)
	if err != nil {
		return ""
	}

	return string(m)
}
