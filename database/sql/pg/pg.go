package pg

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/os"
	pgx "github.com/jackc/pgx/v5/stdlib"
	"github.com/linxGnu/mssqlx"
)

// Register for pg.
func Register() {
	_ = driver.Register("pg", pgx.GetDefaultDriver())
}

// Open for pg.
func Open(lc di.Lifecycle, fs *os.FS, cfg *Config) (*mssqlx.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Open(lc, "pg", fs, cfg.Config)
}
