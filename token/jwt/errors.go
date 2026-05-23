package jwt

import webtoken "github.com/golang-jwt/jwt/v4"

// ValidationError aliases the upstream JWT validation error type.
type ValidationError = webtoken.ValidationError

// ErrTokenMalformed aliases the upstream malformed JWT error.
var ErrTokenMalformed = webtoken.ErrTokenMalformed

// ErrTokenSignatureInvalid aliases the upstream invalid JWT signature error.
var ErrTokenSignatureInvalid = webtoken.ErrTokenSignatureInvalid
