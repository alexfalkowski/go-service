package rsa

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the RSA subsystem into [github.com/alexfalkowski/go-signal/v2].
//
// It provides constructors for:
//   - *[Generator] (via NewGenerator), which generates RSA key pairs,
//   - *[Encryptor] (via NewEncryptor), which encrypts using RSA-OAEP when RSA config is enabled, and
//   - *[Decryptor] (via NewDecryptor), which decrypts using RSA-OAEP when RSA config is enabled.
//
// Disabled behavior: if RSA configuration is disabled (nil *[Config]), NewEncryptor and NewDecryptor
// return (nil, nil) so downstream consumers can treat RSA encryption/decryption as optional.
var Module = signal.Module(
	signal.Constructor(NewGenerator),
	signal.Constructor(NewEncryptor),
	signal.Constructor(NewDecryptor),
)
