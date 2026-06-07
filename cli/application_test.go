package cli_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/stretchr/testify/require"
)

func TestApplicationDuplicateCommand(t *testing.T) {
	var err error

	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				err = runtime.ConvertRecover(recovered)
			}
		}()

		cli.NewApplication(func(c cli.Commander) {
			c.AddServer("server", "Start the server.", test.Options()...)
			c.AddClient("server", "Start the client.", test.Options()...)
		})
	}()

	require.ErrorIs(t, err, cli.ErrCommandRegistered)
	require.ErrorContains(t, err, "server")
}
