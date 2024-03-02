package runtime

import (
	"github.com/KimMachineGun/automemlimit/memlimit"
)

// RegisterMemLimit for runtime.
func RegisterMemLimit() {
	memlimit.SetGoMemLimitWithEnv()
}
