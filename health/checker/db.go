package checker

import (
	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
)

var _ checker.Checker = (*DBChecker)(nil)

// NewDBChecker constructs a DBChecker that verifies database connectivity by pinging
// all configured master and slave pools in db.
//
// The timeout is applied per PingContext invocation (i.e. each individual pool ping
// gets its own derived context with this timeout).
//
// Note: db is expected to be non-nil. If the database subsystem is disabled and no
// pools are available, callers should avoid constructing/registering this checker.
func NewDBChecker(db *mssqlx.DBs, timeout time.Duration) *DBChecker {
	return &DBChecker{db: db, timeout: timeout}
}

// DBChecker is a health checker that verifies database connectivity.
//
// It pings each configured master and slave database pool using PingContext. Any
// ping failures are aggregated and returned from Check.
type DBChecker struct {
	db      *mssqlx.DBs
	timeout time.Duration
}

// Check verifies database health by pinging all configured master and slave databases.
//
// It returns a single aggregated error (via errors.Join) containing all ping errors.
// If all pings succeed, Check returns nil.
func (c *DBChecker) Check(ctx context.Context) error {
	dbs := c.dbs()
	errs := make([]error, 0, len(dbs))
	for _, db := range dbs {
		errs = append(errs, c.ping(ctx, db))
	}
	return errors.Join(errs...)
}

func (c *DBChecker) dbs() []*sqlx.DB {
	masters, _ := c.db.GetAllMasters()
	slaves, _ := c.db.GetAllSlaves()

	dbs := make([]*sqlx.DB, 0, len(masters)+len(slaves))
	dbs = append(dbs, masters...)
	dbs = append(dbs, slaves...)
	return dbs
}

func (o *DBChecker) ping(ctx context.Context, db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	return db.PingContext(ctx)
}
