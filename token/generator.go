package token

// Generator generates authentication tokens for a given audience and subject.
//
// This interface represents the “issuance” side of a token system. Concrete implementations
// may generate different token formats depending on configuration (for example JWT, PASETO,
// or other schemes).
//
// # Parameters
//
// Generate accepts two identity inputs:
//
//   - aud: the intended audience for the token (who the token is meant for).
//   - sub: the subject identifier (who/what the token represents).
//
// In claim-based token formats (for example JWT/PASETO), aud and sub are typically encoded
// into standard claims.
//
// Some token kinds may ignore one or both parameters (for example formats that do not
// carry claims or that encode identity differently). Callers should consult the
// concrete implementation’s documentation for exact semantics.
//
// # Return value
//
// Generate returns the serialized token bytes. For text-based token formats, these bytes
// are typically UTF-8 encoded.
//
// # Errors
//
// Generate returns an error when token issuance fails (for example invalid configuration,
// missing key material, signing failures, or serialization errors).
type Generator interface {
	// Generate creates a new token for the given audience and subject.
	//
	// Implementations should treat aud/sub as logical identity inputs and produce a token
	// suitable for later validation by a corresponding Verifier.
	Generate(aud, sub string) ([]byte, error)
}
