package infrastructure

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	// Prompts for user input
	promptUsername  = "Enter username: "
	promptPassword  = "Enter password: "
	promptRole      = "Enter role (admin/user): "
	promptContinue  = "Do you want to create another user? (y/n): "
	PromptMessage   = "Enter message: "
	PromptLeaveChat = "Do you want to leave the chat? (y/n): "

	// Roles
	roleAdmin = "admin"
	roleUser  = "user"
)

// Prompts for user input
var Prompts = map[string]string{
	"username": promptUsername,
	"password": promptPassword,
	"role":     promptRole,
}

// ReadInput handles user input
func ReadInput(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// ReadUserInput gathers user information
func ReadUserInput(reader *bufio.Reader) (string, string, string, error) {
	username, err := ReadInput(reader, Prompts["username"])
	if err != nil {
		return "", "", "", err
	}

	password, err := ReadInput(reader, Prompts["password"])
	if err != nil {
		return "", "", "", err
	}

	role, err := ReadRole(reader)
	if err != nil {
		return "", "", "", err
	}

	return username, password, role, nil
}

// ReadRole validates user role input
func ReadRole(reader *bufio.Reader) (string, error) {
	role, err := ReadInput(reader, Prompts["role"])
	if err != nil {
		return "", err
	}

	if !IsValidRole(role) {
		return "", fmt.Errorf("invalid role: %s. Please enter 'admin' or 'user'", role)
	}
	return role, nil
}

// isValidRole checks if the provided role is valid
func IsValidRole(role string) bool {
	return role == roleAdmin || role == roleUser
}
