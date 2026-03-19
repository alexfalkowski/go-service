package test

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// StatusURL returns the status service URL for the supplied status name.
//
// The helper respects the `STATUS_PORT` environment variable and otherwise
// falls back to the default test port `6000`.
func StatusURL(status string) string {
	port := cmp.Or(os.Getenv("STATUS_PORT"), "6000")

	return strings.Concat("http://localhost:", port, "/v1/status/", status)
}
