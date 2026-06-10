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

  The repository helper `make kind=status encode-config` uses GNU `base64 -w 0`; on macOS/BSD, use `base64 | tr -d '\n'` for the equivalent single-line payload.

- Otherwise (no `file:`/`env:` prefix), the decoder falls back to **default lookup**, searching for:

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
- `file:/path/to/thing` → read from filesystem
- otherwise → treat as literal string

This is used for secrets and key material (TLS keys, HMAC keys, webhook secrets, SQL DSNs, etc).

Example:

```yaml
hooks:
  secret: env:WEBHOOK_SECRET
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
  options:
    url: env:CACHE_URL
```

> [!NOTE]
> - Built-in driver kinds in this repo are `redis` and `sync`.
> - Unknown `kind` values return `cache/driver.ErrNotFound`.
> - Unknown or empty `compressor` values fall back to `none`.
> - For normal values, unknown or empty `encoder` values fall back to `json`.
> - Cache operations use `plain` for `io.WriterTo`/`io.ReaderFrom` stream values and `proto` for protobuf messages, regardless of the configured `encoder`.
> - `max_size` limits encoded cache values before compression, after compression, and after decompression. A zero value uses the default `4MB`.
> - `options` is backend-specific and decoded as `map[string]any`.

---

## 🚩 Feature flags (OpenFeature)

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
    ca: file:test/certs/rootCA.pem
    server_name: localhost
```

> [!NOTE]
> - Presence enables the feature subsystem configuration-wise, but this repository does not construct a built-in OpenFeature provider from this config.
> - Services that need a remote or custom provider should use `feature.Config` in their own provider constructor and provide the resulting `openfeature.FeatureProvider` in DI; `feature.Module` registers that supplied provider with the OpenFeature SDK lifecycle.

---

## 🪝 Webhooks (Standard Webhooks)

Configured via `hooks.Config`:

```yaml
hooks:
  secret: file:test/secrets/hooks
```

`secret` is a source string.

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

---

## 🚀 Runtime enhancements

Server commands created through `cli.Application.AddServer` include `runtime.Module`, which currently enables:

- [automemlimit](https://github.com/KimMachineGun/automemlimit)

> [!NOTE]
> This registration is best-effort and does not fail startup if a memory limit cannot be applied. Direct Fx compositions and client-style commands should include `runtime.Module` explicitly when they want this behavior.

---

## 🐘 SQL (Postgres)

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

## 🩺 Health

Health checks are based on [go-health](https://github.com/alexfalkowski/go-health).

The framework provides Kubernetes-style endpoints:

- `/<name>/healthz` — general serving health status
- `/<name>/livez` — liveness probe
- `/<name>/readyz` — readiness probe

Successful health responses return HTTP 200 with the plain-text body `SERVING`.
Missing or failing observers return HTTP 503 with the standard go-service error response.

These are modeled after [Kubernetes API health endpoints](https://kubernetes.io/docs/reference/using-api/health-checks/).

---

## 📡 Telemetry

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

> [!NOTE]
> - `headers` values are source strings.
> - Telemetry header maps are resolved during config projection; unset `env:` values and unreadable `file:` values fail fast (panic during startup).

> [!WARNING]
> OTLP exporters reject non-loopback `http://` endpoints when headers are configured. Use HTTPS for remote collectors that require authorization headers; cleartext with headers is accepted only for local loopback endpoints.

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

> [!NOTE]
> Current tracer wiring exports via OTLP/HTTP when `telemetry.tracer.kind` is `otlp` and `url` is configured.

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

Access control is configured inside transport token config:

```yaml
transport:
  http:
    token:
      access:
        model: file:./config/rbac.conf
        policy: file:./config/rbac.csv
```

The model is based on Casbin RBAC:
<https://github.com/casbin/casbin/blob/master/examples/rbac_model.conf>

> [!NOTE]
> `access.model` and `access.policy` are resolved through `os.FS.ReadSource`; use `file:` for files, `env:` for environment-provided content, or literal content.

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

> [!IMPORTANT]
> JWT generation and verification use Ed25519 key material from `jwt.keys`. Keep private key material only on services that mint tokens; verifiers only need public keys.

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
        key: active
        keys:
          active:
            public: file:/keys/ed25519.pub
            private: file:/keys/ed25519
          old:
            public: file:/keys/ed25519-old.pub
```

> [!NOTE]
> The PASETO implementation issues **v4 public** tokens. Generation signs with `paseto.key`, writes that id as footer `kid`, and verification selects the public key from `paseto.keys`.

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
```

> [!NOTE]
> - `interval` is parsed as a Go duration string. Invalid values can fail fast.
> - The built-in limiter is an in-memory, per-process safeguard. Use it as a last resort and prefer an external edge, gateway, ingress, load balancer, or service-mesh limiter for production abuse protection.
> - The `user-id` key uses the verified principal stored in metadata. For JWT/PASETO tokens this is the subject claim; for SSH tokens this is the verified key name. Prefer it when authenticated identity is available.
> - The `transport-service-method` key prefixes the service-method value with the transport name, such as `http:GET /users/{id}` or `grpc:/users.v1.Users/Get`, so HTTP and gRPC operations use separate buckets.
> - The `service-method` key uses HTTP route/path metadata or the gRPC full method name. Prefer `transport-service-method` unless cross-transport operations intentionally share quota.
> - Server-side HTTP and gRPC limiters run after metadata extraction and token verification, so missing, malformed, or invalid authorization is rejected before it reaches the limiter. This is intentional; enforce quotas for those attempts with an external edge, gateway, ingress, load balancer, or service-mesh limiter.
> - gRPC stream limiters consume one token when the stream opens and one token for each `RecvMsg` and `SendMsg` operation. Unary HTTP and gRPC requests consume one token per request/RPC.

---

## 🕒 Time (network time)

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

### HTTP content types

The HTTP REST and RPC helpers resolve encoders from the request `Content-Type`, falling back to the first `Accept` media type when `Content-Type` is absent.

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
> - RSA keys expect PKCS#1 PEM blocks (`RSA PUBLIC KEY` / `RSA PRIVATE KEY`) and must be at least 4096 bits.
> - Ed25519 expects PKIX `PUBLIC KEY` and PKCS#8 `PRIVATE KEY` PEM blocks.

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

<https://github.com/shirou/gopsutil>

---

## 🧑‍💻 Development

### Style

This repo generally follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

Exported Go identifiers should have GoDoc comments, and each comment should start with the identifier name or `Deprecated:`.

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
make sql-benchmarks
make cache-benchmarks
make bytes-benchmarks
make strings-benchmarks
make id-benchmarks
make http-content-benchmarks
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
