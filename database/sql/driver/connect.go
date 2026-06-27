package driver

import (
	"database/sql"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
)

// Connect opens writer/reader [github.com/alexfalkowski/go-service/v2/database/sql] connection pools for a previously
// registered driver name.
//
// It resolves DSNs from cfg using the provided filesystem (DSNs are configured
// as go-service "source strings"), opens writer and reader pools, registers
// OpenTelemetry DB stats metrics for each pool when metrics are enabled, and
// then applies pool settings (connection lifetime, max idle, and max open
// connections).
//
// Like [database/sql.Open], this creates pool handles but does not ping the
// database or verify network reachability. Call [DBs.Ping], [DBs.PingWriter],
// [DBs.PingReader], or a health checker when startup/readiness must prove
// connectivity.
//
// Preconditions:
//   - cfg must be non-nil and already treated as enabled/validated by the caller.
//
// Failure behavior:
//   - returns errors encountered while resolving DSNs or creating pool handles,
//   - returns [ErrEmptyDSN] when any configured DSN source resolves to empty bytes, and
//   - returns [ErrNoDSNs] when neither writers nor readers are configured.
//
// The returned DBs owns repository lifecycle cleanup such as DB stats metric
// unregistration.
func Connect(name string, fs *os.FS, cfg *config.Config, opts ...telemetry.Option) (*DBs, error) {
	return connect(name, fs, cfg, opts...)
}

// Open opens writer/reader [github.com/alexfalkowski/go-service/v2/database/sql] connection pools for a previously
// registered driver name.
//
// Open delegates the connection work to [Connect] and then appends an OnStop hook
// to the provided lifecycle that closes all returned pools by calling [DBs.Destroy].
//
// Preconditions:
//   - cfg must be non-nil and already treated as enabled/validated by the caller.
//   - driver-specific wrappers, such as [github.com/alexfalkowski/go-service/v2/database/sql/pg.Open], own nil/disabled
//     config semantics before delegating here.
//
// The returned type is the same go-service [DBs] wrapper returned by [Connect].
func Open(lc di.Lifecycle, name string, fs *os.FS, cfg *config.Config, opts ...telemetry.Option) (*DBs, error) {
	db, err := Connect(name, fs, cfg, opts...)
	if err != nil {
		return nil, err
	}

	lc.Append(di.Hook{
		OnStop: func(_ context.Context) error {
			return db.Destroy()
		},
	})

	return db, nil
}

// ConnectWritersReaders opens writer/reader [github.com/alexfalkowski/go-service/v2/database/sql] connection pools for a
// previously registered driver name.
//
// It is a low-level helper for callers that already have resolved literal DSNs.
// It creates pool handles with [database/sql.Open] and does not ping the
// databases.
func ConnectWritersReaders(name string, writerDSNs, readerDSNs []string) (*DBs, []error) {
	if len(writerDSNs)+len(readerDSNs) == 0 {
		return nil, []error{ErrNoDSNs}
	}

	writers, writerErrs := open(name, writerDSNs)
	readers, readerErrs := open(name, readerDSNs)
	errs := writerErrs
	errs = append(errs, readerErrs...)
	if err := errors.Join(errs...); err != nil {
		errs = append(errs, closeAll(writers)...)
		errs = append(errs, closeAll(readers)...)

		return nil, errs
	}

	return &DBs{driverName: name, writers: writers, readers: readers}, nil
}

func connect(name string, fs *os.FS, cfg *config.Config, opts ...telemetry.Option) (*DBs, error) {
	writers, err := resolveDSNs(fs, cfg.Writers)
	if err != nil {
		return nil, err
	}

	readers, err := resolveDSNs(fs, cfg.Readers)
	if err != nil {
		return nil, err
	}

	db, err := connectDBs(name, writers, readers, opts...)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

func resolveDSNs(fs *os.FS, dsns []config.DSN) ([]string, error) {
	resolved := make([]string, len(dsns))

	for i, dsn := range dsns {
		url, err := dsn.GetURL(fs)
		if err != nil {
			return nil, err
		}
		if len(url) == 0 {
			return nil, ErrEmptyDSN
		}

		resolved[i] = bytes.String(url)
	}

	return resolved, nil
}

func connectDBs(name string, writerDSNs, readerDSNs []string, opts ...telemetry.Option) (*DBs, error) {
	if len(writerDSNs)+len(readerDSNs) == 0 {
		return nil, ErrNoDSNs
	}

	db, errs := ConnectWritersReaders(name, writerDSNs, readerDSNs)
	if err := errors.Join(errs...); err != nil {
		if db != nil {
			err = errors.Join(err, db.Destroy())
		}

		return nil, err
	}

	var regs []metrics.Registration
	if metrics.IsEnabled() {
		opts := options(name, opts)

		writers := db.Writers()
		regs = append(regs, register(writers, "writer", opts...)...)

		readers := db.Readers()
		regs = append(regs, register(readers, "reader", opts...)...)
	}

	db.registrations = regs

	return db, nil
}

func open(name string, dsns []string) ([]*sql.DB, []error) {
	dbs := make([]*sql.DB, len(dsns))
	errs := make([]error, len(dsns))

	for i, dsn := range dsns {
		dbs[i], errs[i] = sql.Open(name, dsn)
	}

	return dbs, errs
}
