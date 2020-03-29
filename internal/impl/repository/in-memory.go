package repository

import (
	"github.com/ivankatalenic/web-chat/internal/models"
	"sync"
)

// InMemory message repository
type InMemory struct {
	lock sync.RWMutex

	messages models.MessageSlice
	size     int64
	msgCount int64
}

// NewInMemory creates a new in-memory message repository
// The size argument specifies the maximum number of messages which
// are being kept inside a repo: older messages are being overwritten by
// new ones.
func NewInMemory(size int64) *InMemory {
	return &InMemory{
		messages: make(models.MessageSlice, size),
		size:     size,
		msgCount: 0,
	}
}

// GetLast method
func (r *InMemory) GetLast(n int64) (models.MessageSlice, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if n <= 0 {
		return nil, nil
	}

	availableMsgs := min(r.msgCount, r.size)
	returnMessageCount := min(availableMsgs, n)
	if returnMessageCount <= 0 {
		return nil, nil
	}

	returnMessages := make(models.MessageSlice, returnMessageCount)

	startPos := (r.msgCount - returnMessageCount) % r.size
	endPos := (r.msgCount - 1) % r.size
	if startPos <= endPos {
		_ = copy(returnMessages, r.messages[startPos:endPos+1])
	} else {
		_ = copy(returnMessages, r.messages[startPos:])
		copied := r.size - startPos
		leftToCopy := returnMessageCount - copied
		_ = copy(returnMessages[copied:], r.messages[0:leftToCopy])
	}

	return returnMessages, nil
}

func min(a int64, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

// Put method
func (r *InMemory) Put(msg *models.Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	msg.ID = r.msgCount

	msgPos := r.msgCount % r.size
	r.messages[msgPos] = msg

	r.msgCount++

	return nil
}
