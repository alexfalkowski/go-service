package psutil

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
)

// NewHandler for debug.
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

// Response for debug.
type Response struct {
	CPU  *CPU  `yaml:"cpu,omitempty" json:"cpu,omitempty" toml:"cpu,omitempty"`
	Host *Host `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
	Load *Load `yaml:"load,omitempty" json:"load,omitempty" toml:"load,omitempty"`
	Mem  *Mem  `yaml:"mem,omitempty" json:"mem,omitempty" toml:"mem,omitempty"`
	Net  *Net  `yaml:"net,omitempty" json:"net,omitempty" toml:"net,omitempty"`
}
