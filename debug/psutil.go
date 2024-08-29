package debug

import (
	"context"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// RegisterPprof for debug.
func RegisterPsutil(srv *Server, enc *encoding.Map) {
	mux := srv.ServeMux()
	h := content.NewHandler("debug", enc, func(ctx context.Context) any {
		data := make(map[string]any)

		i, err := cpu.InfoWithContext(ctx)
		runtime.Must(err)

		t, err := cpu.TimesWithContext(ctx, true)
		runtime.Must(err)

		data["cpu"] = map[string]any{
			"info":  i,
			"times": t,
		}

		sm, err := mem.SwapMemoryWithContext(ctx)
		runtime.Must(err)

		sd, err := mem.SwapDevicesWithContext(ctx)
		runtime.Must(err)

		vm, err := mem.VirtualMemoryWithContext(ctx)
		runtime.Must(err)

		data["mem"] = map[string]any{
			"swap":    sm,
			"devices": sd,
			"virtual": vm,
		}

		return data
	})

	mux.HandleFunc("/debug/psutil", h)
}
