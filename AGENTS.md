# AGENTS.md

## Repo snapshot

- Module: `github.com/alexfalkowski/go-service/v2`
- Repo type: library of packages; there is no top-level `cmd/`
- Go version: `go 1.26.0`
- DI: Uber Fx/Dig via `di/`
- CLI helpers: `github.com/cristalhq/acmd` via `cli/`
- Most `make` targets come from the `bin/` git submodule

## Setup

- Initialize the submodule before using `make`:

```sh
git submodule sync
git submodule update --init
```

- Equivalent: `make submodule`
- `.gitmodules` uses `git@github.com:alexfalkowski/bin.git`; SSH access is required for submodule fetches.

## Common commands

- `make help`: list targets.
- `make dep`: runs `go mod download`, `go mod tidy`, and `go mod vendor`. Tests use `-mod vendor`, so run this after dependency changes.
- `make specs`, `make lint`, `make fix-lint`, `make format`, `make sec`
- `make benchmarks` and focused variants such as `make http-benchmarks`, `make grpc-benchmarks`, `make bytes-benchmarks`, `make strings-benchmarks`
- `make coverage`, `make html-coverage`, `make func-coverage`
- `make generate`, `make diagrams`, `make start`, `make stop`
- `mkcert -install && make create-certs`
- `make kind=status encode-config`
  `encode-config` uses GNU `base64 -w 0`; on macOS/BSD use `base64 | tr -d '\n'`.

## Layout and wiring

- Feature packages usually follow `config.go`, `module.go`, and implementation files.
- `module/` exports the top-level Fx bundles: `module.Library`, `module.Server`, `module.Client`.
- `config/` defines the top-level config and projections into nested transport, SQL, and telemetry config.
- `net/` contains lower-level HTTP/gRPC, metadata, header, and server helpers.
- `transport/` contains the higher-level transport layer: composed HTTP/gRPC stacks, middleware, and operational endpoints.
- `internal/test/` provides shared test helpers, especially `internal/test/world.go`.
- `test/` stores fixtures such as configs, certs, secrets, and reports.
- Modules are composed with `di.Module(...)`; many constructors consume `di.In` parameter structs.
- `module.Server` includes debug, cache, config, feature, SQL, telemetry, limiter, transport, and health wiring.
- `module.Client` includes cache, config, feature, hooks, SQL, telemetry, and limiter wiring, but not debug, transport, or health by default.

## Configuration rules

- `config.NewDecoder` dispatches `-i` values as `file:<path>`, `env:<ENV_VAR>`, or default lookup for `<serviceName>.{yaml,yml,hjson,toml,json}`.
- Default lookup checks the executable directory, `$XDG_CONFIG_HOME/<serviceName>/`, and `/etc/<serviceName>/`.
- Many fields use go-service source strings resolved by `os.FS.ReadSource`: `env:NAME`, `file:/path`, or a literal value.
- Nil pointer sub-configs usually mean "disabled".

## Gotchas

- Manual transport TLS wiring needs `transport/http.Register(fs)` and `transport/grpc.Register(fs)`; `transport.Module` handles this in the normal path.
- Manual server lifecycle wiring should use `net/server.Register(...)`.
- `telemetry.Register()` installs the global OpenTelemetry propagator.
- `cache.Register(...)` sets the package-level cache used by helpers in `cache/generic.go`.
- `token/access` passes `access.policy` directly to Casbin's file adapter; it needs a real path.
- JWT verification requires both the expected algorithm and a `kid` header.
- `telemetry/header.Map.MustSecrets` can panic during config projection if secret resolution fails.
- `transport/http/health.RegisterParams` and `internal/test.RegisterHealth` both require an `*net/http.ServeMux`.
- Shared metadata, header, and string helpers live under `net/...`, not `transport/...`.
- `vendor/` is gitignored and regenerated via `make dep`.

## Testing, style, and docs

- Tests commonly use `stretchr/testify/require`.
- Go files use tabs; YAML uses 2 spaces (`.editorconfig`).
- All exported identifiers should have GoDoc comments.
- GoDoc comments should start with the identifier name (or `Deprecated:`).
- This standard applies to `internal/test/**` too.

## CI

- CircleCI runs submodule init, `make source-key`, `mkcert -install`, `make create-certs`, waits for services, then runs `make clean`, `make dep`, `make lint`, `make sec`, `make specs`, `make benchmarks`, and `make coverage`.
- CI services:
  - Postgres: `localhost:5432`
  - Valkey: `localhost:6379`
  - `alexfalkowski/status`: `localhost:6000`
  - Grafana Mimir: `localhost:9009`
