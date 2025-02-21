package logger

import "log/slog"

// NewMessage creates a new log message with text and error information.
func NewMessage(text string, err error) Message {
	return Message{Text: text, Error: err}
}

// NewText creates a new log message with text information.
func NewText(text string) Message {
	return Message{Text: text}
}

// Message represents a log message with text and error information.
type Message struct {
	Error error
	Text  string
}

// Level returns the log level for the message.
func (m Message) Level() slog.Level {
	if m.Error != nil {
		return slog.LevelError
	}

	return slog.LevelInfo
}
