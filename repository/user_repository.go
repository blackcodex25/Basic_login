package repository

import (
	"Basic_login/domain"
	"errors"
	"log"
	"sync"
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
	userIDs       map[int64]*domain.User // New map for quick ID lookups
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:   make(map[string]*domain.User),
		userIDs: make(map[int64]*domain.User), // Initialize ID map
	}
}

// GetByID retrieves a user by ID
func (repo *InMemoryUserRepository) GetByID(id int64) (*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.userIDs[id] // Fast lookup by ID
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Create adds a new user to the repository
func (repo *InMemoryUserRepository) Create(user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[user.Username]; exists {
		return ErrUserExists
	}

	repo.userIDCounter++
	user.ID = repo.userIDCounter
	repo.users[user.Username] = user
	repo.userIDs[user.ID] = user // Maintain ID mapping
	log.Printf("Created user: %s with ID: %d", user.Username, user.ID)
	return nil
}

// GetByUsername retrieves a user by username
func (repo *InMemoryUserRepository) GetByUsername(username string) (*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.users[username]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetAll retrieves all users
func (repo *InMemoryUserRepository) GetAll() ([]*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	users := make([]*domain.User, 0, len(repo.users))
	for _, user := range repo.users {
		users = append(users, user)
	}
	return users, nil
}

// Update modifies an existing user
func (repo *InMemoryUserRepository) Update(user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[user.Username]; !exists {
		return ErrUserNotFound
	}

	repo.users[user.Username] = user
	repo.userIDs[user.ID] = user // Ensure the ID mapping is updated
	log.Printf("Updated user: %s", user.Username)
	return nil
}
