package breaker

import "github.com/alexfalkowski/go-service/v2/transport/breaker"

// NewConfig returns an HTTP breaker config from shared breaker mechanics.
//
// If cfg is nil, NewConfig returns nil so client option wiring can preserve disabled breakers.
func NewConfig(cfg *breaker.Config, statusCodes ...int) *Config {
	if cfg == nil {
		return nil
	}

	return &Config{
		Config:      cfg,
		StatusCodes: statusCodes,
	}
}

// Config configures HTTP client circuit breaker behavior.
//
// It embeds shared breaker mechanics and adds HTTP-specific failure classification.
type Config struct {
	// Config carries shared breaker mechanics such as thresholds, intervals, and open-state timeout.
	*breaker.Config `yaml:",inline" json:",inline" toml:",inline"`

	// StatusCodes replaces the default HTTP breaker failure status codes.
	//
	// When empty, HTTP breakers use their default failure classification. When set, only these configured
	// failure status codes are counted as breaker failures; include the defaults here as well when custom
	// configuration should extend rather than replace the default list.
	StatusCodes []int `yaml:"status_codes,omitempty" json:"status_codes,omitempty" toml:"status_codes,omitempty" validate:"omitempty,dive,gte=400,lte=599"`
}

// Options returns circuit breaker options derived from c.
func (c *Config) Options() []Option {
	opts := []Option{WithSettings(c.Settings())}
	if len(c.StatusCodes) > 0 {
		opts = append(opts, WithFailureStatuses(c.StatusCodes...))
	}

	return opts
}
