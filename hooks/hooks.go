package hooks

import (
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// New hook from config.
func New(cfg *Config) (*hooks.Webhook, error) {
	return hooks.NewWebhook(cfg.Secret)
}
