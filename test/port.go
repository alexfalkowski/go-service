package test

import (
	"net"
	"strconv"
)

// Port for test.
func Port() string {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return ""
	}

	defer l.Close()

	p := l.Addr().(*net.TCPAddr).Port

	return strconv.Itoa(p)
}
