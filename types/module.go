package types

import (
	"github.com/alexfalkowski/go-service/types/fields"
	"github.com/alexfalkowski/go-service/types/structs"
	"github.com/alexfalkowski/go-service/types/validator"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(validator.NewValidator),
	fx.Invoke(fields.Register),
	fx.Invoke(structs.Register),
)
