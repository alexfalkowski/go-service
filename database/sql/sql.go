package sql

import (
	"database/sql"

	"github.com/alexfalkowski/go-service/v3/database/sql/driver"
)

// ErrConnDone aliases [database/sql.ErrConnDone].
var ErrConnDone = sql.ErrConnDone

// ErrNoRows aliases [database/sql.ErrNoRows].
var ErrNoRows = sql.ErrNoRows

// ErrTxDone aliases [database/sql.ErrTxDone].
var ErrTxDone = sql.ErrTxDone

// LevelDefault aliases [database/sql.LevelDefault].
const LevelDefault = sql.LevelDefault

// LevelReadUncommitted aliases [database/sql.LevelReadUncommitted].
const LevelReadUncommitted = sql.LevelReadUncommitted

// LevelReadCommitted aliases [database/sql.LevelReadCommitted].
const LevelReadCommitted = sql.LevelReadCommitted

// LevelWriteCommitted aliases [database/sql.LevelWriteCommitted].
const LevelWriteCommitted = sql.LevelWriteCommitted

// LevelRepeatableRead aliases [database/sql.LevelRepeatableRead].
const LevelRepeatableRead = sql.LevelRepeatableRead

// LevelSnapshot aliases [database/sql.LevelSnapshot].
const LevelSnapshot = sql.LevelSnapshot

// LevelSerializable aliases [database/sql.LevelSerializable].
const LevelSerializable = sql.LevelSerializable

// LevelLinearizable aliases [database/sql.LevelLinearizable].
const LevelLinearizable = sql.LevelLinearizable

// ColumnType aliases [database/sql.ColumnType].
type ColumnType = sql.ColumnType

// Conn aliases [database/sql.Conn].
type Conn = sql.Conn

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

// DBStats aliases [database/sql.DBStats].
type DBStats = sql.DBStats

// IsolationLevel aliases [database/sql.IsolationLevel].
type IsolationLevel = sql.IsolationLevel

// NamedArg aliases [database/sql.NamedArg].
type NamedArg = sql.NamedArg

// Null aliases [database/sql.Null].
type Null[T any] = sql.Null[T]

// NullBool aliases [database/sql.NullBool].
type NullBool = sql.NullBool

// NullByte aliases [database/sql.NullByte].
type NullByte = sql.NullByte

// NullFloat64 aliases [database/sql.NullFloat64].
type NullFloat64 = sql.NullFloat64

// NullInt16 aliases [database/sql.NullInt16].
type NullInt16 = sql.NullInt16

// NullInt32 aliases [database/sql.NullInt32].
type NullInt32 = sql.NullInt32

// NullInt64 aliases [database/sql.NullInt64].
type NullInt64 = sql.NullInt64

// NullString aliases [database/sql.NullString].
type NullString = sql.NullString

// NullTime aliases [database/sql.NullTime].
type NullTime = sql.NullTime

// Out aliases [database/sql.Out].
type Out = sql.Out

// RawBytes aliases [database/sql.RawBytes].
type RawBytes = sql.RawBytes

// Result aliases [database/sql.Result].
type Result = sql.Result

// Row aliases [database/sql.Row].
type Row = sql.Row

// Rows aliases [database/sql.Rows].
//
// It is exposed so go-service SQL packages can use the local package namespace
// while preserving standard-library behavior.
type Rows = sql.Rows

// Scanner aliases [database/sql.Scanner].
type Scanner = sql.Scanner

// Stmt aliases [database/sql.Stmt].
type Stmt = sql.Stmt

// Tx aliases [database/sql.Tx].
type Tx = sql.Tx

// TxOptions aliases [database/sql.TxOptions].
type TxOptions = sql.TxOptions

// ConvertAssign wraps [database/sql.ConvertAssign].
func ConvertAssign(scanCtx driver.ScanContext, dest any, src driver.Value) error {
	return sql.ConvertAssign(scanCtx, dest, src)
}

// Drivers wraps [database/sql.Drivers].
func Drivers() []string {
	return sql.Drivers()
}

// Named wraps [database/sql.Named].
func Named(name string, value any) NamedArg {
	return sql.Named(name, value)
}

// Open aliases [database/sql.Open].
//
// The driver name must already be registered with the standard library [database/sql] driver registry.
// The returned DB uses standard-library database/sql behavior.
func Open(driverName, dataSourceName string) (*DB, error) {
	return sql.Open(driverName, dataSourceName)
}

// OpenDB wraps [database/sql.OpenDB].
func OpenDB(c driver.Connector) *DB {
	return sql.OpenDB(c)
}

// Register wraps [driver.Register] using its default telemetry options.
//
// Registration panics are returned as errors, and drivers are wrapped with telemetry
// when tracing or metrics are enabled. For custom telemetry options, use [driver.Register].
func Register(name string, d driver.Driver) error {
	return driver.Register(name, d)
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
