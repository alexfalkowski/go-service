package auth0

import "github.com/form3tech-oss/jwt-go"

// Certificator for Auth0.
type Certificator interface {
	Certificate(token *jwt.Token) (string, error)
}
