package di

import (
	"go.uber.org/dig"
	"go.uber.org/fx"
)

type (
	// Hook is an alias for fx.Hook.
	Hook = fx.Hook

	// Lifecycle is an alias for fx.Lifecycle.
	Lifecycle = fx.Lifecycle

	// In is an alias for fx.In.
	In = fx.In

	// Option is an alias for fx.Option.
	Option = fx.Option

	// Shutdowner is an alias for fx.Shutdowner.
	Shutdowner = fx.Shutdowner

	// ShutdownOption is an alias for fx.ShutdownOption.
	ShutdownOption = fx.ShutdownOption
)

// NoLogger is an alias for fx.NopLogger.
var NoLogger = fx.NopLogger

// Constructor is an alias for fx.Provide.
func Constructor(constructors ...any) Option {
	return fx.Provide(constructors...)
}

// Decorate is an alias for fx.Decorate.
func Decorate(decorators ...any) Option {
	return fx.Decorate(decorators...)
}

// New is an alias for fx.New.
func New(opts ...Option) *fx.App {
	return fx.New(opts...)
}

// Module is an alias for fx.Options.
func Module(opts ...Option) Option {
	return fx.Options(opts...)
}

// Module is an alias for fx.Invoke.
func Register(funcs ...any) Option {
	return fx.Invoke(funcs...)
}

// RootCause is an alias for dig.RootCause.
func RootCause(err error) error {
	return dig.RootCause(err)
}
