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
// It represents the master/slave SQL pool collection used by go-service SQL
// integrations.
//
// The value wraps the upstream master/slave pool collection and embeds it so
// callers can use the usual pool helper methods while go-service-owned cleanup,
// such as metric unregistration, stays attached to [DBs.Destroy].
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

// ConnectMasterSlaves is a thin wrapper around [driver.ConnectMasterSlaves].
//
// It opens a master/slave SQL pool collection for a registered standard library [database/sql] driver
// name and returns any per-DSN connection errors produced by the upstream helper.
//
// The driver name must already be registered with the standard library [database/sql] driver registry.
// The returned [DBs] value may contain zero or more master and slave
// pools depending on the provided DSN lists.
func ConnectMasterSlaves(name string, masterDSNs, slaveDSNs []string) (*DBs, []error) {
	return driver.ConnectMasterSlaves(name, masterDSNs, slaveDSNs)
}
