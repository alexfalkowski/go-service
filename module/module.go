package module

import (
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/crypto"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/types"
	"go.uber.org/fx"
)

// Module for packages in this library.
var Module = fx.Options(
	os.Module,
	env.Module,
	compress.Module,
	encoding.Module,
	crypto.Module,
	time.Module,
	sync.Module,
	id.Module,
	types.Module,
)
