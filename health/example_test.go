package health_test

import (
	"fmt"

	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/time"
)

// ExampleRegistrations shows the Registrations mechanics in isolation, built directly
// rather than through DI. Real services get *server.Server and health.Registrations
// injected via module.Server; register and observe them the same way, wired with
// di.Register (see ExampleClient in module/example_test.go).
func ExampleRegistrations() {
	healthServer := server.NewServer()
	registrations := health.Registrations{
		server.NewRegistration("noop", time.Second.Duration(), checker.NewNoopChecker()),
	}

	healthServer.Register("payments", registrations...)
	for _, probe := range []string{"healthz", "livez", "readyz"} {
		if err := healthServer.Observe("payments", probe, "noop"); err != nil {
			panic(err)
		}
	}

	fmt.Println("registered healthz, livez, and readyz")
	// Output: registered healthz, livez, and readyz
}
