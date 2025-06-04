package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/meta"
)

// UserIDKey for token.
const UserIDKey = "subject"

// Ignored is an alias for meta.Ignored.
var Ignored = meta.Ignored

// WithUserID for token.
func WithUserID(ctx context.Context, id meta.Value) context.Context {
	return meta.WithAttribute(ctx, UserIDKey, id)
}

// UserID for token.
func UserID(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, UserIDKey)
}
