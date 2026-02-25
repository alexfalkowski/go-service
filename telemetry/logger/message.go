package logger

import "log/slog"

// NewMessage constructs a Message with text and an optional error.
//
// When err is non-nil, the message will typically be logged at error level via
// `Message.Level`.
func NewMessage(text string, err error) Message {
	return Message{Text: text, Error: err}
}

// NewText constructs an informational Message with text and no error.
//
// This is equivalent to `NewMessage(text, nil)`.
func NewText(text string) Message {
	return Message{Text: text}
}

// Message represents a structured log message.
//
// It is consumed by `Logger.Log`/`Logger.LogAttrs` and is designed to keep the log
// record text (`Text`) and optional error (`Error`) together, so consistent
// formatting and level selection can be applied.
type Message struct {
	// Error is an optional error associated with the message.
	// When non-nil, it is typically included as an "error" attribute.
	Error error

	// Text is the human-readable log message.
	Text string
}

// Level returns the derived slog level for the message.
//
// If Error is non-nil, it returns `slog.LevelError`; otherwise it returns
// `slog.LevelInfo`.
func (m Message) Level() slog.Level {
	if m.Error != nil {
		return slog.LevelError
	}
	return slog.LevelInfo
}
