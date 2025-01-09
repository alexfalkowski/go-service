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
func RegisterPsutil(srv *Server, cont *content.Content) {
	mux := srv.ServeMux()
	handler := content.NewHandler(cont, "debug", func(ctx context.Context) (*Response, error) {
		res := &Response{}

		info, err := cpu.InfoWithContext(ctx)
		if err != nil {
			return nil, err
		}

		times, err := cpu.TimesWithContext(ctx, true)
		if err != nil {
			return nil, err
		}

		res.CPU = &CPU{
			Info:  info,
			Times: times,
		}

		swapMem, err := mem.SwapMemoryWithContext(ctx)
		if err != nil {
			return nil, err
		}

		swapDev, err := mem.SwapDevicesWithContext(ctx)
		if err != nil {
			return nil, err
		}

		vms, err := mem.VirtualMemoryWithContext(ctx)
		if err != nil {
			return nil, err
		}

		res.Mem = &Mem{
			Swap:    swapMem,
			Devices: swapDev,
			Virtual: vms,
		}

		return res, nil
	})

	mux.HandleFunc("/debug/psutil", handler)
}
