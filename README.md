![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/go-service.svg?style=shield)](https://circleci.com/gh/alexfalkowski/go-service)
[![codecov](https://codecov.io/gh/alexfalkowski/go-service/graph/badge.svg?token=AGP01JOTM0)](https://codecov.io/gh/alexfalkowski/go-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/go-service/v2)](https://goreportcard.com/report/github.com/alexfalkowski/go-service/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/go-service/v2.svg)](https://pkg.go.dev/github.com/alexfalkowski/go-service/v2)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# 🧰 Go Service

`github.com/alexfalkowski/go-service/v2` is an opinionated framework/library for building Go services with consistent wiring for configuration, DI, transports, telemetry, crypto, etc.

This repo is primarily a **library of packages** (no top-level `cmd/` binary). Services built on top typically define their own `main` package elsewhere and import this module.

Most services are expected to be bootstrapped from [`go-service-template`](https://github.com/alexfalkowski/go-service-template) and to compose the high-level module bundles from this repository. That is the primary supported path. Lower-level package-by-package composition is still available, but it is an advanced mode and may require extra manual registration.

---

## 🚀 Install

For a new service, start from `go-service-template` so the application `main`, command wiring, configuration fixtures, and standard module composition are generated together.

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
        serve.AddConfig("file:./config.yml") // adds the `-config` / `-c` config flag with this default
    })

    os.Exit(app.RunCode(context.Background()))
}
```

The `file:./config.yml` default above expects a non-empty config file. A minimal
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

> This repo intentionally does not ship a ready-to-run `main` — it provides the building blocks. In normal usage those building blocks are consumed through `go-service-template` plus `module.Server` / `module.Client`, not by wiring every subsystem manually.

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
  Read config from a file at `<path>`; parser is selected from the file extension (`.json`, `.hjson`, `.yaml`, `.yml`, `.toml`).

- `env:<ENV_VAR>`
  Read config from env var `<ENV_VAR>`. The env var value must be formatted as:

  `"<extension>:<base64-content>"`

  Example format: `yaml:ZW52aXJvbm1lbnQ6IGRldmVsb3BtZW50Cg==`

  Example commands:

  ```sh
  # Linux (GNU base64)
  export SERVICE_CONFIG="yaml:$(base64 -w 0 < ./config.yml)"
  ./your-service serve -config env:SERVICE_CONFIG
  ```

  ```sh
  # macOS/BSD base64
  export SERVICE_CONFIG="yaml:$(base64 < ./config.yml | tr -d '\n')"
  ./your-service serve -c env:SERVICE_CONFIG
  ```

  HJSON works the same way, for example `hjson:<base64-content>`.

  The repository helper `make kind=configs/config encode-config` uses GNU `base64 -w 0`; on macOS/BSD, use `base64 | tr -d '\n'` for the equivalent single-line payload.

- Unsupported explicit `kind:location` prefixes fail startup instead of falling back to another source.

- Unprefixed values, including an empty value, fall back to **default lookup**, searching for:

  `<serviceName>.{yaml,yml,hjson,toml,json}`

  Default lookup checks extensions first (`.yaml`, `.yml`, `.hjson`, `.toml`, `.json`), and for each extension checks:
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

## 🔐 Source strings (secrets, DSNs, paths)

Many fields accept a *source string* rather than only a literal:

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

Encoding kinds used by subsystems that support encoding:

- `json`
- `hjson`
- `toml`
- `yaml`, `yml`
- `msgpack`
- `proto`, `protobuf`, `pb`, `protobin`, `pbbin`
- `protojson`, `pbjson`
- `prototext`, `prototxt`, `pbtxt`
- `gob`
- `plain`, `octet-stream`, `markdown`

> [!NOTE]
> - `plain`, `octet-stream`, and `markdown` all map to the bytes passthrough encoder.
> - Protobuf binary/text/JSON kinds have multiple aliases; the list above reflects the built-in registry.

---

## 💾 Cache

Cache configuration is defined in `cache/config.Config`:

```yaml
cache:
  kind: redis
  compressor: zstd
  encoder: json
  max_size: 4MB
  max_entries: 1024
  options:
    url: env:CACHE_URL
```

> [!NOTE]
> - Built-in driver kinds in this repo are `redis` and `ttlcache`.
> - Unknown `kind` values return `cache/driver/errors.ErrNotFound`.
> - Unknown or empty `compressor` values fall back to `none`.
> - For normal values, unknown or empty `encoder` values fall back to `json`.
> - Configured `compressor` and `encoder` values are part of the cache driver key namespace, so changing either setting creates cache misses for values written with the previous format.
> - Cache operations use `plain` for `io.WriterTo`/`io.ReaderFrom` stream values and `proto` for protobuf messages, regardless of the configured `encoder`.
> - `max_size` limits encoded cache values before compression, after compression, and after decompression. A zero value uses the default `4MB`.
> - `max_entries` limits entries retained by bounded in-memory cache drivers. A zero value uses the default `1024`; negative values are invalid.
> - `options` is backend-specific and decoded as `map[string]any`.
> - Configure each cache backend for a specific service or purpose. For Redis, use a dedicated database, endpoint, or deployment-level key namespace in the connection/configuration instead of sharing one general cache for unrelated data.

> [!WARNING]
> `Cache.Flush` follows backend semantics; for Redis it clears the selected database.

---

## 🚩 Feature flags (OpenFeature)

The `feature.Config` embeds client-side config (`config/client.Config`), so it supports:

- `address`
- `timeout`
- `retry`
- `breaker`
- `limiter`
- `tls`
- `token`
- `options`

Example:

```yaml
feature:
  address: localhost:9000
  timeout: 10s
  breaker:
    max_requests: 2
    interval: 15s
    timeout: 5s
    consecutive_failures: 4
  retry:
    backoff: 100ms
    timeout: 1s
    attempts: 3
  tls:
    cert: file:test/certs/client-cert.pem
    key: file:test/certs/client-key.pem
    ca: file:test/certs/rootCA.pem
    server_name: localhost
```

> [!NOTE]
> - `feature.Config` embeds client config; `IsEnabled` is true only when both the feature config and embedded client config are present. An empty `feature:` block is treated as disabled by feature config helpers.
> - This repository does not construct a built-in OpenFeature provider from this config.
> - Services that need a remote or custom provider should use `feature.Config` in their own provider constructor and provide the resulting `openfeature.FeatureProvider` in DI; `feature.Module` registers that supplied provider with the OpenFeature SDK lifecycle.

---

## 🪝 Webhooks (Standard Webhooks)

Configured via `hooks.Config`:

```yaml
hooks:
  key: current
  secrets:
    current: env:WEBHOOK_SECRET_CURRENT
    previous: env:WEBHOOK_SECRET_PREVIOUS
```

Each `secrets` value is a source string. The resolved value must be accepted by the
Standard Webhooks library, such as a secret generated by `hooks.Generator` with
or without the `whsec_` prefix. Empty resolved secrets fail startup.

Signing uses the active `key`. Verification accepts signatures from every
configured secret, trying the active secret first. Standard Webhooks includes a
message id (`Webhook-Id`) but not a signing key id, so go-service does not extend
the protocol with a custom selector header.

Inbound verification checks Standard Webhooks signatures and timestamps, but
does not store or reject previously seen webhook ids. Receivers that perform
non-idempotent work should deduplicate or process idempotently using
`Webhook-Id` or the event id, backed by durable shared storage when running more
than one receiver instance.

> [!IMPORTANT]
> Webhook-protected CloudEvents must use structured HTTP encoding. Binary-mode
> CloudEvents with `ce-*` headers are rejected before signature verification.

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
> This registration is best-effort and does not fail startup if a memory limit cannot be applied. Direct Fx compositions and client-style commands should include `runtime.Module` explicitly when they want this behavior.

---

## 🐘 SQL (Postgres)

SQL root config is `database/sql.Config`, with Postgres under `sql.pg`.

Postgres config embeds common pool + DSN config (`database/sql/config.Config`), including writer/reader pools. Each role pool owns its `dsns` and `settings`. Enabled SQL pool settings must set positive `max_open_conns` and `max_idle_conns`; `max_idle_conns` must not exceed `max_open_conns`.

`module.Server` and `module.Client` both include `sql.Module`, which currently wires PostgreSQL support via `database/sql/pg.Module`.

Enablement is presence-based: a nil `sql` block or a nil `sql.pg` block disables SQL wiring. When enabled, the pgx stdlib driver is registered under the name `pg`, and reader/writer DSNs are resolved using the source-string rules described above. Enabled PostgreSQL config must provide at least one non-empty `reader.dsns[].url` or `writer.dsns[].url`. Driver instrumentation is installed when tracing or metrics are enabled, OpenTelemetry `database/sql` stats metrics are registered when metrics are enabled, and the resulting pools are closed on lifecycle stop.

SQL wiring creates `database/sql` pool handles and applies pool settings, but it
does not ping PostgreSQL during construction. Call `DBs.Ping`,
`DBs.PingWriter`, `DBs.PingReader`, or register `health/checker.NewDBChecker`
when startup or readiness should verify database reachability.

Example (with source strings for DSNs):

```yaml
sql:
  pg:
    reader:
      dsns:
        - url: env:PG_READER_DSN
      settings:
        max_open_conns: 20
        max_idle_conns: 10
        conn_max_idle_time: 30m
        conn_max_lifetime: 1h
    writer:
      dsns:
        - url: env:PG_WRITER_DSN
      settings:
        max_open_conns: 3
        max_idle_conns: 2
        conn_max_idle_time: 10m
        conn_max_lifetime: 30m
```

Example (literal DSN; not recommended for production secrets):

```yaml
sql:
  pg:
    writer:
      dsns:
        - url: postgres://user:pass@localhost:5432/dbname?sslmode=disable
      settings:
        max_open_conns: 10
        max_idle_conns: 5
```

### Dependencies

![Dependencies](./assets/database.png)

---

## 🩺 Health

Health checks are based on [go-health](https://github.com/alexfalkowski/go-health).

The framework provides Kubernetes-style endpoints:

- `/<name>/healthz` — general serving health status
- `/<name>/livez` — liveness probe
- `/<name>/readyz` — readiness probe

Successful health responses return HTTP 200 with the plain-text body `SERVING`.
Missing or failing observers return HTTP 503 with the standard go-service error response.
During server shutdown, `/readyz` also returns HTTP 503 after the lifecycle starts draining so
orchestrators can stop sending new traffic before the listener fully stops.

Built-in checker helpers under `health/checker` include DB connectivity checks and
cache connectivity checks for pingable cache drivers such as Redis and ttlcache.

When gRPC transport is enabled, `transport/grpc/health` registers the standard
`grpc.health.v1.Health` service on the gRPC server. Named checks use the service
name as the request `service`; an empty service checks overall gRPC health:

```sh
grpcurl -plaintext -d '{"service":"<name>"}' localhost:9000 grpc.health.v1.Health/Check
```

`Check` returns `SERVING` or `NOT_SERVING` for known services and `NotFound` for
unknown services. `List` returns the current statuses for registered services.
`Watch` streams status changes until the client cancels; unknown services stream
`SERVICE_UNKNOWN`. Health operation RPCs bypass token verification. Unary
`Check` and `List` also bypass unary server-side limiting, while health `Watch`
is a stream and still uses stream limiting.

These are modeled after [Kubernetes API health endpoints](https://kubernetes.io/docs/reference/using-api/health-checks/).

---

## 📡 Telemetry

Telemetry config root is `telemetry.Config`:

```yaml
telemetry:
  attributes:
    k8s.namespace.name: payments
    service.instance.id: instance-1
  logger: ...
  metrics: ...
  propagation: ...
  tracer: ...
```

`attributes` are plain OpenTelemetry resource labels attached to logs, metrics,
and traces. They are not source strings. Fixed go-service identity attributes
such as `host.id`, `service.name`, `service.version`, and
`deployment.environment.name` take precedence if the same key is configured.

### Propagation

OpenTelemetry context propagation defaults to W3C Trace Context plus W3C Baggage
for extraction and injection:

```yaml
telemetry:
  propagation:
    formats:
      - tracecontext
      - baggage
```

Mixed tracing estates can enable additional formats:

```yaml
telemetry:
  propagation:
    formats:
      - tracecontext
      - baggage
      - b3
```

Supported propagators are `tracecontext`, `baggage`, `b3`, `b3multi`, and
`none`. Use `none` only as the sole value for `formats`.

B3 uses the upstream B3 propagator, which supports both single-header and
multi-header B3 formats.

### Logging

Logging uses `log/slog`.

Supported built-in logger kinds:

- `json`
- `text`
- `tint`
- `otlp`

Supported logger levels are `debug`, `info`, `warn`, and `error`. When `level`
is unset, logging defaults to `info`; unknown values fail logger construction.

The standard telemetry module writes OpenTelemetry SDK and exporter failures as
JSON to stdout through a handler-owned logger that is independent of the
configured application logger. This keeps OTLP outage diagnostics local without
duplicating normal application logs or feeding failures back into the OTLP log
exporter.

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
    protocol: http
    url: http://localhost:4318/v1/logs
    headers:
      Authorization: env:OTLP_LOGS_AUTH
```

> [!NOTE]
> - `headers` values are source strings.
> - Telemetry header maps are resolved during config projection; unset `env:` values and unreadable `file:` values fail fast (panic during startup).

> [!WARNING]
> OTLP exporters reject non-loopback `http://` endpoints when headers are configured. Use HTTPS for remote collectors that require authorization headers; cleartext with headers is accepted only for local loopback endpoints.
>
> OTLP/gRPC exporters use `protocol: grpc` and a `host:port` endpoint such as `localhost:4317`. Header-bearing remote gRPC endpoints require the signal's `tls` config; loopback gRPC endpoints may still use cleartext.
>
> OTLP exporter endpoints must be set in go-service config fields such as `telemetry.logger.url`, `telemetry.metrics.url`, and `telemetry.tracer.url`. Standard OpenTelemetry endpoint environment variables such as `OTEL_EXPORTER_OTLP_ENDPOINT` are not used as fallback sources.

Remote OTLP/gRPC exporters can use the same TLS source-string model as other go-service clients:

```yaml
telemetry:
  tracer:
    kind: otlp
    protocol: grpc
    url: collector.example.com:4317
    tls:
      ca: file:/etc/otel/ca.pem
      cert: file:/etc/otel/client.crt
      key: file:/etc/otel/client.key
      server_name: collector.example.com
    headers:
      Authorization: env:OTLP_TRACES_AUTH
```

Use the same `tls` shape under `telemetry.logger` or `telemetry.metrics` when those signals export through OTLP/gRPC.

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
    protocol: http
    url: http://localhost:9009/otlp/v1/metrics
    interval: 30s
    timeout: 5s
    headers:
      Authorization: env:OTLP_METRICS_AUTH
```

`interval` and `timeout` apply only to OTLP push metrics. When either value is
unset or zero, the OpenTelemetry SDK default is used.

#### Histogram buckets

Override the default histogram bucket boundaries per instrument with
`telemetry.metrics.views`, keyed by instrument name (OpenTelemetry name matching,
including `*` wildcards):

```yaml
telemetry:
  metrics:
    views:
      http.server.request.duration: [0.005, 0.01, 0.05, 0.1, 0.5, 1, 5]
      "rpc.*.duration": [0.01, 0.1, 1]
```

Boundaries are in the instrument's unit (seconds for duration histograms, bytes
for size histograms) and must be listed in increasing order. Views apply to
histogram instruments regardless of metrics kind; an unset or empty map keeps the
OpenTelemetry SDK default buckets.

### Tracing

Tracing supports OTLP exporter config:

```yaml
telemetry:
  tracer:
    kind: otlp
    protocol: http
    url: http://localhost:4318/v1/traces
    sampler:
      kind: ratio
      ratio: 0.25
    headers:
      Authorization: env:OTLP_TRACES_AUTH
```

> [!NOTE]
> OTLP exporters default to `protocol: http`. Set `protocol: grpc` and use a
> `host:port` `url`, such as `localhost:4317`, to export through OTLP/gRPC.
>
> Supported sampler kinds:
>
> - `always_on`: record every trace.
> - `always_off`: drop every trace.
> - `ratio`: follow an incoming parent span's sampled decision when the request
>   already has trace context; otherwise record the configured fraction of new
>   root traces. Set `ratio` between `0` and `1`, where `0` drops new root
>   traces and `1` records all new root traces.
>
> When `sampler` is omitted, go-service preserves the OpenTelemetry SDK default
> sampler and SDK sampler environment handling.

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

## 🎫 Tokens

Token configuration is rooted at `token.Config`, usually nested under transport config as `transport.http.token` and/or `transport.grpc.token` (via the shared server-side transport config).

Supported token `kind` values:

- `jwt`
- `paseto`
- `ssh`

### Access control (Casbin)

Access control is configured once at the transport level and shared by all
enabled HTTP and gRPC server stacks:

```yaml
transport:
  access:
    model: file:./config/rbac.conf
    policy: file:./config/rbac.csv
```

When `access` is configured, the standard HTTP and gRPC server stacks enforce
the policy after token authentication and before application handlers run. Omit
`access` to leave transport authorization disabled.

The model is based on Casbin RBAC:
<https://github.com/casbin/casbin/blob/master/examples/rbac_model.conf>

Example `rbac.conf`:

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

Policies use the verified user id as `sub`, `meta.TransportServiceMethod` as
`obj`, and `invoke` as `act`. Example `rbac.csv`:

```csv
p, reader, http:GET /users/{id}, invoke
p, writer, http:POST /users, invoke
p, greeter, grpc:/greet.v1.GreeterService/SayHello, invoke
g, frontend, reader
g, admin, reader
g, admin, writer
g, billing-service, greeter
```

The `p` rows define permissions and must match the model's `p = sub, obj, act`
shape, so they include `invoke`. The `g` rows define role membership and match
`g = _, _`, so they only contain `subject, role`.

For HTTP servers the object uses the matched route pattern when available, such
as `http:GET /users/{id}`. HTTP tokens are authenticated against the concrete
request method and path, such as `GET /users/123`; access policy enforcement
uses the canonical route pattern. gRPC tokens are authenticated against the full
method name, such as `/greet.v1.GreeterService/SayHello`; access policy
enforcement uses the transport service-method object, such as
`grpc:/greet.v1.GreeterService/SayHello`.

> [!NOTE]
> `access.model` and `access.policy` are resolved through `os.FS.ReadSource`; use `file:` for files, `env:` for environment-provided content, or literal content.
>
> Access config builds an injectable controller for authorization checks. The built-in HTTP and gRPC server stacks authenticate tokens, store the verified user id, and enforce the configured Casbin policy before application handlers run.

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
        leeway: 30s
        key: active
        keys:
          active:
            public: file:/keys/ed25519.pub
            private: file:/keys/ed25519
          old:
            public: file:/keys/ed25519-old.pub
```

Important behavior:

- JWT generation signs with `jwt.key`; verification requires the token `kid` header to select an entry in `jwt.keys`.
- `exp` is parsed as a Go duration string; invalid values can fail fast.
- `leeway` is optional clock-skew tolerance for verification; keep it small because it extends acceptance around `iat`/`nbf` and `exp`.

> [!IMPORTANT]
> JWT generation and verification use Ed25519 key material from `jwt.keys`. Keep private key material only on services that mint tokens; verifiers only need public keys.

All token `exp` and non-zero `leeway` values are Go duration strings and must be positive whole-second durations. Values such
as `1s`, `15m`, and `24h` validate; sub-second values such as `500ms` do not.

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
        leeway: 30s
        key: active
        keys:
          active:
            public: file:/keys/ed25519.pub
            private: file:/keys/ed25519
          old:
            public: file:/keys/ed25519-old.pub
```

> [!NOTE]
> The PASETO implementation issues **v4 public** tokens. Generation signs with `paseto.key`, writes that id as footer `kid`, and verification selects the public key from `paseto.keys`. `paseto.leeway` is optional clock-skew tolerance for verification.

### SSH tokens

SSH token verification keys are id-addressable and support rotation.

Verification-only example:

```yaml
transport:
  http:
    token:
      kind: ssh
      ssh:
        exp: 5m
        leeway: 30s
        keys:
          active:
            public: file:/keys/active.pub
```

Signing + verification example:

```yaml
transport:
  http:
    token:
      kind: ssh
      ssh:
        exp: 5m
        leeway: 30s
        key: active
        keys:
          active:
            public: file:/keys/active.pub
            private: file:/keys/active
          old:
            public: file:/keys/old.pub
```

> [!NOTE]
> - `ssh.key` is the active key id used for minting tokens (the matching `ssh.keys` entry requires private key material).
> - `ssh.keys` is the trusted key map used for verification (public keys).
> - `ssh.exp` sets the token validity window; SSH keys remain long-lived, while generated tokens are short-lived.
> - `ssh.leeway` is optional clock-skew tolerance for verification; keep it small because it extends acceptance around `iat` and `exp`.
> - SSH tokens carry `sub` equal to `kid`, so the verified subject is the trusted peer key id.

---

## 🚦 Limiter

Limiter config is `transport/limiter.Config` and is typically applied at transport level.

Supported key kinds (built-in):

- `user-id`
- `transport-service-method`
- `service-method`
- `ip`
- `user-agent`

Example:

```yaml
transport:
  http:
    limiter:
      kind: user-agent
      tokens: 10
      interval: 1s
      max_keys: 4096
```

> [!NOTE]
> - `interval` is parsed as a Go duration string. Invalid values can fail fast.
> - `tokens` and `interval` use the underlying in-memory store defaults when set to zero: `1` token per `1s`. Configure positive values for explicit quotas.
> - `max_keys` caps the number of caller-derived keys that receive independent in-memory buckets. A zero value uses the default `4096`; additional distinct keys share one overflow bucket.
> - The built-in limiter is an in-memory, per-process safeguard. Use it as a last resort and prefer an external edge, gateway, ingress, load balancer, or service-mesh limiter for production abuse protection.
> - The `user-id` key uses the verified principal stored in metadata. For JWT/PASETO tokens this is the subject claim; for SSH tokens this is the verified key name. Prefer it when authenticated identity is available.
> - The `transport-service-method` key prefixes the service-method value with the transport name, such as `http:GET /users/{id}` or `grpc:/users.v1.Users/Get`, so HTTP and gRPC operations use separate buckets.
> - The `service-method` key uses HTTP route/path metadata or the gRPC full method name. Prefer `transport-service-method` unless cross-transport operations intentionally share quota.
> - Server-side HTTP and gRPC limiters run after metadata extraction and token verification, so missing, malformed, or invalid authorization is rejected before it reaches the limiter. This is intentional; enforce quotas for those attempts with an external edge, gateway, ingress, load balancer, or service-mesh limiter.
> - Server-side HTTP limiters set `RateLimit` and `RateLimit-Policy` headers; denied HTTP requests also set `Retry-After` when reset timing is available. Server-side gRPC limiters set `ratelimit` and `ratelimit-policy` response metadata; denied gRPC requests also attach a `google.rpc.RetryInfo` detail when reset timing is available.
> - gRPC stream limiters consume one token when the stream opens and one token for each `RecvMsg` and `SendMsg` operation. Unary HTTP and gRPC requests consume one token per request/RPC.

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

CloudEvents HTTP wiring lives under `transport/http/events`: use
`NewReceiver(...).Register(...)` to receive events on a POST route and
`NewSender(...).Send(...)` with `net/http/events.ContextWithTarget(...)` to
send events. The sender uses structured HTTP encoding by default; configure
`WithSenderEncoding(SenderEncodingBinary)` for outbound integrations that
require binary-mode CloudEvents. Webhook-protected receivers require structured
encoding and reject binary-mode CloudEvents with `ce-*` headers before
signature verification. Receiver registration marks the event route as
unauthenticated for transport token/access middleware so webhook verification can
act as the event authentication boundary.

### HTTP content types

The HTTP REST and RPC helpers decode request bodies from the request `Content-Type`, falling back to JSON when `Content-Type` is absent or unknown. Response encoding uses the first `Accept` media type when present, falling back to the request `Content-Type` when `Accept` is absent. Client helpers can set `ContentType` for the request body and `Accept` for an independent response format.

Built-in text/object payload media types include:

- `application/json`
- `application/hjson`
- `application/yaml`, `application/yml`
- `application/toml`
- `application/octet-stream`, `text/plain`, `text/markdown`

Internal binary payload media types include:

- `application/vnd.msgpack`
- `application/gob`

Built-in protobuf-oriented media type aliases include:

- `application/proto`, `application/pb`, `application/protobuf`, `application/protobin`, `application/pbbin`
- `application/protojson`, `application/pbjson`
- `application/prototext`, `application/prototxt`, `application/pbtxt`

> [!NOTE]
> - `application/hjson` maps to the built-in `hjson` encoder kind.
> - Unknown or invalid request media types fall back to JSON selection.
> - `text/error` is reserved for error responses and should not be sent by clients as a request content type.
>
> `application/vnd.msgpack` and `application/gob` can be resolved as media types, but REST/RPC request-body decoding rejects them with HTTP 415.

### HTTP route misses

The HTTP transport wraps the mux with `net/http.NewNotFoundHandler` so generated 404 responses can be rendered consistently while preserving other mux responses such as 405 Method Not Allowed.

- REST/RPC-style missing routes use `net/http/content.NotFoundHandler`, which writes the standard `status.WriteError` response.
- MVC missing routes can use `net/http/mvc.NotFoundHandler` to render the registered MVC not-found view when the request accepts HTML (`Accept: text/html`) or is an HTMX request (`Hx-Request: true`).
- Routes that match and write their own status are not replaced by this mux-level not-found handler.

### HTTP MVC errors

When an MVC controller returns an error, `net/http/mvc.Route` renders the returned view with a client-safe `mvc.Error` model. The model contains the HTTP status `Code` and safe client-visible `Message`.

The raw error string remains available to templates as `mvcModelError` metadata for compatibility. Rendering that metadata can expose diagnostic details, so prefer `.Model.Message` for client-visible error pages.

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

> [!NOTE]
> - Address may use `<network>://<address>` (for example `tcp://:8000`) or a raw listen address such as `:8000`, which defaults to the `tcp` network.
> - If address is omitted, defaults are `tcp://:8080` (HTTP) and `tcp://:9090` (gRPC).
> - `transport.grpc.timeout` bounds unary RPC handlers and feeds gRPC server keepalive/connection defaults; it does not cap stream lifetime. Long-lived streams remain open until client cancellation or stream-specific controls apply.
> - `max_receive_size` limits inbound payload size. A zero value uses the default `4MB`.
> - For HTTP, `max_receive_size` applies per request body. For gRPC, it applies per inbound unary request and per inbound stream message.
> - MVC does not enforce its own body-size caps; supported HTTP server wiring applies `max_receive_size` before MVC handlers run, and go-service HTTP clients apply their configured response-size cap when reading responses.

Receive-limit example:

```yaml
transport:
  http:
    max_receive_size: 2MB
  grpc:
    max_receive_size: 3MB
```

With low-level server options:

```yaml
transport:
  http:
    address: tcp://localhost:8000
    timeout: 10s
    options:
      read_timeout: 10s
      write_timeout: 10s
      idle_timeout: 10s
      read_header_timeout: 10s
  grpc:
    address: tcp://localhost:9000
    timeout: 10s
    options:
      keepalive_enforcement_policy_ping_min_time: 10s
      keepalive_max_connection_idle: 10s
      keepalive_max_connection_age: 10s
      keepalive_max_connection_age_grace: 10s
      keepalive_ping_time: 10s
```

### TLS for transports

TLS config uses `crypto/tls/config.Config` and fields are source strings:

```yaml
transport:
  http:
    tls:
      cert: file:test/certs/cert.pem
      key: file:test/certs/key.pem
      ca: file:test/certs/rootCA.pem
  grpc:
    tls:
      cert: file:test/certs/cert.pem
      key: file:test/certs/key.pem
      ca: file:test/certs/rootCA.pem
```

Set `ca` on server TLS config to require and verify client certificates for mTLS. Set `ca` on client TLS
config to verify server certificates issued by the same local or private CA. `server_name` is only needed
on clients when the dial address differs from the certificate DNS name.

Server-side TLS requires a complete `cert` and `key` pair whenever TLS material is configured. `ca` enables
client-certificate verification for mTLS, but a CA-only server TLS config fails startup.

gRPC clients use insecure transport credentials when TLS is not configured. That default is intended for
local or platform-secured traffic; configure client TLS for calls outside that trusted boundary.

> [!IMPORTANT]
> If you are using `go-service-template` or composing server transport bundles such as `module.Server` or `transport.Module`, the required transport registration is handled for you by DI.
>
> `module.Client` does not wire transports by default. When a client process constructs HTTP or gRPC TLS config from source strings such as `file:`, call the relevant transport-level `Register(...)` functions, such as `transport/http.Register(...)` or `transport/grpc.Register(...)`.
>
> You only need to call transport-level `Register(...)` functions yourself when you intentionally wire transports manually or compose lower-level packages outside the transport module graph.
>
> If you are wiring server lifecycle manually, use `net/server.Register(...)`.

### Forwarded IPs and reflection

> [!WARNING]
> HTTP and gRPC metadata extraction intentionally trusts common forwarded IP headers/metadata such as `X-Forwarded-For`, `X-Real-IP`, `CF-Connecting-IP`, and `True-Client-IP`. Services that rely on extracted IPs for logging, policy, or rate limiting should only receive traffic through trusted edge infrastructure that strips or overwrites client-supplied forwarding headers.

> [!WARNING]
> gRPC server reflection is intentionally always registered by `net/grpc.NewServer` so internal tooling can discover services. Services that should not expose reflection publicly should restrict access with bind addresses, TLS/client authentication, ingress policy, firewall rules, or service-mesh authorization.

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

Client config uses the shared `transport/breaker.Config` shape for breaker mechanics. Any config type that
embeds `config/client.Config` has its own `breaker` block under that client config. This example uses
`feature.Config` only because it is one such client config:

```yaml
feature:
  address: localhost:9000
  breaker:
    max_requests: 2
    interval: 15s
    timeout: 5s
    consecutive_failures: 4
```

When manually constructing HTTP or gRPC clients, pass a transport-specific breaker config to
`transport/http.WithClientBreaker(...)` or `transport/grpc.WithClientBreaker(...)`. These configs
embed the shared breaker mechanics and add protocol-specific failure classification:

```go
httpBreaker := httpbreaker.NewConfig(sharedBreaker, 429, 502, 503)
grpcBreaker := grpcbreaker.NewConfig(sharedBreaker, codes.Unavailable, codes.ResourceExhausted)
```

`NewConfig` returns `nil` when the shared breaker config is `nil`, preserving client-option wiring that
disables breakers by omitting breaker config.

`max_requests` controls half-open probe concurrency. `interval` controls the
closed-state count reset window. `timeout` controls how long the breaker stays
open before allowing half-open probes. `consecutive_failures` controls when the
breaker opens. Zero values keep the package defaults.

HTTP `StatusCodes` and gRPC `Codes` are optional replacement lists for failure
classification. When omitted, the default lists above apply. When set, only the
configured values count as breaker failures, so include the defaults as well
when extending rather than replacing default behavior.

### Client retries

Client config uses the shared `transport/retry.Config` shape for retry mechanics. Any config type that embeds
`config/client.Config` has its own `retry` block under that client config. This example uses `feature.Config`
only because it is one such client config:

```yaml
feature:
  address: localhost:9000
  retry:
    timeout: 1s
    backoff: 100ms
    attempts: 3
    strategy: exponential
```

When manually constructing HTTP or gRPC clients, pass a transport-specific retry config to
`transport/http.WithClientRetry(...)` or `transport/grpc.WithClientRetry(...)`. These configs embed the
shared retry mechanics and add protocol-specific failure classification:

```go
httpRetry := httpretry.NewConfig(sharedRetry, 429, 502, 503)
grpcRetry := grpcretry.NewConfig(sharedRetry, codes.Unavailable, codes.ResourceExhausted)
```

`NewConfig` returns `nil` when the shared retry config is `nil`, preserving client-option wiring that
disables retries by omitting retry config.

`attempts` is the total number of attempts, including the initial call. A value
of `0` or `1` means no retry beyond the first attempt; values above `10` are
rejected during config validation. `backoff` is the base delay between retry
attempts.

`strategy` selects how `backoff` grows between attempts: `constant` (the
default) reuses the base delay for every wait, `exponential` doubles it on each
attempt, and `fibonacci` grows it along the Fibonacci sequence. An unset value
applies `constant`, jitter is applied on top of the chosen strategy, and any
other value is rejected during config validation.

`timeout` is transport-specific. gRPC unary retries apply it per attempt, so
total elapsed time can include multiple attempt timeouts plus backoff unless the
caller context ends first. HTTP retries do not create a retry-owned per-attempt
timeout; bound outbound HTTP calls with the request context or
`http.Client.Timeout`.

HTTP `StatusCodes` and gRPC `Codes` are optional replacement lists for failure
classification. When omitted, the default lists below apply. When set, only the
configured values are retryable, so include the defaults as well when extending
rather than replacing default behavior. HTTP values must be 4xx or 5xx status
codes. gRPC values must be non-OK `codes.Code` values.

Default retry policy is intentionally conservative:

- HTTP retries side-effect-safe methods (`GET`, `HEAD`, `OPTIONS`) or requests with a `Request-Id`.
- HTTP retries response/status failures only for `429 Too Many Requests` and `503 Service Unavailable`, plus selected transport errors classified by `retryablehttp.DefaultRetryPolicy`.
- gRPC retries AIP-style read methods named `Get*` or `List*`, or calls with a `Request-Id`.
- gRPC retries only `Unavailable` by default.

HTTP retryable responses with a valid `Retry-After` delay greater than the
minimum jittered backoff suppress another attempt and return the current
response. gRPC retryable status errors with `google.rpc RetryInfo.retry_delay`
use the same suppression policy.

`Request-Id` identifies the logical request, not an individual wire attempt.
Services that allow retried writes should treat it as the idempotency key and
deduplicate repeated attempts when duplicate processing would be unsafe.

---

## 🔑 Cryptography

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

> [!NOTE]
> - AES keys must be 16/24/32 bytes after resolving the source string.
> - HMAC keys should be high-entropy secrets and must remain private.
> - RSA keys expect PKCS#1 PEM blocks (`RSA PUBLIC KEY` / `RSA PRIVATE KEY`) and must be at least 4096 bits.
> - Ed25519 expects PKIX `PUBLIC KEY` and PKCS#8 `PRIVATE KEY` PEM blocks.
> - SSH keys must be Ed25519 SSH keys: public keys use `authorized_keys` format and private keys use SSH private key format.

AES and RSA encryption APIs accept `crypto.Message`. `Data` is encrypted or
decrypted, while `Meta` is authenticated context that must match during
decryption. AES-GCM uses `Meta` as associated data; RSA-OAEP uses it as the
OAEP label.

### Crypto Dependencies

![Dependencies](./assets/crypto.png)

---

## 🛠️ Debug endpoints

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
    ca: file:test/certs/rootCA.pem
```

Debug TLS uses the same server-side TLS contract as transports: `cert` and `key`
are required whenever TLS material is configured, and `ca` adds client-certificate
verification for mTLS.

All debug endpoints are namespaced by service name: `/<name>/debug/...`.

> [!WARNING]
> If `debug.address` is omitted while debug is enabled, the debug server binds to `tcp://:6060`. Set an explicit address, TLS/mTLS, and network or policy controls appropriate for the deployment.

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

The response is content-negotiated and contains a best-effort snapshot with `cpu`, `host`, `load`, `mem`,
and `net` sections. Individual nested values may be empty or partially populated when the platform or
runtime permissions do not expose a metric.

<https://github.com/shirou/gopsutil>

---

## 🧑‍💻 Development

### Style

This repo generally follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

Exported Go identifiers should have GoDoc comments, and each comment should start with the identifier name or `Deprecated:`.

### Development Dependencies

Common repository targets expect these tools on `PATH`:

- `make`
- `gotestsum` for `make specs`
- `fieldalignment` for `make lint`
- `golangci-lint` for full `make lint` coverage (the wrapper no-ops when it is missing)
- `govulncheck` and `trivy` for `make sec`
- `mkcert` for local TLS fixtures and `make create-certs`
- `buf` for `make generate`
- `goda` and Graphviz `dot` for `make diagrams`

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

`make start` uses the shared Docker-based environment from the sibling
`../docker` repo. It requires Docker and may require GitHub SSH access if that
sibling repo must be fetched.

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
make sql-benchmarks
make cache-benchmarks
make bytes-benchmarks
make strings-benchmarks
make id-benchmarks
make net-http-benchmarks
make http-content-benchmarks
```

### Fuzz tests

```sh
make fuzzes
make bytes-fuzz
make time-fuzz
make encoding-fuzz
make compress-fuzz
make net-fuzz
make package=encoding/json name=FuzzUnmarshal fuzztime=10s fuzz
```

### Coverage reports

```sh
make coverage
make html-coverage
make func-coverage
```

### Code generation (Buf)

Root generation targets are for the `internal/test` protobuf fixtures. After
changing those fixtures, regenerate them. To match the CI stale-output check,
run `make generate-stale` from a clean worktree, or after staging the intended
fixture and generated-file changes:

```sh
make generate
make generate-stale
```

### Architecture diagrams

```sh
make diagrams
make crypto-diagram
make database-diagram
make telemetry-diagram
make transport-diagram
```
