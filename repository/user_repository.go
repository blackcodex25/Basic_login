package repository

import (
	"Basic_login/domain"
	"errors"
	"log"
	"sync"
	"time"
)

// Error messages
var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

// InMemoryUserRepository implements UserRepository
type InMemoryUserRepository struct {
	mu            sync.RWMutex
	userIDCounter int64
	users         map[string]*domain.User
	userIDs       map[int64]*domain.User
	chatRoom      *domain.ChatRoom
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository with a specified chat room buffer size.
func NewInMemoryUserRepository(bufferSize int) *InMemoryUserRepository {
	repo := &InMemoryUserRepository{
		users:    make(map[string]*domain.User),
		userIDs:  make(map[int64]*domain.User),
		chatRoom: domain.NewChatRoom(bufferSize),
	}
	go repo.chatRoom.Run()
	return repo
}

// GetByID retrieves a user by their ID.
func (repo *InMemoryUserRepository) GetByID(id int64) (*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.userIDs[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Create adds a new user to the repository.
func (repo *InMemoryUserRepository) Create(user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[user.Username]; exists {
		return ErrUserExists
	}

	repo.userIDCounter++
	user.ID = repo.userIDCounter
	repo.users[user.Username] = user
	repo.userIDs[user.ID] = user
	log.Printf("Created user: %s with ID: %d", user.Username, user.ID)

	repo.chatRoom.Join <- user.Username
	return nil
}

// GetByUsername retrieves a user by their username.
func (repo *InMemoryUserRepository) GetByUsername(username string) (*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.users[username]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetAll retrieves all users from the repository.
func (repo *InMemoryUserRepository) GetAll() ([]*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	users := make([]*domain.User, 0, len(repo.users))
	for _, user := range repo.users {
		users = append(users, user)
	}
	return users, nil
}

// Update modifies an existing user.
func (repo *InMemoryUserRepository) Update(user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[user.Username]; !exists {
		return ErrUserNotFound
	}

	repo.users[user.Username] = user
	repo.userIDs[user.ID] = user
	log.Printf("Updated user: %s", user.Username)
	return nil
}

// SendChatMessage sends a chat message to the chat room.
func (repo *InMemoryUserRepository) SendChatMessage(sender, message string) {
	chatMessage := domain.ChatMessage{
		Sender:    sender,
		Message:   message,
		TimeStamp: time.Now(),
	}
	repo.chatRoom.Messages <- chatMessage
}

// LeaveChat removes a user from the chat room.
func (repo *InMemoryUserRepository) LeaveChat(username string) {
	repo.chatRoom.Leave <- username
}
