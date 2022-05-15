package test

import (
	"go.uber.org/fx"
)

func NewShutdowner() *Shutdowner {
	return &Shutdowner{}
}

type Shutdowner struct {
	Called bool
}

func (s *Shutdowner) Shutdown(...fx.ShutdownOption) error {
	s.Called = true

	return nil
}
