// Package message defines crypto message payloads with authenticated metadata.
package message

// Message contains data and metadata for encryption and decryption.
//
// Data is the plaintext passed to encryption or the ciphertext passed to decryption.
// Meta is authenticated context that is not encrypted, such as a tenant, purpose, or
// record type. Algorithms that support authenticated metadata bind Meta to Data so
// decryption fails when the wrong Meta is supplied.
type Message struct {
	Data []byte
	Meta []byte
}
