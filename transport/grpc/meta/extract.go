package meta

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// ExtractIncoming for meta.
func ExtractIncoming(ctx context.Context) metadata.MD {
	return extract(metadata.FromIncomingContext(ctx))
}

// ExtractOutgoing for meta.
func ExtractOutgoing(ctx context.Context) metadata.MD {
	return extract(metadata.FromOutgoingContext(ctx))
}

func extract(md metadata.MD, ok bool) metadata.MD {
	if !ok {
		return metadata.MD{}
	}

	return md.Copy()
}
