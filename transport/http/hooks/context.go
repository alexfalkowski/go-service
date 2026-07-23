package hooks

import "github.com/alexfalkowski/go-service/v2/context"

const webhookIDKey = context.Key("webhook-id")

// WithWebhookID returns ctx with the webhook id for one logical delivery.
func WithWebhookID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, webhookIDKey, id)
}

// WebhookID returns the webhook id from ctx.
func WebhookID(ctx context.Context) string {
	id, _ := ctx.Value(webhookIDKey).(string)
	return id
}
