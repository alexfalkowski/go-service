# AGENTS.md

## Shared skills

This repository uses the shared skills from `bin/skills/`. Read
`bin/AGENTS.md` for the canonical shared skill list and use the smallest
matching skill for the task.

## Repo Snapshot

- Library repo; check `go.mod` for module and Go version details.
- DI is under `di/`; CLI helpers are under `cli/`.
- Many `make` targets come from the `bin/` submodule.

## Setup

- Initialize the submodule before using `make`: `git submodule sync` then `git submodule update --init`.
- Equivalent: `make submodule`.
- Submodule fetches use SSH.

## Common Commands

- Discover targets: `make help`.
- Dependencies: `make dep` after dependency changes.
- Quality: `make specs`, `make lint`, `make fix-lint`, `make format`, `make sec`.
- Coverage and benchmarks: use the repo `make` targets shown by `make help`.
- TLS fixtures: `mkcert -install` then `make create-certs`.
- `encode-config` expects GNU `base64 -w 0`; on macOS/BSD use `base64 | tr -d '\n'`.

## Review Workflows

- When the user asks "Do a deep dive secure code review of <package>",
  interpret it as:
  - For each package and subpackage under `<package>`, launch multiple agents.
  - Each agent should perform a thorough and accurate `$code-review` and
    `$security-audit`.
  - After all agents finish, aggregate findings into `FINDINGS.md`.
  - As findings are fixed, remove them from `FINDINGS.md`.
  - Once all findings are resolved, delete `FINDINGS.md`.

## Layout And Wiring

- Feature packages usually use `config.go`, `module.go`, plus implementation files.
- `module/` exports the main Fx bundles.
- `config/` owns top-level config plus projections into transport, SQL, and telemetry config.
- `net/` holds lower-level HTTP/gRPC, metadata, header, and server helpers.
- `transport/` holds the higher-level HTTP/gRPC stacks, middleware, and ops endpoints.
- `internal/test/` contains shared test helpers; `test/` stores fixtures and reports.
- Modules are composed with `di.Module(...)`; many constructors use `di.In`.

## Configuration

- `config.NewDecoder` resolves `-i` as `file:<path>`, `env:<ENV_VAR>`, or a default config file named for the service.
- Default file lookup checks the executable directory, `$XDG_CONFIG_HOME/<serviceName>/`, and `/etc/<serviceName>/`.
- Many config fields use source strings through `os.FS.ReadSource`: `env:NAME`, `file:/path`, or a literal value.
- Nil pointer sub-configs usually mean "disabled".
- Standard module wiring is the supported path. Do not flag hypothetical failures that require hand-wiring an incomplete DI graph unless the public API explicitly promises that custom construction mode.

## Gotchas

- Manual transport TLS setup must call the HTTP and gRPC transport register functions; `transport.Module` does this on the normal path.
- Manual server lifecycle wiring should use `net/server.Register(...)`.
- `transport.Module` is normally consumed through `module.Server`, which also
  wires `debug.Module`; do not flag `transport.NewServers` requiring
  `*debug.Server` unless a public API starts promising standalone
  `transport.Module` composition without the standard `module.Server` bundle.
- `cli.RunCode` returns `os.ExitCodeSuccess` on success, preserves non-zero
  shutdown exit codes requested through `di.ExitCode(...)`, and otherwise
  returns `os.ExitCodeFailure`.
- `net/server.Service` intentionally logs asynchronous `Server.Serve` errors
  and requests shutdown with `di.ExitCode(os.ExitCodeServeFailure)`; it does not
  return the raw serve error from `Stop`.
- `telemetry.Register()` installs the global OpenTelemetry propagator.
- `cache.Register(...)` sets the package-level cache used by generic cache helpers.
- Redis cache config intentionally expects `cache.options.url` to exist and be a string.
- Access model and policy config are resolved through `os.FS.ReadSource`; use `file:` for files or `env:` for content from the environment.
- IP metadata intentionally trusts forwarding headers; deploy behind trusted proxies that strip spoofed headers before using the `"ip"` limiter key.
- `net/header.ForwardedIPs` is intentionally an exported mutable list, similar
  to standard-library package variables such as `os.Args`. Do not flag this
  solely because importing packages could mutate it; only report concrete bugs
  with evidence of accidental mutation, concurrent mutation, or an API promise
  of immutability.
- `Request-Id`/`request-id` is intentionally a logical request identifier, not
  a per-wire-attempt identifier. Client metadata runs before retry middleware,
  so all retry attempts for one logical HTTP/gRPC request share the same value.
  Retry policies intentionally treat a present request id as the idempotency
  key/contract for retryable writes; services that accept retried writes should
  deduplicate by request id when duplicate processing would be unsafe. Do not
  flag the default HTTP/gRPC retry policy merely because metadata injects
  request ids before retry.
- gRPC client constructor options use the package's last-wins functional option
  convention. `WithClientDialOption`, `WithClientUnaryInterceptors`, and
  `WithClientStreamInterceptors` expect all custom values for one client
  construction to be passed in a single call; repeated calls intentionally
  replace earlier values. Do not flag this as dropped configuration unless a
  public API starts promising accumulation across repeated option helpers.
- Transport limiter keys are `"user-agent"`, `"ip"`, and `"user-id"`; `"token"` is intentionally not a limiter key. Server limiters run after metadata extraction and token verification, so `"user-id"` is the verified principal (JWT/PASETO subject or SSH key name), and missing, malformed, or invalid auth is rejected before the limiter by design. Do not flag that bypass; use an external edge/gateway/ingress/load-balancer/service-mesh limiter when those attempts need quota enforcement.
- The built-in transport limiter is intentionally in-memory and per-process. Treat it as a last-resort local safeguard; prefer external edge/gateway/ingress/load-balancer/service-mesh limiting for production abuse protection.
- HTTP telemetry logger service/method derivation may include request URL path
  segments for non-canonical HTTP routes. This is intentional for client and
  server debugging because HTTP clients can call arbitrary paths and route
  patterns are not always available at the logging layer. Do not flag this as a
  query/header/body leakage issue unless the logger starts recording
  `RawQuery`, `RequestURI`, headers, cookies, or bodies, or a specific route
  places secrets in path segments contrary to service policy.
- gRPC telemetry logging intentionally records raw error values for operator
  diagnostics. Client-facing safety is handled by gRPC status/error rendering;
  logs are backend observability data and should be protected by deployment log
  access controls. Do not flag raw gRPC error logging as a data leak unless a
  concrete code path places secrets, credentials, request bodies, or other
  prohibited sensitive values into those errors contrary to service policy.
- gRPC client telemetry intentionally includes the raw `conn.Target()` in client
  log messages to identify the configured downstream endpoint. Targets are
  expected to be configuration-controlled service addresses and must not contain
  credentials, tokens, request data, or other secrets. Do not flag raw target
  logging unless a concrete configuration or call path allows sensitive data in
  the target string.
- Before flagging a nil-pointer panic on embedded pointer configuration types,
  inspect the called method. Go permits calling pointer-receiver methods on nil
  pointers, and methods such as `(*config/server.Config).IsEnabled` are
  intentionally nil-safe.
- gRPC server reflection is intentionally always registered by `net/grpc.NewServer`; restrict public exposure at the bind address, TLS/auth, ingress, firewall, or service-mesh boundary.
- MVC controller errors render a client-safe `mvc.Error` model; `mvcModelError` metadata intentionally remains the raw error string for compatibility and must not be rendered unless diagnostic detail exposure is acceptable.
- MVC not-found handling intentionally uses a simple `Accept` header check for
  `text/html` and does not fully evaluate quality weights such as
  `text/html;q=0`. Do not flag this unless there is concrete evidence of a
  real client, route, or deployment contract that depends on strict weighted
  `Accept` negotiation for MVC 404 responses.
- JWT verification requires both the expected algorithm and a `kid` header.
- `telemetry/header.Map.MustSecrets` can panic if secret resolution fails during config projection.
- Health registration helpers require `*net/http.ServeMux`.
- Shared metadata, header, and string helpers live under `net/...`, not `transport/...`.
- `vendor/` is gitignored and regenerated via `make dep`.
- HTTP webhook verification buffers `req.Body` intentionally for signature checks.
  Under supported server wiring, `transport/http.NewServer` installs the body
  limiter before mux handlers, so inbound webhook bodies are capped by
  `Config.MaxReceiveSize` before verification. Do not flag this as an
  unbounded-read issue unless the code path bypasses the transport server chain
  without an equivalent request-size cap.
- HTTP webhook verification intentionally does not maintain replay state.
  Receivers must deduplicate or process idempotently using `Webhook-Id` or the
  event id, preferably with durable shared storage when duplicate valid
  deliveries would be unsafe. Do not flag missing transport-level replay
  storage unless the code starts promising replay protection.
- HTTP webhook signing intentionally ignores the Standard Webhooks `Sign` error
  because the current vendored implementation always returns nil. Do not flag
  this unless the dependency behavior changes.
- HTTP CloudEvents receiver registration intentionally ignores the current
  CloudEvents constructor errors because supported wiring passes no protocol
  options and uses the typed `ReceiverFunc`, which matches the SDK receive
  handler shape. Do not flag this unless dependency behavior or call arguments
  change.

## Testing, Style, And Docs

- Tests commonly use `stretchr/testify/require`.
- Go files use tabs; YAML uses 2 spaces per `.editorconfig`.
- Every exported identifier, including under `internal/test/**`, needs a GoDoc comment.
- GoDoc comments should start with the identifier name or `Deprecated:`.

## CI

- CI initializes submodules, prepares certs/services, then runs dependency, lint, security, spec, benchmark, and coverage targets.
- Check CI config for exact service images, ports, and command order.
