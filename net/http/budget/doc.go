// Package budget provides a per-value byte budget for a stream decoder, shared by the HTTP server's
// bidirectional streaming request path
// ([github.com/alexfalkowski/go-service/v2/net/http/content.RequestStream.Recv]) and the HTTP client's
// streaming response path ([github.com/alexfalkowski/go-service/v2/net/http/client.ResponseStream.Recv]).
//
// [Reader] wraps an underlying reader with a byte counter the caller resets before every decoded value
// via [Reader.Reset], so a stream decoder bound to it for a whole body gets a per-value cap rather than a
// cumulative one. The reset call has to happen at the exact point a decoded value boundary is known,
// which is why this package only provides the counting reader itself: RequestStream.Recv and
// ResponseStream.Recv remain the callers that know when that boundary occurs. [Reader.Err] and
// [Reader.Exceeds] only ever report [ErrExceeded], a data-free sentinel; each caller already knows the
// configured limit from constructing the Reader with it, so each builds its own status-code mapping (413
// via [github.com/alexfalkowski/go-service/v2/net/http.MaxBytesError] on the server,
// [github.com/alexfalkowski/go-service/v2/net/http/status.SafeError] on the client) once it observes
// ErrExceeded, instead of asking this package to build that error itself.
//
// Start with [NewReader].
package budget
