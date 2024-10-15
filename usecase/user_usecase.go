package usecase

import (
	"Basic_login/domain"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"golang.org/x/crypto/argon2"
)

// Config holds Argon2 parameters
type Config struct {
	SaltLength   int
	ArgonTime    uint32
	ArgonMemory  uint32
	ArgonKeyLen  uint32
	ArgonThreads uint8
}

// DefaultConfig provides default Argon2 settings
func DefaultConfig() *Config {
	return &Config{
		ArgonTime:    1,
		ArgonMemory:  64 * 1024,
		ArgonThreads: 4,
		ArgonKeyLen:  32,
		SaltLength:   16,
	}
}

// User roles and error messages
type Constants struct {
	RoleAdmin                string
	RoleUser                 string
	ErrUserNotFound          error
	ErrInvalidPassword       error
	ErrUsernameTooLong       error
	ErrUsernameAlreadyExists error
}

// NewConstants initializes the constants
func NewConstants() *Constants {
	return &Constants{
		RoleAdmin:                "admin",
		RoleUser:                 "user",
		ErrUserNotFound:          errors.New("user not found"),
		ErrInvalidPassword:       errors.New("invalid password"),
		ErrUsernameTooLong:       errors.New("username must be between 5 and 20 characters long"),
		ErrUsernameAlreadyExists: errors.New("the username is already taken"),
	}
}

// UserRepository is an interface for user operations
type UserRepository interface {
	GetByID(id int64) (*domain.User, error)
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Update(user *domain.User) error
}

// UserUsecase contains methods for user operations
type UserUsecase struct {
	UserRepo  UserRepository
	Config    *Config
	Logger    *log.Logger
	Constants *Constants
}

// NewUserUsecase creates a new UserUsecase
func NewUserUsecase(repo UserRepository, config *Config, logger *log.Logger, constants *Constants) *UserUsecase {
	return &UserUsecase{UserRepo: repo, Config: config, Logger: logger, Constants: constants}
}

// GetUserByID retrieves user data by ID
func (u *UserUsecase) GetUserByID(id int64) (*domain.User, error) {
	return u.UserRepo.GetByID(id)
}

// UpdateUser updates user data
func (u *UserUsecase) Update(user *domain.User) error {
	return u.UserRepo.Update(user)
}

// CreateUser creates a new user
func (u *UserUsecase) CreateUser(user *domain.User, role string) error {
	if err := validateUsername(user.Username, u.Constants); err != nil {
		return err
	}

	user.Role = role
	hashedPassword, salt, err := HashPassword(user.Password, u.Config)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.Salt = salt

	if err := u.UserRepo.Create(user); err != nil {
		if errors.Is(err, u.Constants.ErrUsernameAlreadyExists) {
			return u.Constants.ErrUsernameAlreadyExists
		}
		return err
	}

	u.Logger.Printf("User created: %s\n", user.Username)
	return nil
}

// Login performs user login
func (u *UserUsecase) Login(username, password string) (*domain.User, error) {
	user, err := u.UserRepo.GetByUsername(username)
	if err != nil {
		return nil, u.Constants.ErrUserNotFound
	}

	if !ValidatePassword(password, user.Password, user.Salt, u.Config) {
		return nil, u.Constants.ErrInvalidPassword
	}

	u.Logger.Printf("Login successful: %s with role: %s\n", user.Username, user.Role)
	return user, nil
}

// HashPassword generates a hash for the provided plaintext password
func HashPassword(password string, config *Config) (string, []byte, error) {
	salt, err := generateSalt(config.SaltLength)
	if err != nil {
		return "", nil, err
	}
	hashedPassword := argon2.IDKey([]byte(password), salt, config.ArgonTime, config.ArgonMemory, config.ArgonThreads, config.ArgonKeyLen)
	return base64.RawStdEncoding.EncodeToString(hashedPassword), salt, nil
}

// ValidatePassword compares the provided plaintext password with the stored hashed password
func ValidatePassword(password, hashedPassword string, salt []byte, config *Config) bool {
	hash := argon2.IDKey([]byte(password), salt, config.ArgonTime, config.ArgonMemory, config.ArgonThreads, config.ArgonKeyLen)
	return hashedPassword == base64.RawStdEncoding.EncodeToString(hash)
}

// generateSalt generates a random salt
func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// validateUsername checks if the username meets the requirements
func validateUsername(username string, constants *Constants) error {
	if len(username) < 5 || len(username) > 20 {
		return constants.ErrUsernameTooLong
	}
	return nil
}
