package module_test

import (
	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/module"
)

func ExampleServer() {
	app := cli.NewApplication(func(commander cli.Commander) {
		server := commander.AddServer("serve", "Run the service", module.Server)
		server.AddInput("file:./config.yml")
	})

	_ = app
	// Output:
}

func ExampleClient() {
	app := cli.NewApplication(func(commander cli.Commander) {
		client := commander.AddClient("migrate", "Run client tasks", module.Client)
		client.AddInput("file:./config.yml")
	})

	_ = app
	// Output:
}
