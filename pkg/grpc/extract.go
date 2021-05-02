package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func extractIncoming(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.MD{}
	}

	return md.Copy()
}

func extractOutgoing(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.MD{}
	}

	return md.Copy()
}
