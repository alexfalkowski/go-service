package algo

// Cipher for algo.
type Cipher interface {
	// Encrypt a msg.
	Encrypt(msg string) (string, error)

	// Decrypt a msg.
	Decrypt(msg string) (string, error)
}

// NoCipher for algo.
type NoCipher struct{}

func (*NoCipher) Encrypt(msg string) (string, error) {
	return msg, nil
}

func (*NoCipher) Decrypt(msg string) (string, error) {
	return msg, nil
}
