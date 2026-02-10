// Package time provides time-related helpers and abstractions used by go-service.
//
// This package wraps a subset of the standard library time API, provides convenience helpers such as MustParseDuration,
// and optionally supports fetching the current time from network time providers (for example NTP or NTS) via Network.
package time
