package time

import (
	"github.com/alexfalkowski/go-service/time/ntp"
	"github.com/alexfalkowski/go-service/time/nts"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	ntp.Module,
	nts.Module,
)
