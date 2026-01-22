# AGENTS.md

This repo is a Go module (`github.com/alexfalkowski/go-service/v2`) that provides a framework for building services (DI, config, transport, telemetry, crypto, etc.).

## First-time setup

### Git submodule (required for `make`)

The top-level `Makefile` includes make fragments from the `bin` git submodule (see `.gitmodules`, `Makefile`).

```sh
git submodule sync
git submodule update --init
```

Alternative (same effect, via Make):

```sh
make submodule
```

Gotcha: `.gitmodules` points at an SSH URL (`git@github.com:alexfalkowski/bin.git`). If your environment cannot fetch SSH submodules, `make` will fail until you adjust your Git access.

### Local dev deps (observed)

Some `make` targets assume these tools exist:

- `mkcert` (used by `make create-certs`; also used in CI).
- `gotestsum` (used by `make specs`).
- `govulncheck` (used by `make sec`).
- `codecovcli` (used by `make codecov-upload`).
- `buf` (used by `make generate` via `internal/test`).
- Diagram tooling (used by `make diagrams`): `goda` + `dot` (Graphviz).

## Project type

- Language: Go
- Go version: `go 1.25.0` (see `go.mod`)
- DI: Uber Fx/Dig via wrappers in `di/` (see `di/di.go`)
- CLI command framework: `github.com/cristalhq/acmd` (see `cli/application.go`)
- Linting: `golangci-lint` invoked via `bin/build/go/lint` (see `bin/build/make/go.mak` and `.golangci.yml`)

## Essential commands

Most workflows are driven by `make`. Targets come primarily from `bin/build/make/*.mak`, plus a few top-level targets in `Makefile`.

### Discoverability

```sh
make help
```

### Dependencies

```sh
make dep        # go mod download + tidy + vendor
```

Gotcha: tests/specs are executed with `-mod vendor` (see `bin/build/make/go.mak:61-63`), so `make dep` is typically required after dependency changes.

### Tests

```sh
make specs
```

`make specs` runs with:

- `gotestsum --junitfile test/reports/specs.xml`
- `-race`
- `-mod vendor`
- `-covermode=atomic`
- `-coverpkg=<all repo packages>`
- `-coverprofile=test/reports/profile.cov`

### Lint / format

```sh
make lint       # field-alignment + golangci-lint
make fix-lint   # attempts to auto-fix (field alignment + golangci --fix)
make format     # go fmt ./...
```

Lint configuration is in `.golangci.yml`.

### Security checks

```sh
make sec        # govulncheck -show verbose -test ./...
```

### Coverage

```sh
make coverage       # generates HTML + func coverage (uses covfilter)
make html-coverage
make func-coverage
make codecov-upload # codecovcli upload-process from test/reports/final.cov
```

Coverage-related config and artifacts:

- `.gocov` (coverage tool patterns)
- `.codecov.yml`
- `test/reports/*` (junit xml, cov profiles, HTML)

### Cleanup

```sh
make clean       # cleans build artifacts (via bin/build/go/clean)
make clean-dep   # clears go caches
make clean-lint  # clears golangci-lint cache
make clean-reports
```

### Benchmarks

Top-level `Makefile` defines convenience targets:

```sh
make benchmarks
make http-benchmarks
make grpc-benchmarks
```

Underlying per-package benchmark target (from `bin/build/make/go.mak:65-68`):

```sh
make package=transport/http benchmark
make package=transport/grpc benchmark
```

### Local environment (integration deps)

Targets exist (from `bin/build/make/go.mak:130-136`):

```sh
make start
make stop
```

CI provisions Docker services (see `.circleci/config.yml:4-74`):

- Postgres (waits on `tcp://localhost:5432`)
- Valkey/Redis (waits on `tcp://localhost:6379`)
- `alexfalkowski/status` server (waits on `tcp://localhost:6000`)
- Grafana Mimir (waits on `tcp://localhost:9009`)

### Test fixtures / certs

```sh
mkcert -install
make create-certs
```

TLS fixtures are written to `test/certs/`.

### Code generation (Buf)

Top-level `Makefile` provides:

```sh
make generate
```

This delegates to `make -C internal/test generate` (see `Makefile:30-31`), and that includes the Buf make fragment (`internal/test/Makefile:1`, `bin/build/make/buf.mak:19-21`).

### Diagrams

Top-level `Makefile`:

```sh
make diagrams
make crypto-diagram
make database-diagram
make telemetry-diagram
make transport-diagram
```

These call `make package=<pkg> create-diagram` and write PNGs into `assets/`.

### Git helper targets (optional)

These come from `bin/build/make/git.mak`:

```sh
make latest             # checkout master + pull --rebase + submodule sync/update
make new-feature name=… # creates a user-prefixed branch
make new-fix name=…
make done               # delete current branch after syncing
```

## Code organization (high level)

- `module/` – exported Fx modules for consumers (see `module/module.go`).
  - `module.Library` (baseline functionality)
  - `module.Server` (server-side wiring)
  - `module.Client` (client-side wiring)
- `di/` – thin aliases around `fx`/`dig` (see `di/di.go`).
- `cli/` – application/command wiring (see `cli/application.go`).
- `config/` – decoding, default lookup, validation (see `config/decoder.go`, `config/env.go`, `config/default.go`, `config/validator.go`).
- `transport/` – HTTP + gRPC transports and middleware (`transport/http`, `transport/grpc`).
- `telemetry/` – logger/metrics/tracer wrappers and modules (`telemetry/logger`, `telemetry/metrics`, `telemetry/tracer`).
- `crypto/`, `database/`, `cache/`, `health/`, `feature/`, `hooks/`, `id/`, etc. – pluggable modules.
- `internal/test/` – shared test helpers and generators.
- `test/` – fixtures (configs, certs, secrets, reports).

## Key patterns and conventions (observed)

### Fx module pattern

Modules are defined as `di.Module(...)` values that compose sub-modules and constructors/invocations.

Example: `transport.Module` composes HTTP + gRPC and registers constructors/invocations (see `transport/module.go:10-17`).

### CLI pattern (server vs client)

`cli.Application` registers “server” and “client” commands using `acmd`.

- `AddServer(...)` starts an Fx app, waits on `app.Done()`, then stops on completion (see `cli/application.go:61-88`).
- `AddClient(...)` starts and immediately stops (see `cli/application.go:90-114`).

Errors are wrapped with the command name prefix and `dig.RootCause(...)` (see `cli/application.go:137-139`).

### Config input routing via `-i`

`config.NewDecoder` dispatches based on the `-i` flag value (see `config/decoder.go:21-32`):

- `file:<path>` → file decoder (`config/file.go:21-34`)
- `env:<ENV_VAR>` → env decoder (`config/env.go:12-40`)
- otherwise → default lookup decoder (`config/default.go:24-58`)

Default lookup searches for `<serviceName>.{yaml,yml,toml,json}` in:

- executable directory
- `$XDG_CONFIG_HOME/<serviceName>/` (via `os.UserConfigDir()`)
- `/etc/<serviceName>/`

(see `config/default.go:35-58`).

### Logging

Logging is built on `log/slog` with a project wrapper.

- `telemetry/logger.NewLogger` returns `(*Logger, nil)` when enabled, otherwise `(nil, nil)` (see `telemetry/logger/logger.go:65-78`).
- Callers should treat `*logger.Logger` as optional; many components guard `nil` explicitly (example: `server/service.go:53-59`).

### Tests

- Assertions commonly use `stretchr/testify/require`.
- Many tests use fixtures from `test/` via helpers in `internal/test/`.

## Style / formatting

From `.editorconfig`:

- Go files use tabs (`indent_style = tab`, `indent_size = 4`).
- YAML uses 2-space indentation.

`.golangci.yml` enables formatters such as `gci`, `gofmt`, `gofumpt`, and `goimports`.

## Gotchas

- **Make targets are the source of truth**: CI and local workflows are designed around `make` (see `.circleci/config.yml` and `bin/build/make/go.mak`). Prefer `make` targets over ad-hoc commands for parity.
- **Vendor mode in tests**: `make specs` uses `-mod vendor`; keep `vendor/` up to date via `make dep`.
- **Submodule required**: most targets come from the `bin` submodule; missing submodule manifests as missing `bin/build/make/*.mak` includes.
- **`encode-config` portability**: `make encode-config` uses `base64 -w 0` (`bin/build/make/go.mak:109-112`), which may not work with BSD `base64`.
- **Integration deps**: some tests may assume external services; CI provisions Postgres, Valkey/Redis, a status server, and Mimir (see `.circleci/config.yml`).
- **Certificates for TLS fixtures**: TLS-related tests/configs rely on files in `test/certs/`; CI runs `mkcert -install` + `make create-certs`.
