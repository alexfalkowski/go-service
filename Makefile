include bin/build/make/go.mak
include bin/build/make/git.mak

# Diagrams generated from https://github.com/loov/goda.
diagrams: crypto-diagram cache-diagram database-diagram telemetry-diagram transport-diagram

cache-diagram:
	$(MAKE) package=cache create-diagram

crypto-diagram:
	$(MAKE) package=crypto create-diagram

database-diagram:
	$(MAKE) package=database create-diagram

telemetry-diagram:
	$(MAKE) package=telemetry create-diagram

transport-diagram:
	$(MAKE) package=transport create-diagram

# Run all the benchmarks.
benchmarks:
	$(MAKE) package=transport/http benchmark
