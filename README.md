![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/go-service.svg?style=shield)](https://circleci.com/gh/alexfalkowski/go-service)
[![codecov](https://codecov.io/gh/alexfalkowski/go-service/graph/badge.svg?token=AGP01JOTM0)](https://codecov.io/gh/alexfalkowski/go-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/go-service/v2)](https://goreportcard.com/report/github.com/alexfalkowski/go-service/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/go-service/v2.svg)](https://pkg.go.dev/github.com/alexfalkowski/go-service/v2)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Go Service

`github.com/alexfalkowski/go-service/v2` is an opinionated framework/library for building Go services with consistent wiring for configuration, DI, transports, telemetry, crypto, etc.

This repo is primarily a **library of packages** (no top-level `cmd/` binary). Services built on top typically define their own `main` package elsewhere and import this module.

---

## Dependency Injection (Fx)

The framework is designed around dependency injection and uses [Uber Fx](https://github.com/uber-go/fx) (and Dig under the hood). Most subsystems expose Fx modules that you compose into your service.

If you are new to Fx, their docs/examples are worth reading first.

### Module bundles

The module package exposes three top-level bundles:

- `module.Library` for shared foundations (env, compress, encoding, crypto, time, sync buffer-pool wiring, id)
- `module.Server` for server processes (Library + config, transports, telemetry, debug, health, etc.)
- `module.Client` for short-lived/batch/client processes (Library + config, telemetry, sql, hooks, etc.)

### Minimal CLI bootstrap example

This repository is a library, so your binary is usually in another module. A typical `main` uses `cli.Application` and composes module bundles:

```go
package main

import (
    "context"

    "github.com/alexfalkowski/go-service/v2/cli"
    "github.com/alexfalkowski/go-service/v2/module"
)

func main() {
    app := cli.NewApplication(func(commander cli.Commander) {
       server := commander.AddServer("serve", "Run the service", module.Server)
       server.AddInput("file:./config.yml") // enables the `-i` flag used by config.NewDecoder
    })

    app.ExitOnError(context.Background())
}
```

---

## CLI

Services commonly expose two command shapes:

- **Server**: long-running daemon process
- **Client**: short-lived control/admin process

The framework uses [acmd](https://github.com/cristalhq/acmd). Your service’s `main` typically wires Fx modules + commands.

> This repo intentionally does not ship a ready-to-run `main` — it provides the building blocks.

---

## Repository layout

The repo is intentionally split between high-level service composition and lower-level reusable helpers:

- `module/` exposes the opinionated Fx bundles (`Library`, `Server`, `Client`)
- `config/` defines the standard top-level config shape plus projections used by module wiring
- feature packages such as `cache/`, `crypto/`, `database/sql/`, `feature/`, `telemetry/`, `time/`, and `id/` provide config, constructors, and Fx modules for a subsystem
- `net/...` contains lower-level protocol helpers and reusable primitives (`net/http`, `net/grpc`, metadata/header helpers, gRPC health, and `net/server`)
- `transport/...` contains the higher-level service transport layer: composed HTTP/gRPC stacks, policy middleware, operational endpoints, and transport-specific modules
- `internal/test/` contains the shared test world and fixtures used across packages

As a rule of thumb: if you want protocol primitives or shared helpers, start in `net/...`; if you want service wiring and middleware policy, start in `transport/...`.

---

## Configuration

### Supported config formats

The config decoder supports:

- JSON (`encoding/json`)
- HJSON (`github.com/hjson/hjson-go`)
- TOML (`github.com/BurntSushi/toml`)
- YAML (`go.yaml.in/yaml/v3`)

### Selecting the config source (`-i` flag)

Config input is routed by a flag called `-i`:

- `file:<path>`
  Read config from a file at `<path>`; parser is selected from the file extension (`.json`, `.hjson`, `.yaml`, `.yml`, `.toml`).

- `env:<ENV_VAR>`
  Read config from env var `<ENV_VAR>`. The env var value must be formatted as:

  `"<extension>:<base64-content>"`

  Example format: `yaml:ZW52aXJvbm1lbnQ6IGRldmVsb3BtZW50Cg==`

  Example commands:

  ```sh
  # Linux (GNU base64)
  export SERVICE_CONFIG="yaml:$(base64 -w 0 < ./config.yml)"
  ./your-service serve -i env:SERVICE_CONFIG
  ```

  ```sh
  # macOS/BSD base64
  export SERVICE_CONFIG="yaml:$(base64 < ./config.yml | tr -d '\n')"
  ./your-service serve -i env:SERVICE_CONFIG
  ```

  HJSON works the same way, for example `hjson:<base64-content>`.

- Otherwise (no `file:`/`env:` prefix), the decoder falls back to **default lookup**, searching for:

  `<serviceName>.{yaml,yml,hjson,toml,json}`

  in:
  - executable directory
  - `$XDG_CONFIG_HOME/<serviceName>/` (via `os.UserConfigDir()`)
  - `/etc/<serviceName>/`

### Typed decoding and validation

At runtime, services typically decode into a struct (often embedding `config.Config`) and validate it using `go-playground/validator`.

The library provides a helper `config.NewConfig[T]` which:

- decodes into `*T`
- rejects an “empty” decoded value (guards against starting with a zero-value config)
- validates the decoded config

Example:

```go
type AppConfig struct {
    config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

func loadConfig(decoder config.Decoder, validator *config.Validator) (*AppConfig, error) {
    return config.NewConfig[AppConfig](decoder, validator)
}
```

### The standard top-level config shape

The canonical top-level config type is `config.Config` (in `config/config.go`). It contains:

- `debug`, `cache`, `crypto`, `feature`, `hooks`, `id`, `sql`, `telemetry`, `time`, `transport`, `environment`

Most sub-configs are optional pointers. Conventionally, `nil` means **disabled**.

---

## Source strings (secrets, DSNs, paths)

Many fields accept a *source string* rather than only a literal:

- `env:NAME` → read from environment variable `NAME` (fails if `NAME` is unset; resolves to an empty value if `NAME` is explicitly set to `""`)
- `file:/path/to/thing` → read from filesystem
- otherwise → treat as literal string

This is used for secrets and key material (TLS keys, HMAC keys, webhook secrets, SQL DSNs, etc).

Example:

```yaml
hooks:
  secret: env:WEBHOOK_SECRET
```

---

## Environment

Top-level environment is:

```yaml
environment: development
```

This is an `env.Environment` value used to drive environment-specific behavior in services.

---

## Compression

Compression kinds used by subsystems that support compression:

- `none`
- `zstd`
- `s2`
- `snappy`

---

## Encoders

Encoding kinds used by subsystems that support encoding:

- `json`
- `hjson`
- `toml`
- `yaml`
- `yml`
- `proto`
- `pb`
- `protobuf`
- `protobin`
- `pbbin`
- `protojson`
- `pbjson`
- `prototext`
- `prototxt`
- `pbtxt`
- `gob`
- `plain`
- `octet-stream`
- `markdown`

Notes:

- `plain`, `octet-stream`, and `markdown` all map to the bytes passthrough encoder.
- Protobuf binary/text/JSON kinds have multiple aliases; the list above reflects the built-in registry.

---

## Cache

Cache configuration is defined in `cache/config.Config`:

```yaml
cache:
  kind: redis
  compressor: zstd
  encoder: json
  options:
    url: env:CACHE_URL
```

Notes:

- Built-in driver kinds in this repo are `redis` and `sync`.
- `kind` is still wiring-dependent in practice: services can register additional drivers.
- `options` is backend-specific and decoded as `map[string]any`.

---

## Feature flags (OpenFeature)

The `feature.Config` embeds client-side config (`config/client.Config`), so it supports:

- `address`
- `timeout`
- `retry`
- `limiter`
- `tls`
- `token`
- `options`

Example:

```yaml
feature:
  address: localhost:9000
  timeout: 10s
  retry:
    backoff: 100ms
    timeout: 1s
    attempts: 3
  tls:
    cert: file:test/certs/client-cert.pem
    key: file:test/certs/client-key.pem
```

Notes:

- Presence enables the feature subsystem configuration-wise, but you still need to register an OpenFeature provider in your service wiring.

---

## Webhooks (Standard Webhooks)

Configured via `hooks.Config`:

```yaml
hooks:
  secret: file:test/secrets/hooks
```

`secret` is a source string.

---

## ID generation

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

---

## Runtime enhancements

The runtime is enhanced with:

- [automemlimit](https://github.com/KimMachineGun/automemlimit)

---

## SQL (Postgres)

SQL root config is `database/sql.Config`, with Postgres under `sql.pg`.

Postgres config embeds common pool + DSN config (`database/sql/config.Config`), including master/slave DSNs and pool sizes.

`module.Server` and `module.Client` both include `sql.Module`, which currently wires PostgreSQL support via `database/sql/pg.Module`.

Enablement is presence-based: a nil `sql` block or a nil `sql.pg` block disables SQL wiring. When enabled, the pgx stdlib driver is registered under the name `pg`, master/slave DSNs are resolved using the source-string rules described above, OpenTelemetry `database/sql` stats metrics are registered, and the resulting pools are closed on lifecycle stop.

Example (with source strings for DSNs):

```yaml
sql:
  pg:
    masters:
      - url: env:PG_MASTER_DSN
    slaves:
      - url: env:PG_SLAVE_DSN
    max_open_conns: 5
    max_idle_conns: 5
    conn_max_lifetime: 1h
```

Example (literal DSN; not recommended for production secrets):

```yaml
sql:
  pg:
    masters:
      - url: postgres://user:pass@localhost:5432/dbname?sslmode=disable
    max_open_conns: 10
```

### Dependencies

![Dependencies](./assets/database.png)

---

## Health

Health checks are based on [go-health](https://github.com/alexfalkowski/go-health).

The framework provides Kubernetes-style endpoints:

- `/<name>/healthz` — general serving health status
- `/<name>/livez` — liveness probe
- `/<name>/readyz` — readiness probe

Successful health responses return HTTP 200 with the plain-text body `SERVING`.
Missing or failing observers return HTTP 503 with the standard go-service error response.

These are modeled after [Kubernetes API health endpoints](https://kubernetes.io/docs/reference/using-api/health-checks/).

---

## Telemetry

Telemetry config root is `telemetry.Config`:

```yaml
telemetry:
  logger: ...
  metrics: ...
  tracer: ...
```

### Logging

Logging uses `log/slog`.

Supported built-in logger kinds:

- `json`
- `text`
- `tint`
- `otlp`

#### JSON logger

```yaml
telemetry:
  logger:
    kind: json
    level: info
```

#### Text logger

```yaml
telemetry:
  logger:
    kind: text
    level: info
```

#### OTLP logger

```yaml
telemetry:
  logger:
    kind: otlp
    level: info
    url: http://localhost:4318/v1/logs
    headers:
      Authorization: env:OTLP_LOGS_AUTH
```

Notes:

- `headers` values are source strings.
- Telemetry header maps are resolved during config projection; unset `env:` values and unreadable `file:` values fail fast (panic during startup).

### Metrics

Supported metrics kinds:

- `prometheus`
- `otlp`

#### Prometheus

```yaml
telemetry:
  metrics:
    kind: prometheus
```

When Prometheus is enabled on HTTP transport, metrics are exposed at `/<name>/metrics`.

#### OTLP metrics

```yaml
telemetry:
  metrics:
    kind: otlp
    url: http://localhost:9009/otlp/v1/metrics
    headers:
      Authorization: env:OTLP_METRICS_AUTH
```

### Tracing

Tracing supports OTLP exporter config:

```yaml
telemetry:
  tracer:
    kind: otlp
    url: http://localhost:4318/v1/traces
    headers:
      Authorization: env:OTLP_TRACES_AUTH
```

Note:

- Current tracer wiring exports via OTLP/HTTP when tracer config is present.

### Telemetry libraries used

- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime>
- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/host>
- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp>
- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc>
- <https://github.com/redis/go-redis/tree/master/extra/redisotel>
- <https://github.com/XSAM/otelsql>

### Telemetry Dependencies

![Dependencies](./assets/telemetry.png)

---

## Tokens

Token configuration is rooted at `token.Config`, usually nested under transport config as `transport.http.token` and/or `transport.grpc.token` (via the shared server-side transport config).

Supported token `kind` values:

- `jwt`
- `paseto`
- `ssh`

### Access control (Casbin)

Access control is configured inside transport token config:

```yaml
transport:
  http:
    token:
      access:
        policy: ./config/rbac.csv
```

The model is based on Casbin RBAC:
<https://github.com/casbin/casbin/blob/master/examples/rbac_model.conf>

Note:

- `access.policy` is passed directly to Casbin's file adapter. Use a real file path, or pre-resolve any source string/literal policy handling in your own wiring before constructing the controller.

### JWT

JWT config:

```yaml
transport:
  http:
    token:
      kind: jwt
      jwt:
        iss: my-service
        exp: 1h
        kid: my-key-id
```

Important behavior:

- JWT verification requires the `kid` header to exist and match `kid` in config exactly.
- `exp` is parsed as a Go duration string; invalid values can fail fast.

### Paseto

Paseto config:

```yaml
transport:
  http:
    token:
      kind: paseto
      paseto:
        iss: my-service
        exp: 1h
```

Note:

- The current PASETO implementation issues **v4 public** tokens using Ed25519 key material provided via wiring (not directly from `paseto.secret`). If you want config-driven key material, load it via the crypto subsystem and wire signer/verifier appropriately.

### SSH tokens

SSH token verification keys are name-addressable and support rotation.

Verification-only example:

```yaml
transport:
  http:
    token:
      kind: ssh
      ssh:
        keys:
          - name: active
            public: file:/keys/active.pub
```

Signing + verification example:

```yaml
transport:
  http:
    token:
      kind: ssh
      ssh:
        key:
          name: active
          private: file:/keys/active
        keys:
          - name: active
            public: file:/keys/active.pub
          - name: old
            public: file:/keys/old.pub
```

Notes:

- `ssh.key` is used for minting tokens (requires private key).
- `ssh.keys` is used for verification (public keys).
- The config does not enforce that the signing key name exists in the verification set; include it if you want round-trip.

---

## Limiter

Limiter config is `limiter.Config` and is typically applied at transport level.

Supported key kinds (built-in):

- `user-agent`
- `ip`
- `token`

Example:

```yaml
transport:
  http:
    limiter:
      kind: user-agent
      tokens: 10
      interval: 1s
```

Note:

- `interval` is parsed as a Go duration string. Invalid values can fail fast.

---

## Time (network time)

Time config:

```yaml
time:
  kind: nts
  address: time.cloudflare.com
```

Supported kinds:

- `ntp`
- `nts`

---

## Transport

The transport layer provides higher-level wiring and middleware policy for communication in/out of the service.

At a high level:

- `transport/...` contains the opinionated service transport layer: Fx wiring, composed HTTP/gRPC server and client stacks, retries, breakers, token middleware, health wiring, and related policy.
- `net/...` contains lower-level protocol helpers and reusable primitives such as `net/http`, `net/grpc`, `net/http/meta`, `net/grpc/meta`, `net/http/strings`, `net/grpc/strings`, `net/grpc/health`, `net/header`, and `net/server`.

Supported stacks include:

- gRPC (<https://grpc.io/>)
- HTTP REST abstraction (`net/http/rest`) using content negotiation
- HTTP RPC abstraction (`net/http/rpc`) using content negotiation
- HTTP MVC helpers (`net/http/mvc`)
- CloudEvents (<https://github.com/cloudevents/sdk-go>)

### HTTP content types

The HTTP REST and RPC helpers resolve encoders from the request `Content-Type`.

Built-in text/object payload media types include:

- `application/json`
- `application/hjson`
- `application/yaml`
- `application/yml`
- `application/toml`
- `application/gob`

Built-in protobuf-oriented media type aliases include:

- `application/proto`
- `application/protobuf`
- `application/protojson`
- `application/prototext`

Notes:

- `application/hjson` maps to the built-in `hjson` encoder kind.
- Unknown or invalid request media types fall back to JSON selection.
- `text/error; charset=utf-8` is reserved for error responses and should not be sent by clients as a request content type.

### Transport configuration (servers)

Transport config root is `transport.Config`:

- `transport.http` embeds `config/server.Config`
- `transport.grpc` embeds `config/server.Config`

Minimal example:

```yaml
transport:
  http:
    address: tcp://localhost:8000
    timeout: 10s
  grpc:
    address: tcp://localhost:9000
    timeout: 10s
```

Notes:

- Address format should be `<network>://<address>` (for example `tcp://:8000`).
- If address is omitted, defaults are `tcp://:8080` (HTTP) and `tcp://:9090` (gRPC).

With retry + low-level options map:

```yaml
transport:
  http:
    address: tcp://localhost:8000
    timeout: 10s
    retry:
      backoff: 100ms
      timeout: 1s
      attempts: 3
    options:
      read_timeout: 10s
      write_timeout: 10s
      idle_timeout: 10s
      read_header_timeout: 10s
  grpc:
    address: tcp://localhost:9000
    timeout: 10s
    retry:
      backoff: 100ms
      timeout: 1s
      attempts: 3
    options:
      keepalive_enforcement_policy_ping_min_time: 10s
      keepalive_max_connection_idle: 10s
      keepalive_max_connection_age: 10s
      keepalive_max_connection_age_grace: 10s
      keepalive_ping_time: 10s
```

### TLS for transports

TLS config uses `crypto/tls.Config` and fields are source strings:

```yaml
transport:
  http:
    tls:
      cert: file:test/certs/cert.pem
      key: file:test/certs/key.pem
  grpc:
    tls:
      cert: file:test/certs/cert.pem
      key: file:test/certs/key.pem
```

Important gotcha:

- Some transport packages require that you call a package `Register(...)` function to provide an `os.FS` used to read key material. If you enable TLS and have not registered the FS, TLS construction may not have access to the filesystem.
- If you are wiring server lifecycle manually, use `net/server.Register(...)`.

### Transport Dependencies

![Dependencies](./assets/transport.png)

### Circuit breakers (client-side)

The transport client wrappers include optional circuit breakers:

- HTTP breaker (`transport/http/breaker`):
  - Scope is per `"<METHOD> <HOST>"`.
  - Default failure statuses are `>=500` and `429`.
  - Transport errors are counted as failures.
  - Failure status responses are still returned to callers (while breaker accounting records a failure).

- gRPC breaker (`transport/grpc/breaker`):
  - Scope is per `fullMethod`.
  - Default failure codes are `Unavailable`, `DeadlineExceeded`, `ResourceExhausted`, and `Internal`.
  - Errors with other gRPC codes are treated as successful for breaker accounting.

---

## Cryptography

The crypto root config is `crypto.Config` and supports multiple key types. Most fields are source strings.

Example:

```yaml
crypto:
  aes:
    key: file:test/secrets/aes
  ed25519:
    public: file:test/secrets/ed25519_public
    private: file:test/secrets/ed25519_private
  hmac:
    key: file:test/secrets/hmac
  rsa:
    public: file:test/secrets/rsa_public
    private: file:test/secrets/rsa_private
  ssh:
    public: file:test/secrets/ssh_public
    private: file:test/secrets/ssh_private
```

Notes:

- AES keys must be 16/24/32 bytes after resolving the source string.
- RSA keys expect PKCS#1 PEM blocks (`RSA PUBLIC KEY` / `RSA PRIVATE KEY`).
- Ed25519 expects PKIX `PUBLIC KEY` and PKCS#8 `PRIVATE KEY` PEM blocks.

### Crypto Dependencies

![Dependencies](./assets/crypto.png)

---

## Debug endpoints

Debug server config:

```yaml
debug:
  address: tcp://localhost:6060
  timeout: 10s
```

Enable TLS:

```yaml
debug:
  tls:
    cert: file:test/certs/cert.pem
    key: file:test/certs/key.pem
```

All debug endpoints are namespaced by service name: `/<name>/debug/...`.

### statsviz

```http
GET http://localhost:6060/<name>/debug/statsviz
```

<https://github.com/arl/statsviz>

### pprof

```http
GET http://localhost:6060/<name>/debug/pprof/
GET http://localhost:6060/<name>/debug/pprof/cmdline
GET http://localhost:6060/<name>/debug/pprof/profile
GET http://localhost:6060/<name>/debug/pprof/symbol
GET http://localhost:6060/<name>/debug/pprof/trace
```

<https://pkg.go.dev/net/http/pprof>

### fgprof

```http
GET http://localhost:6060/<name>/debug/fgprof?seconds=10
```

<https://pkg.go.dev/github.com/felixge/fgprof>

### gopsutil

```http
GET http://localhost:6060/<name>/debug/psutil
```

<https://github.com/shirou/gopsutil>

---

## Development

### Style

This repo generally follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

### Development Dependencies

For local TLS fixtures:

- <https://github.com/FiloSottile/mkcert>

### Setup (repo)

This repo uses a `bin/` git submodule for `make` targets.

```sh
git submodule sync
git submodule update --init

mkcert -install
make create-certs

make dep
```

If submodule fetch fails, ensure GitHub SSH access is configured (`.gitmodules` uses `git@github.com:...` URLs).

### Discover targets

```sh
make help
```

### Dependencies (`vendor/` workflow)

```sh
make dep
```

`make dep` runs:

- `go mod download`
- `go mod tidy`
- `go mod vendor`

Tests are run with `-mod vendor`, so after dependency changes run `make dep` before `make specs`.

### Local integration dependencies

Start required services:

```sh
make start
```

Stop them:

```sh
make stop
```

### Tests

Run unit tests with race + coverage:

```sh
make specs
```

Artifacts:

- JUnit XML: `test/reports/specs.xml`
- Coverage profile: `test/reports/profile.cov`

### Lint and format

```sh
make lint
make fix-lint
make format
```

### Security checks

```sh
make sec
```

### Benchmarks

```sh
make benchmarks
make http-benchmarks
make grpc-benchmarks
make bytes-benchmarks
make strings-benchmarks
```

### Coverage reports

```sh
make coverage
make html-coverage
make func-coverage
```

### Code generation (Buf)

```sh
make generate
```

### Architecture diagrams

```sh
make diagrams
make crypto-diagram
make database-diagram
make telemetry-diagram
make transport-diagram
```

### Documentation

All exported identifiers should have GoDoc comments, and each comment should start with the identifier name (or `Deprecated:`).

### Additional gotchas

- `make kind=status encode-config` uses `base64 -w 0` (GNU style). On macOS/BSD use `base64 | tr -d '\n'`.
- If you enable transport TLS and wire transports manually (without transport modules), call:
  - `transport/http.Register(fs)`
  - `transport/grpc.Register(fs)`
- Shared metadata and header helpers live under `net/...`, for example:
  - `net/http/meta`
  - `net/grpc/meta`
  - `net/header`
  - `net/server.Register`
