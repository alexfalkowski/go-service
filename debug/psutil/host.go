package psutil

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/shirou/gopsutil/v4/host"
)

// NewHost for debug.
func NewHost(ctx context.Context) *Host {
	info, _ := host.InfoWithContext(ctx)

	return &Host{
		Info: info,
	}
}

// Host contains host/system details collected for the debug endpoint.
type Host struct {
	// Info contains host information (OS, platform, uptime, etc.).
	Info *host.InfoStat `yaml:"info,omitempty" json:"info,omitempty" toml:"info,omitempty"`
}
