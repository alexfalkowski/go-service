package debug

import (
	"net/http"

	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// RegisterPprof for debug.
func RegisterPsutil(server *Server, json *marshaller.JSON) {
	psutil := func(resp http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		r := make(map[string]any)

		i, _ := cpu.InfoWithContext(ctx)
		t, _ := cpu.TimesWithContext(ctx, true)
		r["cpu"] = map[string]any{
			"info":  i,
			"times": t,
		}

		sw, _ := mem.SwapMemoryWithContext(ctx)
		vi, _ := mem.VirtualMemoryWithContext(ctx)
		r["mem"] = map[string]any{
			"swap":    sw,
			"virtual": vi,
		}

		resp.Header().Add("Content-Type", "application/json")

		b, _ := json.Marshal(r)

		resp.Write(b)
	}

	server.Mux.HandleFunc("/debug/psutil", psutil)
}
