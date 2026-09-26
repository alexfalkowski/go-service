# Accepted Design: Health And Debug

These entries are part of the [accepted design index](../accepted-design.md)
and carry the same mandatory weight.

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
- Build and version provenance is intentionally owned by the container image and
  deployment platform, not a service endpoint. Services built with this framework
  learn their version from the image tag/metadata supplied at build/deploy time,
  surfaced through `service.version`/`deployment.environment.name` telemetry
  resource attributes and static log attributes, and the platform already tracks
  which image a pod runs. Do not flag the absence of a build-info/version debug
  endpoint (for example `runtime/debug.ReadBuildInfo`, `expvar`, or
  `/debug/vars`) as a feature or operability gap; report only concrete bugs where
  configured version/environment attributes are dropped from telemetry or logs.
- Health registration helpers require `*net/http.ServeMux`.
- Health checks intentionally use go-health registration and observer mapping
  directly. Service code may colocate `server.Register` and `server.Observe`
  calls in one DI function for `healthz`, `livez`, `readyz`, and `grpc`
  observers. Do not flag the absence of a standard health probe composition
  helper solely because observer names are hand-mapped; report only concrete
  broken behavior such as missing documented endpoints, ignored observer
  errors, wrong probe names, or a public API promise that standard probes are
  auto-composed.
- The `healthz`/`livez`/`readyz` probes intentionally return only the plain-text
  `SERVING`/`503` probe contract, not a per-check status breakdown. This
  framework targets container platforms where per-dependency health aggregation,
  rollout gating, and outage attribution live at the orchestration and
  observability layer (Kubernetes probes plus events, service mesh, dashboards),
  and a service owner who wants a structured per-check detail route can register
  one. Do not flag the absence of an aggregated health-detail endpoint (per-check
  JSON from the go-health `Observer.Errors()` map) as a feature or operability
  gap; report only concrete bugs such as a documented probe endpoint missing,
  ignored observer errors, or wrong probe names/status codes.
