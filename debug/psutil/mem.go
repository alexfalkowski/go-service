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

// Mem contains memory details collected for the debug endpoint.
type Mem struct {
	// Swap contains overall swap memory statistics.
	Swap *mem.SwapMemoryStat `yaml:"swap,omitempty" json:"swap,omitempty" toml:"swap,omitempty"`

	// Virtual contains virtual memory statistics.
	Virtual *mem.VirtualMemoryStat `yaml:"virtual,omitempty" json:"virtual,omitempty" toml:"virtual,omitempty"`

	// Devices contains per-device swap usage information.
	Devices []*mem.SwapDevice `yaml:"devices,omitempty" json:"devices,omitempty" toml:"devices,omitempty"`
}
