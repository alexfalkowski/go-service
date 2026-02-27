package time

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/beevik/nts"
)

// NewNTSNetwork constructs a Network implementation backed by NTS (Network Time Security).
//
// NTS provides authenticated time over the network, improving on NTP by protecting against
// certain classes of on-path and server spoofing attacks.
// See: https://datatracker.ietf.org/doc/html/rfc8915
//
// The address argument is passed through to the upstream NTS client implementation and is
// typically a hostname (and possibly port) of an NTS-enabled server. If address is empty or
// invalid, calls to Now will typically return an error.
func NewNTSNetwork(address string) *NTSNetwork {
	return &NTSNetwork{address: address}
}

// NTSNetwork implements Network by querying an NTS server and validating the response.
//
// This type is a small adapter around github.com/beevik/nts. It prefixes returned errors
// with "nts" to make failures easier to attribute when multiple time sources are possible.
//
// Note: NTS returns a clock offset relative to the local clock. This implementation returns
// Now().Add(offset), which means the returned value is derived from the local time adjusted
// by the authenticated offset.
type NTSNetwork struct {
	address string
}

// Now returns the current time as reported by the configured NTS server.
//
// This method performs network I/O and validates the NTS response before returning.
//
// The algorithm is:
//   - Establish a session (nts.NewSession).
//   - Query the server (session.Query).
//   - Validate the response (res.Validate).
//   - Apply the server-provided clock offset to the local time.
//
// Any error returned by the underlying NTS library is wrapped/prefixed with "nts" for context.
func (n *NTSNetwork) Now() (Time, error) {
	session, err := nts.NewSession(n.address)
	if err != nil {
		return Time{}, n.prefix(err)
	}

	res, err := session.Query()
	if err != nil {
		return Time{}, n.prefix(err)
	}

	if err := res.Validate(); err != nil {
		return Time{}, n.prefix(err)
	}

	return Now().Add(res.ClockOffset), nil
}

func (n *NTSNetwork) prefix(err error) error {
	return errors.Prefix("nts", err)
}
