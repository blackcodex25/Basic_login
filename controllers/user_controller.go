package controllers

import (
	"Basic_login/domain"
	"Basic_login/infrastructure"
	"Basic_login/usecase"
	"bufio"
	"fmt"
	"os"
)

func CreateUser(usecase *usecase.UserUsecase) error {
	reader := bufio.NewReader(os.Stdin)
	userSet, err := GetExistingUsers(usecase)
	if err != nil {
		return err
	}

	username, password, role, err := infrastructure.ReadUserInput(reader)
	if err != nil {
		return err
	}

	if _, exists := userSet[username]; exists {
		return fmt.Errorf("username '%s' already exists", username)
	}

	if err := usecase.CreateUser(&domain.User{Username: username, Password: password}, role); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	fmt.Println("User created successfully.")
	return nil
}

// GetExistingUsers retrieves all existing users
func GetExistingUsers(usecase *usecase.UserUsecase) (map[string]struct{}, error) {
	existingUsers, err := usecase.UserRepo.GetAll()
	if err != nil {
		return nil, err
	}

	userSet := make(map[string]struct{}, len(existingUsers))
	for _, user := range existingUsers {
		userSet[user.Username] = struct{}{}
	}

	return userSet, nil
}
