package pg

import "github.com/alexfalkowski/go-service/v2/database/sql/config"

// Config contains PostgreSQL SQL database configuration.
//
// It embeds `database/sql/config`.Config to reuse common `database/sql` pool settings and DSN configuration.
//
// # Optional pointers and "enabled" semantics
//
// This type is intentionally optional. By convention across go-service configuration types, a nil *Config is treated
// as "PostgreSQL disabled". The embedded `*config.Config` is also optional; IsEnabled returns true only when both the
// outer *Config and the embedded *config.Config are non-nil.
//
// This allows services to omit either `pg:` or the embedded fields entirely to disable PostgreSQL wiring.
type Config struct {
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled reports whether PostgreSQL configuration is present and enabled.
//
// It returns true only when both the PostgreSQL wrapper config and the embedded shared SQL config are non-nil.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}
