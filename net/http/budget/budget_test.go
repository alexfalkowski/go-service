package budget_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	streamjson "github.com/alexfalkowski/go-service/v2/encoding/stream/json"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http/budget"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestReaderReadDeliversValueAtExactBudget(t *testing.T) {
	t.Parallel()

	r := budget.NewReader(strings.NewReader("ab"), 2)
	buf := make([]byte, 2)

	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	// A second read simulates a decoder's probe past the value's last byte; at exactly the budget, it
	// must still reach the real underlying EOF instead of a synthetic refusal.
	n, err = r.Read(buf)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, 0, n)
}

func TestReaderReadRefusesOnceBudgetExceeded(t *testing.T) {
	t.Parallel()

	// The guard is strictly greater-than (see TestReaderReadDeliversValueAtExactBudget), so one read
	// past the budget still succeeds; a fourth byte is needed to observe an actual refusal at max=2.
	r := budget.NewReader(strings.NewReader("abcd"), 2)
	buf := make([]byte, 1)

	for range 3 {
		_, err := r.Read(buf)
		require.NoError(t, err)
	}

	n, err := r.Read(buf)
	require.ErrorIs(t, err, budget.ErrExceeded)
	require.Equal(t, 0, n)

	// Sticky: a repeated call after the refusal returns the same error without consulting the
	// underlying reader again.
	_, err = r.Read(buf)
	require.ErrorIs(t, err, budget.ErrExceeded)
}

func TestReaderReadNeverRefusesWhenBudgetDisabled(t *testing.T) {
	t.Parallel()

	r := budget.NewReader(strings.NewReader(strings.Repeat("a", 10)), 0)
	buf := make([]byte, 10)

	n, err := r.Read(buf)

	require.NoError(t, err)
	require.Equal(t, 10, n)
}

func TestReaderResetClearsLatchedError(t *testing.T) {
	t.Parallel()

	r := budget.NewReader(strings.NewReader("ab"), 1)
	buf := make([]byte, 1)

	_, err := r.Read(buf)
	require.NoError(t, err)
	_, err = r.Read(buf)
	require.NoError(t, err)

	_, err = r.Read(buf)
	require.ErrorIs(t, err, budget.ErrExceeded)

	r.Reset(0)

	// The underlying reader is now exhausted, so a real io.EOF proves the latched error was cleared
	// rather than merely re-returned.
	_, err = r.Read(buf)
	require.ErrorIs(t, err, io.EOF)
}

func TestReaderResetExcludesReadAheadFromLiveCheck(t *testing.T) {
	t.Parallel()

	r := budget.NewReader(strings.NewReader("a"), 2)
	r.Reset(3)

	// The decoder read these bytes for the next value while handling the previous one. They must not
	// spend the next value's live budget before it reads anything itself.
	require.NoError(t, r.Err())

	n, err := r.Read(make([]byte, 1))
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.NoError(t, r.Err())
}

func TestReaderErrReturnsNilBeforeAnyRefusal(t *testing.T) {
	t.Parallel()

	r := budget.NewReader(strings.NewReader("a"), 10)

	require.NoError(t, r.Err())
}

func TestReaderErrReturnsLatchedError(t *testing.T) {
	t.Parallel()

	r := budget.NewReader(strings.NewReader("ab"), 1)
	buf := make([]byte, 1)

	_, _ = r.Read(buf)
	_, _ = r.Read(buf)
	_, _ = r.Read(buf)

	require.ErrorIs(t, r.Err(), budget.ErrExceeded)
}

func TestReaderExceeds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		max           int64
		read          int64
		bufferedAhead int64
		expected      bool
	}{
		{name: "under budget", max: 10, read: 5, bufferedAhead: 0, expected: false},
		{name: "exactly at budget", max: 10, read: 10, bufferedAhead: 0, expected: false},
		{name: "over budget", max: 10, read: 11, bufferedAhead: 0, expected: true},
		{name: "buffered ahead correction avoids false positive", max: 10, read: 15, bufferedAhead: 5, expected: false},
		{name: "buffered ahead correction still over", max: 10, read: 20, bufferedAhead: 5, expected: true},
		{name: "clamped to zero when buffered ahead exceeds read", max: 10, read: 3, bufferedAhead: 5, expected: false},
		{name: "disabled budget never exceeds", max: 0, read: 1000, bufferedAhead: 0, expected: false},
		{name: "negative budget never exceeds", max: -1, read: 1000, bufferedAhead: 0, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := budget.NewReader(bytes.NewReader(make([]byte, tt.read)), tt.max)

			if tt.read > 0 {
				_, err := r.Read(make([]byte, tt.read))
				require.NoError(t, err)
			}

			require.Equal(t, tt.expected, r.Exceeds(tt.bufferedAhead))
		})
	}
}

func TestBufferedLenReturnsDecoderBuffered(t *testing.T) {
	t.Parallel()

	decoder := streamjson.NewDecoder(strings.NewReader("{} {}"))

	var v struct{}
	require.NoError(t, decoder.Decode(&v))

	require.Equal(t, int64(3), budget.BufferedLen(decoder))
}

type noBufferedDecoder struct{}

func (noBufferedDecoder) Decode(any) error { return nil }
func (noBufferedDecoder) Close() error     { return nil }

func TestBufferedLenReturnsZeroWithoutBufferedMethod(t *testing.T) {
	t.Parallel()

	require.Equal(t, int64(0), budget.BufferedLen(noBufferedDecoder{}))
}

type wrongBufferedDecoder struct{}

func (wrongBufferedDecoder) Decode(any) error    { return nil }
func (wrongBufferedDecoder) Close() error        { return nil }
func (wrongBufferedDecoder) Buffered() io.Reader { return strings.NewReader("") }

func TestBufferedLenReturnsZeroForWrongBufferedType(t *testing.T) {
	t.Parallel()

	require.Equal(t, int64(0), budget.BufferedLen(wrongBufferedDecoder{}))
}
