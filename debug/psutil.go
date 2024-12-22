package debug

import (
	"context"

	"github.com/alexfalkowski/go-service/maps"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// RegisterPprof for debug.
func RegisterPsutil(srv *Server, cont *content.Content) {
	mux := srv.ServeMux()
	h := content.NewHandler(cont, "debug", func(ctx context.Context) (*maps.StringAny, error) {
		data := maps.StringAny{}

		i, err := cpu.InfoWithContext(ctx)
		if err != nil {
			return nil, err
		}

		t, err := cpu.TimesWithContext(ctx, true)
		if err != nil {
			return nil, err
		}

		data["cpu"] = maps.StringAny{
			"info":  i,
			"times": t,
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

		data["mem"] = maps.StringAny{
			"swap":    sm,
			"devices": sd,
			"virtual": vm,
		}

		return &data, nil
	})

	mux.HandleFunc("/debug/psutil", h)
}
