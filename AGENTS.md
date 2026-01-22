# AGENTS.md

This repo is a Go module (`github.com/alexfalkowski/go-service/v2`) that provides a framework for building services (DI, config, transport, telemetry, crypto, etc.).

## Repository bootstrap

### Submodule (required for `make`)
The top-level `Makefile` includes make fragments from the `bin` git submodule (see `.gitmodules` and `Makefile`).

```sh
git submodule sync
git submodule update --init
```

If `make` fails with missing `bin/build/make/*.mak`, the submodule is not initialized.

### Local dev deps (documented)
`README.md` documents installing `mkcert` and creating cert fixtures:

```sh
mkcert -install
make create-certs
```

## Project type

- Language: Go
- Go version: `go 1.25.0` (see `go.mod`)
- DI: Uber Fx/Dig via wrappers in `di/` (see `di/di.go`)
- CLI: `github.com/cristalhq/acmd` (see `cli/application.go`)

## Essential commands

Most workflows are driven by `make` (targets come from `bin/build/make/go.mak`, plus a few top-level targets in `Makefile`).

### Help / discoverability

```sh
make help
```

(`bin/build/make/help.mak` renders target comments.)

### Dependencies / cleanup

```sh
make dep        # go mod download + tidy + vendor
make clean      # runs bin/build/go/clean
make clean-dep  # clears go caches
make clean-lint # clears golangci-lint cache
```

### Lint

```sh
make lint       # field-alignment + golangci-lint
make fix-lint   # attempts to auto-fix (field alignment + golangci --fix)
```

Lint configuration is in `.golangci.yml`.

### Security checks

```sh
make sec        # govulncheck -show verbose -test ./...
```

### Tests (specs) and coverage

```sh
make specs
```

`make specs` uses `gotestsum` and runs tests with:
- `-race`
- `-mod vendor`
- `-covermode=atomic`
- `-coverpkg=<all repo packages>`
- junit output: `test/reports/specs.xml`
- coverage profile: `test/reports/profile.cov`

Coverage post-processing and reporting:

```sh
make coverage       # generates HTML + func coverage (uses covfilter)
make html-coverage
make func-coverage
make codecov-upload # codecovcli upload-process from test/reports/final.cov
```

Additional coverage-related config:
- `.gocov` (patterns used by coverage tooling)
- `.codecov.yml`

### Benchmarks
Top-level `Makefile` defines convenience targets:

```sh
make benchmarks
make http-benchmarks
make grpc-benchmarks
```

Underlying per-package benchmark target is defined in `bin/build/make/go.mak`:

```sh
make package=transport/http benchmark
make package=transport/grpc benchmark
```

### Diagrams
Top-level `Makefile` defines:

```sh
make diagrams
make crypto-diagram
make database-diagram
make telemetry-diagram
make transport-diagram
```

These call `make package=<pkg> create-diagram` (defined in `bin/build/make/go.mak`) which runs `goda graph ... | dot -Tpng` and writes into `assets/`.

### Environment (integration deps)
The make targets exist (defined in `bin/build/make/go.mak`):

```sh
make start
make stop
```

CircleCI starts external services via Docker images (see `.circleci/config.yml`): Postgres, Valkey/Redis, Grafana Mimir, and an `alexfalkowski/status` server.

### Misc helpers

```sh
make encode-config kind=configs/config  # reads test/<kind>.yml and base64 encodes it
make create-certs                       # writes certs into test/certs/
make source-key                         # writes .source-key (used by CI cache keys)
```

## Code organization (high level)

- `module/` – exported Fx modules for consumers (see `module/module.go`).
  - `module.Library` (baseline functionality)
  - `module.Server` (server-side wiring)
  - `module.Client` (client-side wiring)
- `di/` – thin aliases around `fx`/`dig` (see `di/di.go`).
- `cli/` – application/command wiring (see `cli/application.go`).
- `config/` – config decoding, default lookup, validation (see `config/decoder.go`, `config/config.go`).
- `transport/` – HTTP + gRPC transports and middleware (`transport/http`, `transport/grpc`).
- `telemetry/` – logger/metrics/tracer wrappers and modules (`telemetry/logger`, `telemetry/metrics`, `telemetry/tracer`).
- `internal/test/` – shared test helpers and generators (e.g. `internal/test/config.go`).
- `test/` – fixtures (configs, certs, secrets, reports).

## Key patterns and conventions (observed)

### Fx module pattern
Modules are defined as `di.Module(...)` values that compose sub-modules and constructors/invocations.
Example: `transport/Module` composes HTTP + gRPC and registers constructors (see `transport/module.go`).

### Config input routing via `-i`
`config.NewDecoder` dispatches based on the `-i` flag value (see `config/decoder.go`):
- `file:<path>` → file decoder
- `env:<ENV_VAR>` → env decoder
- otherwise → default lookup decoder

README documents env config as `extension:base64(content)` (e.g. `yml:<base64>`).

### Logging
Logging is built on `log/slog` with a project wrapper (`telemetry/logger/logger.go`). Some components accept `*logger.Logger` and no-op safely when it’s nil (example: `server/service.go`).

### Tests
- Assertions commonly use `stretchr/testify/require` (example: `token/jwt/jwt_test.go`).
- Many tests use fixtures from `test/` via helpers in `internal/test/`.

## Style / formatting

- `.editorconfig` indicates:
  - Go files use tabs (`indent_style = tab`, `indent_size = 4`).
  - YAML uses 2-space indentation.
- The repo favors the Uber Go style guide (documented in `README.md`).
- `.golangci.yml` enables many linters and formatters (gofmt/gofumpt/goimports/gci).

## Gotchas for automated agents

- **Make targets are the source of truth**: CI and local workflows are designed around `make` (e.g. `make dep`, `make lint`, `make specs`, `make sec`). Prefer `make` targets over ad-hoc commands for parity with CI.
- **Always init submodules before using `make`**: most targets come from the `bin` submodule.
- **Integration dependencies**: some tests depend on external services; CI starts Postgres and Valkey/Redis (see `.circleci/config.yml`). Use `make start`/`make stop` where appropriate.
- **Cert fixtures**: TLS-related tests/configs rely on cert/key files in `test/certs/`; README documents generating them with `mkcert` + `make create-certs`.
