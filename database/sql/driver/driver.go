package driver

import (
	"database/sql/driver"

	"github.com/alexfalkowski/go-service/v2/errors"
)

// Driver aliases [database/sql/driver.Driver].
//
// It is the concrete driver type expected by Register.
type Driver = driver.Driver

// Conn aliases [database/sql/driver.Conn].
type Conn = driver.Conn

// NamedValue aliases [database/sql/driver.NamedValue].
type NamedValue = driver.NamedValue

// Rows aliases [database/sql/driver.Rows].
type Rows = driver.Rows

// Stmt aliases [database/sql/driver.Stmt].
type Stmt = driver.Stmt

// Tx aliases [database/sql/driver.Tx].
type Tx = driver.Tx

// Value aliases [database/sql/driver.Value].
type Value = driver.Value

// ErrSkip aliases [database/sql/driver.ErrSkip].
var ErrSkip = driver.ErrSkip

// ErrNoDSNs is returned when SQL configuration enables a driver without any writer or reader DSNs.
var ErrNoDSNs = errors.New("driver: no database DSNs configured")

// ErrEmptyDSN is returned when a configured DSN source resolves to an empty string.
var ErrEmptyDSN = errors.New("driver: empty database DSN")
