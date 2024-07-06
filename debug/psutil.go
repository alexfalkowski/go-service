package debug

import (
	"net/http"

	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// RegisterPprof for debug.
func RegisterPsutil(srv *Server, mar *json.Marshaller) {
	mux := srv.ServeMux()

	psutil := func(resp http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		r := make(map[string]any)

		i, _ := cpu.InfoWithContext(ctx)
		t, _ := cpu.TimesWithContext(ctx, true)
		r["cpu"] = map[string]any{
			"info":  i,
			"times": t,
		}

		sm, _ := mem.SwapMemoryWithContext(ctx)
		sd, _ := mem.SwapDevicesWithContext(ctx)
		vm, _ := mem.VirtualMemoryWithContext(ctx)
		r["mem"] = map[string]any{
			"swap":    sm,
			"devices": sd,
			"virtual": vm,
		}

		content.AddJSONHeader(resp.Header())

		b, _ := mar.Marshal(r)

		resp.Write(b)
	}

	mux.HandleFunc("/debug/psutil", psutil)
}
