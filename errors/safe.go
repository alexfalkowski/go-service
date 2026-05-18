package errors

// SafeMessenger allows an error to expose a message that is safe to send outside the process.
//
// Error implementations can use this when Error contains diagnostic details that should be retained for
// logs or wrapping, but clients should receive a smaller, non-sensitive message.
type SafeMessenger interface {
	// SafeMessage returns a non-sensitive message suitable for clients.
	SafeMessage() string
}

// SafeMessage returns the first safe message in err's chain, or fallback when none is available.
//
// Empty safe messages are ignored so callers can always rely on a non-empty fallback.
func SafeMessage(err error, fallback string) string {
	if msg := safeMessage(err); msg != "" {
		return msg
	}

	return fallback
}

func safeMessage(err error) string {
	if err == nil {
		return ""
	}

	if messenger, ok := err.(SafeMessenger); ok {
		if msg := messenger.SafeMessage(); msg != "" {
			return msg
		}
	}

	if unwrapper, ok := err.(interface{ Unwrap() error }); ok {
		return safeMessage(unwrapper.Unwrap())
	}

	if unwrapper, ok := err.(interface{ Unwrap() []error }); ok {
		for _, err := range unwrapper.Unwrap() {
			if msg := safeMessage(err); msg != "" {
				return msg
			}
		}
	}

	return ""
}
