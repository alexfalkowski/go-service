package feature

import "github.com/alexfalkowski/go-service/v2/config/client"

// Config carries optional settings for service-owned OpenFeature integrations.
//
// It embeds [client.Config] to reuse common client-side configuration fields that may be
// shared across feature-related integrations (for example address, timeout, TLS, retry/limiter,
// and key-value options).
//
// The feature package does not consume Config directly or construct a provider from it. Services that
// need a remote/custom OpenFeature provider should consume this config in their own provider constructor
// and provide the resulting [github.com/open-feature/go-sdk/openfeature.FeatureProvider] to the DI graph.
//
// # Optional pointers and "enabled" semantics
//
// This config is intentionally optional. By convention across go-service configuration types, a nil
// *[Config] is treated as "feature disabled". The embedded *[client.Config] is also optional;
// [Config.IsEnabled] returns true only when both the outer *[Config] and the embedded
// *[client.Config] are non-nil/enabled.
//
// Note: provider registration itself is controlled by the presence of an OpenFeature provider in the
// DI graph (see [Register]). A service may have feature config present without wiring
// a provider, in which case OpenFeature behaves with its default provider semantics.
type Config struct {
	*client.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled reports whether feature configuration is present and enabled.
//
// It returns true only when both the feature wrapper config and the embedded client config are non-nil
// and enabled.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}
