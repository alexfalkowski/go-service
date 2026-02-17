package psutil

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
)

// NewHandler builds the HTTP handler that returns a psutil snapshot.
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
type Response struct {
	// CPU contains CPU information and times.
	CPU *CPU `yaml:"cpu,omitempty" json:"cpu,omitempty" toml:"cpu,omitempty"`

	// Host contains host/system information.
	Host *Host `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`

	// Load contains system load averages.
	Load *Load `yaml:"load,omitempty" json:"load,omitempty" toml:"load,omitempty"`

	// Mem contains memory and swap statistics.
	Mem *Mem `yaml:"mem,omitempty" json:"mem,omitempty" toml:"mem,omitempty"`

	// Net contains network I/O counters.
	Net *Net `yaml:"net,omitempty" json:"net,omitempty" toml:"net,omitempty"`
}
