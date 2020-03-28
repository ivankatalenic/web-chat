package models

import "time"

// MessageSlice is a slice of Message pointers
type MessageSlice []*Message

// Message represents a chat message
type Message struct {
	ID        int64
	Author    string
	Content   string
	Timestamp time.Time
}
