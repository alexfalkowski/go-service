// Package client provides a content-aware HTTP client wrapper used by go-service.
//
// This package wraps `net/http.Client` to provide a higher-level API for making HTTP requests with:
//   - request/response body encoding/decoding via `net/http/content.Content`, and
//   - consistent error handling via `net/http/status`.
//
// The wrapper is intended for service-to-service HTTP calls where the payload format is negotiated
// explicitly via Content-Type.
//
// # Options and behavior
//
// `NewClient` accepts optional ClientOption values to configure:
//   - the underlying RoundTripper (e.g. to add retries, breakers, auth, etc.)
//   - the overall client timeout (http.Client.Timeout)
//   - redirect behavior (optionally return redirects instead of following them)
//
// # Encoding/decoding
//
// Each call uses an `Options` value describing:
//   - Request: request payload model (optional)
//   - Response: response payload model (optional)
//   - ContentType: media type used to select an encoder/decoder
//
// If Request is set, it is encoded into the request body.
// If Response is set, it is decoded from the response body on non-error responses.
//
// # Error handling
//
// Response handling follows these rules:
//   - If the response Content-Type indicates an error payload (text/error), the body is treated as a
//     message and returned as a `net/http/status` error.
//   - If the status code is in the 4xx/5xx range and the response is not an error media type,
//     a generic status error is returned.
//   - Otherwise, the response body is decoded into Response (if provided).
//
// Start with `NewClient` and `Options`.
package client
