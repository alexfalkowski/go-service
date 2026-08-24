# 🚩 Feature flags (OpenFeature)

[← Back to README](../README.md)

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
>
> - `feature.Config` embeds client config; `IsEnabled` is true only when both the feature config and embedded client config are present. An empty `feature:` block is treated as disabled by feature config helpers.
> - This repository does not construct a built-in OpenFeature provider from this config.
> - Services that need a remote or custom provider should use `feature.Config` in their own provider constructor and provide the resulting `openfeature.FeatureProvider` in DI; `feature.Module` registers that supplied provider with the OpenFeature SDK lifecycle.
