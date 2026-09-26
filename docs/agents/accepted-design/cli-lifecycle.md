# Accepted Design: CLI And Lifecycle

These entries are part of the [accepted design index](../accepted-design.md)
and carry the same mandatory weight.

- Manual server lifecycle wiring should use `net/server.Register(...)`.
- Module-related behavior is tested through CLI/server application wiring. Do
  not flag missing direct package tests for Fx module provider inventory,
  module composition, or `transport.Module` lifecycle registration solely
  because they are not asserted in the package that declares the module. Report
  only concrete broken behavior through the CLI or supported `module.Server`
  path, or an explicit public promise of lower-level standalone module use.
- `cli.RunCode` returns `os.ExitCodeSuccess` on success, preserves non-zero
  shutdown exit codes requested through `di.ExitCode(...)`, and otherwise
  returns `os.ExitCodeFailure`.
- `cli.Application.AddClient` intentionally models short-lived client command
  work as DI/Fx startup work. Client commands perform their main action from a
  lifecycle `OnStart` hook, then stop the graph immediately after startup
  completes. Constructors and invocations only wire dependencies and register
  lifecycle hooks; they do not perform the command action while the graph is
  being built. Do not flag the absence of a separate post-DI command-task API
  solely because command work lives in `OnStart`; this is the supported pattern,
  as used by downstream client templates. Report only concrete broken behavior
  such as incorrect error propagation, ignored shutdown exit codes, lifecycle
  ordering bugs, or a documented command contract that cannot be expressed with
  the DI lifecycle.
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
- gRPC graceful shutdown delegates to upstream `grpc.Server.GracefulStop`,
  which stops accepting new RPCs and waits for active handlers but does not
  cancel stream handler contexts or interrupt an in-progress `RecvMsg`.
  Long-lived stream handlers that need to finish before the lifecycle stop
  deadline must use a bounded or application-owned cancellation workflow;
  otherwise the lifecycle force-stops the server at the deadline. Do not flag
  this upstream lifecycle boundary as a generic transport drain-interceptor
  gap. Report only a concrete repository-owned handler that fails to honor its
  documented bounded or application-owned shutdown workflow.
- OpenFeature registration uses process-global SDK state. Do not flag hooks or
  provider globals leaking after DI startup failure solely because `OnStop` does
  not run; supported service startup failure exits the process. Report only
  concrete same-process reuse bugs in supported tests/tools, ignored successful
  shutdown cleanup, or an API promise that failed startup is recoverable in the
  same process.
