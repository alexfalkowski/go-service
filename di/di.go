package di

import (
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// Hook is an alias for fx.Hook.
//
// It is used to register lifecycle start/stop callbacks on an fx.Lifecycle.
type Hook = fx.Hook

// Lifecycle is an alias for fx.Lifecycle.
//
// It is used to manage application start/stop hooks (e.g. via Lifecycle.Append).
type Lifecycle = fx.Lifecycle

// In is an alias for fx.In.
//
// Embed In in parameter structs to declare fields that should be injected by Fx/Dig.
type In = fx.In

// Option is an alias for fx.Option.
//
// Options are used to compose modules, provide constructors, decorate values, and register invocations.
type Option = fx.Option

// Shutdowner is an alias for fx.Shutdowner.
//
// It allows triggering application shutdown from within the Fx graph.
type Shutdowner = fx.Shutdowner

// ShutdownOption is an alias for fx.ShutdownOption.
//
// It configures how a shutdown signal is delivered/handled by Fx.
type ShutdownOption = fx.ShutdownOption

// NoLogger is an alias for fx.NopLogger.
//
// It can be included as an option to disable Fx's internal logging.
var NoLogger = fx.NopLogger

// Constructor registers one or more constructors with Fx.
//
// This is a thin wrapper around fx.Provide. Each provided constructor is called by Fx/Dig when its
// outputs are required to satisfy dependencies.
func Constructor(constructors ...any) Option {
	return fx.Provide(constructors...)
}

// Decorate registers one or more decorators with Fx.
//
// This is a thin wrapper around fx.Decorate. Decorators allow you to wrap or modify values after they
// are constructed by constructors.
func Decorate(decorators ...any) Option {
	return fx.Decorate(decorators...)
}

// New constructs a new Fx application from the provided options.
//
// This is a thin wrapper around fx.New.
func New(opts ...Option) *fx.App {
	return fx.New(opts...)
}

// Module composes multiple options into a single option.
//
// This is a thin wrapper around fx.Options and is commonly used to define feature modules as package-level
// variables.
func Module(opts ...Option) Option {
	return fx.Options(opts...)
}

// Register registers one or more invocation functions to run during application startup.
//
// This is a thin wrapper around fx.Invoke. Invocations are typically used for side-effectful wiring such as:
//   - HTTP route registration
//   - telemetry registration
//   - driver initialization/registration
func Register(funcs ...any) Option {
	return fx.Invoke(funcs...)
}

// RootCause returns the underlying root cause of an error returned by Fx/Dig.
//
// This is a thin wrapper around dig.RootCause and is useful for formatting/attributing startup failures.
func RootCause(err error) error {
	return dig.RootCause(err)
}
