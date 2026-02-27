package rsa

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the RSA subsystem into Fx/Dig.
//
// It provides constructors for:
//   - *Generator (via NewGenerator), which generates RSA key pairs,
//   - *Encryptor (via NewEncryptor), which encrypts using RSA-OAEP when RSA config is enabled, and
//   - *Decryptor (via NewDecryptor), which decrypts using RSA-OAEP when RSA config is enabled.
//
// Disabled behavior: if RSA configuration is disabled (nil *Config), NewEncryptor and NewDecryptor
// return (nil, nil) so downstream consumers can treat RSA encryption/decryption as optional.
var Module = di.Module(
	di.Constructor(NewGenerator),
	di.Constructor(NewEncryptor),
	di.Constructor(NewDecryptor),
)
