package config

// Config configures how a go-service gRPC server binds to the network.
//
// The configuration is consumed by the gRPC server wiring in [github.com/alexfalkowski/go-service/v2/net/grpc/server].
// In particular, `server.NewServer` splits Address into a network and address
// and then creates a listener via `net.Listen`.
//
// Address may be in the go-service "network address" format:
//
//	<network>://<address>
//
// or a raw listen address. Raw addresses default to the "tcp" network.
//
// Examples:
//
//	tcp://:9090
//	:9090
type Config struct {
	// Address is the bind address for the gRPC server.
	//
	// It may use the go-service network address format (for example "tcp://:9090")
	// or a raw listen address such as ":9090", which defaults to the "tcp" network.
	Address string
}
