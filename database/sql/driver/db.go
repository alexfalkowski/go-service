package driver

import (
	"database/sql"
	"math/rand/v2"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/time"
)

// DBs contains writer and reader SQL connection pools.
//
// Reader returns a random reader pool when readers are configured and falls back
// to a random writer pool otherwise. Writer returns a random writer pool.
type DBs struct {
	driverName string
	writers    []*sql.DB
	readers    []*sql.DB

	registrations []metrics.Registration
}

// DriverName returns the registered database/sql driver name.
func (d *DBs) DriverName() string {
	return d.driverName
}

// Writers returns all writer pools.
func (d *DBs) Writers() []*sql.DB {
	return d.writers
}

// Readers returns all reader pools.
func (d *DBs) Readers() []*sql.DB {
	return d.readers
}

// Ping pings all writer and reader pools.
func (d *DBs) Ping() error {
	return ping(d.databases())
}

// PingWriter pings all writer pools.
func (d *DBs) PingWriter() error {
	return ping(d.writers)
}

// PingReader pings all reader pools.
func (d *DBs) PingReader() error {
	return ping(d.readers)
}

// Reader returns a database pool suitable for read queries.
func (d *DBs) Reader() (*sql.DB, error) {
	if len(d.readers) > 0 {
		return pick(d.readers), nil
	}

	return d.Writer()
}

// Writer returns a database pool suitable for writes and transactions.
func (d *DBs) Writer() (*sql.DB, error) {
	if len(d.writers) == 0 {
		return nil, ErrNoDSNs
	}

	return pick(d.writers), nil
}

// Destroy unregisters repository-owned DB stats metrics and closes all database
// pools.
func (d *DBs) Destroy() error {
	if d == nil {
		return nil
	}

	regs := d.registrations
	d.registrations = nil

	errs := unregister(regs)
	errs = append(errs, closeAll(d.databases())...)

	return errors.Join(errs...)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
// across all pools.
func (d *DBs) SetConnMaxLifetime(v time.Duration) {
	for _, db := range d.databases() {
		db.SetConnMaxLifetime(v.Duration())
	}
}

// SetMaxIdleConns sets the maximum number of idle connections across all pools.
func (d *DBs) SetMaxIdleConns(v int) {
	for _, db := range d.databases() {
		db.SetMaxIdleConns(v)
	}
}

// SetMaxOpenConns sets the maximum number of open connections across all pools.
func (d *DBs) SetMaxOpenConns(v int) {
	for _, db := range d.databases() {
		db.SetMaxOpenConns(v)
	}
}

func (d *DBs) databases() []*sql.DB {
	dbs := make([]*sql.DB, 0, len(d.writers)+len(d.readers))
	dbs = append(dbs, d.writers...)
	dbs = append(dbs, d.readers...)

	return dbs
}

func pick(databases []*sql.DB) *sql.DB {
	//nolint:gosec // Pool selection is load distribution, not secret generation or security policy.
	return databases[rand.IntN(len(databases))]
}

func ping(databases []*sql.DB) error {
	errs := make([]error, len(databases))
	for i, db := range databases {
		if db != nil {
			errs[i] = db.PingContext(context.Background())
		}
	}

	return errors.Join(errs...)
}

func closeAll(databases []*sql.DB) []error {
	errs := make([]error, len(databases))
	for i, db := range databases {
		if db != nil {
			errs[i] = db.Close()
		}
	}

	return errs
}
