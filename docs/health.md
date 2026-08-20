# 🩺 Health

[← Back to README](../README.md)

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

`module.Server` installs the HTTP/gRPC health transports, but services own the
checks and observer mapping. Create go-health `server.Registration` values,
register them under the service or gRPC service name on `*server.Server`, and
map them to `healthz`, `livez`, `readyz`, or `grpc` with `Observe`. See the
executable [`Registrations` example](../health/example_test.go) and the
[`go-service-template` health module](https://github.com/alexfalkowski/go-service-template/tree/master/internal/health)
for the standard DI pattern. A checker is not exposed by a probe until that
registration and observer mapping exists.

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
