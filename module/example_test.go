package module_test

import (
	"github.com/alexfalkowski/go-service/v2/cli"
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
	application := cli.NewApplication(func(commander cli.Commander) {
		client := commander.AddClient("migrate", "Run client tasks", module.Client)
		client.AddConfig("file:./config.yml")
	})

	_ = application
	// Output:
}
