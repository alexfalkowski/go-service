package jwt

import "github.com/golang-jwt/jwt/v4"

// ValidationError aliases the upstream JWT validation error type.
type ValidationError = jwt.ValidationError

// ErrTokenMalformed aliases the upstream malformed JWT error.
var ErrTokenMalformed = jwt.ErrTokenMalformed

// ErrTokenSignatureInvalid aliases the upstream invalid JWT signature error.
var ErrTokenSignatureInvalid = jwt.ErrTokenSignatureInvalid
