package algo

// Signer for algo.
type Signer interface {
	// Sign a message.
	Sign(msg string) (string, error)

	// Verify sig with msg.
	Verify(sig, msg string) error
}

// NoSigner for algo.
type NoSigner struct{}

func (*NoSigner) Sign(msg string) (string, error) {
	return msg, nil
}

func (*NoSigner) Verify(_, _ string) error {
	return nil
}
