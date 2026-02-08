# AGENTS.md

This repository is a Go module (`github.com/alexfalkowski/go-service/v2`) that provides a reusable framework for building services (DI wiring, config decoding, transports, telemetry, crypto, etc.).

It is primarily a **library of packages** (no top-level `cmd/` binary in this repo).

Most workflows are driven by `make` targets that are defined in the `bin/` git submodule (the top-level `Makefile` includes `bin/build/make/{help,go,git}.mak`).

## First-time setup

### Git submodule (required for `make`)

The top-level `Makefile` includes make fragments from the `bin` submodule (`bin/build/make/*.mak`).

```sh
git submodule sync
git submodule update --init
```

Alternative (same effect):

```sh
make submodule
```

Gotcha: `.gitmodules` uses an SSH URL (`git@github.com:alexfalkowski/bin.git`). If you can’t fetch via SSH, `make` targets will fail until Git access is configured.

## Project type

- Language: Go
- Go version: `go 1.25.0` (`go.mod`)
- DI container: Uber Fx/Dig, wrapped by `di/` (`di/di.go`)
- CLI command framework: `github.com/cristalhq/acmd` (`cli/application.go`)
- Linting: `golangci-lint` plus helper tooling in the `bin` submodule (`bin/build/go/*`, `.golangci.yml`)

## Essential commands

Most targets are defined in `bin/build/make/*.mak` (included by the top-level `Makefile`).

### Discover targets

```sh
make help
```

### Dependencies (keeps `vendor/` in sync)

```sh
make dep
```

`make dep` runs:

- `go mod download`
- `go mod tidy`
- `go mod vendor`

Gotcha: tests run with `-mod vendor`, so after changing dependencies you typically must run `make dep` first.

### Tests

```sh
make specs
```

`make specs` runs `gotestsum` and executes `go test` with race + coverage, using the vendor directory.

Artifacts written under `test/reports/`:

- JUnit XML: `test/reports/specs.xml`
- Coverage profile: `test/reports/profile.cov`

### Lint / format

```sh
make lint
make fix-lint
make format
```

- `make lint` runs field alignment and `golangci-lint`.
- `make fix-lint` runs the same tools with fix mode where supported.
- `make format` runs `go fmt ./...`.

Formatter configuration lives in `.golangci.yml`.

### Security checks

```sh
make sec
```

Runs `govulncheck -show verbose -test ./...`.

### Benchmarks

Convenience targets in the top-level `Makefile`:

```sh
make benchmarks
make http-benchmarks
make grpc-benchmarks
make bytes-benchmarks
make strings-benchmarks
```

These delegate to `make package=<pkg> benchmark`.

### Coverage

```sh
make coverage
make html-coverage
make func-coverage
```

Coverage commands expect that `make specs` has already produced `test/reports/profile.cov`.

`make html-coverage` / `make func-coverage` operate on `test/reports/final.cov`, which is generated from `test/reports/profile.cov` by filtering (see `bin/quality/go/covfilter`).

### Code generation (Buf)

```sh
make generate
```

Delegates to `make -C internal/test generate`.

### Diagrams

Top-level targets:

```sh
make diagrams
make crypto-diagram
make database-diagram
make telemetry-diagram
make transport-diagram
```

Under the hood this uses `goda graph ... | dot -Tpng` and writes PNGs into `assets/`.

### Local environment (integration deps)

```sh
make start
make stop
```

Uses `bin/build/docker/env`.

### TLS fixtures / certs

```sh
mkcert -install
make create-certs
```

Creates fixtures under `test/certs/`.

### Encoding configs for env-based loading

```sh
make kind=status encode-config
```

Gotcha: `encode-config` uses `base64 -w 0` (GNU coreutils); on macOS/BSD `base64`, `-w` may not exist.

## Code organization

This repo is organized as packages under the module root.

Common conventions:

- Many subsystems are “feature packages” with `config.go`, `module.go`, and implementation files.
- `module/` exports top-level Fx modules (`module/module.go`):
  - `module.Library`
  - `module.Server`
  - `module.Client`
- `internal/test/` contains shared test helpers.
- `test/` contains fixtures used by tests (configs, certs, secrets, reports).

## Dependency injection patterns (Fx)

### Module composition

Modules are typically defined as `di.Module(...)` values composing submodules, constructors, and invocations.

Example: `module.Server` composes most server-side subsystems (`module/module.go`).

### Injected parameter structs

Packages frequently use `di.In` structs to declare injected dependencies.

Example: `config.DecoderParams` (`config/decoder.go`).

### Registrations / init-time wiring

Some packages use package-level registration for wiring globals in tests and/or when a feature is used.

Examples:

- Telemetry setup via `telemetry.Register()`.
- Transport packages set a package-level `fs` via `transport/http.Register(...)` and `transport/grpc.Register(...)` (used to load TLS key material via `os.FS.ReadSource`).

Gotcha: if you construct the HTTP/gRPC servers or clients with TLS enabled without having called those `Register(...)` functions, TLS config construction may not have an `*os.FS` available.

## Configuration

### Config input routing via `-i`

`config.NewDecoder` dispatches based on the `-i` flag value:

- `file:<path>` → file decoder
- `env:<ENV_VAR>` → env decoder
- otherwise → default lookup decoder

The default lookup searches for `<serviceName>.{yaml,yml,toml,json}` in:

- executable directory
- `$XDG_CONFIG_HOME/<serviceName>/` (via `os.UserConfigDir()`)
- `/etc/<serviceName>/`

### Env config format

Env configs expect `extension:content` where `content` is base64-encoded (see `config/env.go`, README).

### “Source strings” pattern (`file:` / `env:` / raw)

Several configs accept a “source string” that can be:

- `env:NAME` (read from environment)
- `file:/path/to/secret` (read from filesystem)
- otherwise treated as the literal value

This is implemented by `os.FS.ReadSource` (`os/fs.go`).

Gotcha: some paths are expanded (leading `~`) and cleaned via `os.FS.CleanPath`.

## HTTP / gRPC patterns

### HTTP handler instrumentation

The wrapper in `net/http` instruments handlers by wrapping them with `otelhttp.NewHandler` when registered via `net/http.Handle(...)` (`net/http/http.go`).

### Path patterns

`net/http.Pattern(name, pattern)` builds routes as `/<name><pattern>` (`net/http/http.go`).

### Circuit breakers (client-side)

- gRPC breaker uses per-`fullMethod` circuit breakers and counts only selected gRPC status codes as failures (`transport/grpc/breaker`).
- HTTP breaker uses per-`method + host` circuit breakers and counts failures via status code classification (defaults to `>= 500` and `429`) (`transport/http/breaker`).
- Both breakers treat non-response/transport errors (e.g., cancellations) as successful for breaker accounting by default.

## Cache API gotcha

The `cache/` package contains both an instance API and package-level generic helpers.

- Instance methods are on `*cache.Cache` (see `cache/cache.go`).
- Package-level generic helpers use a package-global set via `cache.Register(...)` (see `cache/generic.go`, wired by `cache/module.go`).

Gotcha: callers are expected to tolerate cache being nil/disabled; helpers are designed to be nil-safe.

## Testing

- Tests commonly use `stretchr/testify/require`.
- Shared helpers live under `internal/test/`.
  - `internal/test/world.go` builds a test “world” with `fxtest.NewLifecycle` and configures telemetry, transports, and helpers.
- Fixtures are under `test/` (configs, certs, secrets).

## Style / formatting

From `.editorconfig`:

- Go files use tabs (`indent_style = tab`).
- YAML uses 2-space indentation.

## Observed gotchas

- JWT verification enforces both algorithm + key id (see `token/jwt/jwt.go`). If you mint test tokens directly via `github.com/golang-jwt/jwt/v4`, remember to set `kid` or verification will fail.
- Telemetry header secrets may panic during config projection: `telemetry/header/header.go` uses `header.Map.MustSecrets`, which panics on secret read errors.
- `vendor/` is ignored by git and is expected to be generated by the dependency workflow (see `.gitignore`).

## CI notes

CircleCI (`.circleci/config.yml`) runs, in order:

- submodule init
- `make source-key`
- `mkcert -install` / `make create-certs`
- wait for dependent services
- `make clean`
- `make dep`
- `make lint`
- `make sec`
- `make specs`
- `make benchmarks`
- `make coverage`

CI sets `GOEXPERIMENT=greenteagc`.

CI provisions and waits for these services:

- Postgres (`tcp://localhost:5432`)
- Valkey/Redis (`tcp://localhost:6379`)
- `alexfalkowski/status` server (`tcp://localhost:6000`)
- Grafana Mimir (`tcp://localhost:9009`)
