package checker

import (
	"fmt"

	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
	"github.com/jmoiron/sqlx"
)

var _ checker.Checker = (*DBChecker)(nil)

// ErrNoConnections is returned when a DBChecker has no master or slave pools to verify.
var ErrNoConnections = errors.New("db: no connections")

// ErrPingTimeout is the cause recorded when a DB health-check ping times out.
var ErrPingTimeout = fmt.Errorf("db: ping timeout: %w", sync.ErrTimeout)

// NewDBChecker constructs a DBChecker that verifies database connectivity by pinging
// all configured master and slave pools in db.
//
// The timeout is applied per PingContext invocation (i.e. each individual pool ping
// gets its own derived context with this timeout).
//
// Note: db is expected to be non-nil. If the database subsystem is disabled and no
// pools are available, callers should avoid constructing/registering this checker.
func NewDBChecker(db *sql.DBs, timeout time.Duration) *DBChecker {
	return &DBChecker{db: db, timeout: timeout}
}

// DBChecker is a health checker that verifies database connectivity.
//
// It pings each configured master and slave database pool using PingContext. Any
// ping failures are aggregated and returned from Check.
type DBChecker struct {
	db      *sql.DBs
	timeout time.Duration
}

// Check verifies database health by pinging all configured master and slave databases.
//
// It returns a single aggregated error (via [errors.Join]) containing all ping errors.
// If no database pools are configured, Check returns ErrNoConnections.
// If all pings succeed, Check returns nil.
func (c *DBChecker) Check(ctx context.Context) error {
	databases := c.databases()
	if len(databases) == 0 {
		return ErrNoConnections
	}

	errs := make([]error, 0, len(databases))
	for _, database := range databases {
		err := c.ping(ctx, database.db)
		if err != nil {
			errs = append(errs, fmt.Errorf("db %s[%d]: %w", database.role, database.index, err))
		}
	}
	return errors.Join(errs...)
}

type database struct {
	db    *sqlx.DB
	role  string
	index int
}

func (c *DBChecker) databases() []database {
	masters, _ := c.db.GetAllMasters()
	slaves, _ := c.db.GetAllSlaves()

	databases := make([]database, 0, len(masters)+len(slaves))
	for index, db := range masters {
		databases = append(databases, database{role: "master", index: index, db: db})
	}
	for index, db := range slaves {
		databases = append(databases, database{role: "slave", index: index, db: db})
	}
	return databases
}

func (c *DBChecker) ping(ctx context.Context, db *sqlx.DB) error {
	ctx, cancel := context.WithTimeoutCause(ctx, c.timeout, ErrPingTimeout)
	defer cancel()

	err := db.PingContext(ctx)
	if err != nil && ctx.Err() != nil {
		return context.Cause(ctx)
	}

	return err
}
