package debug

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
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
//   - `Register`, the front door for optional debug endpoint registration:
//
//   - statsviz under /debug/statsviz
//
//   - pprof under /debug/pprof
//
//   - fgprof under /debug/fgprof
//
//   - psutil under /debug/psutil
//
// Register attaches handlers to the debug mux only when the debug server is enabled. The mux is then
// used by the debug server when it is enabled via configuration.
var Module = di.Module(
	di.Constructor(http.NewServeMux),
	di.Constructor(NewServer),
	di.Register(Register),
)
