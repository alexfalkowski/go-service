package module_test

import (
	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/module"
)

func ExampleServer() {
	application := cli.NewApplication(func(commander cli.Commander) {
		server := commander.AddServer("serve", "Run the service", module.Server)
		server.AddConfig("file:./config.yml")
	})

	_ = application
	// Output:
}

func ExampleClient() {
	task := di.Register(func(lifecycle di.Lifecycle) {
		lifecycle.Append(di.Hook{
			OnStart: func(context.Context) error {
				// Perform the short-lived command action here.
				return nil
			},
		})
	})

	application := cli.NewApplication(func(commander cli.Commander) {
		client := commander.AddClient("migrate", "Run client tasks", module.Client, task)
		client.AddConfig("file:./config.yml")
	})

	_ = application
	// Output:
}
