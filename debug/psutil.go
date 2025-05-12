package debug

import (
	"context"

	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type (
	// CPU for debug.
	CPU struct {
		Info  []cpu.InfoStat  `yaml:"info,omitempty" json:"info,omitempty" toml:"info,omitempty"`
		Times []cpu.TimesStat `yaml:"times,omitempty" json:"times,omitempty" toml:"times,omitempty"`
	}

	// Mem for debug.
	Mem struct {
		Swap    *mem.SwapMemoryStat    `yaml:"swap,omitempty" json:"swap,omitempty" toml:"swap,omitempty"`
		Virtual *mem.VirtualMemoryStat `yaml:"virtual,omitempty" json:"virtual,omitempty" toml:"virtual,omitempty"`
		Devices []*mem.SwapDevice      `yaml:"devices,omitempty" json:"devices,omitempty" toml:"devices,omitempty"`
	}

	// Response for debug.
	Response struct {
		CPU *CPU `yaml:"cpu,omitempty" json:"cpu,omitempty" toml:"cpu,omitempty"`
		Mem *Mem `yaml:"mem,omitempty" json:"mem,omitempty" toml:"mem,omitempty"`
	}
)

// RegisterPprof for debug.
func RegisterPsutil(mux *ServeMux, cont *content.Content) {
	handler := content.NewHandler(cont, func(ctx context.Context) (*Response, error) {
		info, _ := cpu.InfoWithContext(ctx)
		times, _ := cpu.TimesWithContext(ctx, true)
		swapMem, _ := mem.SwapMemoryWithContext(ctx)
		swapDev, _ := mem.SwapDevicesWithContext(ctx)
		vms, _ := mem.VirtualMemoryWithContext(ctx)
		res := &Response{
			CPU: &CPU{
				Info:  info,
				Times: times,
			},
			Mem: &Mem{
				Swap:    swapMem,
				Devices: swapDev,
				Virtual: vms,
			},
		}

		return res, nil
	})

	mux.HandleFunc("/debug/psutil", handler)
}
