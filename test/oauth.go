package test

import (
	"github.com/alexfalkowski/go-service/security/oauth"
)

// NewOAuthConfig for test.
func NewOAuthConfig() *oauth.Config {
	return &oauth.Config{
		URL:           "http://localhost:5000/v1/oauth/token",
		ClientID:      "e1602e185cba2a90d8bbcfc3f3c5530c",
		ClientSecret:  "uC?MxwKO+r1@0RX[q8V5s4F|3oQ)yZ7TYDlUHmIfeNn9E&ScL2Pk{g$pi]z6bBta",
		Audience:      "standort",
		Issuer:        "https://auth.falkowski.io",
		Algorithm:     "EdDSA",
		JSONWebKeySet: "http://localhost:5000/v1/.well-known/jwks.json",
	}
}
