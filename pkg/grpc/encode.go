package grpc

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func encode(message interface{}) string {
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
