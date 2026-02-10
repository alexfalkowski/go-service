package time

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrNotFound is returned when a configured network time provider kind is not supported.
var ErrNotFound = errors.New("time: network not found")

// NewNetwork constructs a Network time provider based on cfg.
//
// If cfg is disabled, it returns (nil, nil).
// Supported kinds include "ntp" and "nts". For any other kind it returns ErrNotFound.
func NewNetwork(cfg *Config) (Network, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	switch cfg.Kind {
	case "ntp":
		return NewNTPNetwork(cfg.Address), nil
	case "nts":
		return NewNTSNetwork(cfg.Address), nil
	default:
		return nil, ErrNotFound
	}
}

// Network provides the current time from a network time provider (for example NTP or NTS).
type Network interface {
	// Now returns the current time from the network provider.
	Now() (Time, error)
}
