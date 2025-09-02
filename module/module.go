package module

import (
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/crypto"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
)

var (
	// Library module.
	Library = di.Module(
		env.Module,
		compress.Module,
		encoding.Module,
		crypto.Module,
		time.Module,
		sync.Module,
		id.Module,
	)

	// Server module.
	Server = di.Module(
		Library,
		debug.Module,
		cache.Module,
		config.Module,
		feature.Module,
		sql.Module,
		telemetry.Module,
		limiter.Module,
		transport.Module,
		health.Module,
	)

	// Client module.
	Client = di.Module(
		Library,
		cache.Module,
		config.Module,
		feature.Module,
		hooks.Module,
		sql.Module,
		telemetry.Module,
		limiter.Module,
	)
)
