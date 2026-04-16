package time

import "time"

// RFC3339 is the RFC3339 time format layout.
//
// It is an alias of time.RFC3339.
const RFC3339 = time.RFC3339

// Ticker is the go-service ticker type used across the repository.
//
// It is a type alias of time.Ticker, meaning it has identical semantics and method
// set to the standard library type.
type Ticker = time.Ticker

// Time is the go-service time type used across the repository.
//
// It is a type alias of time.Time, meaning it has identical semantics and method
// set to the standard library type.
type Time = time.Time

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
//
// This is a thin wrapper around time.After and does not change semantics.
func After(d Duration) <-chan Time {
	return time.After(d.Duration())
}

// NewTicker returns a new [Ticker] containing a channel that will send the current time on the channel after each tick.
//
// This is a thin wrapper around time.NewTicker and does not change semantics.
func NewTicker(d Duration) *Ticker {
	return time.NewTicker(d.Duration())
}

// Now returns the current local time.
//
// This is a thin wrapper around time.Now and does not change semantics.
func Now() Time {
	return time.Now()
}

// Since returns the time elapsed since t.
//
// This is a thin wrapper around time.Since and does not change semantics.
func Since(t Time) Duration {
	return Duration(time.Since(t))
}

// Sleep pauses the current goroutine for at least the duration d.
//
// This is a thin wrapper around time.Sleep and does not change semantics.
func Sleep(d Duration) {
	time.Sleep(d.Duration())
}
