package repository

import (
	"github.com/ivankatalenic/web-chat/internal/models"
	"sync"
)

// InMemory message repository
type InMemory struct {
	lock     sync.Mutex
	messages models.MessageSlice
}

// NewInMemory creates a new in-memory message repository
func NewInMemory() *InMemory {
	return &InMemory{
		messages: make(models.MessageSlice, 0),
	}
}

// GetLast method
func (r *InMemory) GetLast(n int) (models.MessageSlice, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	messageLen := len(r.messages)
	if n <= 0 {
		return nil, nil
	} else if n > messageLen {
		return r.messages, nil
	}

	return r.messages[messageLen-n:], nil
}

// Put method
func (r *InMemory) Put(msg *models.Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	msg.ID = int64(len(r.messages))

	r.messages = append(r.messages, msg)

	return nil
}
