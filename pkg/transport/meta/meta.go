package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/meta"
)

const (
	requestID = "app.request_id"
)

// WithRequestID for transport.
func WithRequestID(ctx context.Context, id string) context.Context {
	return meta.WithAttribute(ctx, requestID, id)
}

// RequestID for transport.
func RequestID(ctx context.Context) string {
	return meta.Attribute(ctx, requestID)
}
