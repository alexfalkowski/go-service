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

---

## CLI

Services commonly expose two command shapes:

- **Server**: long-running daemon process
- **Client**: short-lived control/admin process

The framework uses [acmd](https://github.com/cristalhq/acmd). Your service’s `main` typically wires Fx modules + commands.

> This repo intentionally does not ship a ready-to-run `main` — it provides the building blocks.

---

## Configuration

### Supported config formats

The config decoder supports:

- JSON (`encoding/json`)
- TOML (`github.com/BurntSushi/toml`)
- YAML (`go.yaml.in/yaml/v3`)

### Selecting the config source (`-i` flag)

Config input is routed by a flag called `-i`:

- `file:<path>`  
  Read config from a file at `<path>`; parser is selected from the file extension (`.yaml`, `.yml`, `.toml`, `.json`).

- `env:<ENV_VAR>`  
  Read config from env var `<ENV_VAR>`. The env var value must be formatted as:

  `"<extension>:<base64-content>"`

  Example format: `yaml:ZW52aXJvbm1lbnQ6IGRldmVsb3BtZW50Cg==`

- Otherwise (no `file:`/`env:` prefix), the decoder falls back to **default lookup**, searching for:

  `<serviceName>.{yaml,yml,toml,json}`

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

### The standard top-level config shape

The canonical top-level config type is `config.Config` (in `config/config.go`). It contains:

- `debug`, `cache`, `crypto`, `feature`, `hooks`, `id`, `sql`, `telemetry`, `time`, `transport`, `environment`

Most sub-configs are optional pointers. Conventionally, `nil` means **disabled**.

---

## Source strings (secrets, DSNs, paths)

Many fields accept a *source string* rather than only a literal:

- `env:NAME` → read from environment variable `NAME`
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
- `toml`
- `yaml`
- `proto`
- `gob`

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
- `kind` is implementation-dependent (for example `redis`, `valkey`, `noop`) based on what your service wires/compiles in.
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
    cert: file:test/certs/client.pem
    key: file:test/certs/client-key.pem
```

Notes:
- Presence enables the feature subsystem configuration-wise, but you still need to register an OpenFeature provider in your service wiring.

---

## Webhooks (Standard Webhooks)

Configured via `hooks.Config`:

```yaml
hooks:
  secret: file:test/secrets/webhook_secret.txt
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

- `/<name>/healthz` — detailed dependency breakdown (verification/debugging)
- `/<name>/livez` — liveness probe
- `/<name>/readyz` — readiness probe

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

#### OTLP-ish logger (Loki push)

```yaml
telemetry:
  logger:
    kind: otlp
    level: info
    url: http://localhost:3100/loki/api/v1/push
    headers:
      Authorization: env:LOKI_AUTH
```

Notes:
- `headers` values are source strings. If a header value cannot be resolved, some header wiring can fail fast.

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

### Telemetry libraries used

- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime>
- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/host>
- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp>
- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc>
- <https://github.com/redis/go-redis/tree/master/extra/redisotel>
- <https://github.com/XSAM/otelsql>

### Dependencies

![Dependencies](./assets/telemetry.png)

---

## Tokens

Token configuration is rooted at `token.Config`, usually nested under transport config as `transport.http.token` and/or `transport.grpc.token` (via the shared server-side transport config).

Supported token `kind` values:

- `jwt`
- `paseto`
- `ssh`

### Access control (Casbin)

Optional access control is configured under:

```yaml
token:
  access:
    policy: file:/path/to/policy.csv
```

In transport configs this becomes:

```yaml
transport:
  http:
    token:
      access:
        policy: file:test/policy.csv
```

The model is based on Casbin RBAC:
<https://github.com/casbin/casbin/blob/master/examples/rbac_model.conf>

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
            public: file:test/keys/active.pub
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
          private: file:test/keys/active
        keys:
          - name: active
            public: file:test/keys/active.pub
          - name: old
            public: file:test/keys/old.pub
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

The transport layer provides common wiring and helpers for communication in/out of the service.

Supported stacks include:

- gRPC (<https://grpc.io/>)
- HTTP REST abstraction (`net/http/rest`) using content negotiation
- HTTP RPC abstraction (`net/http/rpc`) using content negotiation
- HTTP MVC helpers (`net/http/mvc`)
- CloudEvents (<https://github.com/cloudevents/sdk-go>)

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
      cert: file:test/certs/server.pem
      key: file:test/certs/server-key.pem
  grpc:
    tls:
      cert: file:test/certs/server.pem
      key: file:test/certs/server-key.pem
```

Important gotcha:
- Some transport packages require that you call a package `Register(...)` function to provide an `os.FS` used to read key material. If you enable TLS and have not registered the FS, TLS construction may not have access to the filesystem.

### Dependencies

![Dependencies](./assets/transport.png)

---

## Cryptography

The crypto root config is `crypto.Config` and supports multiple key types. Most fields are source strings.

Example:

```yaml
crypto:
  aes:
    key: env:AES_KEY
  ed25519:
    public: file:test/keys/ed25519.pub.pem
    private: file:test/keys/ed25519.priv.pem
  hmac:
    key: env:HMAC_KEY
  rsa:
    public: file:test/keys/rsa.pub.pem
    private: file:test/keys/rsa.priv.pem
  ssh:
    public: file:test/keys/id_ed25519.pub
    private: file:test/keys/id_ed25519
```

Notes:
- AES keys must be 16/24/32 bytes after resolving the source string.
- RSA keys expect PKCS#1 PEM blocks (`RSA PUBLIC KEY` / `RSA PRIVATE KEY`).
- Ed25519 expects PKIX `PUBLIC KEY` and PKCS#8 `PRIVATE KEY` PEM blocks.

### Dependencies

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
    cert: file:test/certs/debug.pem
    key: file:test/certs/debug-key.pem
```

### statsviz

```http
GET http://localhost:6060/debug/statsviz
```

<https://github.com/arl/statsviz>

### pprof

```http
GET http://localhost:6060/debug/pprof/
GET http://localhost:6060/debug/pprof/cmdline
GET http://localhost:6060/debug/pprof/profile
GET http://localhost:6060/debug/pprof/symbol
GET http://localhost:6060/debug/pprof/trace
```

<https://pkg.go.dev/net/http/pprof>

### fgprof

```http
GET http://localhost:6060/debug/fgprof?seconds=10
```

<https://pkg.go.dev/github.com/felixge/fgprof>

### gopsutil

```http
GET http://localhost:6060/debug/psutil
```

<https://github.com/shirou/gopsutil>

---

## Development

### Style

This repo generally follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

### Dependencies

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

Run unit tests + race + coverage:

```sh
make specs
```
