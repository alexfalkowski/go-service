# AGENTS.md

This repository is a Go module (`github.com/alexfalkowski/go-service/v2`) that provides a framework for building services (DI wiring, config decoding, transport, telemetry, crypto, etc.).

The project is driven by `make` targets that are mostly defined in the `bin/` git submodule.

## First-time setup

### Git submodule (required for `make`)

The top-level `Makefile` includes make fragments from the `bin` submodule (see `Makefile:1-3`, `.gitmodules:1-3`).

```sh
git submodule sync
git submodule update --init
```

Alternative (same effect):

```sh
make submodule
```

Gotcha: `.gitmodules` points at an SSH URL (`git@github.com:alexfalkowski/bin.git`). If you can’t fetch via SSH, `make` will fail until Git access is configured.

### Local dev dependencies observed in this repo

These tools are used by `make` targets and/or CI:

- `mkcert` (used by `make create-certs`; also used in CI).
- `gotestsum` (used by `make specs`).
- `govulncheck` (used by `make sec`).
- `codecovcli` (used by `make codecov-upload`).
- `buf` (used by `make generate` via `internal/test/Makefile`).
- Diagram tooling: `goda` + `dot` (Graphviz) (used by `make diagrams`).

## Project type

- Language: Go
- Go version: `go 1.25.0` (see `go.mod:1-4`)
- DI: Uber Fx/Dig, wrapped in `di/` (see `di/di.go`)
- CLI command framework: `github.com/cristalhq/acmd` (see `cli/application.go`)
- Linting/formatting: `golangci-lint` + formatters (see `.golangci.yml`)

## Essential commands

Most workflows are driven by `make`.

### Discover targets

```sh
make help
```

### Dependencies (keeps `vendor/` in sync)

```sh
make dep
```

`make dep` runs `go mod download`, `go mod tidy`, and `go mod vendor` (see `bin/build/make/go.mak:9-26`).

Gotcha: tests run with `-mod vendor`, so after changing dependencies you typically must run `make dep`.

### Tests

```sh
make specs
```

`make specs` runs:

- `gotestsum --junitfile test/reports/specs.xml -- ...`
- `-race`
- `-mod vendor`
- `-covermode=atomic`
- `-coverpkg=<all repo packages>`
- `-coverprofile=test/reports/profile.cov`

(see `bin/build/make/go.mak:61-64`).

### Lint / format

```sh
make lint
make fix-lint
make format
```

- `make lint` runs field-alignment + `golangci-lint` (see `bin/build/make/go.mak:39-56`).
- Formatters enabled via golangci config: `gci`, `gofmt`, `gofumpt`, `goimports` (see `.golangci.yml:44-50`).

### Security checks

```sh
make sec
```

Runs `govulncheck -show verbose -test ./...` (see `bin/build/make/go.mak:95-98`).

### Coverage

```sh
make coverage
make html-coverage
make func-coverage
```

Coverage artifacts live under `test/reports/` (see `bin/build/make/go.mak:73-86`).

### Benchmarks

Convenience targets in the top-level `Makefile`:

```sh
make benchmarks
make http-benchmarks
make grpc-benchmarks
make bytes-benchmarks
make strings-benchmarks
```

Underlying target:

```sh
make package=<pkg> benchmark
```

(see `Makefile:20-34` and `bin/build/make/go.mak:65-71`).

### Local environment (integration deps)

```sh
make start
make stop
```

The `start/stop` targets shell out to `bin/build/docker/env` (see `bin/build/make/go.mak:130-136`).

CI provisions and waits for these services (see `.circleci/config.yml:9-30`):

- Postgres (`tcp://localhost:5432`)
- Valkey/Redis (`tcp://localhost:6379`)
- `alexfalkowski/status` server (`tcp://localhost:6000`)
- Grafana Mimir (`tcp://localhost:9009`)

### TLS fixtures / certs

```sh
mkcert -install
make create-certs
```

Generates fixtures into `test/certs/` (see `bin/build/make/go.mak:113-117`).

### Code generation (Buf)

Top-level target:

```sh
make generate
```

Delegates to `make -C internal/test generate` (see `Makefile:35-37`), which includes `../../bin/build/make/buf.mak` (see `internal/test/Makefile:1`).

Buf config:

- `internal/test/buf.yaml`
- `internal/test/buf.gen.yaml`

### Diagrams

```sh
make diagrams
make crypto-diagram
make database-diagram
make telemetry-diagram
make transport-diagram
```

These call `make package=<pkg> create-diagram` and write PNGs into `assets/` (see `Makefile:5-18`, `bin/build/make/go.mak:118-121`).

## Code organization (high level)

The code is structured as a library of packages (no single `cmd/` binary in this repo).

Common patterns:

- Most features are packages with `config.go`, `module.go`, and implementation files.
- `module/` exports top-level Fx modules (see `module/module.go`).
  - `module.Library`
  - `module.Server`
  - `module.Client`
- `di/` provides thin aliases/wrappers around `fx` / `dig` (`di.In`, `di.Module`, `di.Constructor`, `di.Register`, etc.) (see `di/di.go`).
- `cli/` wires commands via `acmd` (see `cli/application.go`).
- `config/` handles decoding/validation and input routing (see `config/decoder.go`, `config/default.go`, `config/env.go`, `config/file.go`).
- `transport/` provides HTTP + gRPC transport packages (`transport/http`, `transport/grpc`) plus common transport config.
- `telemetry/` provides logger/metrics/tracer packages and modules.
- `internal/test/` contains shared test helpers, fixtures, templates, and protobuf generation config.
- `test/` contains fixtures used by tests (configs, certs, secrets, reports).

## Key patterns and conventions (observed)

### Fx module composition

Modules are typically defined as `di.Module(...)` values that compose submodules plus constructors/invocations.

Example: `transport.Module` composes HTTP + gRPC wiring (see `transport/module.go`).

### Dependency injection parameter structs

Packages frequently use `di.In` structs to declare injected dependencies (example: `config.DecoderParams` in `config/decoder.go:12-19`).

### CLI pattern (server vs client)

`cli.Application` registers “server” and “client” subcommands:

- `AddServer(...)` starts an Fx app and blocks until `app.Done()` then stops (see `cli/application.go:61-88`).
- `AddClient(...)` starts and immediately stops (see `cli/application.go:90-114`).

Errors are prefixed with the command name and use `dig.RootCause(...)` via `di.RootCause` (see `cli/application.go:137-139`).

### Config input routing via `-i`

`config.NewDecoder` dispatches based on the `-i` flag value (see `config/decoder.go:21-32`):

- `file:<path>` → file decoder
- `env:<ENV_VAR>` → env decoder
- otherwise → default lookup decoder

Default lookup searches for `<serviceName>.{yaml,yml,toml,json}` in (see `config/default.go:35-58`):

- executable directory
- `$XDG_CONFIG_HOME/<serviceName>/` (via `os.UserConfigDir()`)
- `/etc/<serviceName>/`

### Optional logger

Logging is based on `log/slog` and the project’s wrapper (`telemetry/logger`).

- `telemetry/logger.NewLogger` returns `(*Logger, nil)` when enabled, otherwise `(nil, nil)` (see `telemetry/logger/logger.go:65-78`).
- Callers frequently treat `*logger.Logger` as optional and guard `nil` (example: `server/service.go:53-59`).

## Testing

- Tests use `stretchr/testify/require` in many places (search in `*_test.go`).
- Shared helpers and fixtures live under `internal/test/` and `test/`.
- Test reports and coverage artifacts are written to `test/reports/` by `make specs` and `make coverage`.

## Style / formatting

From `.editorconfig`:

- Go files use tabs (`indent_style = tab`, `indent_size = 4`).
- YAML uses 2-space indentation.

Golangci-lint is configured in `.golangci.yml`.

## CI notes / gotchas

- Many `make` targets come from the `bin/` submodule; if it’s missing, `make` will fail due to missing includes.
- CI sets `GOEXPERIMENT=greenteagc` (see `.circleci/config.yml:6-9`).
- `make encode-config` uses `base64 -w 0` (see `bin/build/make/go.mak:109-112`), which may not work on BSD/macOS `base64`.
- CI generates a `.source-key` file via `make source-key` for caching (see `.circleci/config.yml:27`, `bin/build/make/git.mak:175-177`). This file is ignored by Git (see `.gitignore:11`).
