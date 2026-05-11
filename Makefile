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
benchmarks: http-benchmarks grpc-benchmarks sql-benchmarks cache-benchmarks bytes-benchmarks strings-benchmarks id-benchmarks http-content-benchmarks

http-benchmarks:
	@make package=transport/http benchtime=100x benchmark

grpc-benchmarks:
	@make package=transport/grpc benchtime=100x benchmark

sql-benchmarks:
	@make package=database/sql/driver benchtime=100x benchmark

cache-benchmarks:
	@make package=cache/driver benchtime=100x benchmark

bytes-benchmarks:
	@make package=bytes benchtime=100x benchmark

strings-benchmarks:
	@make package=strings benchtime=100x benchmark

id-benchmarks:
	@make package=id benchtime=100x benchmark

http-content-benchmarks:
	@make package=net/http/content benchtime=100x benchmark

# Generate for tests.
generate:
	@make -C internal/test generate
