package psutil

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/shirou/gopsutil/v4/mem"
)

// NewMem for debug.
func NewMem(ctx context.Context) *Mem {
	swapMem, _ := mem.SwapMemoryWithContext(ctx)
	swapDev, _ := mem.SwapDevicesWithContext(ctx)
	vms, _ := mem.VirtualMemoryWithContext(ctx)

	return &Mem{
		Swap:    swapMem,
		Devices: swapDev,
		Virtual: vms,
	}
}

// Mem for debug.
type Mem struct {
	Swap    *mem.SwapMemoryStat    `yaml:"swap,omitempty" json:"swap,omitempty" toml:"swap,omitempty"`
	Virtual *mem.VirtualMemoryStat `yaml:"virtual,omitempty" json:"virtual,omitempty" toml:"virtual,omitempty"`
	Devices []*mem.SwapDevice      `yaml:"devices,omitempty" json:"devices,omitempty" toml:"devices,omitempty"`
}
