package aes

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the AES subsystem into Fx/Dig.
//
// It provides constructors for:
//   - *Generator (via NewGenerator), which can be used to generate AES-256 key material, and
//   - *Cipher (via NewCipher), which provides AES-GCM Encrypt/Decrypt when AES config is enabled.
//
// Disabled behavior: if AES configuration is disabled (nil *Config), NewCipher returns (nil, nil) so
// downstream consumers can treat AES as optional.
var Module = di.Module(
	di.Constructor(NewGenerator),
	di.Constructor(NewCipher),
)
