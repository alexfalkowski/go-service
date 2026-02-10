package checker

import (
	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
)

var _ checker.Checker = (*DBChecker)(nil)

// NewDBChecker constructs a DBChecker that pings all configured master and slave databases.
//
// timeout is applied per PingContext invocation.
func NewDBChecker(db *mssqlx.DBs, timeout time.Duration) *DBChecker {
	return &DBChecker{db: db, timeout: timeout}
}

// DBChecker is a health checker that verifies database connectivity.
type DBChecker struct {
	db      *mssqlx.DBs
	timeout time.Duration
}

// Check verifies database health by pinging all configured master and slave databases.
func (c *DBChecker) Check(ctx context.Context) error {
	dbs, _ := c.db.GetAllMasters()
	for _, db := range dbs {
		if err := c.ping(ctx, db); err != nil {
			return err
		}
	}

	dbs, _ = c.db.GetAllSlaves()
	for _, db := range dbs {
		if err := c.ping(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

func (o *DBChecker) ping(ctx context.Context, db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	return db.PingContext(ctx)
}
