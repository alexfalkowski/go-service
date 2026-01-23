# AGENTS.md

This repository is a Go module (`github.com/alexfalkowski/go-service/v2`) that provides a framework for building services (DI wiring, config decoding, transport, telemetry, crypto, etc.).

It is primarily a **library of packages** (there is no `cmd/` binary in this repo).

Most workflows are driven by `make` targets that are mostly defined in the `bin/` git submodule.

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

## Project type

- Language: Go
- Go version: `go 1.25.0` (see `go.mod:1-4`)
- DI container: Uber Fx/Dig, wrapped in `di/` (see `di/di.go`)
- CLI command framework: `github.com/cristalhq/acmd` (see `cli/application.go`)
- Linting/formatting: `golangci-lint` with multiple formatters enabled (see `.golangci.yml`)

## Essential commands

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

`make specs` runs `gotestsum` and executes `go test` with (see `bin/build/make/go.mak:61-64`):

- `-vet=off`
- `-race`
- `-mod vendor`
- `-covermode=atomic`
- `-coverpkg=<all repo packages>`
- `-coverprofile=test/reports/profile.cov`

It computes the package list from tracked Go sources (see `bin/build/make/go.mak:5-7`).

Artifacts:

- JUnit XML: `test/reports/specs.xml`
- Coverage profile: `test/reports/profile.cov`

### Lint / format

```sh
make lint
make fix-lint
make format
```

- `make lint` runs field alignment + `golangci-lint` (see `bin/build/make/go.mak:39-56`).
- `make fix-lint` runs the same tools with fixes enabled where supported.
- `make format` runs `go fmt ./...` (see `bin/build/make/go.mak:58-60`).

Golangci formatters enabled (see `.golangci.yml:44-49`): `gci`, `gofmt`, `gofumpt`, `goimports`.

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

These shell out to `bin/build/docker/env` (see `bin/build/make/go.mak:130-136`).

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

Generates fixtures into `test/certs/` (see `bin/build/make/go.mak:113-117`, `.circleci/config.yml:28-29`).

### Code generation (Buf)

Top-level target:

```sh
make generate
```

Delegates to `make -C internal/test generate` (see `Makefile:35-37`).

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

This repo is organized as packages under the module root.

Common conventions:

- Most features are packages with `config.go`, `module.go`, and implementation files.
- `module/` exports top-level Fx modules (see `module/module.go`):
  - `module.Library`
  - `module.Server`
  - `module.Client`
- `internal/test/` contains shared test helpers/fixtures and Buf generation config.
- `test/` contains fixtures used by tests (configs, certs, secrets, reports).

Major subsystems:

- `config/`: decoding/validation and input routing (`config/decoder.go`, `config/default.go`, `config/env.go`, `config/file.go`).
- `transport/`: HTTP + gRPC wiring (`transport/http`, `transport/grpc`) and shared transport config.
- `telemetry/`: logger/metrics/tracer packages + modules.
- `crypto/`, `token/`, `database/sql/`, `health/`, `cache/`, etc.

## Dependency injection patterns

### Fx module composition

Modules are typically defined as `di.Module(...)` values composing submodules and constructors/invocations.

Example: `transport.Module` composes HTTP + gRPC wiring (see `transport/module.go`).

### Injected parameter structs

Packages frequently use `di.In` structs to declare injected dependencies (example: `config.DecoderParams` in `config/decoder.go:12-19`).

### Lifecycle hooks

Construction frequently registers cleanup via `Lifecycle.Append(di.Hook{...})` (example: `cache/NewCache` appends an `OnStop` hook in `cache/cache.go`).

## CLI wiring pattern

`cli.Application` registers “server” and “client” subcommands:

- `AddServer(...)` starts an Fx app and blocks until `app.Done()` then stops (see `cli/application.go:61-88`).
- `AddClient(...)` starts and immediately stops (see `cli/application.go:90-114`).

Errors are prefixed with the command name and use `dig.RootCause(...)` via `di.RootCause` (see `cli/application.go:137-139`).

## Configuration

### Config input routing via `-i`

`config.NewDecoder` dispatches based on the `-i` flag value (see `config/decoder.go:21-32`):

- `file:<path>` → file decoder
- `env:<ENV_VAR>` → env decoder
- otherwise → default lookup decoder

The default lookup searches for `<serviceName>.{yaml,yml,toml,json}` in (see `config/default.go:35-58`):

- executable directory
- `$XDG_CONFIG_HOME/<serviceName>/` (via `os.UserConfigDir()`)
- `/etc/<serviceName>/`

### Env config format

The README documents that env configs expect `extension:content` where `content` is base64-encoded (see `README.md` “Configuration”).

CI demonstrates this with `-i env:CONFIG` and `CONFIG=yml:<base64...>` (see `.circleci/config.yml:16-20`).

## Cache API gotcha

The `cache/` package contains both an instance API and package-level helpers.

- Instance methods are on `*cache.Cache` (see `cache/cache.go`).
- Package-level generic helpers use a package-global set via `cache.Register(...)` (see `cache/generic.go`).

Current helper signatures (see `cache/generic.go`):

- `cache.Get[T](ctx, key) (*T, error)`
- `cache.Persist[T](ctx, key, value, ttl) error`

Gotcha: `cache.Register(nil)` is used when cache is disabled; helpers are designed to be nil-safe and return the zero `*T` and `nil` error.

## Testing

- Tests commonly use `stretchr/testify/require` (see many `*_test.go`).
- Shared helpers live under `internal/test/`.
  - `internal/test/world.go` builds a test “world” using `fxtest.NewLifecycle` and registers multiple subsystems.
  - `internal/test/world.go:init` performs package registrations for tests.
- Fixtures are under `test/` (configs, certs, secrets).

## Style / formatting

From `.editorconfig`:

- Go files use tabs (`indent_style = tab`, `indent_size = 4`).
- YAML uses 2-space indentation.

Golangci-lint is configured in `.golangci.yml`.

## CI notes / gotchas

- Many `make` targets come from the `bin/` submodule; if it’s missing, `make` will fail due to missing includes.
- CI sets `GOEXPERIMENT=greenteagc` (see `.circleci/config.yml:6-9`).
- `make encode-config` uses `base64 -w 0` (see `bin/build/make/go.mak:109-112`), which may not work on BSD/macOS `base64`.
- CI generates a `.source-key` file via `make source-key` for caching (see `.circleci/config.yml:27`, `bin/build/make/git.mak`). This file is ignored by Git (see `.gitignore:11`).
- `vendor/` is ignored by git (see `.gitignore:2`) and is expected to be generated by `make dep`.
