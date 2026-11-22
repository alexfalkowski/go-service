# 🐘 SQL (Postgres)

[← Back to README](../README.md)

SQL root config is `database/sql.Config`, with Postgres under `sql.pg`.

Postgres config embeds common pool + DSN config (`database/sql/config.Config`), including writer/reader pools. Each role pool owns its `dsns` and `settings`. Enabled SQL pool settings must set positive `max_open_conns` and `max_idle_conns`; `max_idle_conns` must not exceed `max_open_conns`.

`module.Server` and `module.Client` both include `sql.Module`, which currently wires PostgreSQL support via `database/sql/pg.Module`.

Enablement is presence-based: a nil `sql` block or a nil `sql.pg` block disables SQL wiring. When enabled, the pgx stdlib driver is registered under the name `pg`, and reader/writer DSNs are resolved using the source-string rules described in [Source strings](../README.md#-source-strings-secrets-dsns-paths). Enabled PostgreSQL config must provide at least one non-empty `reader.dsns[].source` or `writer.dsns[].source`. Driver instrumentation is installed when tracing or metrics are enabled, OpenTelemetry `database/sql` stats metrics are registered when metrics are enabled, and the resulting pools are closed on lifecycle stop.

When migrating an existing configuration, rename `url` to `source` in each DSN
entry. In Go, use `DSN.Source` and `DSN.Resolve`.

SQL wiring creates `database/sql` pool handles and applies pool settings, but it
does not ping PostgreSQL during construction. Call `DBs.Ping`, `DBs.PingWriter`,
or `DBs.PingReader` with an SLO-appropriate
deadline, or register `health/checker.NewDBChecker`, when startup or readiness
should verify database reachability.

Example (with source strings for DSNs):

```yaml
sql:
  pg:
    reader:
      dsns:
        - source: env:PG_READER_DSN
      settings:
        max_open_conns: 20
        max_idle_conns: 10
        conn_max_idle_time: 30m
        conn_max_lifetime: 1h
    writer:
      dsns:
        - source: env:PG_WRITER_DSN
      settings:
        max_open_conns: 3
        max_idle_conns: 2
        conn_max_idle_time: 10m
        conn_max_lifetime: 30m
```

A literal DSN can use a PostgreSQL URL or a keyword/value connection string.

```yaml
sql:
  pg:
    writer:
      dsns:
        - source: host=localhost dbname=example
      settings:
        max_open_conns: 10
        max_idle_conns: 5
```

## Dependencies

![Dependencies](../assets/database.png)
