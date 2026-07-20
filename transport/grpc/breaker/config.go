package breaker

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/transport/breaker"
)

// NewConfig returns a gRPC breaker config from shared breaker mechanics.
//
// If cfg is nil, NewConfig returns nil so client option wiring can preserve disabled breakers.
func NewConfig(cfg *breaker.Config, codes ...codes.Code) *Config {
	if cfg == nil {
		return nil
	}

	return &Config{
		Config: cfg,
		Codes:  codes,
	}
}

// Config configures gRPC unary client circuit breaker behavior.
//
// It embeds shared breaker mechanics and adds gRPC-specific failure classification.
type Config struct {
	// Config carries shared breaker mechanics such as thresholds, intervals, and open-state timeout.
	*breaker.Config `yaml:",inline" json:",inline" toml:",inline"`

	// Codes replaces the default gRPC breaker failure codes.
	//
	// When empty, gRPC breakers use their default failure classification. When set, only these configured
	// non-OK status codes are counted as breaker failures; include the defaults here as well when custom
	// configuration should extend rather than replace the default list. [codes.Canceled] is always counted
	// as a success even when included here, so a caller aborting an in-flight call never trips the breaker.
	Codes []codes.Code `yaml:"codes,omitempty" json:"codes,omitempty" toml:"codes,omitempty" validate:"omitempty,dive,gt=0,lte=16"`
}

// Options returns circuit breaker options derived from c.
func (c *Config) Options() []Option {
	opts := []Option{WithSettings(c.Settings())}
	if len(c.Codes) > 0 {
		opts = append(opts, WithFailureCodes(c.Codes...))
	}

	return opts
}
