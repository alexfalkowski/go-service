package stream

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Options configures the streaming route helpers built by [NewHandler] and [NewRequestHandler].
//
// The zero value disables every knob: no deadline extension, no per-value receive cap, and no drain
// signal, matching today's behavior when no options are supplied.
type Options struct {
	// Drain closes when the server starts draining. A nil channel disables drain handling, preserving
	// the behavior of manually registered routes that do not supply a server lifecycle signal.
	Drain <-chan struct{}

	// ReadTimeout is the per-message inactivity budget applied to [RequestStream.Recv]'s request read
	// deadline. Zero disables read deadline extension. Unused by [NewHandler]'s send-only stream.
	ReadTimeout time.Duration

	// WriteTimeout is the per-message inactivity budget applied to [Stream.Send]'s response write
	// deadline. Zero disables write deadline extension.
	WriteTimeout time.Duration

	// MaxReceiveSize bounds each value decoded by [RequestStream.Recv], not the request stream as a
	// whole (see [github.com/alexfalkowski/go-service/v2/net/http/quota.Reader]). Zero disables the
	// per-value cap. Unused by [NewHandler]'s send-only stream.
	MaxReceiveSize bytes.Size
}
