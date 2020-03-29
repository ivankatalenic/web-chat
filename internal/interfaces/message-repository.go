package interfaces

import "github.com/ivankatalenic/web-chat/internal/models"

// MessageRepository is a message repository
type MessageRepository interface {
	GetLast(n int64) (models.MessageSlice, error)
	Put(msg *models.Message) error
}
