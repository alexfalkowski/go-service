package token

// Verifier verifies authentication tokens for an expected audience and returns the subject identifier.
//
// This interface represents the “verification” side of a token system. Concrete implementations may
// verify different token formats (for example JWT, PASETO, or other schemes) and may impose additional
// checks such as issuer matching, algorithm constraints, key ID matching, or time validity.
//
// # Parameters
//
// Verify accepts:
//   - token: the serialized token bytes to verify.
//   - aud: the expected audience value for which the token must be valid.
//
// In claim-based token formats (for example JWT/PASETO), aud is typically validated against an "aud"
// claim. Some token kinds may ignore aud entirely (for example formats without claims). Callers should
// consult the concrete implementation’s documentation for exact semantics.
//
// # Return value
//
// On success, Verify returns the subject identifier represented by the token (commonly the "sub" claim).
// If the token kind does not carry a subject claim, the implementation may return an alternate identifier.
//
// # Errors
//
// Verify returns an error when validation fails (for example malformed token, signature mismatch, wrong
// issuer/audience, expired/not-yet-valid token, key mismatch, or missing key material). When an error is
// returned, the subject return value should not be trusted.
type Verifier interface {
	// Verify validates token for the given audience and returns the subject identifier.
	Verify(token []byte, aud string) (string, error)
}
