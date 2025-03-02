package runtime

import (
	"log/slog"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// RegisterMemLimit for runtime.
func RegisterMemLimit(logger *slog.Logger) error {
	_, err := memlimit.SetGoMemLimitWithOpts(memlimit.WithLogger(logger))

	return err
}
