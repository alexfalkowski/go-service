package token

import (
	"context"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/runtime"
)

// RegisterKeys for limiter.
func RegisterKeys() {
	limiter.RegisterKey("token", Key)
}

// Key for token.
func Key(ctx context.Context) meta.Valuer {
	t, err := ExtractToken(ctx)
	runtime.Must(err)

	return meta.Redacted(t)
}
