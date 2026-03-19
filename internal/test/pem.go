package test

import "github.com/alexfalkowski/go-service/v2/crypto/pem"

// PEM is the shared PEM decoder wired to the test filesystem.
var PEM = pem.NewDecoder(FS)
