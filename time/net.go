package time

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrNotFound is returned when [Config.Kind] does not match a supported network
// time provider.
//
// This error is returned by [NewNetwork] when cfg is enabled (non-nil) but Kind
// is not recognized by this package.
var ErrNotFound = errors.New("time: network not found")

// NewNetwork constructs a network time provider based on cfg.
//
// Enablement is modeled by presence: if cfg is nil (disabled), NewNetwork
// returns (nil, nil). A non-nil cfg is enabled even when cfg.Kind is empty.
//
// Supported kinds:
//   - "ntp": constructs an NTP-backed provider (see [NewNTPNetwork])
//   - "nts": constructs an NTS-backed provider (see [NewNTSNetwork])
//
// If cfg.Kind is empty or not recognized, NewNetwork returns (nil, ErrNotFound).
//
// Note: Address validation is provider-specific. NewNetwork does not validate
// cfg.Address; providers typically return an error from [Network.Now] when the
// address is empty or invalid.
func NewNetwork(cfg *Config) (Network, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	switch cfg.Kind {
	case "ntp":
		return NewNTPNetwork(cfg.Address, cfg.Timeout), nil
	case "nts":
		return NewNTSNetwork(cfg.Address, cfg.Timeout), nil
	default:
		return nil, ErrNotFound
	}
}

// Network provides the current time from a network time provider (for example NTP or NTS).
//
// Implementations should return the current time as reported by the configured provider.
// Errors returned by Now should include enough context for callers to diagnose the failure
// (for example connection failures, protocol errors, or validation failures).
type Network interface {
	// Now returns the current time from the provider.
	//
	// Implementations may perform network I/O and may return an error if the provider
	// cannot be reached, the response is invalid, or the configured address is incorrect.
	Now() (Time, error)
}
