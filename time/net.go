package time

import "errors"

// ErrNotFound for metrics.
var ErrNotFound = errors.New("time: network not found")

// NewNetwork for time.
func NewNetwork(cfg *Config) (Network, error) {
	switch {
	case !IsEnabled(cfg):
		return nil, nil
	case cfg.IsNTP():
		return NewNTPNetwork(cfg.Address), nil
	case cfg.IsNTS():
		return NewNTSNetwork(cfg.Address), nil
	}

	return nil, ErrNotFound
}

// Network for time.
type Network interface {
	// Now from the network.
	Now() (Time, error)
}
