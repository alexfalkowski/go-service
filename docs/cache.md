# 💾 Cache

[← Back to README](../README.md)

Cache configuration is defined in `cache/config.Config`:

```yaml
cache:
  kind: redis
  compressor: zstd
  encoder: json
  max_size: 4MB
  max_entries: 1024
  options:
    url: env:CACHE_URL
```

> [!NOTE]
>
> - Built-in driver kinds in this repo are `redis` and `ttlcache`.
> - Unknown `kind` values return `cache/driver/errors.ErrNotFound`.
> - Unknown or empty `compressor` values fall back to `none`.
> - For normal values, unknown or empty `encoder` values fall back to `json`.
> - Configured `compressor` and `encoder` values are part of the cache driver key namespace, so changing either setting creates cache misses for values written with the previous format.
> - Cache operations use `bytes` for `io.WriterTo`/`io.ReaderFrom` stream values and `protobuf` for protobuf messages, regardless of the configured `encoder`.
> - `max_size` limits encoded cache values before compression, after compression, and after decompression. A zero value uses the default `4MB`.
> - `max_entries` limits entries retained by bounded in-memory cache drivers. A zero value uses the default `1024`; negative values are invalid.
> - `options` is backend-specific and decoded as `map[string]any`.
> - Redis-backed `GetOrPersist` requires [Redis 7.0-compatible `SET` semantics](https://redis.io/docs/latest/commands/set/) because atomic publication uses `SET ... NX GET`.
> - Configure each cache backend for a specific service or purpose. For Redis, use a dedicated database, endpoint, or deployment-level key namespace in the connection/configuration instead of sharing one general cache for unrelated data.

> [!WARNING]
> `Cache.Flush` follows backend semantics; for Redis it clears the selected database.
