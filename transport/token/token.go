package token

import (
	"context"

	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport/meta"
)

// Verify token from context.
func Verify(ctx context.Context, verifier token.Verifier) (context.Context, error) {
	token := meta.Authorization(ctx).Value()

	return verifier.Verify(ctx, []byte(token))
}
