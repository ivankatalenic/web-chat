package logger

import (
	"fmt"
	"github.com/ivankatalenic/web-chat/interfaces"
)

type logger struct {
}

// NewConsoleLogger returns a new console logger
func NewConsoleLogger() interfaces.Logger {
	return logger{}
}

// Informational message
func (logger) Info(msg string) {
	fmt.Println("INFO: " + msg)
}

// Warning message
func (logger) Warning(msg string) {
	fmt.Println("WARNING: " + msg)
}

// Error message
func (logger) Error(msg string) {
	fmt.Println("ERROR: " + msg)
}
