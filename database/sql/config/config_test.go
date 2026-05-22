package config_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestConfigRejectsNegativeConnMaxLifetime(t *testing.T) {
	cfg := &config.Config{ConnMaxLifetime: -time.Second}
	require.Error(t, test.Validator.Struct(cfg))
}
