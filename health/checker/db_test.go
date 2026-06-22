package checker_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestDBCheckerPingsWorldDB(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	check := checker.NewDBChecker(world.DB, time.Second)
	require.NoError(t, check.Check(t.Context()))
}

func TestDBCheckerReturnsPingError(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	check := checker.NewDBChecker(world.DB, time.Second)
	require.NoError(t, world.DB.Destroy())
	require.Error(t, check.Check(t.Context()))
}

func TestDBCheckerPingsWritersAndReaders(t *testing.T) {
	db := newPingDB(t, []string{"writer"}, []string{"reader"})

	check := checker.NewDBChecker(db, time.Second)
	err := check.Check(t.Context())
	require.ErrorIs(t, err, errWriterPing)
	require.ErrorIs(t, err, errReaderPing)
	require.Contains(t, err.Error(), "db writer[0]")
	require.Contains(t, err.Error(), "db reader[0]")
}

func TestDBCheckerReturnsNoConnections(t *testing.T) {
	check := checker.NewDBChecker(&sql.DBs{}, time.Second)

	require.ErrorIs(t, check.Check(t.Context()), checker.ErrNoConnections)
}

func TestDBCheckerReturnsTimeoutCause(t *testing.T) {
	db := newPingDB(t, []string{"timeout"}, nil)

	check := checker.NewDBChecker(db, time.Millisecond)
	require.ErrorIs(t, check.Check(t.Context()), checker.ErrPingTimeout)
}

func TestDBCheckerReturnsTimeoutCauseWhenWaitingForConnection(t *testing.T) {
	db := newPingDB(t, []string{"writer"}, nil)
	writers := db.Writers()
	require.Len(t, writers, 1)

	writer := writers[0]
	writer.SetMaxOpenConns(1)
	conn, err := writer.Conn(t.Context())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	check := checker.NewDBChecker(db, time.Millisecond)
	require.ErrorIs(t, check.Check(t.Context()), checker.ErrPingTimeout)
}

var (
	errWriterPing = errors.New("writer ping")
	errReaderPing = errors.New("reader ping")
)

func newPingDB(t *testing.T, writers, readers []string) *sql.DBs {
	t.Helper()

	driverName := registerPingDriver(t)
	db, errs := sql.ConnectWritersReaders(driverName, writers, readers)
	require.Empty(t, errs)
	t.Cleanup(func() {
		require.NoError(t, db.Destroy())
	})

	return db
}

func registerPingDriver(t *testing.T) string {
	t.Helper()

	return test.RegisterSQLDriver(t, "health-checker-", pingDriver{})
}

type pingDriver struct{}

func (pingDriver) Open(name string) (driver.Conn, error) {
	return pingConn{name: name}, nil
}

type pingConn struct {
	name string
}

func (pingConn) Prepare(string) (driver.Stmt, error) {
	return nil, driver.ErrSkip
}

func (pingConn) Close() error {
	return nil
}

func (pingConn) Begin() (driver.Tx, error) {
	return nil, driver.ErrSkip
}

func (c pingConn) Ping(ctx context.Context) error {
	switch c.name {
	case "writer":
		return errWriterPing
	case "reader":
		return errReaderPing
	case "timeout":
		<-ctx.Done()
		return context.Cause(ctx)
	default:
		return nil
	}
}
