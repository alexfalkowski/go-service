package time

import (
	"github.com/alexfalkowski/go-service/errors"
	"github.com/beevik/nts"
)

// NewNTSNetwork creates a new NTS network with an address.
// https://datatracker.ietf.org/doc/html/rfc8915
func NewNTSNetwork(address string) *NTSNetwork {
	return &NTSNetwork{address: address}
}

// NTSNetwork implements the Network interface for NTS.
type NTSNetwork struct {
	address string
}

// Now returns the current time from the NTS server.
func (n *NTSNetwork) Now() (Time, error) {
	session, err := nts.NewSession(n.address)
	if err != nil {
		return Time{}, errors.Prefix("nts", err)
	}

	res, err := session.Query()
	if err != nil {
		return Time{}, errors.Prefix("nts", err)
	}

	if err := res.Validate(); err != nil {
		return Time{}, errors.Prefix("nts", err)
	}

	return Now().Add(res.ClockOffset), nil
}
