package interfaces

// Logger is a logger interface
type Logger interface {
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}
