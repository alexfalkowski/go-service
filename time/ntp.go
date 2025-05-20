package time

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/beevik/ntp"
)

// NewNTPNetwork creates a new NTP network with an address.
// https://en.wikipedia.org/wiki/Network_Time_Protocol
func NewNTPNetwork(address string) *NTPNetwork {
	return &NTPNetwork{address: address}
}

// NTPNetwork implements the Network interface for NTP.
type NTPNetwork struct {
	address string
}

// Now returns the current time from the NTP server.
func (n *NTPNetwork) Now() (Time, error) {
	t, err := ntp.Time(n.address)

	return t, errors.Prefix("ntp", err)
}
