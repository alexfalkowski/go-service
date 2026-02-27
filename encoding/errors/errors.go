package errors

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrInvalidType indicates that an Encoder cannot operate on the provided value type.
//
// It is returned by encoding adapters that require the input value to implement a specific interface
// (for example io.WriterTo/io.ReaderFrom for the bytes passthrough encoder, or proto.Message for protobuf
// encoders).
var ErrInvalidType = errors.New("encoding: invalid type")
