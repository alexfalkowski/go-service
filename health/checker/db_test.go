package checker_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
)

func TestDBCheckerWithoutConnections(t *testing.T) {
	db, errs := sql.ConnectMasterSlaves("pg", nil, nil)
	require.Empty(t, errs)

	check := checker.NewDBChecker(db, time.Second)
	require.ErrorIs(t, check.Check(t.Context()), checker.ErrNoConnections)
}

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

func TestDBCheckerPingsMastersAndSlaves(t *testing.T) {
	driverName := registerPingDriver(t)
	db, errs := sql.ConnectMasterSlaves(driverName, []string{"master"}, []string{"slave"})
	require.Empty(t, errs)
	defer func() {
		require.NoError(t, db.Destroy())
	}()

	check := checker.NewDBChecker(db, time.Second)
	err := check.Check(t.Context())
	require.ErrorIs(t, err, errMasterPing)
	require.ErrorIs(t, err, errSlavePing)
}

func TestDBCheckerReturnsTimeoutCause(t *testing.T) {
	driverName := registerPingDriver(t)
	db, errs := sql.ConnectMasterSlaves(driverName, []string{"timeout"}, nil)
	require.Empty(t, errs)
	defer func() {
		require.NoError(t, db.Destroy())
	}()

	check := checker.NewDBChecker(db, time.Millisecond)
	require.ErrorIs(t, check.Check(t.Context()), checker.ErrPingTimeout)
}

var (
	errMasterPing = errors.New("master ping")
	errSlavePing  = errors.New("slave ping")
	pingDriverID  sync.Uint64
)

func registerPingDriver(t *testing.T) string {
	t.Helper()

	name := "health-checker-" + strconv.FormatUint(pingDriverID.Add(1), 10)
	require.NoError(t, driver.Register(name, pingDriver{}))
	return name
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
	case "master":
		return errMasterPing
	case "slave":
		return errSlavePing
	case "timeout":
		<-ctx.Done()
		return context.Cause(ctx)
	default:
		return nil
	}
}
