package prometheus

type grpcType string

const (
	unary        grpcType = "unary"
	clientStream grpcType = "client_stream"
	serverStream grpcType = "server_stream"
	bidiStream   grpcType = "bidi_stream"
)
