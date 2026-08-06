package budget

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
)

// ErrExceeded is the sentinel [Reader.Read] latches, and [Reader.Err] returns, once a value's bytes
// exceed the configured budget.
//
// ErrExceeded carries no data of its own: the budget's limit is already known to whichever caller
// constructed the [Reader] with it, so that caller builds whatever concrete, contract-specific error it
// needs to return (for example a
// [github.com/alexfalkowski/go-service/v2/net/http.MaxBytesError] carrying the limit, or a
// [github.com/alexfalkowski/go-service/v2/net/http/status.SafeError]) once [Reader.Err] or
// [Reader.Exceeds] tells it the budget was spent, rather than asking Reader to build that error itself.
var ErrExceeded = errors.New("budget: exceeded")

// BufferedLen returns the number of bytes decoder has already pulled from its underlying reader for a
// value it has not decoded yet, so a caller can subtract them from the raw bytes a [Reader] counted for
// the value it just decoded (see [Reader.Exceeds]).
//
// This only recognizes decoders whose Buffered method matches [encoding/json.Decoder.Buffered]'s shape
// and returns a concrete [bytes.Reader]. It returns 0 for any decoder that does not match, which is
// always safe: the correction is skipped, not a wrong non-zero answer, so the check falls back to the
// same "bound on reads attributed to one value" behavior documented on [Reader] rather than any
// incorrect result.
func BufferedLen(decoder stream.Decoder) int64 {
	buffered, ok := decoder.(interface{ Buffered() io.Reader })
	if !ok {
		return 0
	}

	reader, ok := buffered.Buffered().(*bytes.Reader)
	if !ok {
		return 0
	}

	return int64(reader.Len())
}

// NewReader constructs a Reader wrapping r with a per-value byte budget of limit. limit <= 0 disables
// the budget entirely.
func NewReader(r io.Reader, limit int64) *Reader {
	return &Reader{r: r, max: limit}
}

// Reader wraps an [io.Reader] with a resettable per-value byte budget, so a stream decoder bound to it
// for a whole body gets a per-value cap rather than a cumulative one. The caller resets the budget via
// [Reader.Reset] at the exact point a decoded value boundary is known.
//
// Reader never truncates or discards a Read call's data: every byte pulled from the underlying reader is
// always returned to the caller, so a decoder's own internal buffering can never be desynchronized by
// this guard. Read refuses to pull more data once the budget from a previous Read within the same value
// is already exceeded — a coarse, live backstop bounding a single pathological value across repeated
// Read calls, not the authoritative bound. The guard is strictly greater-than, not at-or-above: a value
// landing exactly on the budget still reaches the underlying reader for its terminating read (which may
// legitimately be [io.EOF]), rather than a synthetic refusal. A caller wanting the precise, corrected
// bound uses [Reader.Exceeds] after a successful decode; see its bufferedAhead correction.
//
// A budget of zero or less disables it outright: Read never refuses to read and Exceeds always reports
// false.
type Reader struct {
	r    io.Reader
	max  int64
	read int64
}

// Reset rearms the byte counter for the next decoded value. bufferedAhead has already been read from the
// underlying reader by the decoder while handling the previous value, so it starts this value's
// accounting rather than being discarded when the counter resets.
func (r *Reader) Reset(bufferedAhead int64) {
	r.read = bufferedAhead
}

// Exceeds reports whether the bytes read since the last [Reader.Reset], minus bufferedAhead, exceed the
// configured budget. bufferedAhead is the decoder's own read-ahead for a later value; subtracting it
// corrects for a decoder that pulled more than one value's worth of bytes out of the reader in a single
// underlying Read. The corrected count is clamped to zero, and this always reports false when the budget
// is disabled (max <= 0). This is the post-decode check; [Reader.Read] applies its own uncorrected, live
// check to bound a single pathological value across repeated Read calls.
func (r *Reader) Exceeds(bufferedAhead int64) bool {
	if r.max <= 0 {
		return false
	}

	attributed := max(r.read-bufferedAhead, 0)

	return attributed > r.max
}

// Read implements [io.Reader]. Once the bytes read since the last [Reader.Reset] already exceed the
// configured budget, Read refuses to read more and returns [ErrExceeded], and every subsequent call
// returns that same error until Reset is called. Read never refuses to read when the budget is disabled
// (max <= 0).
func (r *Reader) Read(p []byte) (int, error) {
	if r.exceededLive() {
		return 0, ErrExceeded
	}

	n, err := r.r.Read(p)
	r.read += int64(n)

	return n, err
}

// Err returns [ErrExceeded] if the budget is already spent — the same condition a subsequent
// [Reader.Read] would refuse on — or nil otherwise.
func (r *Reader) Err() error {
	if r.exceededLive() {
		return ErrExceeded
	}

	return nil
}

// exceededLive reports the same live, uncorrected condition [Reader.Read] and [Reader.Err] both answer:
// the budget is enabled and the raw count since the last [Reader.Reset] already exceeds it. This needs
// no cached "latched" state of its own — r.read only advances when Read actually reads more from r,
// which it refuses to do once this already reports true, so re-evaluating it on every call is exactly as
// sticky as caching its first result would be, without an extra field to keep in sync with Reset.
func (r *Reader) exceededLive() bool {
	return r.max > 0 && r.read > r.max
}
