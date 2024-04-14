package metrics

import (
	"google.golang.org/grpc"
)

type kind string

const (
	unary        kind = "unary"
	clientStream kind = "client_stream"
	serverStream kind = "server_stream"
	bidiStream   kind = "bidi_stream"
)

func streamKind(info *grpc.StreamServerInfo) kind {
	if info.IsClientStream && !info.IsServerStream {
		return clientStream
	} else if !info.IsClientStream && info.IsServerStream {
		return serverStream
	}

	return bidiStream
}

func clientStreamKind(desc *grpc.StreamDesc) kind {
	if desc.ClientStreams && !desc.ServerStreams {
		return clientStream
	} else if !desc.ClientStreams && desc.ServerStreams {
		return serverStream
	}

	return bidiStream
}
