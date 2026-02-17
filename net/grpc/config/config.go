package config

// Config holds the bind address for the gRPC server.
type Config struct {
	// Address is the bind address for the gRPC server (for example "tcp://:9090").
	Address string
}
