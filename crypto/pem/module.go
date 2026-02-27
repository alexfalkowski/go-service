package pem

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the PEM decoding subsystem into Fx/Dig.
//
// It provides a constructor for `*Decoder` via `NewDecoder`, which resolves PEM-encoded sources using the
// go-service "source string" pattern and extracts raw bytes for a requested PEM block kind.
//
// This module does not require configuration.
var Module = di.Module(
	di.Constructor(NewDecoder),
)
