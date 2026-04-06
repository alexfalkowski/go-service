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
- `make dep`: runs `go mod download`, `go mod tidy`, and `go mod vendor`.
  Tests use `-mod vendor`, so run this after dependency changes.
- `make specs`: race + coverage via `gotestsum`.
  Outputs `test/reports/specs.xml` and `test/reports/profile.cov`.
- `make lint`, `make fix-lint`, `make format`
- `make sec`
- `make benchmarks`, `make http-benchmarks`, `make grpc-benchmarks`, `make bytes-benchmarks`, `make strings-benchmarks`
- `make coverage`, `make html-coverage`, `make func-coverage`
- `make generate`
- `make diagrams`, `make crypto-diagram`, `make database-diagram`, `make telemetry-diagram`, `make transport-diagram`
- `make start`, `make stop`
- `mkcert -install && make create-certs`
- `make kind=status encode-config`
  `encode-config` uses GNU `base64 -w 0`; on macOS/BSD use `base64 | tr -d '\n'`.

## Repo layout

- Feature packages usually follow `config.go`, `module.go`, and implementation files.
- `module/` exports the top-level Fx bundles: `module.Library`, `module.Server`, `module.Client`.
- `config/` defines the standard top-level config plus projections into nested transport/SQL/telemetry config.
- `net/` contains lower-level protocol helpers: `net/http`, `net/grpc`, metadata helpers, `net/header`, `net/server`, and gRPC health.
- `transport/` contains the higher-level service transport layer: composed HTTP/gRPC stacks, middleware policy, operational endpoints, and transport modules.
- `internal/test/` provides shared test helpers, especially `internal/test/world.go`.
- `test/` stores fixtures (configs, certs, secrets, reports).

## Wiring patterns

- Modules are composed with `di.Module(...)`.
- Many constructors consume `di.In` parameter structs.
- `module.Server` already includes config, telemetry, transports, health, debug, cache, feature, limiter, and SQL.
- `module.Client` includes config, telemetry, feature, hooks, cache, limiter, and SQL, but not debug/transport/health by default.

## Configuration rules

- `config.NewDecoder` dispatches `-i` values as:
  - `file:<path>`
  - `env:<ENV_VAR>`
  - otherwise default lookup for `<serviceName>.{yaml,yml,hjson,toml,json}`
- Default lookup checks:
  - executable directory
  - `$XDG_CONFIG_HOME/<serviceName>/`
  - `/etc/<serviceName>/`
- Many fields use go-service “source strings” resolved by `os.FS.ReadSource`:
  - `env:NAME`
  - `file:/path`
  - otherwise literal value
- Nil pointer sub-configs usually mean “disabled”.

## Gotchas

- Transport TLS requires `transport/http.Register(fs)` and `transport/grpc.Register(fs)` if you wire transports manually; `transport.Module` handles normal Fx wiring.
- Manual server lifecycle wiring should use `net/server.Register(...)`.
- `telemetry.Register()` installs the global OpenTelemetry propagator.
- `cache.Register(...)` sets the package-level cache used by generic helpers in `cache/generic.go`.
- `token/access` passes `access.policy` directly to Casbin’s file adapter; it needs a real path.
- JWT verification requires both the expected algorithm and a `kid` header.
- `telemetry/header` uses `header.Map.MustSecrets`; secret-resolution failures can panic during config projection.
- `transport/http/health.RegisterParams` requires an `*net/http.ServeMux`; `internal/test.RegisterHealth` does too.
- Shared metadata/header/string helpers live under `net/...`, not `transport/...`.
- `vendor/` is ignored by git and regenerated via `make dep`.

## Testing, style, and docs

- Tests commonly use `stretchr/testify/require`.
- Go files use tabs; YAML uses 2 spaces (`.editorconfig`).
- All exported identifiers should have GoDoc comments.
- GoDoc comments should start with the identifier name (or `Deprecated:`).
- This standard applies to `internal/test/**` too.

## CI

- CircleCI runs: submodule init, `make source-key`, `mkcert -install`, `make create-certs`, waits for services, then `make clean`, `make dep`, `make lint`, `make sec`, `make specs`, `make benchmarks`, `make coverage`.
- CI services:
  - Postgres: `localhost:5432`
  - Valkey: `localhost:6379`
  - `alexfalkowski/status`: `localhost:6000`
  - Grafana Mimir: `localhost:9009`
