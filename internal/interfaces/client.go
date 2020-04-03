package interfaces

import "github.com/ivankatalenic/web-chat/internal/models"

// Client is an interface for a chat client
type Client interface {
	GetAddress() string

	SendMessage(message *models.Message) error
	GetMessage() (*models.Message, error)

	Disconnect() error
	IsDisconnected() bool
}
