package sql

import (
	"github.com/alexfalkowski/go-service/pkg/sql/pg"
)

// Config for sql.
type Config struct {
	PG pg.Config `yaml:"pg"`
}
