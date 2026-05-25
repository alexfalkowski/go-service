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
  When default lookup reaches the user config directory candidate, HOME or
  XDG_CONFIG_HOME is expected to be available; missing both is treated as a
  misconfigured runtime. Do not flag the resulting `os.UserConfigDir` panic as
  a config issue. Services that do not want this environment contract should
  pass an explicit `-i file:<path>` or `-i env:<ENV_VAR>` source.
- Many config fields use source strings through `os.FS.ReadSource`: `env:NAME`, `file:/path`, or a literal value.
- Service configuration files should contain configuration values and secret
  source references, not raw passwords or credentials.
- Nil pointer sub-configs usually mean "disabled".
- Standard module wiring is the supported path. Do not flag hypothetical failures that require hand-wiring an incomplete DI graph unless the public API explicitly promises that custom construction mode.

## Gotchas

- Manual transport TLS setup must call the HTTP and gRPC transport register functions; `transport.Module` does this on the normal path.
- Manual server lifecycle wiring should use `net/server.Register(...)`.
- `transport.Module` is normally consumed through `module.Server`, which also
  wires `debug.Module`; do not flag `transport.NewServers` requiring
  `*debug.Server` unless a public API starts promising standalone
  `transport.Module` composition without the standard `module.Server` bundle.
- Server defaults that use `net.DefaultAddress`, including the debug server's
  `tcp://:6060` fallback, intentionally bind all interfaces so containerized
  workloads remain reachable through Kubernetes Services, probes, ingress
  controllers, and sidecars. Do not flag this solely because debug endpoints can
  be reached remotely when the debug server is enabled; restrict exposure with
  explicit addresses, TLS/mTLS, NetworkPolicy, ingress/firewall, or
  service-mesh policy. Only report concrete bugs such as accidental listener
  changes, ignored explicit addresses, missing documented protections, or a
  public API promise of localhost-only debug binding.
- `cli.RunCode` returns `os.ExitCodeSuccess` on success, preserves non-zero
  shutdown exit codes requested through `di.ExitCode(...)`, and otherwise
  returns `os.ExitCodeFailure`.
- `net/server.Service` intentionally logs asynchronous `Server.Serve` errors
  and requests shutdown with `di.ExitCode(os.ExitCodeServeFailure)`; it does not
  return the raw serve error from `Stop`.
- `telemetry.Register()` installs the global OpenTelemetry propagator.
- `cache.Register(...)` sets the package-level cache used by generic cache helpers.
  It is intentionally called by the supported wiring path during startup/test
  setup, not as a concurrent runtime reconfiguration API. Do not flag
  unsynchronized global-state or nil-pointer issues based solely on hypothetical
  concurrent manual `cache.Register`, `cache.Get`, or `cache.Persist` calls
  unless a concrete public API path starts promising concurrent manual
  re-registration or the repository adds such a runtime path.
- `encoding/base64.EncodedLen` intentionally follows the standard library's
  `EncodedLen(int)` contract. Do not flag hypothetical integer overflow from
  passing absurd configured `bytes.Size` values directly to it; callers that
  compare configuration-sized limits, such as cache decode size guards, must
  guard those limits before calling it.
- `bytes.ParseSize` intentionally delegates human-readable size parsing to
  `github.com/docker/go-units` for compatibility with existing configuration
  values. Do not flag upstream float-to-int range behavior from absurdly large
  size strings as a local issue unless this repository adds a public promise to
  reject every representability edge case or a concrete supported path shows an
  unsafe limit being applied from such a value. Do not flag accepted suffix
  spellings such as `MiB` merely because `go-units.FromHumanSize` treats them
  as decimal multipliers; that compatibility is documented local behavior. Do
  not flag `bytes.Size` marshal/unmarshal round-trip failures for exabyte-scale
  values solely because `go-units.HumanSize` can format suffixes that
  `go-units.FromHumanSize` does not parse; that is accepted upstream behavior
  unless this repository adds a strict round-trip promise for those values.
- Redis cache config intentionally expects `cache.options.url` to exist and be a string.
- PostgreSQL DSN security options, including TLS/`sslmode`, are intentionally
  part of the DSN supplied by the service configuration. `database/sql/pg`
  passes resolved DSNs through to pgx and does not impose repository-level DSN
  construction policy. Do not flag pass-through DSN handling merely because an
  insecure DSN could be supplied; only report concrete bugs such as accidental
  DSN rewriting, secret leakage, or a public API promise to enforce secure DSN
  policy.
- RSA public-key loading intentionally validates the repository's key-size
  policy and delegates deeper RSA parameter checks, such as public exponent
  usability, to the standard library operations that consume the key. Do not
  flag `crypto/rsa.Config.PublicKey` merely because `Encrypt` can later return a
  standard-library error for an unusual parsed key parameter; only report
  concrete bugs such as panics, accepted weak key sizes, secret leakage, or a
  public API promise to fully validate every RSA parameter at load time.
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
- `net/http/status.Code(err)` usually returns a valid HTTP status code because
  repository code should construct HTTP status errors with the constants exposed
  by `net/http` (for example `http.StatusBadRequest`) or intentionally supported
  valid custom codes such as `499`. Do not flag hypothetical invalid
  `WriteHeader` panics from manually constructed bogus codes unless a concrete
  public API path accepts untrusted status codes or starts promising validation.
- gRPC client constructor options use the package's last-wins functional option
  convention. `WithClientDialOption`, `WithClientUnaryInterceptors`, and
  `WithClientStreamInterceptors` expect all custom values for one client
  construction to be passed in a single call; repeated calls intentionally
  replace earlier values. Do not flag this as dropped configuration unless a
  public API starts promising accumulation across repeated option helpers.
- Transport limiter keys are `"user-agent"`, `"ip"`, `"user-id"`, and
  `"service-method"`; `"token"` is intentionally not a limiter key. Server
  limiters run after metadata extraction and token verification, so `"user-id"`
  is the verified principal (JWT/PASETO subject or SSH key name), and missing,
  malformed, or invalid auth is rejected before the limiter by design. Do not
  flag that bypass; use an external edge, gateway, ingress, load balancer, or
  service mesh limiter when those attempts need quota enforcement.
- The built-in transport limiter is intentionally in-memory and per-process. Treat it as a last-resort local safeguard; prefer external edge/gateway/ingress/load-balancer/service-mesh limiting for production abuse protection.
- HTTP telemetry logger service/method derivation may include request URL path
  segments for non-canonical HTTP routes. This is intentional for client and
  server debugging because HTTP clients can call arbitrary paths and route
  patterns are not always available at the logging layer. Do not flag this as a
  query/header/body leakage issue unless the logger starts recording
  `RawQuery`, `RequestURI`, headers, cookies, or bodies, or a specific route
  places secrets in path segments contrary to service policy.
- OpenTelemetry HTTP client instrumentation is delegated to upstream
  `otelhttp.NewTransport`, which currently records semantic-convention URL
  attributes from the request URL and may include `RawQuery` in `url.full`.
  Treat this as a documented third-party instrumentation behavior, not a local
  `net/http` finding. Do not flag it unless this repository adds local URL
  attribute construction, logging/export code that records queries directly, or
  a supported upstream option that can sanitize this behavior without mutating
  the outbound request.
- HTTP client `RoundTripper` implementations that can return locally before
  delegating to another `RoundTripper` must make request-body ownership explicit
  with `net/http.ClosingRoundTripper`. Return the closing adapter's close-body
  flag as true only for local rejection paths, and false after delegating
  because the delegated transport owns `req.Body` closure.
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
- gRPC service and method names are generated from Buf-managed proto files such
  as `internal/test/greet/v1/service.proto`, which require package-qualified
  service names and valid RPC method names. Do not flag
  `net/grpc.ParseServiceMethod` merely because a hypothetical manually
  constructed method string like `/pkg.Service/` or `/svc/Get.Name` could split
  or fall back to `root`. Only report a concrete issue if untrusted,
  non-generated method strings are used for a security decision, or a public API
  promises strict validation of arbitrary method strings.
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
- JSON decoding intentionally keeps the standard library's duplicate object key
  behavior, where later values replace earlier values. Do not flag this as a
  finding unless a public API starts promising duplicate-key rejection or this
  repository adds a strict JSON decoder mode.
- HJSON duplicate-key errors intentionally preserve upstream diagnostic detail,
  including duplicate decoded values, to keep bad configuration files
  debuggable. Configuration files are not expected to contain raw passwords or
  credentials; they should contain source references such as `env:NAME` or
  `file:/path`. Do not flag this as a secret leak unless a concrete code path
  starts placing raw secrets into HJSON configuration contrary to that policy.

## Testing, Style, And Docs

- Tests commonly use `stretchr/testify/require`.
- Go files use tabs; YAML uses 2 spaces per `.editorconfig`.
- Every exported identifier, including under `internal/test/**`, needs a GoDoc comment.
- GoDoc comments should start with the identifier name or `Deprecated:`.

## CI

- CI initializes submodules, prepares certs/services, then runs dependency, lint, security, spec, benchmark, and coverage targets.
- Check CI config for exact service images, ports, and command order.
