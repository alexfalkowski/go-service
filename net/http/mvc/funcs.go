package mvc

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/go-sprout/sprout"
	"github.com/go-sprout/sprout/registry/conversion"
	"github.com/go-sprout/sprout/registry/maps"
	"github.com/go-sprout/sprout/registry/numeric"
	"github.com/go-sprout/sprout/registry/regexp"
	"github.com/go-sprout/sprout/registry/slices"
	"github.com/go-sprout/sprout/registry/std"
	"github.com/go-sprout/sprout/registry/strings"
	"github.com/go-sprout/sprout/registry/time"
)

// FunctionMapParams defines dependencies used to build a sprout.FunctionMap.
type FunctionMapParams struct {
	di.In
	Logger     *slog.Logger
	Registries []sprout.Registry `optional:"true"`
}

// NewFunctionMap builds a FunctionMap with common sprout registries enabled.
// List of registries can be found at https://docs.atom.codes/sprout/registries/list-of-all-registries
func NewFunctionMap(params FunctionMapParams) sprout.FunctionMap {
	handler := sprout.New(sprout.WithLogger(params.Logger), sprout.WithSafeFuncs(true))

	runtime.Must(handler.AddRegistries(params.Registries...))
	runtime.Must(handler.AddRegistry(conversion.NewRegistry()))
	runtime.Must(handler.AddRegistry(std.NewRegistry()))
	runtime.Must(handler.AddRegistry(maps.NewRegistry()))
	runtime.Must(handler.AddRegistry(numeric.NewRegistry()))
	runtime.Must(handler.AddRegistry(regexp.NewRegistry()))
	runtime.Must(handler.AddRegistry(slices.NewRegistry()))
	runtime.Must(handler.AddRegistry(strings.NewRegistry()))
	runtime.Must(handler.AddRegistry(time.NewRegistry()))

	return handler.Build()
}
