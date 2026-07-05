# AGENTS.md

## Shared guidance

Use `bin/AGENTS.md` for shared skills and cross-repository defaults.

## Repo Snapshot

- Library repo; check `go.mod` for module and toolchain details.
- DI is under `di/`; CLI helpers are under `cli/`.
- Many `make` targets come from the `bin/` submodule.

## Setup

- Initialize the submodule before using shared Make targets. Use
  `make submodule` once the shared `bin` checkout is present; see
  `bin/AGENTS.md` for fresh-clone bootstrap details.
- Submodule fetches use SSH.

## Common Commands

- Discover targets: `make help`.
- Dependencies: `make dep` after dependency changes.
- Quality: `make specs`, `make lint`, `make fix-lint`, `make format`, `make sec`.
- Coverage and benchmarks: use the repo `make` targets shown by `make help`.
- TLS fixtures: `make create-certs` when the target is relevant.

## Layout And Wiring

- Feature packages usually use `config.go`, `module.go`, plus implementation files.
- `module/` exports the main Fx bundles.
- `config/` owns top-level config plus projections into transport, SQL, and telemetry config.
- `net/` holds lower-level HTTP/gRPC, metadata, header, and server helpers.
- `transport/` holds the higher-level HTTP/gRPC stacks, middleware, and ops endpoints.
- `internal/test/` contains shared test helpers; `test/` stores fixtures and reports.
- `internal/test` protobufs are test fixtures, not an external API contract.
  Use `make -C internal/test generate` / `stale` when changing them, but do
  not treat the inherited `breaking` target as applicable there; the shared
  Buf target assumes an `api/` contract directory.
- Modules are composed with `di.Module(...)`; many constructors use `di.In`.

## Configuration

- `config.NewDecoder` resolves `-config`/`-c` as `file:<path>`, `env:<ENV_VAR>`, or a default config file named for the service.
- Default file lookup checks the executable directory, `$XDG_CONFIG_HOME/<serviceName>/`, and `/etc/<serviceName>/`.
  Default lookup may resolve the user config directory before probing file
  candidates, so HOME or XDG_CONFIG_HOME is expected to be available when
  default lookup starts; missing both is treated as a misconfigured runtime. Do
  not flag the resulting `os.UserConfigDir` panic as a config issue. Services
  that do not want this environment contract should pass an explicit
  `-config file:<path>` or `-config env:<ENV_VAR>` source.
- Many config fields use source strings through `os.FS.ReadSource`: `env:NAME`, `file:/path`, or a literal value.
- Source strings are administrator-supplied configuration. This repository
  accepts that admins must point file/env/literal sources at appropriate
  material, including reasonably sized secrets, keys, certificates, access
  models, policies, DSNs, and service config files. Do not flag unbounded
  `ReadSource`/`ReadFile` behavior solely because a misconfigured admin could
  point a source string at an unexpectedly large file, env var, or literal.
  Report only concrete bugs where untrusted runtime input controls the source,
  documented size limits are ignored, or a public API promises bounded source
  reads.
- Low-level `config/options.Map` values are administrator-supplied tuning
  knobs. Do not flag size-valued options such as HTTP `max_header_bytes` or
  gRPC `max_header_list_size`, `initial_window_size`,
  `initial_conn_window_size`, or `max_send_msg_size` solely because an admin can
  configure values above `bytes.MaxConfigSize`. `bytes.MaxConfigSize` applies
  only where a typed config field or public API explicitly promises that
  repository-owned cap, such as `MaxReceiveSize` validation. Report only
  concrete bugs such as ignored typed validation, untrusted runtime input
  controlling the option, destination-type overflow, or documented bounds being
  bypassed.
- Service configuration files should contain configuration values and secret
  source references, not raw passwords or credentials.
- Nil pointer sub-configs usually mean "disabled".
- Downstream services that need application-specific configuration compose the
  standard `module.Server` or `module.Client` bundle with a service-local
  `internal/config.Module`. The supported pattern is to provide
  `config.NewConfig[ServiceConfig]`, decorate the embedded shared
  `*config.Config`, and expose service-specific projections with constructors.
  Do not flag the absence of generic helpers such as `module.ServerWithConfig[T]`
  or `module.ClientWithConfig[T]` as a feature gap solely because services need
  custom typed config. Report only concrete bugs where the documented/template
  pattern fails, duplicate config decoding is forced by supported wiring,
  projections cannot be supplied through `di.Decorate`, or the support boundary
  changes.
- Standard module wiring is the supported path. Do not flag hypothetical failures that require hand-wiring an incomplete DI graph unless the public API explicitly promises that custom construction mode.
- Recommend Fx `optional:"true"` dependency tags only with concrete evidence
  that a supported DI graph may omit that dependency. A nil check, "optional"
  prose, or guarded hook installation is a lead, not enough by itself; verify
  standard `module.Server`/`module.Client` wiring, CLI/server application
  wiring, a config-disabled path, or existing supported test/user wiring where
  the dependency can genuinely be absent. Directly composing an exported
  lower-level module is not enough evidence unless a public contract explicitly
  promises that module works standalone without the standard bundle. If the
  standard module bundle always resolves the dependency and no such lower-level
  contract exists, do not flag the missing optional tag.
- `*os.FS` dependencies are provided by the supported DI wiring path and are
  expected to be non-nil there. Do not flag nil-`*os.FS` panics in token,
  crypto, config, or source-string loading paths based solely on manually
  calling constructors with nil filesystem dependencies unless a public API
  explicitly promises nil-FS tolerance or a supported path can provide nil.

## Gotchas

- Manual transport TLS setup must call the HTTP and gRPC transport register functions; `transport.Module` does this on the normal path.
- Manual server lifecycle wiring should use `net/server.Register(...)`.
- `transport.Module` is normally consumed through `module.Server`, which also
  wires `debug.Module`; do not flag `transport.NewServers` requiring
  `*debug.Server`. Mentions that `transport.Module` handles transport
  registration, TLS filesystem registration, or lifecycle registration do not
  promise a complete standalone server bundle without `module.Server`. Report
  this only if a public API explicitly promises standalone `transport.Module`
  composition without the standard `module.Server` bundle.
- Module-related behavior is tested through CLI/server application wiring. Do
  not flag missing direct package tests for Fx module provider inventory,
  module composition, or `transport.Module` lifecycle registration solely
  because they are not asserted in the package that declares the module. Report
  only concrete broken behavior through the CLI or supported `module.Server`
  path, or an explicit public promise of lower-level standalone module use.
- Server defaults that use `net.DefaultAddress`, including the debug server's
  `tcp://:6060` fallback, intentionally bind all interfaces so containerized
  workloads remain reachable through Kubernetes Services, probes, ingress
  controllers, and sidecars. Do not flag this solely because debug endpoints can
  be reached remotely when the debug server is enabled; restrict exposure with
  explicit addresses, TLS/mTLS, NetworkPolicy, ingress/firewall, or
  service-mesh policy. Only report concrete bugs such as accidental listener
  changes, ignored explicit addresses, missing documented protections, or a
  public API promise of localhost-only debug binding.
- Debug profiling endpoints are administrator/operator diagnostics. Do not flag
  `net/http/pprof` or `fgprof` duration parameters solely because an authorized
  debug caller can request a long profile. Long captures are an intentional
  diagnostic capability; restrict debug exposure with bind addresses, TLS/mTLS,
  ingress/firewall, NetworkPolicy, or service-mesh policy. Report only concrete
  bugs such as ignored explicit debug server limits, accidental public exposure,
  missing documented protections, or repository-owned profiling wrapper logic
  that violates its own duration/admission contract.
- `debug/internal/fgprof` intentionally delegates request handling to upstream
  `github.com/felixge/fgprof.Handler`. Treat its cancellation behavior and
  zero-sample profile export edge cases as upstream behavior, not local debug
  findings, unless this repository adds local fgprof handler logic or promises
  cancellation-aware profiling semantics.
- `cli.RunCode` returns `os.ExitCodeSuccess` on success, preserves non-zero
  shutdown exit codes requested through `di.ExitCode(...)`, and otherwise
  returns `os.ExitCodeFailure`.
- `cli.Application.AddClient` intentionally models short-lived client command
  work as DI/Fx startup work. Client commands may perform their main action
  from constructors or lifecycle `OnStart` hooks, then stop the graph
  immediately after startup completes. Do not flag the absence of a separate
  post-DI command-task API solely because command work lives in `OnStart`; this
  is the supported pattern, as used by downstream client templates. Report only
  concrete broken behavior such as incorrect error propagation, ignored
  shutdown exit codes, lifecycle ordering bugs, or a documented command
  contract that cannot be expressed with the DI lifecycle.
- `cli.Application.Run` intentionally sanitizes Go test harness `-test.*`
  arguments before handing `os.Args` to the command runner because this
  repository commonly exercises CLI applications through Go test binaries. Do
  not flag this solely because a hypothetical downstream command could define a
  user-facing flag with the reserved `test.` prefix; report only concrete
  breakage in a supported CLI contract or an explicit public promise to
  preserve `-test.*` command flags.
- `net/server.Service` intentionally logs asynchronous `Server.Serve` errors
  and requests shutdown with `di.ExitCode(os.ExitCodeServeFailure)`; it does not
  return the raw serve error from `Stop`.
- gRPC `Server.Serve` returns `nil` when `Stop` or `GracefulStop` is called
  after serving has started. `grpc.ErrServerStopped` is only returned when
  `Serve` is called after the server was already stopped. Do not flag normal
  DI-managed gRPC shutdown as a serve-failure bug based solely on
  `ErrServerStopped` speculation; report only a concrete supported lifecycle
  path that demonstrates `Serve` is invoked after `Stop`/`GracefulStop` and
  causes an incorrect exit code.
- `telemetry.Register()` installs the global OpenTelemetry propagator.
- OTLP exporter endpoints intentionally come from explicit go-service config
  fields such as `telemetry.logger.url`, `telemetry.metrics.url`, and
  `telemetry.tracer.url`. Standard OpenTelemetry endpoint environment variables
  such as `OTEL_EXPORTER_OTLP_ENDPOINT` are not fallback sources and should not
  be projected into config by default. Operators that want env-managed endpoints
  should set the go-service config values through their deployment/config
  source; do not flag missing automatic `OTEL_*` endpoint projection as a
  feature gap unless the documented support boundary changes.
- Prometheus pull metrics are intentionally exposed through the service HTTP
  transport at `/<name>/metrics` when HTTP transport and Prometheus metrics are
  enabled. Do not flag the absence of a Prometheus scrape endpoint on the debug
  server, or the absence of a built-in scrape endpoint for gRPC-only/non-HTTP
  services, as a feature, reliability, or operator gap. Services that want
  Prometheus pull metrics should enable the HTTP transport endpoint, expose a
  service-owned scrape route, or use OTLP push metrics. Report only concrete
  bugs such as the documented HTTP metrics route not being registered, ignored
  explicit metrics config, or a changed public support boundary.
- Telemetry logger, metrics, and tracer setup installs process-global
  OpenTelemetry providers. Do not flag provider globals leaking after DI startup
  failure solely because lifecycle `OnStop` does not run; supported service
  startup failure exits the process. Report only concrete same-process reuse bugs
  in supported tests/tools, ignored successful shutdown cleanup, or an API
  promise that failed startup is recoverable in the same process.
- Process-global telemetry logger state, including the OpenTelemetry logger
  provider and the process-wide `slog` default installed by the OTLP logger, is
  intentionally not reset during normal process shutdown. Do not flag stale
  logger/provider globals after a clean lifecycle stop solely because a
  hypothetical same-process app instance could be started after shutdown.
  Supported service shutdown exits the process; report only concrete reuse bugs
  in supported tests/tools that actually continue running after shutdown and
  require isolated global logger state.
- Telemetry logger, metrics, and tracer shutdown hooks intentionally ignore
  provider/exporter shutdown errors so one telemetry flush failure does not stop
  later lifecycle shutdown hooks. Do not flag swallowed telemetry shutdown
  export errors as reliability gaps solely because operators will not see those
  final flush failures; report only concrete bugs such as shutdown hooks not
  running, globals not resetting after successful shutdown, or a public API
  promise to surface telemetry shutdown errors.
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
- Configured byte-size limits that drive buffering paths are capped by
  `bytes.MaxConfigSize` at the supported config validation boundary, including
  cache value sizes and server receive sizes used by HTTP and gRPC transports.
  Do not flag `MaxInt`, allocator, or `limit+1` overflow speculation from
  absurd configured sizes unless a supported config/DI path bypasses that
  validation or the public API starts promising safe arbitrary direct
  constructor limits.
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
- Redis cache isolation should use Redis URL database selection, such as
  `/0` through `/15`, a dedicated endpoint, or deployment-level isolation. Do
  not flag missing service-name key namespacing or prefix-scoped Redis
  `Cache.Flush`; implementing that cleanup requires client-side key iteration,
  while the supported Redis flush behavior is `FLUSHDB` against the selected
  database.
- PostgreSQL DSN security options, including TLS/`sslmode`, are intentionally
  part of the DSN supplied by the service configuration. `database/sql/pg`
  passes resolved DSNs through to pgx and does not impose repository-level DSN
  construction policy. Do not flag pass-through DSN handling merely because an
  insecure DSN could be supplied; only report concrete bugs such as accidental
  DSN rewriting, secret leakage, or a public API promise to enforce secure DSN
  policy.
- SQL OpenTelemetry query text capture is disabled by the supported
  `database/sql/pg` DI wiring, and ping spans are disabled by the upstream
  zero-value `otelsql.SpanOptions`. Do not flag lower-level manual
  `database/sql/telemetry.WithSpanOptions` use solely because a caller could
  replace those defaults unless a supported config/DI path exposes that option,
  raw SQL text is emitted by repository code, or the public API starts promising
  safe merging of arbitrary span options.
- RSA public-key loading intentionally validates the repository's key-size
  policy and delegates deeper RSA parameter checks, such as public exponent
  usability, to the standard library operations that consume the key. Do not
  flag `crypto/rsa.Config.PublicKey` merely because `Encrypt` can later return a
  standard-library error for an unusual parsed key parameter; only report
  concrete bugs such as panics, accepted weak key sizes, secret leakage, or a
  public API promise to fully validate every RSA parameter at load time.
- RSA key loading intentionally supports the PKCS#1 PEM formats generated by
  this repository's RSA generator: `RSA PUBLIC KEY` and `RSA PRIVATE KEY`.
  Do not flag missing support for generic PKIX `PUBLIC KEY` or PKCS#8
  `PRIVATE KEY` RSA PEM blocks as a feature gap solely for external-tool
  compatibility. Report only concrete bugs where repository-generated RSA keys
  fail to load, weak keys are accepted, secrets leak, or the documented support
  boundary changes.
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
- HTTP retry intentionally does not apply `transport/retry.Config.Timeout` as a
  per-attempt timeout. HTTP response bodies are caller-owned after `RoundTrip`
  returns, and tying retry-owned attempt contexts to returned bodies requires
  response body wrapping that can hide optional body interfaces. Bound outbound
  HTTP calls with the request context or `http.Client.Timeout`; do not flag the
  absence of HTTP retry per-attempt timeout unless a public API starts promising
  that timeout or the retry layer reintroduces owned response-body lifecycle
  handling.
- HTTP retry intentionally runs the first attempt directly after its initial
  request-context cancellation check, then uses the retry/backoff helper only
  for later attempts. This keeps request-body ownership simple: an already
  canceled request is closed locally, while any non-canceled request body is
  handed to the inner `RoundTripper` on the first attempt. Do not reintroduce
  per-attempt ownership flags or synthetic tests for the tiny cancellation
  window before the first attempt unless the retry architecture changes. Local
  redirect sentinels such as `net/http.ErrUseLastResponse` are terminal retry
  outcomes, not transport failures to retry.
- HTTP and gRPC client retry/load-control ordering intentionally match at the
  supported transport stack level: metadata is outside retry so `Request-Id`
  stays logical-request scoped, retry wraps the client limiter and breaker so
  each attempt consumes local quota and breaker capacity, and token generation
  remains inside retry so each wire attempt gets a fresh token. HTTP local
  limiter and breaker rejections are marked with `net/http/status.LocalError`
  and are terminal, not retried as upstream 429/503 responses. Do not flag
  retry/load-control ordering unless a supported stack regresses from this
  contract, local load-control rejections become retryable, or a public API
  starts promising a different composition.
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
- Transport limiter keys are `"user-id"`, `"transport-service-method"`,
  `"service-method"`, `"ip"`, and `"user-agent"`; `"token"` is
  intentionally not a limiter key. Server limiters run after metadata
  extraction and token verification, so `"user-id"` is the verified principal
  (JWT/PASETO subject or SSH key name), and missing, malformed, or invalid auth
  is rejected before the limiter by design. Do not flag that bypass; use an
  external edge, gateway, ingress, load balancer, or service mesh limiter when
  those attempts need quota enforcement.
- The built-in transport limiter is intentionally in-memory and per-process. Treat it as a last-resort local safeguard; prefer external edge/gateway/ingress/load-balancer/service-mesh limiting for production abuse protection.
- go-service is a microservices framework; browser-facing concerns such as
  CORS are expected to live at a BFF, API gateway, ingress, CDN, or other edge
  layer. Do not flag missing built-in CORS/preflight support solely because
  browser clients cannot call authenticated service endpoints cross-origin
  through the standard HTTP stack. Report only concrete bugs where a public API
  promises browser-direct support, an existing edge/BFF integration is broken,
  or the repo adds first-class CORS/pre-auth middleware semantics and violates
  them.
- Transport limiter `max_keys` intentionally caps the number of
  caller-derived keys that get independent in-memory buckets. Additional
  distinct keys share one overflow bucket; do not flag this as accidental key
  collision unless explicit limiter config is ignored, the overflow bucket is
  bypassed, or documented status/header behavior is wrong.
- gRPC stream limiters intentionally consume one token when a stream opens and
  one token for each `RecvMsg` and `SendMsg` operation. Do not flag this as
  accidental double-counting; report only concrete bugs such as missing message
  limiting, ignored explicit limiter config, or incorrect status/header
  behavior.
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
- MVC views should be constructed during startup or route registration so
  missing, unreadable, or malformed templates fail fast before serving traffic.
  Do not flag `mvc.NewFullView`, `mvc.NewPartialView`, or `mvc.NewViewPair`
  panics as request-path reliability gaps solely because those constructors
  panic; report only concrete supported paths that construct views per request
  contrary to the documented lifecycle.
- MVC static files are expected to be served from embedded or otherwise stable
  application assets. Do not flag ignored mid-stream `io.Copy` errors in
  `mvc.writeStaticFile` as reliability gaps solely because an artificial
  `fs.FS` can return partial data after successful `Open`/`Stat`; report only
  concrete supported filesystem paths where static reads can fail mid-stream
  and operators need different behavior.
- MVC not-found handling intentionally uses a simple `Accept` header check for
  `text/html` and does not fully evaluate quality weights such as
  `text/html;q=0`. Do not flag this unless there is concrete evidence of a
  real client, route, or deployment contract that depends on strict weighted
  `Accept` negotiation for MVC 404 responses.
- JWT verification requires both the expected algorithm and a `kid` header.
- JWT, PASETO, and SSH token key material is intentionally loaded and checked
  by the runtime `Generate` and `Verify` paths. Do not flag missing startup
  warmup/validation for token key sources solely because bad, missing,
  unreadable, or malformed administrator-supplied key material can surface
  during token issuance or verification. Do not recommend duplicating runtime
  token checks at startup merely to reject an empty or incomplete trusted key
  set, including SSH `keys`, unless a supported config contract explicitly
  promises that structural key usability is validated before runtime token
  operations. Report only concrete bugs such as panics, accepted weak or wrong
  key types/sizes, secret leakage, ignored documented token config validation,
  or a public API promise that token key material is fully resolved before
  runtime token operations.
- Token tests should focus on repository-owned wrapper behavior. Do not flag or
  add tests that require hand-crafting upstream JWT/PASETO internals solely to
  prove library-owned parsing, signature, expiration, not-before, or
  missing-claim behavior. Prefer tests through this repository's
  `Generate`/`Verify` APIs. Only test crafted tokens when the repository adds
  local validation logic or a public contract beyond upstream library behavior.
- `telemetry/header.Map.MustSecrets` can panic if secret resolution fails during config projection.
- Health registration helpers require `*net/http.ServeMux`.
- Health checks intentionally use go-health registration and observer mapping
  directly. Service code may colocate `server.Register` and `server.Observe`
  calls in one DI function for `healthz`, `livez`, `readyz`, and `grpc`
  observers. Do not flag the absence of a standard health probe composition
  helper solely because observer names are hand-mapped; report only concrete
  broken behavior such as missing documented endpoints, ignored observer
  errors, wrong probe names, or a public API promise that standard probes are
  auto-composed.
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
- OpenFeature trace event attributes are produced by the upstream
  `hooks.NewTracesHook` implementation and may include semantic-convention
  fields such as the evaluation context targeting key and evaluated flag value.
  Treat this as documented upstream instrumentation behavior; protect telemetry
  exporter and backend access controls. Do not flag it as a local `feature`
  issue unless this repository adds local feature trace attribute construction,
  logs/exports those values independently, or starts promising sanitized feature
  trace events by default.
- OpenFeature registration uses process-global SDK state. Do not flag hooks or
  provider globals leaking after DI startup failure solely because `OnStop` does
  not run; supported service startup failure exits the process. Report only
  concrete same-process reuse bugs in supported tests/tools, ignored successful
  shutdown cleanup, or an API promise that failed startup is recoverable in the
  same process.
- The UUIDv7 generator intentionally calls `google/uuid.EnableRandPool` at
  package init time. UUID is the default ID generator and sits on request
  metadata hot paths; the process-wide heap-backed random pool tradeoff is
  accepted for these operational identifiers, which are not secrets or bearer
  tokens. Do not flag this as a global-state or security issue unless the
  generator starts being used for secret material, the upstream pool semantics
  change materially, or a public API starts promising no `google/uuid` global
  mutation.

## Testing, Style, And Docs

- Tests commonly use `stretchr/testify/require`.
- When adding test coverage, first follow the existing test shape in the
  package. Prefer extending current fixture/table/assertion helpers over adding
  standalone tests for behavior already exercised nearby.
- Config tests usually use fixture-driven `config.NewConfig` coverage plus
  `verifyConfig`; do not add separate decoder-routing or Fx projection tests
  unless they cover a distinct repository-owned behavior not already exercised
  by those fixture tests.
- Do not add build-tagged or architecture-specific tests unless CI actually
  runs that build tag or architecture.
- Go files use tabs; YAML uses 2 spaces per `.editorconfig`.
- Every exported identifier, including under `internal/test/**`, needs a GoDoc comment.
- GoDoc comments should start with the identifier name or `Deprecated:`.

## CI

- CI initializes submodules, prepares certs/services, then runs dependency, lint, security, spec, benchmark, and coverage targets.
- Check CI config for exact service dependencies, ports, and command order.
- CircleCI owns the selected dependency sidecars. Do not flag the sidecar
  selection as a reliability, release-safety, or reproducibility gap solely
  because it is CI-owned; report only concrete breakage such as CI no longer
  starting the required service, waiting on the wrong port, or using a
  dependency that no longer satisfies the documented test dependency.
- The `time` package intentionally exercises multiple live public NTP/NTS
  providers in normal tests as smoke coverage for the network time adapters.
  Do not flag those tests or the `make specs` CI path solely because they
  perform internet/network time queries. Report only concrete breakage such as
  all configured providers currently failing in CI, ignored timeouts, removed
  provider redundancy, or a documented promise of hermetic/offline specs.
