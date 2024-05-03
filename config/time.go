package config

import (
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/time/ntp"
	"github.com/alexfalkowski/go-service/time/nts"
)

// NTPConfig from time.
func (cfg *Config) NTPConfig() *ntp.Config {
	if !time.IsEnabled(cfg.Time) {
		return nil
	}

	return cfg.Time.NTP
}

// NTSConfig from time.
func (cfg *Config) NTSConfig() *nts.Config {
	if !time.IsEnabled(cfg.Time) {
		return nil
	}

	return cfg.Time.NTS
}
