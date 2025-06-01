package psutil

import (
	"context"

	"github.com/shirou/gopsutil/v4/cpu"
)

// NewCPU for debug.
func NewCPU(ctx context.Context) *CPU {
	info, _ := cpu.InfoWithContext(ctx)
	times, _ := cpu.TimesWithContext(ctx, true)

	return &CPU{
		Info:  info,
		Times: times,
	}
}

// CPU for debug.
type CPU struct {
	Info  []cpu.InfoStat  `yaml:"info,omitempty" json:"info,omitempty" toml:"info,omitempty"`
	Times []cpu.TimesStat `yaml:"times,omitempty" json:"times,omitempty" toml:"times,omitempty"`
}
