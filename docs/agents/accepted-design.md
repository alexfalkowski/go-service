# Accepted Design

These entries are mandatory, not advisory. They record accepted designs,
intentional tradeoffs, and support boundaries that agents repeatedly re-flag as
bugs. Before recording a finding, review candidate, audit entry, or proposed
change, agents MUST read this index and every area file below that is relevant
to their scope, and MUST NOT report a finding in these areas without
accounting for the entry that covers it. When a scope spans several areas or
relevance is unclear, read every candidate area file.

- [`accepted-design/transport.md`](accepted-design/transport.md): HTTP/gRPC transport stacks, TLS setup and runtime TLS material, `transport.Module` composition, forwarding headers and request IDs, gRPC user-agent, retry and load-control ordering, HTTP status/route/content-type contracts, gRPC client options, limiters and limiter keys, CORS, HTTP and gRPC streaming, gRPC reflection and service names, terminal response-write errors, webhooks, CloudEvents, and `net/...` helper placement.
- [`accepted-design/telemetry.md`](accepted-design/telemetry.md): OpenTelemetry propagation, OTLP endpoints and headers, Prometheus metrics, process-global telemetry providers and logger state, telemetry shutdown hooks, log level, SQL OpenTelemetry, HTTP/gRPC telemetry logging and instrumentation, `telemetry/header` secrets, and OpenFeature trace attributes.
- [`accepted-design/crypto-token.md`](accepted-design/crypto-token.md): RSA key loading and validation, crypto key generators, access model and Casbin policy loading, JWT/PASETO verification and key material, token tests, and UUIDv7 ID generation.
- [`accepted-design/config.md`](accepted-design/config.md): `config.NewConfig` decode and validation, `bytes` size limits and parsing, cache registration and Redis config, PostgreSQL DSNs, nil-safe pointer configs, and `vendor/` regeneration.
- [`accepted-design/encoding.md`](accepted-design/encoding.md): Base64 encoded length, decoder admissibility for untrusted input (including `msgpack` and `gob`), and JSON/HJSON duplicate keys.
- [`accepted-design/mvc.md`](accepted-design/mvc.md): MVC error rendering, view construction, static file serving and path handling, conditional caching, and not-found handling.
- [`accepted-design/health-debug.md`](accepted-design/health-debug.md): Default bind addresses and debug exposure, pprof/fgprof profiling, build/version provenance, health registration, and health probes.
- [`accepted-design/cli-lifecycle.md`](accepted-design/cli-lifecycle.md): Manual server lifecycle wiring, module-wiring tests, `cli.RunCode`, `cli.Application.AddClient` and `Run`, `net/server.Service` serve errors, gRPC serve and graceful shutdown, and OpenFeature process-global state.
