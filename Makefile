include bin/build/make/help.mak
include bin/build/make/go.mak
include bin/build/make/git.mak

# Diagrams generated from https://github.com/loov/goda.
diagrams: crypto-diagram database-diagram telemetry-diagram transport-diagram

crypto-diagram:
	@$(MAKE) package=crypto create-diagram

database-diagram:
	@$(MAKE) package=database create-diagram

telemetry-diagram:
	@$(MAKE) package=telemetry create-diagram

transport-diagram:
	@$(MAKE) package=transport create-diagram

# Run all the benchmarks.
benchmarks: http-benchmarks grpc-benchmarks sql-benchmarks cache-benchmarks bytes-benchmarks strings-benchmarks id-benchmarks http-content-benchmarks

http-benchmarks:
	@$(MAKE) package=transport/http benchtime=100x benchmark

grpc-benchmarks:
	@$(MAKE) package=transport/grpc benchtime=100x benchmark

sql-benchmarks:
	@$(MAKE) package=database/sql/driver benchtime=100x benchmark

cache-benchmarks:
	@$(MAKE) package=cache/driver benchtime=100x benchmark

bytes-benchmarks:
	@$(MAKE) package=bytes benchtime=100x benchmark

strings-benchmarks:
	@$(MAKE) package=strings benchtime=100x benchmark

id-benchmarks:
	@$(MAKE) package=id benchtime=100x benchmark

http-content-benchmarks:
	@$(MAKE) package=net/http/content benchtime=100x benchmark

# Generate for tests.
generate:
	@$(MAKE) -C internal/test generate

# Check generated test protobuf outputs are fresh.
generate-stale:
	@$(MAKE) -C internal/test stale
