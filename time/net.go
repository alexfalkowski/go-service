package time

import (
	"errors"

	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/runtime"
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

func (n *ntsNetwork) Now() (t Time, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = se.Prefix("nts", runtime.ConvertRecover(r))
		}
	}()

	se, err := nts.NewSession(n.c.Address)
	runtime.Must(err)

	res, err := se.Query()
	runtime.Must(err)

	err = res.Validate()
	runtime.Must(err)

	t = Now().Add(res.ClockOffset)

	return
}
