package config

// Config configures how a go-service gRPC server binds to the network.
//
// The configuration is consumed by the gRPC server wiring in `net/grpc/server`.
// In particular, `server.NewServer` splits Address into a network and address
// and then creates a listener via `net.Listen`.
//
// Address is expected to be in the go-service "network address" format:
//
//	<network>://<address>
//
// Example:
//
//	tcp://:9090
type Config struct {
	// Address is the bind address for the gRPC server, expressed in the go-service
	// network address format (for example "tcp://:9090").
	//
	// This value is split into network/address components and passed to net.Listen.
	Address string
}
