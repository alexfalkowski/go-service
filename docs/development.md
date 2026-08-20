# 🧑‍💻 Development

[← Back to README](../README.md)

## Style

This repo generally follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

Exported Go identifiers should have GoDoc comments, and each comment should start with the identifier name or `Deprecated:`.

## Development Dependencies

Common repository targets expect these tools on `PATH`:

- `make`
- `gotestsum` for `make specs`
- `fieldalignment` for `make lint`
- `golangci-lint` for full `make lint` coverage (the wrapper no-ops when it is missing)
- `govulncheck` and `trivy` for `make sec`
- `mkcert` for local TLS fixtures and `make create-certs`
- `buf` for `make generate`
- `goda` and Graphviz `dot` for `make diagrams`

## Setup (repo)

This repo uses a `bin/` git submodule for `make` targets.

```sh
git submodule sync
git submodule update --init

mkcert -install
make create-certs

make dep
```

If submodule fetch fails, ensure GitHub SSH access is configured (`.gitmodules` uses `git@github.com:...` URLs).

## Discover targets

```sh
make help
```

## Dependencies (`vendor/` workflow)

```sh
make dep
```

`make dep` runs:

- `go mod download`
- `go mod tidy`
- `go mod vendor`

Tests are run with `-mod vendor`, so after dependency changes run `make dep` before `make specs`.

## Local integration dependencies

`make start` uses the shared Docker-based environment from the sibling
`../docker` repo. It requires Docker and may require GitHub SSH access if that
sibling repo must be fetched.

Start required services:

```sh
make start
```

Stop them:

```sh
make stop
```

## Tests

Run unit tests with race + coverage:

```sh
make specs
```

Artifacts:

- JUnit XML: `test/reports/specs.xml`
- Coverage profile: `test/reports/profile.cov`

## Lint and format

```sh
make lint
make fix-lint
make format
```

## Security checks

```sh
make sec
```

## Benchmarks

```sh
make benchmarks
make http-benchmarks
make grpc-benchmarks
make limiter-benchmarks
make sql-benchmarks
make cache-benchmarks
make bytes-benchmarks
make strings-benchmarks
make id-benchmarks
make net-http-benchmarks
make http-content-benchmarks
```

## Fuzz tests

```sh
make fuzzes
make bytes-fuzz
make time-fuzz
make encoding-fuzz
make compress-fuzz
make net-fuzz
make package=encoding/json name=FuzzUnmarshal fuzztime=10s fuzz
```

## Coverage reports

```sh
make coverage
make html-coverage
make func-coverage
```

## Code generation (Buf)

Root generation targets are for the `internal/test` protobuf fixtures. After
changing those fixtures, regenerate them. To match the CI stale-output check,
run `make generate-stale` from a clean worktree, or after staging the intended
fixture and generated-file changes:

```sh
make generate
make generate-stale
```

## Architecture diagrams

```sh
make diagrams
make crypto-diagram
make database-diagram
make telemetry-diagram
make transport-diagram
```
