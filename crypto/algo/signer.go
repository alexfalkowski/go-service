package algo

// Signer for algo.
type Signer interface {
	// Sign a message.
	Sign(msg string) string

	// Verify sig with msg.
	Verify(sig, msg string) error
}

// NoSigner for algo.
type NoSigner struct{}

func (*NoSigner) Sign(msg string) string {
	return msg
}

func (*NoSigner) Verify(_, _ string) error {
	return nil
}
