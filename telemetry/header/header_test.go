package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/stretchr/testify/require"
)

func TestSecrets(t *testing.T) {
	require.NoError(t, header.Map{"test": test.FilePath("secrets/hooks")}.Secrets(test.FS))
	require.Error(t, header.Map{"test": test.FilePath("none")}.Secrets(test.ErrFS))
}
