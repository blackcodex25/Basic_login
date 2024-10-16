package main

import (
	"Basic_login/controllers"
	"Basic_login/repository"
	"Basic_login/usecase"
	"log"
)

// Main function
func main() {
	userRepo := repository.NewInMemoryUserRepository(1000)
	userUsecase := usecase.NewUserUsecase(
		userRepo,
		usecase.DefaultConfig(),
		log.Default(),
		usecase.NewConstants(),
	)

	// Create a new user (for demonstration)
	err := controllers.CreateUser(userUsecase)
	if err != nil {
		log.Fatalf("Failed to create user: %v\n", err)
	}

	// Start chat session
	controllers.StartChat(userUsecase)
}
