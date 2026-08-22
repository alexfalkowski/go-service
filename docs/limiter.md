# 🚦 Limiter

[← Back to README](../README.md)

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
>
> - Omitting the `limiter` block disables limiting entirely; a nil config is treated as disabled.
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
