include bin/build/make/go.mak
include bin/build/make/git.mak

# Diagrams generated from https://github.com/loov/goda.
diagrams: crypto-diagram database-diagram telemetry-diagram transport-diagram

crypto-diagram:
	$(MAKE) package=crypto create-diagram

database-diagram:
	$(MAKE) package=database create-diagram

telemetry-diagram:
	$(MAKE) package=telemetry create-diagram

transport-diagram:
	$(MAKE) package=transport create-diagram

# Run all the benchmarks.
benchmarks: http-benchmarks grpc-benchmarks

http-benchmarks:
	$(MAKE) package=transport/http benchmark

grpc-benchmarks:
	$(MAKE) package=transport/grpc benchmark
