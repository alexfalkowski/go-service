package time

import (
	"errors"

	se "github.com/alexfalkowski/go-service/errors"
	"github.com/beevik/ntp"
	"github.com/beevik/nts"
)

// ErrNotFound for metrics.
var ErrNotFound = errors.New("time: network not found")

// Network for time.
type Network interface {
	// Now from the network.
	Now() (Time, error)
}

// NewNetwork for time.
func NewNetwork(cfg *Config) (Network, error) {
	switch {
	case !IsEnabled(cfg):
		return nil, nil
	case cfg.IsNTP():
		return &ntpNetwork{c: cfg}, nil
	case cfg.IsNTS():
		return &ntsNetwork{c: cfg}, nil
	}

	return nil, ErrNotFound
}

type ntpNetwork struct {
	c *Config
}

func (n *ntpNetwork) Now() (Time, error) {
	t, err := ntp.Time(n.c.Address)

	return t, se.Prefix("ntp", err)
}

type ntsNetwork struct {
	c *Config
}

func (n *ntsNetwork) Now() (Time, error) {
	session, err := nts.NewSession(n.c.Address)
	if err != nil {
		return Time{}, se.Prefix("nts", err)
	}

	res, err := session.Query()
	if err != nil {
		return Time{}, se.Prefix("nts", err)
	}

	if err := res.Validate(); err != nil {
		return Time{}, se.Prefix("nts", err)
	}

	return Now().Add(res.ClockOffset), nil
}
