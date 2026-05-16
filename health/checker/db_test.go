package checker_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
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
	require.NoError(t, errors.Join(world.DB.Destroy()...))
	require.Error(t, check.Check(t.Context()))
}
