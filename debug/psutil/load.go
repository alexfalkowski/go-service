package psutil

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/shirou/gopsutil/v4/load"
)

// NewLoad for debug.
func NewLoad(ctx context.Context) *Load {
	avg, _ := load.AvgWithContext(ctx)

	return &Load{
		Avg: avg,
	}
}

// Load contains system load details collected for the debug endpoint.
type Load struct {
	// Avg contains load averages (for example 1m/5m/15m).
	Avg *load.AvgStat `yaml:"avg,omitempty" json:"avg,omitempty" toml:"avg,omitempty"`
}
