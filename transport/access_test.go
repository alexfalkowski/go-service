package transport_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/stretchr/testify/require"
)

func TestNewAccessController(t *testing.T) {
	for _, tc := range []struct {
		name           string
		config         *access.Config
		wantController bool
	}{
		{name: "without access config"},
		{name: "with access config", config: test.NewAccessConfig(), wantController: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			controller, err := transport.NewAccessController(tc.config, test.FS)
			require.NoError(t, err)
			require.Equal(t, tc.wantController, controller != nil)
		})
	}
}
