package module

import (
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/crypto"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	"github.com/alexfalkowski/go-service/v2/types"
	"go.uber.org/fx"
)

var (
	// Module for packages in this library.
	Module = fx.Options(
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

	// Server module.
	Server = fx.Options(
		Module,
		debug.Module,
		cache.Module,
		cli.Module,
		config.Module,
		feature.Module,
		sql.Module,
		telemetry.Module,
		token.Module,
		transport.Module,
		health.Module,
	)

	// Client module.
	Client = fx.Options(
		Module,
		cache.Module,
		cli.Module,
		config.Module,
		feature.Module,
		hooks.Module,
		sql.Module,
		telemetry.Module,
		token.Module,
	)
)
