package domain

import (
	"log"
	"sync"
	"time"
)

// ChatRoom represents a chat room where users can join, leave, and send messages.
type ChatRoom struct {
	Messages chan ChatMessage
	Join     chan string
	Leave    chan string
	Users    map[string]struct{} // Use struct{} to save memory since we only need existence.
	mu       sync.Mutex          // Mutex for thread-safe access to Users.
}

// NewChatRoom initializes a new ChatRoom with buffered channels.
func NewChatRoom(bufferSize int) *ChatRoom {
	return &ChatRoom{
		Messages: make(chan ChatMessage, bufferSize),
		Join:     make(chan string, bufferSize),
		Leave:    make(chan string, bufferSize),
		Users:    make(map[string]struct{}),
	}
}

// Run starts the chat room and listens for messages, joins, and leaves.
func (c *ChatRoom) Run() {
	for {
		select {
		case message := <-c.Messages:
			go c.processMessage(message) // Process messages concurrently.
		case user := <-c.Join:
			go c.processJoin(user) // Process user joining concurrently.
		case user := <-c.Leave:
			go c.processLeave(user) // Process user leaving concurrently.
		}
	}
}

// processMessage handles and logs a received message.
func (c *ChatRoom) processMessage(message ChatMessage) {
	log.Printf("[%s] %s: %s\n", message.TimeStamp.Format(time.RFC3339), message.Sender, message.Message)
}

// processJoin manages a user joining the chat.
func (c *ChatRoom) processJoin(user string) {
	c.mu.Lock() // Lock to ensure thread-safe access to Users.
	defer c.mu.Unlock()

	c.Users[user] = struct{}{} // Only store existence
	log.Printf("%s joined the chat.\n", user)
}

// processLeave manages a user leaving the chat.
func (c *ChatRoom) processLeave(user string) {
	c.mu.Lock() // Lock to ensure thread-safe access to Users.
	defer c.mu.Unlock()

	delete(c.Users, user)
	log.Printf("%s left the chat.\n", user)
}
