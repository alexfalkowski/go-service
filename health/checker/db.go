package checker

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/linxGnu/mssqlx"
)

// NewDBChecker for health.
func NewDBChecker(db *mssqlx.DBs, timeout time.Duration) *DBChecker {
	return &DBChecker{db: db, timeout: timeout}
}

// DBChecker for health.
type DBChecker struct {
	db      *mssqlx.DBs
	timeout time.Duration
}

// Check db health.
func (c *DBChecker) Check(ctx context.Context) error {
	_, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return errors.Join(c.db.Ping()...)
}
