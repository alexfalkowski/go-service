package psutil

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
)

// NewHandler constructs the HTTP handler that returns a psutil snapshot response.
//
// The handler collects a point-in-time view of system information using gopsutil-backed helpers in this package
// (CPU, host info, load averages, memory/swap, and network counters).
//
// Collection is best-effort: underlying collection helpers intentionally ignore gopsutil errors, so individual
// sections may be partially populated or empty depending on platform support and runtime permissions.
//
// The handler is typically registered under /debug/psutil (namespaced by service name via debug/http.Pattern).
func NewHandler(cont *content.Content) http.HandlerFunc {
	return content.NewHandler(cont, func(ctx context.Context) (*Response, error) {
		res := &Response{
			CPU:  NewCPU(ctx),
			Host: NewHost(ctx),
			Load: NewLoad(ctx),
			Mem:  NewMem(ctx),
			Net:  NewNet(ctx),
		}

		return res, nil
	})
}

// Response is the response body returned by the psutil debug endpoint.
//
// All fields are optional and may be nil if collection failed or if the platform does not support the underlying
// metric. Fields are tagged for common config/encoding formats (YAML/JSON/TOML) to support standard response
// encoders used in go-service.
type Response struct {
	// CPU contains CPU information and time statistics.
	CPU *CPU `yaml:"cpu,omitempty" json:"cpu,omitempty" toml:"cpu,omitempty"`

	// Host contains host/system information (OS, platform, uptime, etc.).
	Host *Host `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`

	// Load contains system load averages (where supported by the OS).
	Load *Load `yaml:"load,omitempty" json:"load,omitempty" toml:"load,omitempty"`

	// Mem contains memory and swap statistics.
	Mem *Mem `yaml:"mem,omitempty" json:"mem,omitempty" toml:"mem,omitempty"`

	// Net contains network I/O counters.
	Net *Net `yaml:"net,omitempty" json:"net,omitempty" toml:"net,omitempty"`
}
