package sql_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/stretchr/testify/require"
)

func TestOpenUnknownDriver(t *testing.T) {
	db, err := sql.Open("go-service-missing-driver", "benchmark")

	require.Error(t, err)
	require.Nil(t, db)
}

func TestConfigIsEnabled(t *testing.T) {
	require.False(t, (*sql.Config)(nil).IsEnabled())
	require.True(t, (&sql.Config{}).IsEnabled())
}
