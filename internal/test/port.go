package test

import (
	"net"
	"strings"

	"github.com/alexfalkowski/go-service/runtime"
)

// Address returns a random address for tests.
func Address() string {
	l, err := net.Listen("tcp", "localhost:0")
	runtime.Must(err)

	defer l.Close()

	addr := l.Addr().String()
	addr = strings.ReplaceAll(addr, "127.0.0.1", "localhost")

	return addr
}
