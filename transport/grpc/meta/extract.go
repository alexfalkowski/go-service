package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"google.golang.org/grpc/metadata"
)

// ExtractIncoming extracts incoming gRPC metadata from ctx.
//
// If no incoming metadata is present, it returns an empty metadata map.
// The returned metadata is a copy and is safe to mutate by the caller.
func ExtractIncoming(ctx context.Context) metadata.MD {
	return extract(metadata.FromIncomingContext(ctx))
}

// ExtractOutgoing extracts outgoing gRPC metadata from ctx.
//
// If no outgoing metadata is present, it returns an empty metadata map.
// The returned metadata is a copy and is safe to mutate by the caller.
func ExtractOutgoing(ctx context.Context) metadata.MD {
	return extract(metadata.FromOutgoingContext(ctx))
}

func extract(md metadata.MD, ok bool) metadata.MD {
	if !ok {
		return metadata.MD{}
	}

	return md.Copy()
}
