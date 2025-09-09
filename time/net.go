package time

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrNotFound for metrics.
var ErrNotFound = errors.New("time: network not found")

// NewNetwork for time.
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

// Network for time.
type Network interface {
	// Now from the network.
	Now() (Time, error)
}
