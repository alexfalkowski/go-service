package debug

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewServeMux),
	fx.Provide(NewServer),
	fx.Invoke(RegisterStatsviz),
	fx.Invoke(RegisterPprof),
	fx.Invoke(RegisterFgprof),
	fx.Invoke(RegisterPsutil),
)
