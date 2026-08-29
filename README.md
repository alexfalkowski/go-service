![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/go-service.svg?style=shield)](https://circleci.com/gh/alexfalkowski/go-service)
[![codecov](https://codecov.io/gh/alexfalkowski/go-service/graph/badge.svg?token=AGP01JOTM0)](https://codecov.io/gh/alexfalkowski/go-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/go-service/v2)](https://goreportcard.com/report/github.com/alexfalkowski/go-service/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/go-service/v2.svg)](https://pkg.go.dev/github.com/alexfalkowski/go-service/v2)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# 🧰 Go Service

`github.com/alexfalkowski/go-service/v2` is an opinionated framework/library for building Go services with consistent wiring for configuration, DI, transports, telemetry, crypto, etc.

This repo is primarily a **library of packages** (no top-level `cmd/` binary). Services built on top typically define their own `main` package elsewhere and import this module.

Long-running services are expected to start from [`go-service-template`](https://github.com/alexfalkowski/go-service-template), while short-lived client commands start from [`go-client-template`](https://github.com/alexfalkowski/go-client-template). Both compose the high-level module bundles from this repository. These are the primary supported paths. Lower-level package-by-package composition is still available, but it is an advanced mode and may require extra manual registration.

---

## 🚀 Install

For a new long-running service, start from `go-service-template` so the application `main`, server command wiring, configuration fixtures, and standard module composition are generated together. For a short-lived control, migration, or batch command, start from `go-client-template`; it demonstrates `cli.Application.AddClient`, `module.Client`, and lifecycle `OnStart` command work.

For direct package use in an existing module, add the library dependency with the versioned module path:

```sh
go get github.com/alexfalkowski/go-service/v2
```

Use the Go version declared in `go.mod` or newer when installing or building this module.

---

## 🧩 Dependency Injection (Fx)

The framework is designed around dependency injection and uses [Uber Fx](https://github.com/uber-go/fx) (and Dig under the hood). Most subsystems expose Fx modules that you compose into your service.

If you are new to Fx, their docs/examples are worth reading first.

### Module bundles

The module package exposes three top-level bundles:

- `module.Library` for shared foundations (env, compress, encoding, crypto, time, sync buffer-pool wiring, id)
- `module.Server` for server processes (Library + config, transports, telemetry, debug, health, etc.)
- `module.Client` for short-lived/batch/client processes (Library + config, telemetry, sql, hooks, etc.)

These bundles are the intended default for services generated from `go-service-template`. They handle the internal registration expected by the framework so most services do not need to wire lower-level transport or lifecycle helpers manually.

### Minimal CLI bootstrap example

This repository is a library, so your binary is usually in another module. A typical `main` uses `cli.Application` and composes module bundles:

```go
package main

import (
    "github.com/alexfalkowski/go-service/v2/cli"
    "github.com/alexfalkowski/go-service/v2/context"
    "github.com/alexfalkowski/go-service/v2/module"
    "github.com/alexfalkowski/go-service/v2/os"
)

func main() {
    app := cli.NewApplication(func(commander cli.Commander) {
        serve := commander.AddServer("serve", "Run the service", module.Server)
        serve.AddConfig("file:./config.yaml") // adds the `-config` / `-c` config flag with this default
    })

    os.Exit(app.RunCode(context.Background()))
}
```

The `file:./config.yaml` default above expects a non-empty config file. A minimal
server config can start with the environment plus one enabled transport:

```yaml
environment: development
transport:
  http:
    address: tcp://localhost:8000
    timeout: 10s
```

Use `app.RunCode(context.Background())` from `main` when exiting the process. It
returns `os.ExitCodeSuccess` on success, returns a requested non-zero shutdown
exit code such as `os.ExitCodeServeFailure`, and returns `os.ExitCodeFailure`
for other errors. Use `app.Run(context.Background())` in tests or embedding code
that needs to inspect the returned error.

---

## 🖥️ CLI

Services commonly expose two command shapes:

- **Server**: long-running daemon process
- **Client**: short-lived control/admin process

The framework uses [acmd](https://github.com/cristalhq/acmd). Your service’s `main` typically wires Fx modules + commands.

> This repo intentionally does not ship a ready-to-run `main` — it provides the building blocks. In normal usage server applications consume them through `go-service-template` plus `module.Server`, while short-lived commands use `go-client-template` plus `module.Client`, rather than wiring every subsystem manually.

---

## 🗂️ Repository layout

The repo is intentionally split between high-level service composition and lower-level reusable helpers:

- `module/` exposes the opinionated Fx bundles (`Library`, `Server`, `Client`)
- `config/` defines the standard top-level config shape plus projections used by module wiring
- feature packages such as `cache/`, `crypto/`, `database/sql/`, `feature/`, `telemetry/`, `time/`, and `id/` provide config, constructors, and Fx modules for a subsystem
- `net/...` contains lower-level protocol helpers and reusable primitives (`net/http`, `net/grpc`, metadata/header helpers, gRPC health protocol aliases, and `net/server`)
- `transport/...` contains the higher-level service transport layer: composed HTTP/gRPC stacks, policy middleware, operational endpoints, and transport-specific modules
- `internal/test/` contains the shared test world and fixtures used across packages

As a rule of thumb: if you want protocol primitives or shared helpers, start in `net/...`; if you want service wiring and middleware policy, start in `transport/...`. Shared metadata, header, and lifecycle helpers live under `net/...`, including `net/http/meta`, `net/grpc/meta`, `net/header`, and `net/server.Register`.

For most service authors, the right starting point is still the high-level module bundles rather than these lower-level packages directly.

---

## ⚙️ Configuration

### Supported config formats

The config decoder supports:

- JSON
- HJSON (`github.com/hjson/hjson-go/v4`)
- TOML (`github.com/BurntSushi/toml`)
- YAML (`go.yaml.in/yaml/v3`)

### Selecting the config source (`-config` / `-c` flags)

Config input is routed by flags called `-config` and `-c`:

- `file:<path>`
  Read config from a file at `<path>`; parser is selected from the file extension (`.json`, `.hjson`, `.yaml`, `.toml`).

- `env:<ENV_VAR>`
  Read config from env var `<ENV_VAR>`. The env var value must be formatted as:

  `"<extension>:<base64-content>"`

  Example format: `yaml:ZW52aXJvbm1lbnQ6IGRldmVsb3BtZW50Cg==`

  Example commands:

  ```sh
  # Linux (GNU base64)
  export SERVICE_CONFIG="yaml:$(base64 -w 0 < ./config.yaml)"
  ./your-service serve -config env:SERVICE_CONFIG
  ```

  ```sh
  # macOS/BSD base64
  export SERVICE_CONFIG="yaml:$(base64 < ./config.yaml | tr -d '\n')"
  ./your-service serve -c env:SERVICE_CONFIG
  ```

  HJSON works the same way, for example `hjson:<base64-content>`.

  The repository helper `make kind=configs/config encode-config` uses GNU `base64 -w 0`; on macOS/BSD, use `base64 | tr -d '\n'` for the equivalent single-line payload.

- Unsupported explicit `kind:location` prefixes fail startup instead of falling back to another source.

- Unprefixed values, including an empty value, fall back to **default lookup**, searching for:

  `<serviceName>.{yaml,hjson,toml,json}`

  Default lookup checks extensions first (`.yaml`, `.hjson`, `.toml`, `.json`), and for each extension checks:
  - executable directory
  - `$XDG_CONFIG_HOME/<serviceName>/` (via `os.UserConfigDir()`)
  - `/etc/<serviceName>/`

> [!IMPORTANT]
> Because the user config directory is part of that search, runtimes using default lookup are expected to provide `HOME` or `XDG_CONFIG_HOME`. Services that cannot rely on those environment variables should pass an explicit `-config file:<path>` or `-config env:<ENV_VAR>` source.

### Typed decoding and validation

At runtime, services typically decode into a struct (often embedding `config.Config`) and validate it using `go-playground/validator`.

The library provides a helper `config.NewConfig[T]` which:

- decodes into `*T`
- rejects an “empty” decoded value (guards against starting with a zero-value config)
- validates the decoded config

Empty detection uses zero-value semantics and supports config types containing maps, slices, or other
non-comparable fields.

Example:

```go
type WorkerConfig struct {
    Queue string `yaml:"queue" json:"queue" toml:"queue" validate:"required"`
}

type AppConfig struct {
    Worker         *WorkerConfig `yaml:"worker" json:"worker" toml:"worker" validate:"required"`
    *config.Config `yaml:",inline" json:",inline" toml:",inline" validate:"required"`
}

func loadConfig(decoder config.Decoder, validator *config.Validator) (*AppConfig, error) {
    return config.NewConfig[AppConfig](decoder, validator)
}

func sharedConfig(cfg *AppConfig) *config.Config {
    return cfg.Config
}

func workerConfig(cfg *AppConfig) *WorkerConfig {
    return cfg.Worker
}

var AppConfigModule = di.Module(
    di.Constructor(config.NewConfig[AppConfig]),
    di.Decorate(sharedConfig),
    di.Constructor(workerConfig),
)
```

Compose `AppConfigModule` alongside `module.Server` or `module.Client`. The
decorator projects the embedded shared `*config.Config` into the standard graph,
so the service-specific config is decoded once while existing transport, SQL,
and telemetry projections continue to work. Add constructors like
`workerConfig` for service-owned sub-configs.

### The standard top-level config shape

The canonical top-level config type is `config.Config` (in `config/config.go`). It contains:

- `debug`, `cache`, `crypto`, `feature`, `hooks`, `id`, `sql`, `telemetry`, `time`, `transport`, `environment`

Most sub-configs are optional pointers. Conventionally, `nil` means **disabled**.

---

## 🔐 Source strings (secrets, DSNs, paths)

Many fields accept a _source string_ rather than only a literal:

- `env:NAME` → read from environment variable `NAME` (fails if `NAME` is unset; resolves to an empty value if `NAME` is explicitly set to `""`)
- `file:/path/to/thing` → read from filesystem after path cleaning; returned bytes are trimmed of leading and trailing whitespace
- otherwise → treat as literal string

This is used for secrets and key material (TLS keys, HMAC keys, webhook secrets, SQL DSNs, etc).
`env:` values and literal values are returned exactly as provided; they are not
trimmed.

Example:

```yaml
hooks:
  key: current
  secrets:
    current: env:WEBHOOK_SECRET
```

---

## 🌍 Environment

Top-level environment is:

```yaml
environment: development
```

This is an `env.Environment` value used to drive environment-specific behavior in services.

---

## 🗜️ Compression

Compression kinds used by subsystems that support compression:

- `none`
- `zstd`
- `s2`
- `snappy`

---

## 🧾 Encoders

Encoding kinds used by subsystems that support encoding. `encoding.Map` registers each encoder under
exactly one canonical kind (no aliases):

- `json`
- `hjson`
- `toml`
- `yaml`
- `msgpack`
- `protobuf`
- `protojson`
- `prototext`
- `gob`
- `bytes`

> [!NOTE]
>
> - `bytes` is the passthrough encoder for `io.ReaderFrom`/`io.WriterTo` payloads.
> - HTTP media-type aliases such as `pb`, `proto`, `protobin`, `pbbin`, `pbtxt`, `prototxt`,
>   `pbjson`, `octet-stream`, `plain`, and `yml` are resolved to the canonical kinds above by
>   `net/http/content/unary` before they ever reach this registry. See [HTTP content types](docs/transport.md#http-content-types).
> - `encoding/stream.Map` is a separate registry for streaming (multi-value) encoding — `json`, `msgpack`,
>   `gob`, `yaml` — used by [HTTP streaming (NDJSON)](docs/transport.md#http-streaming-ndjson), not by this single-value
>   registry.
> - Not every kind in this registry is interchangeable for HTTP request-body decoding: `msgpack` and
>   `gob` remain valid response codecs but are rejected as a request `Content-Type`. See
>   [HTTP content types](docs/transport.md#http-content-types).

---

## 💾 Cache

Cache configuration is defined in `cache/config.Config`. Built-in driver kinds are `redis` and `ttlcache`.

See [docs/cache.md](docs/cache.md) for the config shape, driver semantics, size/entry limits, key-namespace implications of `compressor`/`encoder`, and `Cache.Flush` behavior.

---

## 🚩 Feature flags (OpenFeature)

`feature.Config` embeds client-side config (`config/client.Config`): address, timeout, retry, breaker, limiter, TLS, token, and options. This repository does not construct a built-in OpenFeature provider from this config — supply your own `openfeature.FeatureProvider` and `feature.Module` registers it with the SDK lifecycle.

See [docs/feature-flags.md](docs/feature-flags.md) for a full config example and provider-wiring notes.

---

## 🪝 Webhooks (Standard Webhooks)

Configured via `hooks.Config`, using Standard Webhooks signing/verification with key rotation and an optional clock-skew `leeway`.

See [docs/webhooks.md](docs/webhooks.md) for the config shape, signing/verification behavior, `leeway`, idempotency guidance, and the CloudEvents structured-encoding requirement.

---

## 🆔 ID generation

Supported ID kinds:

- `uuid`
- `ksuid`
- `nanoid`
- `ulid`
- `xid`

Config:

```yaml
id:
  kind: uuid
```

> [!NOTE]
> ID generators produce operational identifiers such as request ids, webhook ids, and token `jti` values. They are not a secret-material API and should not be used as passwords, bearer tokens, or other credentials. Omit `id` entirely to select the `uuid` default. If `id` is present, `kind` must be one of the supported registered kinds. Sortable kinds such as `ksuid`, `ulid`, and `xid` expose ordering characteristics.

---

## 🚀 Runtime enhancements

Server commands created through `cli.Application.AddServer` include `runtime.Module`, which currently enables:

- [automemlimit](https://github.com/KimMachineGun/automemlimit)

> [!NOTE]
> This registration is best-effort and does not fail startup if a memory limit cannot be applied. When automemlimit detects an unlimited cgroup and `GOMEMLIMIT` is not already configured, it sets Go's runtime memory limit to `math.MaxInt64`, replacing any programmatic limit applied before startup. Direct Fx compositions and client-style commands should include `runtime.Module` explicitly when they want this behavior.

---

## 🐘 SQL (Postgres)

SQL root config is `database/sql.Config`, with Postgres under `sql.pg`. `module.Server` and `module.Client` both wire PostgreSQL support via `database/sql/pg.Module`; a nil `sql` or `sql.pg` block disables it.

See [docs/sql.md](docs/sql.md) for the config shape, reader/writer pool settings, DSN resolution, ping/health-check guidance, and dependencies.

---

## 🩺 Health

Health checks are based on [go-health](https://github.com/alexfalkowski/go-health) and expose Kubernetes-style `/<name>/healthz`, `/<name>/livez`, and `/<name>/readyz` endpoints, plus the standard gRPC `grpc.health.v1.Health` service when gRPC transport is enabled.

`module.Server` installs the HTTP/gRPC health transports, but services own the checks and observer mapping — see the executable [`Registrations` example](health/example_test.go).

See [docs/health.md](docs/health.md) for endpoint behavior, checker wiring, and the gRPC health protocol.

---

## 📡 Telemetry

Telemetry config root is `telemetry.Config`, covering resource attributes, metadata size limits, propagation formats, and the logger/metrics/tracer signals.

See [docs/telemetry.md](docs/telemetry.md) for the config shape, metadata, propagation, logging (JSON/text/tint/OTLP), metrics (Prometheus/OTLP, histogram buckets), tracing, libraries used, and dependencies.

---

## 🎫 Tokens

Token configuration is rooted at `token.Config`, usually nested under `transport.http.token` and/or `transport.grpc.token`. Supported kinds are `jwt` and `paseto`, plus a shared Casbin-based access-control layer configured at `transport.access`.

See [docs/tokens.md](docs/tokens.md) for Casbin RBAC access control, and the JWT/PASETO config shapes, key rotation, and verification semantics.

---

## 🚦 Limiter

Limiter config is `transport/limiter.Config`, typically applied at transport level. Built-in key kinds are `user-id`, `transport-service-method`, `service-method`, `ip`, and `user-agent`.

See [docs/limiter.md](docs/limiter.md) for the config shape, defaults, key semantics, and response headers/metadata.

---

## 🕒 Time (network time)

Time config:

```yaml
time:
  kind: nts
  address: time.cloudflare.com
  timeout: 2s
```

Supported kinds:

- `ntp`
- `nts`

Omit the `time` block to disable network time. If the block is present, `kind`
must be `ntp` or `nts`; empty or unknown kinds fail startup with the time
provider not found error. `address` is provider-specific and is used when the
network time provider performs I/O. `timeout` bounds network operations for the
selected provider; a zero value uses the upstream client's default timeout, and
negative values are invalid.

---

## 🌐 Transport

The transport layer provides higher-level wiring and middleware policy for communication in/out of the service: composed HTTP/gRPC server and client stacks, retries, breakers, token middleware, and health wiring. `net/...` holds the lower-level protocol primitives those stacks build on.

Supported stacks include gRPC, HTTP REST/RPC abstractions with content negotiation, HTTP MVC helpers, and CloudEvents (`transport/http/events`).

See [docs/transport.md](docs/transport.md) for server configuration and TLS, HTTP content types and streaming (NDJSON), route policy, route-miss and MVC error handling, forwarded-IP/reflection posture, and client-side circuit breakers and retries.

---

## 🔑 Cryptography

The crypto root config is `crypto.Config` and supports AES, Ed25519, HMAC, and RSA key types. Most fields are source strings.

See [docs/crypto.md](docs/crypto.md) for the config shape, key format requirements, the `crypto.Message` encryption API, and dependencies.

---

## 🛠️ Debug endpoints

Debug server config exposes `statsviz`, `pprof`, and `fgprof` under `/<name>/debug/...`, optionally behind TLS.

See [docs/debug.md](docs/debug.md) for the config shape, endpoint paths, and TLS setup.

---

## 🧑‍💻 Development

This repo generally follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) and uses a `bin/` git submodule for `make` targets; run `make help` to discover them.

See [docs/development.md](docs/development.md) for repo setup, development dependencies, tests, lint/format, security checks, benchmarks, fuzz tests, coverage reports, code generation, and architecture diagrams.
