package debug

import (
	"github.com/alexfalkowski/go-service/v2/debug/fgprof"
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/debug/pprof"
	"github.com/alexfalkowski/go-service/v2/debug/psutil"
	"github.com/alexfalkowski/go-service/v2/debug/statsviz"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires the debug subsystem into Fx/Dig.
//
// It provides:
//
//   - a debug router (`*debug/http.ServeMux`) via `http.NewServeMux`,
//
//   - the debug server (`*debug.Server`) via `NewServer` (returns nil when disabled), and
//
//   - registrations for optional debug endpoints:
//
//   - statsviz under /debug/statsviz
//
//   - pprof under /debug/pprof
//
//   - fgprof under /debug/fgprof
//
//   - psutil under /debug/psutil
//
// The endpoint registrations attach handlers to the debug mux. The mux is then used by the debug server
// when it is enabled via configuration.
var Module = di.Module(
	di.Constructor(http.NewServeMux),
	di.Constructor(NewServer),
	di.Register(statsviz.Register),
	di.Register(pprof.Register),
	di.Register(fgprof.Register),
	di.Register(psutil.Register),
)
