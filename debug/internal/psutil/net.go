package psutil

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/shirou/gopsutil/v4/net"
)

// NewNet collects network I/O counters for the debug endpoint.
func NewNet(ctx context.Context) *Net {
	counters, _ := net.IOCountersWithContext(ctx, true)

	return &Net{Counters: counters}
}

// Net contains network details collected for the debug endpoint.
type Net struct {
	// Counters contains per-interface network I/O counters.
	Counters []net.IOCountersStat `yaml:"counters,omitempty" json:"counters,omitempty" toml:"counters,omitempty"`
}
