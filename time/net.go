package time

import (
	"time"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/beevik/ntp"
	"github.com/beevik/nts"
)

// Network for time.
type Network interface {
	// Now from the network.
	Now() (time.Time, error)
}

// NewNetwork for time.
func NewNetwork(cfg *Config) Network {
	switch {
	case !IsEnabled(cfg):
		return &sysNetwork{}
	case cfg.IsNTP():
		return &ntpNetwork{c: cfg}
	case cfg.IsNTS():
		return &ntsNetwork{c: cfg}
	default:
		return &sysNetwork{}
	}
}

type sysNetwork struct{}

func (*sysNetwork) Now() (time.Time, error) {
	return time.Now(), nil
}

type ntpNetwork struct {
	c *Config
}

func (n *ntpNetwork) Now() (time.Time, error) {
	t, err := ntp.Time(n.c.Address)

	return t, errors.Prefix("ntp", err)
}

type ntsNetwork struct {
	c *Config
}

func (n *ntsNetwork) Now() (time.Time, error) {
	se, err := nts.NewSession(n.c.Address)
	if err != nil {
		return time.Now(), errors.Prefix("nts", err)
	}

	r, err := se.Query()
	if err != nil {
		return time.Now(), errors.Prefix("nts", err)
	}

	err = r.Validate()
	if err != nil {
		return time.Now(), errors.Prefix("nts", err)
	}

	return time.Now().Add(r.ClockOffset), nil
}
