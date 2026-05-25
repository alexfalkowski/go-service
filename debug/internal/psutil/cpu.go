package psutil

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/shirou/gopsutil/v4/cpu"
)

// NewCPU collects CPU information and times for the debug endpoint.
func NewCPU(ctx context.Context) *CPU {
	info, _ := cpu.InfoWithContext(ctx)
	times, _ := cpu.TimesWithContext(ctx, true)

	return &CPU{
		Info:  info,
		Times: times,
	}
}

// CPU contains CPU details collected for the debug endpoint.
type CPU struct {
	// Info contains static CPU information (model, cores, cache sizes, etc.).
	Info []cpu.InfoStat `yaml:"info,omitempty" json:"info,omitempty" toml:"info,omitempty"`

	// Times contains per-CPU time statistics (user/system/idle/etc.).
	Times []cpu.TimesStat `yaml:"times,omitempty" json:"times,omitempty" toml:"times,omitempty"`
}
