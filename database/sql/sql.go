package sql

import "github.com/linxGnu/mssqlx"

// DBs is an alias of mssqlx.DBs.
//
// It represents the master/slave SQL pool collection used by go-service SQL
// integrations and preserves the upstream behavior exactly.
//
// The value groups master and replica `sqlx.DB` pools behind a single type with
// helper methods for querying masters, slaves, and running operations against
// the configured pools.
type DBs = mssqlx.DBs

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
