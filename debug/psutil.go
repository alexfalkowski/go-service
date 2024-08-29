package debug

import (
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// RegisterPprof for debug.
func RegisterPsutil(srv *Server, enc *encoding.Map) {
	mux := srv.ServeMux()

	psutil := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		defer func() {
			if r := recover(); r != nil {
				err := errors.Prefix("health", runtime.ConvertRecover(r))
				nh.WriteError(ctx, res, err, status.Code(err))
			}
		}()

		ct := content.NewFromRequest(req)
		res.Header().Add(content.TypeKey, ct.Media)

		e := ct.Encoder(enc)

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

		err = e.Encode(res, data)
		runtime.Must(err)
	}

	mux.HandleFunc("/debug/psutil", psutil)
}
