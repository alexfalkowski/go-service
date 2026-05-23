package errors

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
)

// ErrInvalidType indicates that an Encoder cannot operate on the provided value type.
//
// It is returned by encoding adapters that require the input value to implement a specific interface
// (for example io.WriterTo/io.ReaderFrom for the bytes passthrough encoder, or proto.Message for protobuf
// encoders).
var ErrInvalidType = errors.New("encoding: invalid type")

// ErrTrailingData indicates that a decoder found another value after the first decoded payload.
var ErrTrailingData = errors.New("encoding: trailing data")

// TrailingData returns nil for EOF and ErrTrailingData for any extra decoded value or parse error.
func TrailingData(err error) error {
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err == nil {
		return ErrTrailingData
	}

	return fmt.Errorf("%w: %w", ErrTrailingData, err)
}
