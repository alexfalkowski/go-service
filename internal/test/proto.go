package test

import "google.golang.org/protobuf/reflect/protoreflect"

// ErrProto for test.
type ErrProto struct{}

func (ErrProto) ProtoReflect() protoreflect.Message {
	return nil
}
