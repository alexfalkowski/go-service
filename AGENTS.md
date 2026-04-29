# AGENTS.md

## Shared skill

Use the shared `coding-standards` skill from `bin/skills/coding-standards` for code changes, bug fixes, refactors, reviews, tests, linting, documentation, PR summaries, commits, Makefile changes, CI validation, and verification.

## Repo snapshot

- Module: `github.com/alexfalkowski/go-service/v2`
- Library repo; no top-level `cmd/`
- Go: `1.26.0`
- DI: Uber Fx/Dig via `di/`
- CLI helpers: `github.com/cristalhq/acmd` via `cli/`
- Many `make` targets come from the `bin/` submodule

## Setup

- Initialize the submodule before using `make`:

```sh
git submodule sync
git submodule update --init
```

- Equivalent: `make submodule`
- Submodule fetches use SSH via `git@github.com:alexfalkowski/bin.git`

## Common commands

- `make help`
- `make dep`: run after dependency changes; it does `go mod download`, `go mod tidy`, and `go mod vendor`. Tests use `-mod vendor`.
- Quality: `make specs`, `make lint`, `make fix-lint`, `make format`, `make sec`
- Coverage: `make coverage`, `make html-coverage`, `make func-coverage`
- Benchmarks: `make benchmarks`, `make http-benchmarks`, `make grpc-benchmarks`, `make bytes-benchmarks`, `make strings-benchmarks`
- Other: `make generate`, `make diagrams`, `make start`, `make stop`
- TLS setup: `mkcert -install && make create-certs`
- `make kind=status encode-config`
- `encode-config` expects GNU `base64 -w 0`; on macOS/BSD use `base64 | tr -d '\n'`

## Layout and wiring

- Feature packages usually use `config.go`, `module.go`, plus implementation files.
- `module/` exports the main Fx bundles: `module.Library`, `module.Server`, `module.Client`.
- `config/` owns top-level config plus projections into transport, SQL, and telemetry config.
- `net/` holds lower-level HTTP/gRPC, metadata, header, and server helpers.
- `transport/` holds the higher-level HTTP/gRPC stacks, middleware, and ops endpoints.
- `internal/test/` contains shared test helpers, especially `internal/test/world.go`.
- `test/` stores fixtures such as configs, certs, secrets, and reports.
- Modules are composed with `di.Module(...)`; many constructors use `di.In`.
- `module.Server` wires debug, cache, config, feature, SQL, telemetry, limiter, transport, and health.
- `module.Client` wires cache, config, feature, hooks, SQL, telemetry, and limiter, but not debug, transport, or health by default.

## Configuration

- `config.NewDecoder` resolves `-i` as `file:<path>`, `env:<ENV_VAR>`, or `<serviceName>.{yaml,yml,hjson,toml,json}`.
- Default file lookup checks the executable directory, `$XDG_CONFIG_HOME/<serviceName>/`, and `/etc/<serviceName>/`.
- Many config fields use go-service source strings through `os.FS.ReadSource`: `env:NAME`, `file:/path`, or a literal value.
- Nil pointer sub-configs usually mean "disabled".

## Gotchas

- Manual transport TLS setup must call `transport/http.Register(fs)` and `transport/grpc.Register(fs)`; `transport.Module` does this on the normal path.
- Manual server lifecycle wiring should use `net/server.Register(...)`.
- `telemetry.Register()` installs the global OpenTelemetry propagator.
- `cache.Register(...)` sets the package-level cache used by `cache/generic.go`.
- `token/access` passes `access.policy` directly to Casbin's file adapter, so it must be a real path.
- JWT verification requires both the expected algorithm and a `kid` header.
- `telemetry/header.Map.MustSecrets` can panic if secret resolution fails during config projection.
- `transport/http/health.RegisterParams` and `internal/test.RegisterHealth` both require `*net/http.ServeMux`.
- Shared metadata, header, and string helpers live under `net/...`, not `transport/...`.
- `vendor/` is gitignored and regenerated via `make dep`.

## Testing, style, and docs

- Tests commonly use `stretchr/testify/require`.
- Go files use tabs; YAML uses 2 spaces per `.editorconfig`.
- Every exported identifier, including under `internal/test/**`, needs a GoDoc comment.
- GoDoc comments should start with the identifier name or `Deprecated:`.

## CI

- CircleCI runs submodule init, `make source-key`, `mkcert -install`, `make create-certs`, waits for services, then runs `make clean`, `make dep`, `make lint`, `make sec`, `make specs`, `make benchmarks`, and `make coverage`.
- Services: Postgres `localhost:5432`, Valkey `localhost:6379`, `alexfalkowski/status` `localhost:6000`, Grafana Mimir `localhost:9009`
