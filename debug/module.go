package debug

import (
	"github.com/alexfalkowski/go-service/v2/debug/fgprof"
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/debug/pprof"
	"github.com/alexfalkowski/go-service/v2/debug/psutil"
	"github.com/alexfalkowski/go-service/v2/debug/statsviz"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(http.NewServeMux),
	fx.Provide(NewServer),
	fx.Invoke(statsviz.Register),
	fx.Invoke(pprof.Register),
	fx.Invoke(fgprof.Register),
	fx.Invoke(psutil.Register),
)
