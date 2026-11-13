package aes

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the AES subsystem into [github.com/alexfalkowski/go-signal/v2].
//
// It provides constructors for:
//   - *[Generator] (via NewGenerator), which can be used to generate AES-256 key material, and
//   - *[Cipher] (via NewCipher), which provides AES-GCM Encrypt/Decrypt when AES config is enabled.
//
// Disabled behavior: if AES configuration is disabled (nil *[Config]), NewCipher returns (nil, nil) so
// downstream consumers can treat AES as optional.
var Module = signal.Module(
	signal.Constructor(NewGenerator),
	signal.Constructor(NewCipher),
)
