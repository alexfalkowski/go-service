# Accepted Design: Transport

These entries are part of the [accepted design index](../accepted-design.md)
and carry the same mandatory weight.

- Manual transport TLS setup must call the HTTP and gRPC transport register functions; `transport.Module` does this on the normal path.
- Runtime TLS material, like other resolved configuration, is intentionally
  materialized once at startup; the leaf certificate and CA pools are pinned for
  the process lifetime. This framework targets container platforms where
  certificate rotation and other configuration changes are rolled out by
  restarting or recreating the workload (cert-manager renewal plus rolling
  restart, mounted-secret refresh plus pod recreation), so there is no
  `GetCertificate`/`GetClientCertificate` callback, fsnotify watch, or SIGHUP
  reload. Do not flag the absence of TLS certificate hot-reload as a reliability
  or operability gap; report only concrete bugs such as certs not being resolved
  at startup, ignored explicit TLS config, or a public API promise of in-process
  certificate rotation.
- `transport.Module` is normally consumed through `module.Server`, which also
  wires `debug.Module`; do not flag `transport.NewServers` requiring
  `*debug.Server`. Mentions that `transport.Module` handles transport
  registration, TLS filesystem registration, or lifecycle registration do not
  promise a complete standalone server bundle without `module.Server`. Report
  this only if a public API explicitly promises standalone `transport.Module`
  composition without the standard `module.Server` bundle.
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
- gRPC user-agent is intentionally a connection-level protocol header managed
  by grpc-go and configured with `grpc.WithUserAgent`. Do not introduce a
  request-scoped or logical user-agent metadata header, or otherwise propagate
  context/outgoing user-agent values across gRPC calls. Do not flag the
  resulting difference between local client metadata and downstream gRPC
  user-agent attribution as a code issue unless the supported contract changes.
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
- HTTP operation route patterns registered through
  `net/http.Router.HandleRoute` with `net/http.WithRouteOperation` are expected
  to include an HTTP method prefix, such as `"GET /<name>/metrics"`, matching
  the supported callers in `transport/http/health` and
  `transport/http/telemetry/metrics`. `HandleRoute` stores only the path portion
  for operations, so a bare method-less
  pattern registers under the empty key and never matches; operation matching is
  intentionally path-only so method mismatches still reach the mux for normal
  method handling. Do not flag method-less operation patterns failing to match,
  or the `WithRouteOperation`-versus-`WithRouteUnauthenticated` handling asymmetry, as a bug;
  method-prefixing is the supported registration form. Report only concrete bugs
  where a method-prefixed operation pattern fails to match or the supported
  callers stop prefixing.
- Unary HTTP REST and RPC response encoding
  (`net/http/content/unary.Content.NewFromAccept`) intentionally uses the first
  `Accept` media type, falling back to `Content-Type` and then JSON, as
  documented in `docs/transport.md` and the `net/http/rest` and `net/http/rpc`
  package docs. It does not evaluate quality weights, so a leading excluded
  range such as `application/yaml;q=0, application/json` still selects YAML.
  RFC 9110 §12.5.1 permits a server to disregard `Accept`, and go-service
  clients set it from one configured `Options.Accept` value. Streaming routes
  (`net/http/content/stream.Content.NewFromAccept`) evaluate the full list,
  including `q=0`, because an unsatisfiable stream `Accept` is rejected rather
  than answered in a different wire format. Do not flag the unary behavior or the
  unary/stream asymmetry; report only concrete bugs where the first `Accept`
  media type is not the one used, or a public API starts promising weighted
  unary negotiation.
- HTTP `net/http/client.Options.ContentType` is expected to be a real encodable
  request media type. Error media types (`text/error` and other `*/error`
  subtypes) are internal error-response media with no encoder, so supplying one
  together with a request body is rejected rather than encoded: `Client.Do`
  returns `http: encode: unary: unsupported media` from `newRequest`'s
  nil-encoder guard and sends no request. This is intentional: error media are
  not valid request content types and callers control `ContentType`. Do not flag
  that rejection, or the matching nil-encoder guard on `Client.Do`'s
  response-decode path, as a bug based on supplying an error media type; report
  only concrete bugs where a valid request media type fails to encode.
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
  (JWT/PASETO subject), and missing, malformed, or invalid auth
  is rejected before the limiter by design. Do not flag that bypass; use an
  external edge, gateway, ingress, load balancer, or service mesh limiter when
  those attempts need quota enforcement.
- The built-in transport limiter is intentionally in-memory and per-process. Treat it as a last-resort local safeguard; prefer external edge/gateway/ingress/load-balancer/service-mesh limiting for production abuse protection.
- The built-in transport limiter's memory store returns an error from `Take`
  only after the store has been closed by its lifecycle hook. Do not report
  error classification, retry, or status behavior demonstrated solely by
  explicitly closing a limiter and then reusing it. Before recording such a
  finding, prove either a non-shutdown `Take` error or a concrete supported
  request path that can reach the limiter after or concurrently with its
  lifecycle close; exported `Close`/`Take` methods alone do not establish that
  post-close reuse is supported.
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
- HTTP streaming routes (`net/http/content/stream.Stream`/`RequestStream`) mirror the
  gRPC stream limiter behavior above: `transport/http/limiter.Handler` still
  charges one token when the request/stream opens (its existing single
  `TakeDecision` call, unchanged for every route), and on top of that it stores
  the underlying `*transport/limiter.Limiter` on the request context via
  `net/http/meta.WithLimiter` when the request is allowed. `Stream.Send` and
  `RequestStream.Recv` retrieve it via `net/http/meta.Limiter` and each charge
  one additional token per message, so a streaming route pays the same
  stream-open-plus-per-message cost gRPC streams do; non-streaming routes never
  read that context value, so it is harmless overhead for them. A denial
  discovered by `Send`/`Recv` before the first successful `Send` is an ordinary
  pre-commit `429` through `status.Error`, matching the middleware's own
  rejection; a denial discovered after the response is committed cannot be a
  `429` (headers are already sent), so it aborts the response like any other
  post-commit streaming error. The
  `RateLimit`/`RateLimit-Policy` response headers are set once, from the
  stream-open decision only, and are not re-emitted per message, since HTTP
  headers cannot change after a streaming response is committed — unlike
  gRPC, which can fall back to trailers. Do not flag the missing per-message
  header updates or the context-value plumbing as bugs; report only concrete
  issues such as a streaming route never being charged, a charge applied
  without a captured limiter, or a mid-stream denial incorrectly writing a
  status after commit.
- HTTP streaming responses (`net/http/content/stream.NewHandler`,
  `NewRequestHandler`) intentionally opt out of gzip by setting
  `gzhttp.HeaderNoCompression` before the first byte is written. This is a
  deliberate cost/interop decision, not an oversight: per-value flush emits a
  separate gzip flush block per message (poor value for small NDJSON records),
  and `gzhttp`'s deferred `Close` on an aborted stream can still emit a
  complete, valid gzip footer, which is a residual interop risk for a client
  whose gzip reader stops at that footer rather than reading ahead into the
  truncated chunked stream (Go's own client does read ahead and still observes
  the abort). Do not flag the missing compression for streaming responses as a
  performance regression; report only concrete bugs such as the header not
  being set before commit, or a streaming route that is compressed despite it.
- Bidirectional HTTP streaming routes (`net/http/content/stream.NewRequestHandler`,
  and the route helpers that use it: `net/http/rest`'s `StreamRouteRequest`,
  `StreamPost`, `StreamPut`, `StreamPatch`, and `net/http/rpc.StreamRoute`)
  require HTTP/2 (including h2c). `NewRequestHandler` rejects a request
  with `req.ProtoMajor < 2` with `505 HTTP Version Not Supported` before the
  handler runs. This is intentional and measured, not a missing feature: an
  HTTP/1.x request body is buffered ahead of the handler by intermediaries and
  the Go transport, so a bidi handler that both reads the request stream and
  writes the response stream hangs rather than failing outright over h1 — the
  pre-handler rejection is what turns that hang into an immediate, diagnosable
  error. Send-only streaming routes (`net/http/content/stream.NewHandler`, and
  `net/http/rest`'s `StreamRoute`/`StreamGet`, `net/http/rpc` has no send-only
  helper) have no such requirement and stay fully supported on HTTP/1.1
  chunked responses, since only the bidi handle interleaves reads and writes.
  Do not flag the 505 rejection, or the asymmetry between the bidi and
  send-only helpers, as a bug; report only concrete issues such as the gate
  firing for a send-only route, failing to fire for a bidi route, or an h2/h2c
  deployment still hanging instead of failing.
- HTTP client `RoundTripper` implementations that can return locally before
  delegating to another `RoundTripper` must make request-body ownership explicit
  with `net/http.ClosingRoundTripper`. Return the closing adapter's close-body
  flag as true only for local rejection paths, and false after delegating
  because the delegated transport owns `req.Body` closure.
- gRPC server reflection is intentionally always registered by `net/grpc.NewServer`; restrict public exposure at the bind address, TLS/auth, ingress, firewall, or service-mesh boundary.
- gRPC service and method names are generated from Buf-managed proto files such
  as `internal/test/greet/v1/service.proto`, which require package-qualified
  service names and valid RPC method names. Do not flag
  `net/grpc.ParseServiceMethod` merely because a hypothetical manually
  constructed method string like `/pkg.Service/` or `/svc/Get.Name` could split
  or fall back to `root`. Only report a concrete issue if untrusted,
  non-generated method strings are used for a security decision, or a public API
  promises strict validation of arbitrary method strings.
- Terminal writes and closes after a buffered HTTP response has already been
  produced are best-effort transport and cleanup signals, not delivery
  acknowledgements. Go may accept a `ResponseWriter.Write` into internal
  buffers without proving that the client received it, and the standard
  library itself discards equivalent terminal response-copy errors. Do not flag
  ignored final `ResponseWriter.Write`, `bytes.Buffer.WriteTo`, response/body or
  codec `Close` errors solely to observe client disconnects, write-deadline
  expiry, or cleanup failure. Report only a concrete repository-owned contract
  that requires the error to change control flow, preserve data, or feed an
  actionable required diagnostic. This does not apply to active streaming
  `Send`/`Recv`/`Flush` errors, which participate in the documented stream
  abort, limiting, and truncation contracts.
- Shared metadata, header, and string helpers live under `net/...`, not `transport/...`.
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
- HTTP CloudEvents receiver registration intentionally ignores the current
  CloudEvents constructor errors because supported wiring passes no protocol
  options and uses the typed `ReceiverFunc`, which matches the SDK receive
  handler shape. Do not flag this unless dependency behavior or call arguments
  change.
