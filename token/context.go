package token

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/meta"
)

// SubjectKey for token.
const SubjectKey = "subject"

// WithSubject for token.
func WithSubject(ctx context.Context, id meta.Value) context.Context {
	return meta.WithAttribute(ctx, SubjectKey, id)
}

// Subject for token.
func Subject(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, SubjectKey)
}
