package pem

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the PEM decoding subsystem into [github.com/alexfalkowski/go-signal/v2].
//
// It provides a constructor for [*Decoder] via [NewDecoder], which resolves PEM-encoded sources using the
// go-service "source string" pattern and extracts raw bytes for a requested PEM block kind.
//
// This module does not require configuration.
var Module = signal.Module(
	signal.Constructor(NewDecoder),
)
