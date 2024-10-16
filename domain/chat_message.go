package domain

import (
	"time"
)

// ChatMessage represents a message in the chat room.
type ChatMessage struct {
	TimeStamp time.Time
	Sender    string
	Message   string
}
