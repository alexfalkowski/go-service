[![CircleCI](https://circleci.com/gh/alexfalkowski/go-service.svg?style=shield)](https://circleci.com/gh/alexfalkowski/go-service)
[![codecov](https://codecov.io/gh/alexfalkowski/go-service/graph/badge.svg?token=AGP01JOTM0)](https://codecov.io/gh/alexfalkowski/go-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/go-service)](https://goreportcard.com/report/github.com/alexfalkowski/go-service)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/go-service.svg)](https://pkg.go.dev/github.com/alexfalkowski/go-service)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Go Service

A framework to build services in go. This came out of out building services over the years and what I have considered good practices in building services. Hence it is highly subjective and opinionated.

This framework [stands on the shoulder of giants](https://en.wikipedia.org/wiki/Standing_on_the_shoulders_of_giants) so we don't reinvent the wheel!

## Dependency Injection

This framework heavily relies on [DI](https://en.wikipedia.org/wiki/Dependency_injection). We have chosen to use [Uber FX](https://github.com/uber-go/fx). So there is great information online to get you up to speed.

## Commands

A service has commands that are configured using [Cobra](github.com/spf13/cobra). Each service has the following commands (you can add more):
- `Server` - This will provide your server needs.
- `Client` - This will provide your client needs.

These are configured in the main function.

## Configuration

The supported configuration kinds are as follows:
- [JSON](https://github.com/goccy/go-json)
- [TOML](https://github.com/BurntSushi/toml)
- [YAML](https://github.com/go-yaml/yaml)

The configuration can be read from multiple sources by specifying a flag called `--input` or `-i`. As per the following:
- `env:CONFIG_FILE` - Read from an env variable called `CONFIG_FILE`. This is the default if nothing is passed. The env variable can be file path or the configuration. If it is the config, we expect the format of `extension:ENV_VARIABLE`, where extension is the supported kinds and `ENV_VARIABLE` contains the contents of the config that are *base64 encoded*. **This can be overridden.**
- `file:path` - Read from the path.

The reason for this is that we want to be able to separate how configuration is retrieved. This way we can use and [application configuration system](https://github.com/alexfalkowski/konfig).

This is the [configuration](config/config.go). We will outline the config required in each section. The following configuration examples will use YAML.

## Environment

You can specify the environment of the service.

### Configuration

To configure, please specify the following:

```yaml
environment: development
```

## Compression

We support the following:
- None
- [Zstd](https://github.com/klauspost/compress/tree/master/zstd)
- [S2](https://github.com/klauspost/compress/tree/master/s2)
- [Snappy](https://github.com/klauspost/compress/tree/master/snappy)

## Encoders

We support the following:
- [JSON](https://github.com/goccy/go-json)
- [TOML](https://github.com/BurntSushi/toml)
- [YAML](https://gopkg.in/yaml.v3)
- [Proto](https://google.golang.org/protobuf/proto)
- [GOB](https://pkg.go.dev/encoding/gob)


## Caching

The framework currently supports the following caching solutions:
- [Redis Cache](https://github.com/go-redis/cache)

### Configuration

To configure, please specify the following:

```yaml
cache:
  redis:
    compressor: snappy
    encoder: proto
    addresses:
      server: localhost:6379
    url: path to url
```

### Dependencies

![Dependencies](./assets/cache.png)

## Feature

The framework supports [OpenFeature](https://openfeature.dev/).

### Configuration

To configure, please specify the following:

```yaml
feature:
  address: localhost:9000
  retry:
    backoff: 100ms
    timeout: 1s
    attempts: 3
  timeout: 10s
```

## Hooks

The framework supports [Standard Webhooks](https://www.standardwebhooks.com/).

### Configuration

To configure, please specify the following:

```yaml
hooks:
  secret: path to secret
```

## Runtime

We enhance the runtime with the following:
- [Automaxprocs](https://github.com/uber-go/automaxprocs)
- [Automemlimit](https://github.com/KimMachineGun/automemlimit)

## SQL

For SQL databases we support the following:
- [Postgres](https://github.com/jackc/pgx)

We also support master, slave combinations with the awesome [mssqlx](https://github.com/linxGnu/mssqlx).

### Configuration

To configure, please specify the following:

```yaml
sql:
  pg:
    masters:
      -
        url: path to url
    slaves:
      -
        url: path to url
    max_open_conns: 5
    max_idle_conns: 5
    conn_max_lifetime: 1h
```

### Dependencies

![Dependencies](./assets/database.png)

## Health

The health package is based on [go-health](https://github.com/alexfalkowski/go-health). This package allows us to create all sorts of ways to check external and internal systems.

We also provide ways to integrate into container integration systems. So we provide the following endpoints:
- `/healthz` - This allows us to check any external dependency and provide a breakdown of what is not functioning. This should only be used for verification.
- `/livez`: Can be used for k8s [liveness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-liveness-command).
- `/readyz`: Can be used for k8s [readiness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-readiness-probes).

This is modelled around [Kubernetes API health endpoints](https://kubernetes.io/docs/reference/using-api/health-checks/).

## Telemetry

Telemetry is broken down in the following sections:

### Logging

For logging we use [Uber Zap](https://github.com/uber-go/zap).

#### Configuration

To configure, please specify the following:

```yaml
telemetry:
  logger:
    level: info
```

### Metrics

For metrics we support the following:
- [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go)
- [Prometheus](https://github.com/prometheus/client_golang)

#### Configuration

Below is the configuration for each system.

##### Prometheus

To configure, please specify the following:

```yaml
telemetry:
  metrics:
    kind: prometheus
```

##### OTLP

To configure, please specify the following:

```yaml
telemetry:
  metrics:
    kind: otlp
    url: http://localhost:9009/otlp/v1/metrics
    headers:
      Authorization: path to key
```

### Trace

For distributed tracing we support the following:
- [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go)

#### Configuration

Below is the configuration for each system.

##### OTLP

To configure, please specify the following:

```yaml
telemetry:
  tracer:
    kind: otlp
    url: localhost:4318
    headers:
      Authorization: path to key
```

### Dependencies

![Dependencies](./assets/telemetry.png)

## Token

The framework allows you to define different token generators and verifiers. This is left up to you!

We recommend that you look at better auth strategies, such as:
- https://github.com/supertokens/supertokens-core
- https://github.com/ory/hydra

## Limiter

The framework allows you to define a [limiter](https://github.com/sethvargo/go-limiter). This will be applied to the different transports.

The different kinds are:
- [user-agent](meta/meta.go)
- [ip](meta/meta.go)
- [token](transport/grpc/security/token/token.go)

### Configuration

To configure, please specify the following:

```yaml
limiter:
  kind: user-agent
  tokens: 10
  interval: 1s
```

## Time

The framework allows you use network time services. We use:
- [ntp](https://github.com/beevik/ntp)
- [nts](https://github.com/beevik/nts)

### Configuration

To configure, please specify the following:

```yaml
time:
  kind: nts
  address: time.cloudflare.com
```

## Transport

The transport layer provides ways to abstract communication for in/out of the service. So we have the following integrations:
- [gRPC](https://grpc.io/) - The author truly believes in [IDLs](https://en.wikipedia.org/wiki/Interface_description_language).
- [REST](https://github.com/alexfalkowski/go-service/tree/master/net/http/rest) - An abstraction using [content negotiation](https://github.com/elnormous/) and the awesome [resty](https://github.com/go-resty/resty).
- [RPC](https://github.com/alexfalkowski/go-service/tree/master/net/http/rpc) - abstraction using [content negotiation](https://github.com/elnormous/contenttype).
- [MVC](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller) - We have a simple [framework](https://github.com/alexfalkowski/go-service/tree/master/net/http/mvc).
- [CloudEvents](https://github.com/cloudevents/sdk-go) - A specification for describing event data in a common way.

### gRPC

Below is list of the provided interceptors:
- [Limiter](https://github.com/sethvargo/go-limiter)

### REST

Below is list of the provided handlers:
- [Limiter](https://github.com/sethvargo/go-limiter)

### Configuration

To configure, please specify the following:

```yaml
transport:
  http:
    address: :8000
    retry:
      backoff: 100ms
      timeout: 1s
      attempts: 3
    timeout: 10s
  grpc:
    address: :9000
    retry:
      backoff: 100ms
      timeout: 1s
      attempts: 3
    timeout: 10s
```

If you would like to enable TLS, do the following:

```yaml
transport:
  http:
    tls:
      cert: path of cert
      key: path of key
  grpc:
    tls:
      cert: path of cert
      key: path of key
```

### Dependencies

![Dependencies](./assets/transport.png)

## Cryptography

The crypto package provides sensible defaults for symmetric, asymmetric, hashing and randomness.

We rely on the following libraries:
- [argon2](https://github.com/matthewhartstonge/argon2)
- [crypto](https://pkg.go.dev/golang.org/x/crypto)

### Configuration

To configure, please specify the following:

```yaml
crypto:
  aes:
    key: path to the key
  ed25519:
    public: path to the public
    private: path to the private
  hmac:
    key: path to the key
  rsa:
    public: path to the public
    private: path to the private
  ssh:
    public: path to the public
    private: path to the private
```

### Dependencies

![Dependencies](./assets/crypto.png)

## Debug

This section outlines all utilities added for you troubleshooting abilities.

### statsviz

```http
GET http://localhost:6060/debug/statsviz
```

Check out [statsviz](https://github.com/arl/statsviz).

### pprof

```http
GET http://localhost:6060/debug/pprof/
GET http://localhost:6060/debug/pprof/cmdline
GET http://localhost:6060/debug/pprof/profile
GET http://localhost:6060/debug/pprof/symbol
GET http://localhost:6060/debug/pprof/trace
```

Check out [pprof](https://pkg.go.dev/net/http/pprof).

### fgprof

```http
GET http://localhost:6060/debug/fgprof?seconds=10
```

Check out [fgprof](https://pkg.go.dev/github.com/felixge/fgprof).

### gopsutil

```http
GET http://localhost:6060/debug/psutil
```

Check out [gopsutil](https://github.com/shirou/gopsutil).

### Configuration

To configure, please specify the following:

```yaml
debug:
  address: :6060
  timeout: 10s
```

If you would like to enable TLS, do the following:

```yaml
debug:
  tls:
    cert: path of cert
    key: path of key
```

## Development

This section describes how to run and contribute to the project, if you are interested.

### Style

We favour what is defined in the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

### Dependencies

Please setup the following:
- https://github.com/FiloSottile/mkcert

### Setup

To get yourself setup, please run:

```sh
git submodule sync
git submodule update --init

mkcert -install
make create-certs

make dep
```

### Environment

As we rely on external services these need to be configured:

#### Starting

Please run:

```sh
make start
```

#### Stopping

Please run:

```sh
make stop
```

### Testing

To be able to test locally, please run:

```sh
make specs
```
