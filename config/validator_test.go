package config_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

type invalidConfigSize struct {
	Size string `validate:"config_size"`
}

func TestValidatorRejectsConfigSizeWithNonInteger(t *testing.T) {
	cfg := &invalidConfigSize{Size: "64B"}

	require.Error(t, test.Validator.Struct(cfg))
}
