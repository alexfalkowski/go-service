package test

import (
	"net"
	"strings"

	"github.com/alexfalkowski/go-service/v2/runtime"
)

// Address returns a random address for tests.
func Address() string {
	l, err := net.Listen("tcp", "localhost:0")
	runtime.Must(err)

	defer l.Close()

	return strings.ReplaceAll(l.Addr().String(), "127.0.0.1", "localhost")
}
