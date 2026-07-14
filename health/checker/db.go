package checker

import (
	"fmt"

	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
)

var _ checker.Checker = (*DBChecker)(nil)

// ErrNoConnections is returned when a DBChecker has no writer or reader pools to verify.
var ErrNoConnections = errors.New("db: no connections")

// ErrPingTimeout is the cause recorded when a DB health-check ping times out.
var ErrPingTimeout = fmt.Errorf("db: ping timeout: %w", sync.ErrTimeout)

// NewDBChecker constructs a DBChecker that verifies database connectivity by pinging
// all configured writer and reader pools in db.
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
// It pings each configured writer and reader database pool using PingContext. Any
// ping failures are aggregated and returned from Check.
type DBChecker struct {
	db      *sql.DBs
	timeout time.Duration
}

// Check verifies database health by pinging all configured writer and reader databases.
//
// It wraps ping failures with the pool role and index (for example `db writer[0]`) and returns them as
// a single aggregated error via [errors.Join]. If a ping context finishes before PingContext returns
// successfully, the wrapped cause is [ErrPingTimeout] for this checker's timeout, or the parent context's
// cause when the parent context is canceled first. If no database pools are configured, Check returns
// [ErrNoConnections]. If all pings succeed, Check returns nil.
func (c *DBChecker) Check(ctx context.Context) error {
	databases := c.databases()
	if len(databases) == 0 {
		return ErrNoConnections
	}

	group := sync.ErrorsGroup{}
	group.SetLimit(len(databases))

	for _, database := range databases {
		group.Go(func() error {
			err := c.ping(ctx, database.db)
			if err != nil {
				return fmt.Errorf("db %s[%d]: %w", database.role, database.index, err)
			}

			return nil
		})
	}

	return group.Wait()
}

func (c *DBChecker) ping(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeoutCause(ctx, c.timeout, ErrPingTimeout)
	defer cancel()

	err := db.PingContext(ctx)
	if err != nil && ctx.Err() != nil {
		return context.Cause(ctx)
	}

	return err
}

type database struct {
	db    *sql.DB
	role  string
	index int
}

func (c *DBChecker) databases() []database {
	writers := c.db.Writers()
	readers := c.db.Readers()

	databases := make([]database, 0, len(writers)+len(readers))
	for index, db := range writers {
		databases = append(databases, database{role: "writer", index: index, db: db})
	}
	for index, db := range readers {
		databases = append(databases, database{role: "reader", index: index, db: db})
	}
	return databases
}
