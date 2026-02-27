package header

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
)

// Map is a set of telemetry exporter/request headers.
//
// The map keys are header names and the map values are header values.
//
// Values are often configured as go-service “source strings” so that secrets can be
// supplied at runtime rather than embedded directly in config files. See Secrets
// for the supported forms and resolution behavior.
type Map map[string]string

// Secrets resolves configured header values using the go-service “source string” convention.
//
// It traverses m and replaces each value in place by reading it through fs.ReadSource.
// fs.ReadSource supports these forms:
//
//   - "env:NAME"    reads the value of environment variable NAME.
//   - "file:/path"  reads bytes from the file at /path (including path cleaning and trimming).
//   - otherwise     treats the value as a literal string.
//
// After successful resolution, every entry in m contains the final literal header
// value that should be sent by exporters/clients.
//
// Secrets returns the first error encountered while resolving any value.
//
// Note: Secrets mutates the map in place. If you need to preserve the original
// configured source strings, copy the map before calling Secrets.
func (m Map) Secrets(fs *os.FS) error {
	for key, name := range m {
		data, err := fs.ReadSource(name)
		if err != nil {
			return err
		}

		m[key] = bytes.String(data)
	}

	return nil
}

// MustSecrets resolves configured header values like Secrets, but panics on error.
//
// It calls Secrets and panics if any value cannot be resolved (via runtime.Must).
// This is intended for strict startup/config projection paths where missing or
// unreadable secret material should fail fast.
func (m Map) MustSecrets(fs *os.FS) {
	runtime.Must(m.Secrets(fs))
}
