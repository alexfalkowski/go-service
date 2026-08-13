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
- An exported constructor, setter, helper, option, or other symbol does not
  make every direct invocation or argument a supported contract. Before
  recording an issue that cannot occur through standard module, CLI, template,
  or integration wiring, identify a current supported consumer or an explicit
  contract for the exact construction and triggering input. Package-local tests
  of nil, zero, invalid, or defensive branches do not establish supported usage
  by themselves. If no supported path supplies the triggering input, reject the
  candidate as unsupported misuse.
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

## Compatibility

- Treat compatibility-affecting changes that have landed in repository history
  as accepted decisions, even when they appear semver-incompatible or require
  downstream source changes. Do not flag or reopen them solely by comparing
  current code with older tags, releases, documentation, examples, or exported
  surfaces. Review current code against the current supported behavior, and
  revisit a landed compatibility decision only when the user explicitly asks.

## Accepted Design

Before recording a finding, review candidate, audit entry, or proposed change,
agents MUST read `docs/agents/accepted-design.md`. It records the accepted
designs, intentional tradeoffs, and support boundaries that agents repeatedly
re-flag as bugs, covering transport, TLS, config, telemetry, limiters, tokens,
crypto, debug endpoints, encoding, MVC, and health surfaces. Those rules are
mandatory and carry the same weight as this file. Agents MUST NOT report a
finding in those areas without checking that file first.

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
- Agents verifying changes in a sandboxed shell should run `make specs`
  (optionally `package=<path>` to scope it) rather than a raw `go test`
  invocation, especially for anything touching or depending on the `time`
  package. In at least one sandboxed agent environment, raw `go test` against
  `time`'s NTP/NTS tests failed reproducibly (`dial udp/tcp ...: operation not
  permitted`) while `package=time make specs` passed reproducibly for the same
  code, same machine, same session — `make specs` wraps the same `go test`
  invocation (see `bin/build/go/test`), so the difference is in how that
  sandbox's tooling permission/network layer treats the two invocations, not
  in what the repo's tests do. Do not conclude live NTP/NTS coverage is broken
  from a raw `go test` failure alone; reproduce through `make specs` before
  treating it as a real breakage.
