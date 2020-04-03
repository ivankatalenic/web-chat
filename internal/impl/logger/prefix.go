package logger

import (
	"github.com/ivankatalenic/web-chat/internal/interfaces"
)

type prefixLogger struct {
	baseLogger interfaces.Logger
	prefix     string
}

// NewPrefix returns a new logger which appends a log message to a common prefix.
// Example: INFO: BROADCASTER: <message> or ERROR: BROADCASTER: <message>
func NewPrefix(baseLogger interfaces.Logger, prefix string) interfaces.Logger {
	return prefixLogger{
		baseLogger: baseLogger,
		prefix:     prefix,
	}
}

// Informational message
func (p prefixLogger) Info(msg string) {
	p.baseLogger.Info(p.prefix + ": " + msg)
}

// Warning message
func (p prefixLogger) Warning(msg string) {
	p.baseLogger.Warning(p.prefix + ": " + msg)
}

// Error message
func (p prefixLogger) Error(msg string) {
	p.baseLogger.Error(p.prefix + ": " + msg)
}
