package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
)

const (
	securityAZP = "security.azp"
)

// WithAuthorizedParty for security.
func WithAuthorizedParty(ctx context.Context, azp string) context.Context {
	return meta.WithAttribute(ctx, securityAZP, azp)
}

// AuthorizedParty for security.
func AuthorizedParty(ctx context.Context) string {
	return meta.Attribute(ctx, securityAZP)
}
