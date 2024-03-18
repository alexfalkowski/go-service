package debug

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewServer),
	fx.Invoke(RegisterStatsviz),
	fx.Invoke(RegisterPprof),
	fx.Invoke(RegisterPsutil),
)
