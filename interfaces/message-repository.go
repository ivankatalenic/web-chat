package interfaces

import "github.com/ivankatalenic/web-chat/models"

// MessageRepository is a message repository
type MessageRepository interface {
	GetLast(n int) (models.MessageSlice, error)
	Put(msg *models.Message) error
}
