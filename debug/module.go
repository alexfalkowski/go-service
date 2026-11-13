package debug

import (
	"github.com/alexfalkowski/go-service/v3/debug/http"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires the debug subsystem into go-signal.
//
// It provides:
//
//   - a debug router (*[github.com/alexfalkowski/go-service/v3/debug/http.ServeMux]) via [http.NewServeMux],
//   - the debug server (*[Server]) via [NewServer] (returns nil when disabled), and
//   - [Register], the front door for optional debug endpoint registration.
//
// Register installs statsviz, pprof, and fgprof handlers under their
// service-namespaced /<name>/debug/... routes.
//
// Register attaches handlers to the debug mux only when the debug server is enabled. The mux is then
// used by the debug server when it is enabled via configuration.
var Module = signal.Module(
	signal.Constructor(http.NewServeMux),
	signal.Constructor(NewServer),
	signal.Register(Register),
)
