[![CircleCI](https://circleci.com/gh/alexfalkowski/go-service.svg?style=svg)](https://circleci.com/gh/alexfalkowski/go-service)
[![Coverage Status](https://coveralls.io/repos/github/alexfalkowski/go-service/badge.svg?branch=master)](https://coveralls.io/github/alexfalkowski/go-service?branch=master)

# Go Service

A framework to build services in go. This came out of out building services over the years and what I have considered good practices in building services. Hence it is highly subjective and opinionated.

This framework [stands on the shoulder of giants](https://en.wikipedia.org/wiki/Standing_on_the_shoulders_of_giants) so we don't reinvent the wheel!

## Dependency Injection

This framework heavily relies on [DI](https://en.wikipedia.org/wiki/Dependency_injection). We have chosen to use [Uber FX](https://github.com/uber-go/fx). So there is great information online to get you up to speed.

## Commands

A service has commands that are configured using [Cobra](github.com/spf13/cobra). Each service has the following commands (you can add more):
- `Server` - This will host your API.
- `Worker` - This will host your background processing.
- `Client` - This will have a command that starts and finishes.

These are configured in the main function.

## Configuration

The supported configuration kinds are as follows:
- [YAML](https://github.com/go-yaml/yaml)
- [TOML](https://github.com/BurntSushi/toml)

The configuration can be read from multiple sources by specifying a flag called `--input`. As per the following:
- `env:CONFIG_FILE` - Read from an env variable called `CONFIG_FILE`. This is the default if nothing is passed. The env variable can be file path or the configuration. If it is the config, we expect the format of `extension:ENV_VARIABLE`, where extension is the supported kinds and ENV_VARIABLE contains the contents of the config.
- `file:path` - Read from the path.

The reason for this is that we want to be able to separate how configuration is retrieved. This way we can use and [application configuration system](https://github.com/alexfalkowski/konfig).

This is the [configuration](config/config.go). We will outline the config required in each section.

## Caching

The framework currently supports the following caching solutions:
- [Redis Cache](https://github.com/go-redis/cache)
- [Ristretto](https://github.com/dgraph-io/ristretto)

We also support the following compressions to optimize cache size:
- [Snappy](https://github.com/golang/snappy)

### Configuration

To configure, please specify the following:

```yaml
cache:
  redis:
    addresses:
      server: localhost:6379
    username: test
    password: test
    db: 0
  ristretto:
    max_counters: 10000000
    max_cost: 100000000
    buffer_items: 64
```

```toml
[cache.redis]
username = "test"
password = "test"
db = 0

[cache.redis.addresses]
server = "localhost:6379"

[cache.ristretto]
max_counters = 10_000_000
max_cost = 100_000_000
buffer_items = 64
```

## Health

The health package is based on [go-health](https://github.com/alexfalkowski/go-health). This package allows us to create all sorts of ways to check external and internal systems.

We also provide ways to integrate into container integration systems. So we provide the following endpoints:
- `/healthz` - This allows us to check any external dependency and provide a breakdown of what is not functioning. This should only be used for verification.
- `/livez`: Can be used for k8s [liveness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-liveness-command).
- `/readyz`: Can be used for k8s [readiness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-readiness-probes).

This is modelled around [Kubernetes API health endpoints](https://kubernetes.io/docs/reference/using-api/health-checks/).

## Logging

For logging we use [Uber Zap](https://github.com/uber-go/zap).

## Metrics

For metrics we use [Prometheus](https://github.com/prometheus/client_golang).

## Security

For security we support the following:
- [Auth0](https://auth0.com/)

### Configuration

To configure, please specify the following:

```yaml
security:
  auth0:
    url: test_url
    client_id: test_client_id
    client_secret: test_client_secret
    audience: test_audience
    issuer: test_issuer
    algorithm: test_algorithm
    json_web_key_set: test_json_web_key_set
```

```toml
[security.auth0]
url = "test_url"
client_id = "test_client_id"
client_secret = "test_client_secret"
audience = "test_audience"
issuer = "test_issuer"
algorithm = "test_algorithm"
json_web_key_set = "test_json_web_key_set"
```

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
        url: postgres://test:test@localhost:5432/test?sslmode=disable
    slaves:
      -
        url: postgres://test:test@localhost:5432/test?sslmode=disable
    max_open_conns: 5
    max_idle_conns: 5
    conn_max_lifetime: 1h
```

```toml
[sql.pg]
max_open_conns = 5
max_idle_conns = 5
conn_max_lifetime = "1h"

[[sql.pg.masters]]
url = "postgres://test:test@localhost:5432/test?sslmode=disable"

[[sql.pg.slaves]]
url = "postgres://test:test@localhost:5432/test?sslmode=disable"
```

## Telemetry

For distributed tracing we support the following:
- [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go)
- [Jaeger](https://github.com/jaegertracing/jaeger)

### Configuration

To configure, please specify the following:

```yaml
otel:
  kind: jaeger
  host: localhost:6831
```

```toml
[otel]
kind = "jaeger"
host = "localhost:6831"
```

## Transport

The transport layer provides ways to abstract communication for in/out of the service. So we have the following integrations:
- [gRPC](https://grpc.io/) - The author truly believes in [IDLs](https://en.wikipedia.org/wiki/Interface_description_language).
- [REST](https://en.wikipedia.org/wiki/Representational_state_transfer) - This is achieved with [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway).
- [NSQ](https://github.com/nsqio/go-nsq)

### gRPC

Below is list of the provided interceptors:
- [Limiter](https://github.com/ulule/limiter)

### REST

Below is list of the provided handlers:
- [CORS](https://github.com/rs/cors)

### Configuration

To configure, please specify the following:

```yaml
transport:
  port: 8080
  http:
    retry:
      timeout: 1s
      attempts: 3
  grpc:
    retry:
      timeout: 1s
      attempts: 3
  nsq:
    lookup_host: localhost:4161
    host: localhost:4150
    retry:
      timeout: 1s
      attempts: 3
```

```toml
[transport]
port = "8080"

[transport.http.retry]
timeout = "1s"
attempts = 3

[transport.grpc.retry]
timeout = "1s"
attempts = 3

[transport.nsq]
lookup_host = "localhost:4161"
host = "localhost:4150"

[transport.nsq.retry]
timeout = "1s"
attempts = 3
```

## Development

This section describes how to run and contribute to the project, if you are interested.

### Style

We favour what is defined in the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

### Setup

To get yourself setup, please run:

```sh
make setup
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
