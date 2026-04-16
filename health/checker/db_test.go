package checker_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestDBCheckerWithoutConnections(t *testing.T) {
	db, errs := sql.ConnectMasterSlaves("pg", nil, nil)
	require.Empty(t, errs)

	check := checker.NewDBChecker(db, time.Second)
	err := check.Check(t.Context())
	require.ErrorIs(t, err, checker.ErrNoConnections)
}
