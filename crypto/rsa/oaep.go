package rsa

import (
	"crypto/rsa"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/message"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
)

// EncryptOAEP encrypts msg.Data using RSA-OAEP with SHA-512.
//
// msg.Meta is used as the OAEP label and must be supplied unchanged to DecryptOAEP.
// msg.Data must fit RSA-OAEP's plaintext limit for publicKey: modulus bytes
// minus two SHA-512 digest lengths minus two bytes.
func EncryptOAEP(generator *rand.Generator, publicKey *rsa.PublicKey, msg message.Message) ([]byte, error) {
	return rsa.EncryptOAEP(sha512.New(), generator.Reader(), publicKey, msg.Data, msg.Meta)
}

// DecryptOAEP decrypts msg.Data using RSA-OAEP with SHA-512.
//
// msg.Meta is used as the OAEP label and must match the metadata supplied to EncryptOAEP.
// The generator parameter is retained for API consistency with EncryptOAEP. The standard library's
// RSA-OAEP decryption treats the randomness parameter as legacy and ignores it.
func DecryptOAEP(generator *rand.Generator, privateKey *rsa.PrivateKey, msg message.Message) ([]byte, error) {
	return rsa.DecryptOAEP(sha512.New(), generator, privateKey, msg.Data, msg.Meta)
}
