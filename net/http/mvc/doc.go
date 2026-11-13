// Package mvc provides a small MVC-style HTML rendering layer for go-service HTTP servers.
//
// This package provides helpers to register controller-driven routes that render HTML templates.
// [Server] (constructed via [NewServer], typically via dependency injection) holds the HTTP router,
// template function map, filesystem, and layout used by those helpers.
//
// Definition requirements:
//
// Route/Get/Post/etc. return false when MVC is not defined on the Server (for example when no
// filesystem or layout is configured); see [Server.IsDefined].
//
// Context requirements:
//
// Controllers and views use net/http/meta to retrieve the request and response writer from context.
// The routing helpers set these values before invoking the controller.
//
// View lifecycle:
//
// Construct views during startup or route registration so template read and parse failures fail fast before
// the server starts handling requests.
//
// Localization:
//
// Provide an application-owned *[github.com/alexfalkowski/go-service/v3/i18n.Catalog]
// alongside the filesystem and layout to enable [Template.Localizer]. Each render
// selects a language from [WithLanguage], then Accept-Language, then the catalog
// default. Set WithLanguage in middleware before the MVC handler runs; a
// controller cannot replace its caller's context by deriving a new one.
// Full, partial, error, and not-found views share this behavior. Localized
// responses add Vary: Accept-Language and Content-Language for the selected
// catalog language. User-specific language preferences also require an
// application-appropriate cache policy, such as Cache-Control: private.
//
// Localizer methods return plain text, which html/template escapes. Missing
// translations fall back to the default language; missing messages or variable
// errors fail rendering before the response body is written. Without a catalog,
// Localizer is nil and localization does not change response headers.
package mvc
