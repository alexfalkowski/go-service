module github.com/alexfalkowski/go-service

go 1.19

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/alexfalkowski/go-health v1.11.0
	github.com/avast/retry-go/v3 v3.1.1
	github.com/dgraph-io/ristretto v0.1.1
	github.com/go-redis/cache/v8 v8.4.4
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v4 v4.4.3
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.0
	github.com/hashicorp/go-retryablehttp v0.7.2
	github.com/jackc/pgx/v4 v4.17.2
	github.com/jmoiron/sqlx v1.3.5
	github.com/klauspost/compress v1.15.14
	github.com/linxGnu/mssqlx v1.1.8
	github.com/ngrok/sqlmw v0.0.0-20220520173518-97c9c04efc79
	github.com/nsqio/go-nsq v1.1.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.14.0
	github.com/rs/cors v1.8.3
	github.com/smartystreets/goconvey v1.7.2
	github.com/soheilhy/cmux v0.1.5
	github.com/sony/gobreaker v0.5.0
	github.com/spf13/cobra v1.6.1
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/ulule/limiter/v3 v3.10.0
	github.com/vmihailenco/msgpack/v5 v5.3.5
	go.uber.org/fx v1.18.2
	go.uber.org/multierr v1.9.0
	go.uber.org/zap v1.24.0
	golang.org/x/net v0.5.0
	google.golang.org/genproto v0.0.0-20221207170731-23e4bf6bdc37
	google.golang.org/grpc v1.52.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/DataDog/dd-trace-go.v1 v1.46.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.38.2 // indirect
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state v0.42.0-rc.1 // indirect
	github.com/DataDog/datadog-go/v5 v5.1.1 // indirect
	github.com/DataDog/go-tuf v0.3.0--fix-localmeta-fork // indirect
	github.com/DataDog/sketches-go v1.4.1 // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.4.0 // indirect
	github.com/smartystreets/assertions v1.13.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/vmihailenco/go-tinylfu v0.2.2 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/dig v1.16.0 // indirect
	go4.org/intern v0.0.0-20211027215823-ae77deb06f29 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20220617031537-928513b29760 // indirect
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/sync v0.0.0-20220819030929-7fc1605a5dde // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	inet.af/netaddr v0.0.0-20220617031823-097006376321 // indirect
)
