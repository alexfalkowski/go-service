# 🛠️ Debug endpoints

[← Back to README](../README.md)

Debug server config:

```yaml
debug:
  address: tcp://localhost:6060
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

## statsviz

```http
GET http://localhost:6060/<name>/debug/statsviz
```

<https://github.com/arl/statsviz>

## pprof

```http
GET http://localhost:6060/<name>/debug/pprof/
GET http://localhost:6060/<name>/debug/pprof/cmdline
GET http://localhost:6060/<name>/debug/pprof/profile
GET http://localhost:6060/<name>/debug/pprof/symbol
GET http://localhost:6060/<name>/debug/pprof/trace
```

<https://pkg.go.dev/net/http/pprof>

## fgprof

```http
GET http://localhost:6060/<name>/debug/fgprof?seconds=10
```

<https://pkg.go.dev/github.com/felixge/fgprof>
