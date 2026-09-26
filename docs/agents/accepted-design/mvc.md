# Accepted Design: MVC

These entries are part of the [accepted design index](../accepted-design.md)
and carry the same mandatory weight.

- MVC controller errors render a client-safe `mvc.Error` model; `mvcModelError` metadata intentionally remains the raw error string for compatibility and must not be rendered unless diagnostic detail exposure is acceptable.
- MVC views should be constructed during startup or route registration so
  missing, unreadable, or malformed templates fail fast before serving traffic.
  Do not flag `mvc.NewFullView`, `mvc.NewPartialView`, or `mvc.NewViewPair`
  panics as request-path reliability gaps solely because those constructors
  panic; report only concrete supported paths that construct views per request
  contrary to the documented lifecycle.
- MVC static files are expected to be served from embedded or otherwise stable
  application assets. Do not flag ignored mid-stream `io.Copy` errors in
  `mvc.writeStaticFile` as reliability gaps solely because an artificial
  `fs.FS` can return partial data after successful `Open`/`Stat`; report only
  concrete supported filesystem paths where static reads can fail mid-stream
  and operators need different behavior.
- `mvc.StaticPathValue` intentionally treats the decoded
  `Request.PathValue` as an `fs.FS` path beneath the configured prefix. Go's
  `ServeMux` can decode `%2F` inside a single wildcard to `/`, so a pattern such
  as `/{file}` can resolve descendants under that prefix. Do not flag this
  solely as path traversal; report only a concrete prefix escape, a supported
  top-level-only contract, or unintended exposure of a non-public asset.
- MVC static serving intentionally does not generate `ETag`, `Last-Modified`,
  or conditional 304 responses. The supported embedded filesystem path has no
  reliable content identity, so do not propose reintroducing validators unless
  the supported Go toolchain provides a trustworthy constant-time content hash
  or MVC gains explicit, reliable build metadata.
- MVC not-found handling intentionally uses a simple `Accept` header check for
  `text/html` and does not fully evaluate quality weights such as
  `text/html;q=0`. Do not flag this unless there is concrete evidence of a
  real client, route, or deployment contract that depends on strict weighted
  `Accept` negotiation for MVC 404 responses.
