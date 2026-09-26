# Accepted Design: Encoding

These entries are part of the [accepted design index](../accepted-design.md)
and carry the same mandatory weight.

- `encoding/base64.EncodedLen` intentionally follows the standard library's
  `EncodedLen(int)` contract. Do not flag hypothetical integer overflow from
  passing absurd configured `bytes.Size` values directly to it; callers that
  compare configuration-sized limits, such as cache decode size guards, must
  guard those limits before calling it.
- A wire format is admissible for decoding untrusted input only when its decoder
  is both ratio-bounded, meaning it never allocates from a declared count
  without validating that count against bytes actually received, and
  depth-bounded. This is a decoder rule and applies per direction: encoding a
  response is not gated by it, because the encoder's input is a service-owned
  domain object rather than caller-controlled bytes. `max_receive_size`,
  `net/http/body.NewHandler`, and the streaming per-value `capReader` bound
  input bytes only, so they do not mitigate amplification and are not a
  substitute for these bounds. `json`, `yaml`, and protobuf are admissible:
  stdlib `encoding/json` has `maxNestingDepth` plus incremental append,
  `yaml.v3` has `max_flow_level`/`max_indents` plus `allowedAliasRatio`
  alias-bomb protection, and protobuf validates every declared length against
  the remaining buffer in `protowire.ConsumeBytes` and applies
  `proto.UnmarshalOptions.RecursionLimit`.
- `msgpack` and `gob` fail the decoder-bounds rule above, so they are
  intentionally rejected for request-body decoding by
  `net/http/content/unary.Media.CanDecodeRequest` with HTTP 415 and intentionally
  absent from `net/http/content/stream`'s `streamKinds`. Both remain fully supported
  for response encoding, which is why they still resolve as media types.
  `encoding/gob`'s own `Decoder` documents that it does only basic sanity
  checking on decoded input sizes and that its limits are not configurable, so
  there is no hook to add one. The vendored `Basekick-Labs/msgpack/v6` has
  default allocation limits (`bytesAllocLimit`, `sliceAllocLimit`,
  `maxMapSize`) but they are per-container constants rather than ratio bounds,
  and the package has no depth limit or option for one: `decodeSlice` recurses
  through `decodeInterfaceCond` while holding a clamped `[]interface{}`, so
  nested container headers decoded into an `any` or `map[string]any` target
  amplify per level and can exhaust the goroutine stack, which is a fatal
  runtime error rather than a recoverable panic. Do not propose exposing
  `msgpack` or `gob` for request-body decoding or streaming negotiation, do not
  generalize this into a "binary formats are internal" or "unsupported over
  HTTP" policy, and do not recommend `GOMEMLIMIT`, watchdog goroutines, or
  lower input-size caps as mitigations. Report only concrete bugs such as a
  rejected media type being decoded anyway, an admissible format being
  rejected, or new evidence that a listed decoder's bounds changed upstream.
- JSON decoding intentionally keeps the standard library's duplicate object key
  behavior, where later values replace earlier values. Do not flag this as a
  finding unless a public API starts promising duplicate-key rejection or this
  repository adds a strict JSON decoder mode.
- HJSON duplicate-key errors intentionally preserve upstream diagnostic detail,
  including duplicate decoded values, to keep bad configuration files
  debuggable. Configuration files are not expected to contain raw passwords or
  credentials; they should contain source references such as `env:NAME` or
  `file:/path`. Do not flag this as a secret leak unless a concrete code path
  starts placing raw secrets into HJSON configuration contrary to that policy.
