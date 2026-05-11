package sql

import (
	"database/sql"

	"github.com/linxGnu/mssqlx"
)

// DB aliases `database/sql`.DB.
//
// It is exposed so go-service SQL packages can use the local package namespace
// while preserving standard-library behavior.
type DB = sql.DB

// DBs is an alias of mssqlx.DBs.
//
// It represents the master/slave SQL pool collection used by go-service SQL
// integrations and preserves the upstream behavior exactly.
//
// The value groups master and replica `sqlx.DB` pools behind a single type with
// helper methods for querying masters, slaves, and running operations against
// the configured pools.
type DBs = mssqlx.DBs

// Rows aliases `database/sql`.Rows.
//
// It is exposed so go-service SQL packages can use the local package namespace
// while preserving standard-library behavior.
type Rows = sql.Rows

// Open aliases `database/sql`.Open.
//
// The driver name must already be registered with the global `database/sql`
// registry. The returned DB uses standard-library database/sql behavior.
func Open(driverName, dataSourceName string) (*DB, error) {
	return sql.Open(driverName, dataSourceName)
}

// ConnectMasterSlaves is a thin wrapper around mssqlx.ConnectMasterSlaves.
//
// It opens a master/slave SQL pool collection for a registered `database/sql`
// driver name and returns any per-DSN connection errors produced by the
// upstream helper.
//
// The driver name must already be registered with the global `database/sql`
// registry. The returned `DBs` value may contain zero or more master and slave
// pools depending on the provided DSN lists.
func ConnectMasterSlaves(name string, masterDSNs, slaveDSNs []string) (*DBs, []error) {
	return mssqlx.ConnectMasterSlaves(name, masterDSNs, slaveDSNs)
}
