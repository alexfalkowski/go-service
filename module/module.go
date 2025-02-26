package module

import (
	"github.com/alexfalkowski/go-service/compress"
	"github.com/alexfalkowski/go-service/crypto"
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/sync"
	"github.com/alexfalkowski/go-service/time"
	"go.uber.org/fx"
)

// Module for packages in this library.
var Module = fx.Options(
	os.Module,
	env.Module,
	runtime.Module,
	compress.Module,
	encoding.Module,
	crypto.Module,
	time.Module,
	sync.Module,
	id.Module,
)
