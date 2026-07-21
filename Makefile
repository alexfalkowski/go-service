fuzztime ?= 1000x

include bin/build/make/help.mak
include bin/build/make/go.mak
include bin/build/make/git.mak
include bin/build/make/claude.mak
include bin/build/make/codex.mak

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
benchmarks: http-benchmarks grpc-benchmarks limiter-benchmarks sql-benchmarks cache-benchmarks bytes-benchmarks \
	strings-benchmarks id-benchmarks net-http-benchmarks http-content-benchmarks

http-benchmarks:
	@$(MAKE) package=transport/http benchtime=100x benchmark

grpc-benchmarks:
	@$(MAKE) package=transport/grpc benchtime=100x benchmark

limiter-benchmarks:
	@$(MAKE) package=transport/limiter benchtime=100x benchmark

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

net-http-benchmarks:
	@$(MAKE) package=net/http benchtime=100x benchmark

http-content-benchmarks:
	@$(MAKE) package=net/http/content benchtime=100x benchmark

# Run bounded fuzz tests. Set fuzztime=<duration-or-count> to override the default 1000 executions per target.
fuzzes: bytes-fuzz time-fuzz encoding-fuzz compress-fuzz net-fuzz

bytes-fuzz:
	@$(MAKE) package=bytes name=FuzzSizeTextRoundTrip fuzz
	@$(MAKE) package=bytes name=FuzzSizeJSONRoundTrip fuzz

time-fuzz:
	@$(MAKE) package=time name=FuzzDurationTextRoundTrip fuzz
	@$(MAKE) package=time name=FuzzDurationJSONRoundTrip fuzz

encoding-fuzz:
	@$(MAKE) package=encoding/bytes name=FuzzEncoder fuzz
	@$(MAKE) package=encoding/gob name=FuzzUnmarshal fuzz
	@$(MAKE) package=encoding/hjson name=FuzzUnmarshal fuzz
	@$(MAKE) package=encoding/json name=FuzzUnmarshal fuzz
	@$(MAKE) package=encoding/msgpack name=FuzzUnmarshal fuzz
	@$(MAKE) package=encoding/toml name=FuzzUnmarshal fuzz
	@$(MAKE) package=encoding/yaml name=FuzzUnmarshal fuzz

compress-fuzz:
	@$(MAKE) package=compress/none name=FuzzCompressor fuzz
	@$(MAKE) package=compress/s2 name=FuzzCompressor fuzz
	@$(MAKE) package=compress/s2 name=FuzzDecompress fuzz
	@$(MAKE) package=compress/snappy name=FuzzCompressor fuzz
	@$(MAKE) package=compress/snappy name=FuzzDecompress fuzz
	@$(MAKE) package=compress/zstd name=FuzzCompressor fuzz
	@$(MAKE) package=compress/zstd name=FuzzDecompress fuzz

net-fuzz:
	@$(MAKE) package=net/grpc name=FuzzParseServiceMethod fuzz
	@$(MAKE) package=net/header name=FuzzParseBearer fuzz
	@$(MAKE) package=net/http name=FuzzParseServiceMethod fuzz
	@$(MAKE) package=net/http/media name=FuzzParse fuzz
	@$(MAKE) package=net/url name=FuzzSplitPath fuzz

# Generate for tests.
generate:
	@$(MAKE) -C internal/test generate

# Check generated test protobuf outputs are fresh.
generate-stale:
	@$(MAKE) -C internal/test stale
