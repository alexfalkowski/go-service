package meta

import "github.com/alexfalkowski/go-service/v2/context"

// ExtractIncoming extracts incoming gRPC metadata from ctx.
//
// If no incoming metadata is present, it returns an empty metadata map.
// The returned metadata is a copy and is safe to mutate by the caller.
func ExtractIncoming(ctx context.Context) Map {
	return extract(FromIncomingContext(ctx))
}

// ExtractOutgoing extracts outgoing gRPC metadata from ctx.
//
// If no outgoing metadata is present, it returns an empty metadata map.
// The returned metadata is a copy and is safe to mutate by the caller.
func ExtractOutgoing(ctx context.Context) Map {
	return extract(FromOutgoingContext(ctx))
}

func extract(md Map, ok bool) Map {
	if !ok {
		return Map{}
	}

	return md.Copy()
}
