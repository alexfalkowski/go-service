package debug

import (
	"context"

	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type (
	// Request for debug.
	Request struct{}

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
	h := content.NewQueryHandler(cont, "debug", func(ctx context.Context, _ *Request) (*Response, error) {
		res := &Response{}

		i, err := cpu.InfoWithContext(ctx)
		if err != nil {
			return nil, err
		}

		t, err := cpu.TimesWithContext(ctx, true)
		if err != nil {
			return nil, err
		}

		res.CPU = &CPU{
			Info:  i,
			Times: t,
		}

		sm, err := mem.SwapMemoryWithContext(ctx)
		if err != nil {
			return nil, err
		}

		sd, err := mem.SwapDevicesWithContext(ctx)
		if err != nil {
			return nil, err
		}

		vm, err := mem.VirtualMemoryWithContext(ctx)
		if err != nil {
			return nil, err
		}

		res.Mem = &Mem{
			Swap:    sm,
			Devices: sd,
			Virtual: vm,
		}

		return res, nil
	})

	mux.HandleFunc("/debug/psutil", h)
}
