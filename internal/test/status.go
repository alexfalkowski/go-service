package test

import (
	"cmp"

	"github.com/alexfalkowski/go-service/os"
)

// StatusURL for test.
func StatusURL(status string) string {
	port := cmp.Or(os.Getenv("STATUS_PORT"), "6000")

	return "http://localhost:" + port + "/v1/status/" + status
}
