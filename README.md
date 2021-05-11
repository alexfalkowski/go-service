[![CircleCI](https://circleci.com/gh/alexfalkowski/go-service.svg?style=svg)](https://circleci.com/gh/alexfalkowski/go-service)
[![Coverage Status](https://coveralls.io/repos/github/alexfalkowski/go-service/badge.svg?branch=master)](https://coveralls.io/github/alexfalkowski/go-service?branch=master)

# Go Service

A framework to build services in go

## Dependencies

This framework [stands on the shoulder of giants](https://en.wikipedia.org/wiki/Standing_on_the_shoulders_of_giants), therefore we have added them here. These are as following:
- [Uber FX](https://github.com/uber-go/fx)
- [Uber Zap](https://github.com/uber-go/zap)
- [OpenTracing](https://github.com/opentracing/opentracing-go)
- [Prometheus](https://github.com/prometheus/client_golang)
- [Jaeger](https://github.com/jaegertracing/jaeger)
- [DataDog](https://github.com/DataDog/dd-trace-go)
- [Redis Cache](https://github.com/go-redis/cache)
- [Snappy](https://github.com/golang/snappy)
- [Ristretto](https://github.com/dgraph-io/ristretto)
- [NSQ](https://github.com/nsqio/go-nsq)

## Testing

To be able to test things locally you have to setup the environment.

### Starting

Please run:

```sh
make start
```

### Stopping

Please run:

```sh
make stop
```

## Usage

The best way to learn is to check out some projects built with this framework. These are the following:
- [go-nonnative-example](https://github.com/alexfalkowski/go-nonnative-example)
