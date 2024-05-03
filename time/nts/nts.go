package nts

import (
	"github.com/alexfalkowski/go-service/errors"
	"github.com/beevik/ntp"
	"github.com/beevik/nts"
)

// NewService for NTS.
func NewService(cfg *Config) (*Service, error) {
	s := &Service{}

	if !IsEnabled(cfg) {
		return s, nil
	}

	se, err := nts.NewSession(cfg.Host)
	if err != nil {
		return nil, err
	}

	s.session = se

	return s, nil
}

// Service for NTS.
type Service struct {
	session *nts.Session
}

// Query for NTS.
func (s *Service) Query() (*ntp.Response, error) {
	if s.session == nil {
		return &ntp.Response{}, nil
	}

	q, err := s.session.Query()

	return q, errors.Prefix("nts query", err)
}
