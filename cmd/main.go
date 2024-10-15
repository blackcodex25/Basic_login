package main

import (
	"Basic_login/domain"
	"Basic_login/repository"
	"Basic_login/usecase"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	promptUsername = "Enter username: "
	promptPassword = "Enter password: "
	promptRole     = "Enter role (admin/user): "
	promptContinue = "Do you want to create another user? (y/n): "
	roleAdmin      = "admin"
	roleUser       = "user"
)

func main() {
	userRepo := repository.NewInMemoryUserRepository()
	userUsecase := usecase.NewUserUsecase(
		userRepo,
		usecase.DefaultConfig(),
		log.Default(),
		usecase.NewConstants(),
	)

	for {
		if err := createUserFromInput(userUsecase); err != nil {
			fmt.Println("Error:", err)
		}

		if !askToContinue() {
			break
		}
	}
}

func createUserFromInput(usecase *usecase.UserUsecase) error {
	reader := bufio.NewReader(os.Stdin)

	// โหลดผู้ใช้ที่มีอยู่แล้วในแฮชเซต
	userSet := make(map[string]struct{})
	existingUsers, err := usecase.UserRepo.GetAll() // สมมติว่ามีฟังก์ชันนี้
	if err != nil {
		return err
	}
	for _, user := range existingUsers {
		userSet[user.Username] = struct{}{}
	}

	// อ่านข้อมูลผู้ใช้
	username, password, role, err := readUserInput(reader)
	if err != nil {
		return err
	}

	// ตรวจสอบการมีอยู่ของผู้ใช้
	if _, exists := userSet[username]; exists {
		return fmt.Errorf("error creating user: username '%s' already exists", username)
	}

	// สร้างผู้ใช้ใหม่
	if err := usecase.CreateUser(&domain.User{Username: username, Password: password}, role); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	fmt.Println("User created successfully.")
	return nil
}

func readUserInput(reader *bufio.Reader) (string, string, string, error) {
	username, err := readInput(reader, promptUsername)
	if err != nil {
		return "", "", "", err
	}

	password, err := readInput(reader, promptPassword)
	if err != nil {
		return "", "", "", err
	}

	role, err := readRole(reader)
	if err != nil {
		return "", "", "", err
	}

	return username, password, role, nil
}

func readRole(reader *bufio.Reader) (string, error) {
	role, err := readInput(reader, promptRole)
	if err != nil {
		return "", err
	}

	if role != roleAdmin && role != roleUser {
		return "", fmt.Errorf("invalid role: %s. Please enter 'admin' or 'user'", role)
	}
	return role, nil
}

func readInput(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func askToContinue() bool {
	cont, err := readInput(bufio.NewReader(os.Stdin), promptContinue)
	if err != nil {
		log.Println("Error reading input:", err)
		return false
	}
	return strings.EqualFold(cont, "y")
}
