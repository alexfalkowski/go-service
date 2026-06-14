// Package options provides helpers for working with low-level configuration
// option key-value pairs.
//
// Options are administrator-supplied startup tuning knobs used by transport and
// feature integrations for backend-specific settings that do not belong in the
// strongly typed configuration schema. Helper methods return caller-provided
// fallbacks when a key is absent and panic when a present value cannot be parsed
// for the requested type.
//
// Byte-size helpers parse the same human-readable decimal size strings as
// [github.com/alexfalkowski/go-service/v2/bytes.ParseSize]. They do not apply
// the typed configuration cap [github.com/alexfalkowski/go-service/v2/bytes.MaxConfigSize];
// callers that need a repository-owned bound should use typed config fields
// validated with the config package's `config_size` rule.
package options
