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

The configuration is based on YAML and is read from an env variable called `CONFIG_FILE`. The reason for this is that we want to be able to separate how configuration is retrieved. This way we can use and [application configuration system](https://github.com/alexfalkowski/konfig).

The configuration can be [watched](https://github.com/fsnotify/fsnotify) for write changes. If it changes the application is stopped. This way an [orchestration](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy) system can just restart the process.

This is the [configuration](config/config.go).

## Caching

The framework currently supports the following caching solutions:
- [Redis Cache](https://github.com/go-redis/cache)
- [Ristretto](https://github.com/dgraph-io/ristretto)

We also support the following compressions to optimize cache size:
- [Snappy](https://github.com/golang/snappy)

## Health

The health package is based on [go-health](https://github.com/alexfalkowski/go-health). This package allows us to create all sorts of ways to check external and internal systems.

We also provide ways to integrate into container integration systems. So we provide the following endpoints:
- `/health` - This allows us to check any external dependency and provide a breakdown of what is not functioning. This should only be used for verification.
- `/liveness`: Can be used for k8s [liveness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-liveness-command).
- `/readiness`: Can be used for k8s [readiness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-readiness-probes).

## Logging

For logging we use [Uber Zap](https://github.com/uber-go/zap).

## Metrics

For metrics we use [Prometheus](https://github.com/prometheus/client_golang).

## Security

For security we support the following:
- [Auth0](https://auth0.com/)

## SQL

For SQL databases we support the following:
- [Postgres](https://github.com/jackc/pgx)

We also support master, slave combinations with the awesome [mssqlx](https://github.com/linxGnu/mssqlx).

## Tracing

For distributed tracing we support the following:
- [OpenTracing](https://github.com/opentracing/opentracing-go)
- [Jaeger](https://github.com/jaegertracing/jaeger)
- [DataDog](https://github.com/DataDog/dd-trace-go)

## Transport

The transport layer provides ways to abstract communication for in/out of the service. So we have the following integrations:
- [gRPC](https://grpc.io/) - The author truly believes in [IDLs](https://en.wikipedia.org/wiki/Interface_description_language).
- [REST](https://en.wikipedia.org/wiki/Representational_state_transfer) - This is achieved with [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway).
- [NSQ](https://github.com/nsqio/go-nsq)

### REST

Below is list of the provided handlers:
- [CORS](https://github.com/rs/cors)

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

## Projects

Below is a list of projects using this framework:
- [Konfig](https://github.com/alexfalkowski/konfig)
