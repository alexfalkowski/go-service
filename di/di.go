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

var (
	// Constructor is an alias for fx.Provide.
	Constructor = fx.Provide

	// Decorate is an alias for fx.Decorate.
	Decorate = fx.Decorate

	// New is an alias for fx.New.
	New = fx.New

	// NoLogger is an alias for fx.NopLogger.
	NoLogger = fx.NopLogger

	// Module is an alias for fx.Options.
	Module = fx.Options

	// Module is an alias for fx.Invoke.
	Register = fx.Invoke

	// RootCause is an alias for dig.RootCause.
	RootCause = dig.RootCause
)
