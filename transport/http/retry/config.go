package retry

import "github.com/alexfalkowski/go-service/v2/transport/retry"

// Config configures retry behavior for outbound HTTP requests.
//
// It embeds shared retry mechanics and adds HTTP-specific failure classification.
type Config struct {
	// Config carries shared retry mechanics such as attempts, timeout, and backoff.
	*retry.Config `yaml:",inline" json:",inline" toml:",inline"`

	// StatusCodes replaces the default retryable HTTP response/status error codes.
	//
	// When empty, HTTP retry uses its default failure classification. When set, only these
	// configured failure status codes are retryable; include the defaults here as well when
	// custom configuration should extend rather than replace the default list.
	//
	// Decoded configuration accepts only 4xx and 5xx status codes. Direct public API
	// construction ignores non-failure status codes so successful responses cannot become retryable.
	StatusCodes []int `yaml:"status_codes,omitempty" json:"status_codes,omitempty" toml:"status_codes,omitempty" validate:"omitempty,dive,gte=400,lte=599"`
}
