package rsa

import (
	"crypto/rsa"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
)

// EncryptOAEP encrypts msg using RSA-OAEP with SHA-512 and a nil label.
func EncryptOAEP(generator *rand.Generator, publicKey *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha512.New(), generator.Reader(), publicKey, msg, nil)
}

// DecryptOAEP decrypts msg using RSA-OAEP with SHA-512 and a nil label.
func DecryptOAEP(generator *rand.Generator, privateKey *rsa.PrivateKey, msg []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha512.New(), generator, privateKey, msg, nil)
}
