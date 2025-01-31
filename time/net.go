package time

import (
	"time"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/runtime"
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

func (n *ntsNetwork) Now() (t time.Time, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("nts", runtime.ConvertRecover(r))
		}
	}()

	se, err := nts.NewSession(n.c.Address)
	runtime.Must(err)

	res, err := se.Query()
	runtime.Must(err)

	err = res.Validate()
	runtime.Must(err)

	t = time.Now().Add(res.ClockOffset)

	return
}
