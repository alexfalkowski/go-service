package time

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/beevik/ntp"
)

// NewNTPNetwork constructs a Network implementation backed by NTP (Network Time Protocol).
//
// NTP is a widely deployed protocol for synchronizing clocks over packet-switched networks.
// See: https://en.wikipedia.org/wiki/Network_Time_Protocol
//
// The address argument is passed through to the upstream NTP client implementation and is
// typically a hostname (for example an NTP pool name) or host:port depending on the client
// library behavior. If address is empty or invalid, calls to Now will typically return an error.
//
// If timeout is zero, the upstream client's default timeout is used.
func NewNTPNetwork(address string, timeout Duration) *NTPNetwork {
	return &NTPNetwork{address: address, timeout: timeout}
}

// NTPNetwork implements Network by querying an NTP server for the current time.
//
// This type is a small adapter around [github.com/beevik/ntp]. It prefixes returned errors
// with "ntp" to make failures easier to attribute when multiple time sources are possible.
type NTPNetwork struct {
	address string
	timeout Duration
}

// Now returns the current time as reported by the configured NTP server.
//
// This method performs network I/O. Any error returned by the underlying NTP library is
// wrapped/prefixed with "ntp" for context.
func (n *NTPNetwork) Now() (Time, error) {
	res, err := ntp.QueryWithOptions(n.address, ntp.QueryOptions{Timeout: n.timeout.Duration()})
	if err != nil {
		return Time{}, errors.Prefix("ntp", err)
	}

	if err := res.Validate(); err != nil {
		return Time{}, errors.Prefix("ntp", err)
	}

	return Now().Add(res.ClockOffset), nil
}
