# AGENTS.md

## Shared Skill

Use the shared `coding-standards` skill from `bin/skills/coding-standards` for code changes, bug fixes, refactors, reviews, tests, linting, documentation, PR summaries, commits, Makefile changes, CI validation, and verification.

Before answering any covered task, read `bin/skills/coding-standards/SKILL.md` and the relevant reference file under `bin/skills/coding-standards/references/`.

- Commit messages and PR summaries: read `references/pr.md` and follow its required format exactly.
- Reviews: read `references/review.md`.
- Verification: read `references/verification.md`.
- Go changes: read `references/go.md`, plus workflow or change-safety references when relevant.

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
- Many config fields use source strings through `os.FS.ReadSource`: `env:NAME`, `file:/path`, or a literal value.
- Nil pointer sub-configs usually mean "disabled".

## Gotchas

- Manual transport TLS setup must call the HTTP and gRPC transport register functions; `transport.Module` does this on the normal path.
- Manual server lifecycle wiring should use `net/server.Register(...)`.
- `telemetry.Register()` installs the global OpenTelemetry propagator.
- `cache.Register(...)` sets the package-level cache used by generic cache helpers.
- Redis cache config intentionally expects `cache.options.url` to exist and be a string.
- Access policy config is passed to Casbin's file adapter, so it must be a real path.
- IP metadata intentionally trusts forwarding headers; deploy behind trusted proxies that strip spoofed headers before using the `"ip"` limiter key.
- JWT verification requires both the expected algorithm and a `kid` header.
- `telemetry/header.Map.MustSecrets` can panic if secret resolution fails during config projection.
- Health registration helpers require `*net/http.ServeMux`.
- Shared metadata, header, and string helpers live under `net/...`, not `transport/...`.
- `vendor/` is gitignored and regenerated via `make dep`.

## Testing, Style, And Docs

- Tests commonly use `stretchr/testify/require`.
- Go files use tabs; YAML uses 2 spaces per `.editorconfig`.
- Every exported identifier, including under `internal/test/**`, needs a GoDoc comment.
- GoDoc comments should start with the identifier name or `Deprecated:`.

## CI

- CI initializes submodules, prepares certs/services, then runs dependency, lint, security, spec, benchmark, and coverage targets.
- Check CI config for exact service images, ports, and command order.
