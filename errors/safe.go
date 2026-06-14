package errors

// SafeMessenger allows an error to expose a message that is safe to send outside the process.
//
// Error implementations can use this when Error contains diagnostic details that should be retained for
// logs or wrapping, but clients should receive a smaller, non-sensitive message.
type SafeMessenger interface {
	// SafeMessage returns a non-sensitive message suitable for clients.
	//
	// Returning an empty string means this error has no safe message of its own; [SafeMessage] will continue
	// walking wrapped errors or use the caller-provided fallback.
	SafeMessage() string
}

// SafeMessage returns the first safe message in err's tree, or fallback when none is available.
//
// SafeMessage checks err itself, then follows single-error wrappers and joined-error children in order.
// Empty safe messages are ignored. The fallback is returned as supplied, so callers that need a non-empty
// client message should pass a non-empty fallback.
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
		for _, wrapped := range unwrapper.Unwrap() {
			if msg := safeMessage(wrapped); msg != "" {
				return msg
			}
		}
	}

	return ""
}
