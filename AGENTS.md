# AGENTS.md

This repository is a Go module (`github.com/alexfalkowski/go-service/v2`) that provides a reusable framework for building services (DI wiring, config decoding, transports, telemetry, crypto, etc.).

It is primarily a **library of packages** (no top-level `cmd/` binary in this repo).

Most workflows are driven by `make` targets that are defined in the `bin/` git submodule (see `Makefile:1-3`).

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

Gotcha: `.gitmodules` points at an SSH URL (`git@github.com:alexfalkowski/bin.git`) (`.gitmodules:1-3`). If you can’t fetch via SSH, `make` targets will fail until Git access is configured.

## Project type

- Language: Go
- Go version: `go 1.25.0` (`go.mod:1-4`)
- DI container: Uber Fx/Dig, wrapped by `di/` (`di/di.go:8-55`)
- CLI command framework: `github.com/cristalhq/acmd` (`cli/application.go`)
- Linting: `golangci-lint` plus additional tooling via `bin/build/go/*` (`bin/build/make/go.mak:24-55`, `.golangci.yml`)

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

`make dep` runs `go mod download`, `go mod tidy`, and `go mod vendor` (`bin/build/make/go.mak:10-20`).

Gotcha: tests run with `-mod vendor` (`bin/build/make/go.mak:46-47`), so after changing dependencies you typically must run `make dep` first.

### Tests

```sh
make specs
```

`make specs` runs `gotestsum` and executes `go test` with race + coverage, using the vendor directory (`bin/build/make/go.mak:46-47`).

Artifacts written under `test/reports/`:

- JUnit XML: `test/reports/specs.xml` (`bin/build/make/go.mak:46-47`)
- Coverage profile: `test/reports/profile.cov` (`bin/build/make/go.mak:46-47`)

### Lint / format

```sh
make lint
make fix-lint
make format
```

- `make lint` runs field alignment and `golangci-lint` (`bin/build/make/go.mak:36-42`).
- `make fix-lint` runs the same tools with fix mode where supported (`bin/build/make/go.mak:32-35`, `bin/build/make/go.mak:43-44`).
- `make format` runs `go fmt ./...` (`bin/build/make/go.mak:37-38`).

Golangci configuration is in `.golangci.yml` (formatters are enabled via `gci`, `gofmt`, `gofumpt`, `goimports` in `.golangci.yml:44-49`).

### Security checks

```sh
make sec
```

Runs `govulncheck -show verbose -test ./...` (`bin/build/make/go.mak:74-75`).

### Benchmarks

Convenience targets in `Makefile` (`Makefile:20-34`):

```sh
make benchmarks
make http-benchmarks
make grpc-benchmarks
make bytes-benchmarks
make strings-benchmarks
```

These delegate to `make package=<pkg> benchmark` (`bin/build/make/go.mak:49-50`).

### Coverage

```sh
make coverage
make html-coverage
make func-coverage
```

Coverage processing uses `test/reports/final.cov` and writes `test/reports/coverage.html` (`bin/build/make/go.mak:56-63`).

### Code generation (Buf)

```sh
make generate
```

Delegates to `make -C internal/test generate` (`Makefile:35-37`).

### Diagrams

Top-level targets (`Makefile:5-18`):

```sh
make diagrams
make crypto-diagram
make database-diagram
make telemetry-diagram
make transport-diagram
```

Under the hood this uses `goda graph ... | dot -Tpng` and writes PNGs into `assets/` (`bin/build/make/go.mak:84-85`).

### Local environment (integration deps)

```sh
make start
make stop
```

Uses `bin/build/docker/env` (`bin/build/make/go.mak:94-99`).

### TLS fixtures / certs

```sh
mkcert -install
make create-certs
```

Creates fixtures under `test/certs/` (`bin/build/make/go.mak:80-82`).

## Code organization

This repo is organized as packages under the module root.

Common conventions:

- Many subsystems are “feature packages” with `config.go`, `module.go`, and implementation files.
- `module/` exports top-level Fx modules (`module/module.go:24-61`):
  - `module.Library`
  - `module.Server`
  - `module.Client`
- `internal/test/` contains shared test helpers.
- `test/` contains fixtures used by tests (configs, certs, secrets, reports).

## Dependency injection patterns (Fx)

### Module composition

Modules are typically defined as `di.Module(...)` values composing submodules, constructors, and invocations.

Example: `module.Server` composes most server-side subsystems (`module/module.go:36-48`).

### Injected parameter structs

Packages frequently use `di.In` structs to declare injected dependencies.

Example: `config.DecoderParams` (`config/decoder.go:12-19`).

### Registrations / init-time wiring

Tests register some package globals in `internal/test/world.go:init` (`internal/test/world.go:33-39`):

- `telemetry.Register()`
- `transport/grpc.Register(FS)`
- `transport/http.Register(FS)`

Gotcha: both `transport/http` and `transport/grpc` use a package-level `var fs *os.FS` set via `Register(...)` (`transport/http/register.go:7-12`, `transport/grpc/register.go:7-12`). If you construct those servers without calling `Register`, TLS-related config construction may not have an FS available.

## Configuration

### Config input routing via `-i`

`config.NewDecoder` dispatches based on the `-i` flag value (`config/decoder.go:21-32`):

- `file:<path>` → file decoder
- `env:<ENV_VAR>` → env decoder
- otherwise → default lookup decoder

The default lookup searches for `<serviceName>.{yaml,yml,toml,json}` in (`config/default.go:35-58`):

- executable directory
- `$XDG_CONFIG_HOME/<serviceName>/` (via `os.UserConfigDir()`)
- `/etc/<serviceName>/`

### Env config format

Env configs expect `extension:content` where `content` is base64-encoded (`README.md:35-44`, `config/env.go`).

### “Source strings” pattern (`file:` / `env:` / raw)

Several configs accept a “source string” that can be:

- `env:NAME` (read from environment)
- `file:/path/to/secret` (read from filesystem)
- otherwise treated as the literal value

This is implemented by `os.FS.ReadSource` (`os/fs.go:100-110`) and is used for secrets/keys in multiple subsystems.

Gotcha: some paths are expanded (leading `~`) and cleaned via `os.FS.CleanPath` (`os/fs.go:75-88`).

Gotcha: `make encode-config` uses `base64 -w 0` (`bin/build/make/go.mak:77-78`), which may not work on BSD/macOS `base64`.

## HTTP / gRPC patterns

### HTTP handler instrumentation

The wrapper in `net/http` instruments handlers by wrapping them with `otelhttp.NewHandler` when registered via `net/http.Handle(...)` (`net/http/http.go:140-143`).

### Path patterns

`net/http.Pattern(name, pattern)` builds routes as `/<name><pattern>` (`net/http/http.go:181-184`).

## Cache API gotcha

The `cache/` package contains both an instance API and package-level generic helpers.

- Instance methods are on `*cache.Cache` (see `cache/cache.go`).
- Package-level generic helpers use a package-global set via `cache.Register(...)` (see `cache/generic.go`).

Gotcha: `cache.Register(nil)` is used when cache is disabled; helpers are designed to be nil-safe.

## Testing

- Tests commonly use `stretchr/testify/require`.
- Shared helpers live under `internal/test/`.
  - `internal/test/world.go` builds a test “world” with `fxtest.NewLifecycle` and configures telemetry, transports, and helpers.
- Fixtures are under `test/` (configs, certs, secrets).

## Style / formatting

From `.editorconfig`:

- Go files use tabs (`indent_style = tab`, `indent_size = 4`) (`.editorconfig:14-17`).
- YAML uses 2-space indentation (`.editorconfig:18-20`).

## Observed gotchas

- JWT verification enforces both algorithm + key id:
  - `token/jwt.Token.validate` rejects non-EdDSA tokens (`token/errors.ErrInvalidAlgorithm`) and requires `kid` to exist and match configured `jwt.kid` (`token/errors.ErrInvalidKeyID`) (`token/jwt/jwt.go:72-87`).
  - If you mint test tokens directly via `github.com/golang-jwt/jwt/v4`, remember to set `kid` or verification will fail.
- Telemetry header secrets may panic during config projection: `header.Map.MustSecrets` panics on read errors (`telemetry/header/header.go:26-29`). Treat secret loading failures as startup-fatal.
- `vendor/` is ignored by git and is expected to be generated by the dependency workflow (see `.gitignore`).

## CI notes

CircleCI (`.circleci/config.yml`) runs, in order (`.circleci/config.yml:25-74`):

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

CI sets `GOEXPERIMENT=greenteagc` (`.circleci/config.yml:6-9`).

CI provisions and waits for these services (`.circleci/config.yml:5-31`):

- Postgres (`tcp://localhost:5432`)
- Valkey/Redis (`tcp://localhost:6379`)
- `alexfalkowski/status` server (`tcp://localhost:6000`)
- Grafana Mimir (`tcp://localhost:9009`)
