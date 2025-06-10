package debug

import (
	"github.com/alexfalkowski/go-service/v2/debug/fgprof"
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/debug/pprof"
	"github.com/alexfalkowski/go-service/v2/debug/psutil"
	"github.com/alexfalkowski/go-service/v2/debug/statsviz"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module for fx.
var Module = di.Module(
	di.Constructor(http.NewServeMux),
	di.Constructor(NewServer),
	di.Register(statsviz.Register),
	di.Register(pprof.Register),
	di.Register(fgprof.Register),
	di.Register(psutil.Register),
)
