package debug

import (
	"net/http"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func (s *server) psutil(resp http.ResponseWriter, req *http.Request) {
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

	b, _ := s.json.Marshal(resp)

	resp.Write(b)
}
