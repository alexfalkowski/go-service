package test

import (
	"net"
	"strconv"

	"github.com/alexfalkowski/go-service/runtime"
)

// Port for test.
func Port() string {
	l, err := net.Listen("tcp", "localhost:0")
	runtime.Must(err)

	defer l.Close()

	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}
