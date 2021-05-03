package test

import (
	"go.uber.org/fx"
)

func NewShutdowner() *Shutdowner {
	return &Shutdowner{}
}

type Shutdowner struct{}

func (*Shutdowner) Shutdown(...fx.ShutdownOption) error {
	return nil
}
