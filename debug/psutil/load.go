package psutil

import (
	"context"

	"github.com/shirou/gopsutil/v4/load"
)

// NewLoad for debug.
func NewLoad(ctx context.Context) *Load {
	avg, _ := load.AvgWithContext(ctx)

	return &Load{
		Avg: avg,
	}
}

// Load for debug.
type Load struct {
	Avg *load.AvgStat `yaml:"avg,omitempty" json:"avg,omitempty" toml:"avg,omitempty"`
}
