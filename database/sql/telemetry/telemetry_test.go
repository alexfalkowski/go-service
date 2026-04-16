package telemetry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	db, err := telemetry.Open("missing", "dsn")
	require.Nil(t, db)
	require.Error(t, err)
	require.ErrorContains(t, err, `sql: unknown driver "missing"`)
}
