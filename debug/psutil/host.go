package psutil

import (
	"context"

	"github.com/shirou/gopsutil/v4/host"
)

// NewHost for debug.
func NewHost(ctx context.Context) *Host {
	info, _ := host.InfoWithContext(ctx)

	return &Host{
		Info: info,
	}
}

// Host for debug.
type Host struct {
	Info *host.InfoStat `yaml:"info,omitempty" json:"info,omitempty" toml:"info,omitempty"`
}
