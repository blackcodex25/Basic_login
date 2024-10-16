package controllers

import (
	"Basic_login/infrastructure"
	"Basic_login/usecase"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// StartChat manages the chat session
func StartChat(usecase *usecase.UserUsecase) {
	reader := bufio.NewReader(os.Stdin)
	username, err := infrastructure.ReadInput(reader, infrastructure.Prompts["username"])
	if err != nil {
		log.Println("Error:", err)
		return
	}

	fmt.Printf("%s has joined the chat.\n", username)

	for {
		message, err := infrastructure.ReadInput(reader, infrastructure.PromptMessage)
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		usecase.SendChatMessage(username, message)
		fmt.Printf("[%s] %s: %s\n", time.Now().Format(time.RFC3339), username, message)

		if Confirm(infrastructure.PromptLeaveChat) {
			usecase.LeaveChat(username)
			fmt.Printf("%s has left the chat.\n", username)
			break
		}
	}
}

// Confirm prompts for yes/no input
func Confirm(prompt string) bool {
	response, err := infrastructure.ReadInput(bufio.NewReader(os.Stdin), prompt)
	if err != nil {
		log.Println("Error reading input:", err)
		return false
	}
	return strings.EqualFold(response, "y")
}
