package sql

import (
	"database/sql"

	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
)

// DB aliases [database/sql.DB].
//
// It is exposed so go-service SQL packages can use the local package namespace
// while preserving standard-library behavior.
type DB = sql.DB

// DBs is an alias of [driver.DBs].
//
// It represents the writer/reader SQL pool collection used by go-service SQL
// integrations. Callers choose a pool with [DBs.Reader] or [DBs.Writer] and use
// standard-library [database/sql] methods on the returned pool.
type DBs = driver.DBs

// Rows aliases [database/sql.Rows].
//
// It is exposed so go-service SQL packages can use the local package namespace
// while preserving standard-library behavior.
type Rows = sql.Rows

// Open aliases [database/sql.Open].
//
// The driver name must already be registered with the standard library [database/sql] driver registry.
// The returned DB uses standard-library database/sql behavior.
func Open(driverName, dataSourceName string) (*DB, error) {
	return sql.Open(driverName, dataSourceName)
}

// ConnectWritersReaders is a thin wrapper around [driver.ConnectWritersReaders].
//
// It opens a writer/reader SQL pool collection for a registered standard library [database/sql] driver
// name and returns any per-DSN open errors.
//
// The driver name must already be registered with the standard library [database/sql] driver registry.
// At least one writer or reader DSN must be configured.
func ConnectWritersReaders(name string, writerDSNs, readerDSNs []string) (*DBs, []error) {
	return driver.ConnectWritersReaders(name, writerDSNs, readerDSNs)
}
