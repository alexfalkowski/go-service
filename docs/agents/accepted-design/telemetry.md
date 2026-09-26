# Accepted Design: Telemetry

These entries are part of the [accepted design index](../accepted-design.md)
and carry the same mandatory weight.

- `telemetry.RegisterPropagation(...)` installs the global OpenTelemetry
  propagator.
- OTLP exporter endpoints intentionally come from explicit go-service config
  fields such as `telemetry.logger.url`, `telemetry.metrics.url`, and
  `telemetry.tracer.url`. Standard OpenTelemetry endpoint environment variables
  such as `OTEL_EXPORTER_OTLP_ENDPOINT` are not fallback sources and should not
  be projected into config by default. Operators that want env-managed endpoints
  should set the go-service config values through their deployment/config
  source; do not flag missing automatic `OTEL_*` endpoint projection as a
  feature gap unless the documented support boundary changes.
- OTLP header maps are intentionally passed to the selected exporter after
  source resolution without repository-owned HTTP or gRPC syntax validation.
  Missing `env:` values and unreadable `file:` values still fail startup, but
  malformed administrator-supplied names or resolved values may be rejected by
  the exporter only when it attempts an export. Do not flag that upstream
  protocol rejection as a local code issue. Report only concrete local bugs such
  as valid headers being altered or dropped, explicit headers being ignored,
  secret values being exposed, or a public promise of local syntax validation.
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
- The effective log level is intentionally resolved once from
  `telemetry.logger.level` at startup and installed as the process-wide `slog`
  default; there is no runtime log-level mutation endpoint or `slog.LevelVar`
  toggle. The debug server is off by default and not run in production, so any
  workflow that could reach such an endpoint is already a configuration change
  where `telemetry.logger.level` can be set directly. Do not flag the absence of
  a runtime/debug log-level control endpoint or dynamic level toggle as a feature
  or operability gap; operators change verbosity through `telemetry.logger.level`
  config. Report only concrete bugs where the configured level is ignored or a
  public API promises runtime level mutation.
- SQL OpenTelemetry query text capture is disabled by the supported
  `database/sql/pg` DI wiring, and ping spans are disabled by the upstream
  zero-value `otelsql.SpanOptions`. Do not flag lower-level manual
  `database/sql/telemetry.WithSpanOptions` use solely because a caller could
  replace those defaults unless a supported config/DI path exposes that option,
  raw SQL text is emitted by repository code, or the public API starts promising
  safe merging of arbitrary span options.
- SQL OpenTelemetry error recording intentionally preserves upstream `otelsql`
  error events and span status. Those events use `err.Error()` as
  `exception.message`, and PostgreSQL may echo bound parameter values in error
  text even while raw query capture is disabled. Trace telemetry is privileged
  operator diagnostic data and should be protected by backend access controls.
  Do not flag upstream SQL error messages solely because they may contain bound
  values; report only concrete local bugs such as raw query text, credentials,
  or DSNs emitted independently by repository code, ignored explicit telemetry
  options, or a public API promise of sanitized SQL span errors.
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
- `telemetry/header.Map.MustSecrets` can panic if secret resolution fails during config projection.
- OpenFeature trace event attributes are produced by the upstream
  `hooks.NewTracesHook` implementation and may include semantic-convention
  fields such as the evaluation context targeting key and evaluated flag value.
  Treat this as documented upstream instrumentation behavior; protect telemetry
  exporter and backend access controls. Do not flag it as a local `feature`
  issue unless this repository adds local feature trace attribute construction,
  logs/exports those values independently, or starts promising sanitized feature
  trace events by default.
