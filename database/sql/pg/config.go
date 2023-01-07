package pg

import "github.com/alexfalkowski/go-service/database/sql/config"

// Config for SQL.
type Config struct {
	config.Config `yaml:",inline" json:",inline" toml:",inline"`
}
