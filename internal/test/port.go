package test

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// RandomAddress address for tests.
func RandomAddress() string {
	network, address := RandomNetworkHost()

	return strings.Concat(network, "://", address)
}

// RandomHost for tests.
func RandomHost() string {
	_, address := RandomNetworkHost()

	return address
}

// RandomNetworkHost for tests.
func RandomNetworkHost() (string, string) {
	l, err := net.Listen(context.Background(), "tcp", "localhost:0")
	runtime.Must(err)

	defer l.Close()

	addr := l.Addr().String()
	addr = strings.ReplaceAll(addr, "127.0.0.1", "localhost")

	return l.Addr().Network(), addr
}
