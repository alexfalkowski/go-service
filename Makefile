include bin/build/make/help.mak
include bin/build/make/go.mak
include bin/build/make/git.mak

# Diagrams generated from https://github.com/loov/goda.
diagrams: crypto-diagram database-diagram telemetry-diagram transport-diagram

crypto-diagram:
	@make package=crypto create-diagram

database-diagram:
	@make package=database create-diagram

telemetry-diagram:
	@make package=telemetry create-diagram

transport-diagram:
	@make package=transport create-diagram

# Run all the benchmarks.
benchmarks: http-benchmarks grpc-benchmarks bytes-benchmarks

http-benchmarks:
	@make package=transport/http benchmark

grpc-benchmarks:
	@make package=transport/grpc benchmark

bytes-benchmarks:
	@make package=bytes benchmark

# Generate for tests.
generate:
	@make -C internal/test generate
