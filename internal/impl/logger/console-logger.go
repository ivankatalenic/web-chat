package logger

import (
	"fmt"
	"github.com/ivankatalenic/web-chat/internal/interfaces"
)

type console struct {
}

// NewConsoleLogger returns a new console prefixLogger
func NewConsoleLogger() interfaces.Logger {
	return console{}
}

// Informational message
func (console) Info(msg string) {
	fmt.Println("INFO: " + msg)
}

// Warning message
func (console) Warning(msg string) {
	fmt.Println("WARNING: " + msg)
}

// Error message
func (console) Error(msg string) {
	fmt.Println("ERROR: " + msg)
}
