package stream

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrUnsupportedMedia is returned when a streaming media type has no registered streaming codec.
var ErrUnsupportedMedia = errors.New("stream: unsupported media")

// ErrBidiRequiresHTTP2 is returned when a bidirectional streaming route is called over HTTP/1.x.
var ErrBidiRequiresHTTP2 = errors.New("stream: bidirectional streaming requires HTTP/2")

// ErrDraining is returned by [Stream.Send] and [RequestStream.Recv] after the configured drain starts.
var ErrDraining = errors.New("stream: draining")
