package test

import "github.com/alexfalkowski/go-service/v2/net"

const (
	localNetwork = "tcp"
	localAddress = "localhost:0"
)

// RandomAddress returns a local transport address in `<network>://<host:port>` form that lets the server pick an ephemeral port.
func RandomAddress() string {
	return net.NetworkAddress(localNetwork, localAddress)
}

// RandomHost returns a local `host:port` listen address that lets the server pick an ephemeral port.
func RandomHost() string {
	return localAddress
}

// RandomNetworkHost returns the local network name and listen address that let the server pick an ephemeral port.
func RandomNetworkHost() (string, string) {
	return localNetwork, localAddress
}

// BoundAddress rewrites configured to the bound listener address while preserving the configured network and host when possible.
func BoundAddress(configured, actual string) string {
	network, address := net.ListenNetworkAddress(configured)
	host, _, err := net.SplitHostPort(address)
	if err != nil || host == "" {
		return net.NetworkAddress(network, actual)
	}

	_, port, err := net.SplitHostPort(actual)
	if err != nil {
		return net.NetworkAddress(network, actual)
	}

	return net.NetworkAddress(network, net.JoinHostPort(host, port))
}
