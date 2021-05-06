package meta

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// ExtractIncoming for meta.
func ExtractIncoming(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.MD{}
	}

	return md.Copy()
}

// ExtractOutgoing for meta.
func ExtractOutgoing(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.MD{}
	}

	return md.Copy()
}
