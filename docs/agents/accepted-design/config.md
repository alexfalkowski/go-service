# Accepted Design: Config

These entries are part of the [accepted design index](../accepted-design.md)
and carry the same mandatory weight.

- `config.NewConfig[T]` intentionally decodes and validates configuration on the
  normal startup path, so a service started with bad configuration fails fast at
  run time. There is no separate config validate/preflight subcommand (for
  example an `AddCheck` command or `svc check`), and its absence is not a feature
  gap: running the service is the check. Do not flag the missing side-effect-free
  config-only command solely because CI or pre-deploy config linting would prefer
  not to start subsystems; report only concrete bugs where documented config
  validation is skipped, invalid config is accepted, or the decode/validate path
  does not run at startup.
- `cache.Register(...)` sets the package-level cache used by generic cache helpers.
  It is intentionally called by the supported wiring path during startup/test
  setup, not as a concurrent runtime reconfiguration API. Do not flag
  unsynchronized global-state or nil-pointer issues based solely on hypothetical
  concurrent manual `cache.Register`, `cache.Get`, or `cache.Persist` calls
  unless a concrete public API path starts promising concurrent manual
  re-registration or the repository adds such a runtime path.
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
- Before flagging a nil-pointer panic on embedded pointer configuration types,
  inspect the called method. Go permits calling pointer-receiver methods on nil
  pointers, and methods such as `(*config/server.Config).IsEnabled` are
  intentionally nil-safe.
- `vendor/` is gitignored and regenerated via `make dep`.
