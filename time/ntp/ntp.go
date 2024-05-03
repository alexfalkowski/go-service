package ntp

import (
	"time"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/beevik/ntp"
)

// NewService for NTP.
func NewService(cfg *Config) *Service {
	return &Service{cfg: cfg}
}

// Service for NTP.
type Service struct {
	cfg *Config
}

// Time for NTP.
func (s *Service) Time() (time.Time, error) {
	if !IsEnabled(s.cfg) {
		return time.Now(), nil
	}

	t, err := ntp.Time(s.cfg.Host)

	return t, errors.Prefix("ntp time", err)
}

// Query for NTP.
func (s *Service) Query() (*ntp.Response, error) {
	if !IsEnabled(s.cfg) {
		return &ntp.Response{}, nil
	}

	q, err := ntp.Query(s.cfg.Host)

	return q, errors.Prefix("ntp query", err)
}
