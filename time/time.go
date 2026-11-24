package time

import "time"

// Ticker is the go-service ticker type used across the repository.
//
// It is a type alias of [time.Ticker], meaning it has identical semantics and method
// set to the standard library type.
type Ticker = time.Ticker

// Timer is the go-service timer type used across the repository.
//
// It is a type alias of [time.Timer], meaning it has identical semantics and method
// set to the standard library type.
type Timer = time.Timer

// Time is the go-service time type used across the repository.
//
// It is a type alias of [time.Time], meaning it has identical semantics and method
// set to the standard library type.
type Time = time.Time

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
//
// This is a thin wrapper around [time.After] and does not change semantics.
func After(d Duration) <-chan Time {
	return time.After(d.Duration())
}

// AfterFunc schedules f to run in its own goroutine after d and returns a [Timer]
// whose Stop method can cancel a pending call. It wraps [time.AfterFunc].
func AfterFunc(d Duration, f func()) *Timer {
	return time.AfterFunc(d.Duration(), f)
}

// NewTimer returns a new [Timer] that will send the current time on its channel
// after at least duration d.
//
// This is a thin wrapper around [time.NewTimer] and does not change semantics.
func NewTimer(d Duration) *Timer {
	return time.NewTimer(d.Duration())
}

// NewTicker returns a new [Ticker] containing a channel that will send the
// current time on the channel after each tick.
//
// This is a thin wrapper around [time.NewTicker] and does not change semantics.
func NewTicker(d Duration) *Ticker {
	return time.NewTicker(d.Duration())
}

// Tick wraps [time.Tick] and returns the channel for ticks spaced by d.
// It returns a nil channel when d is zero or negative.
func Tick(d Duration) <-chan Time {
	return time.Tick(d.Duration())
}

// Now returns the current local time.
//
// This is a thin wrapper around [time.Now] and does not change semantics.
func Now() Time {
	return time.Now()
}

// Unix returns the local time sec seconds after the Unix epoch, with zero
// fractional seconds. It wraps [time.Unix].
func Unix(sec int64) Time {
	return time.Unix(sec, 0)
}

// UnixMilli wraps [time.UnixMilli] and returns the local time msec milliseconds
// after the Unix epoch.
func UnixMilli(msec int64) Time {
	return time.UnixMilli(msec)
}

// UnixMicro wraps [time.UnixMicro] and returns the local time usec microseconds
// after the Unix epoch.
func UnixMicro(usec int64) Time {
	return time.UnixMicro(usec)
}

// Since returns the time elapsed since t.
//
// This is a thin wrapper around [time.Since] and does not change semantics.
func Since(t Time) Duration {
	return Duration(time.Since(t))
}

// Until returns the duration until t.
//
// This is a thin wrapper around [time.Until] and does not change semantics.
func Until(t Time) Duration {
	return Duration(time.Until(t))
}

// Sleep pauses the current goroutine for at least the duration d.
//
// This is a thin wrapper around [time.Sleep] and does not change semantics.
func Sleep(d Duration) {
	time.Sleep(d.Duration())
}
