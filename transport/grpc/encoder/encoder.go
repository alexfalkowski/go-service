package encoder

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Message to be encoded.
func Message(message interface{}) string {
	m, err := protojson.Marshal(message.(proto.Message))
	if err != nil {
		return ""
	}

	return string(m)
}
