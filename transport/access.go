package transport

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/token/access"
)

// NewAccessController constructs the shared transport access controller.
//
// If transport config or access config is omitted, it returns (nil, nil) so server stacks can treat access
// control as disabled. When configured, the returned controller is shared by the HTTP and gRPC server stacks.
func NewAccessController(cfg *access.Config, fs *os.FS) (access.Controller, error) {
	return access.NewController(cfg, fs)
}
