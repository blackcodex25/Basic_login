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

// StartChat จัดการเซสชันแชท
func StartChat(usecase *usecase.UserUsecase) {
	reader := bufio.NewReader(os.Stdin)                                                   // สร้าง reader สำหรับอ่านข้อมูลจาก stdin
	username, err := infrastructure.ReadInput(reader, infrastructure.Prompts["username"]) // อ่านชื่อผู้ใช้
	if err != nil {
		log.Println("Error:", err)
		return
	}

	fmt.Printf("%s has joined the chat.\n", username) // บันทึกข้อผิดพลาดถ้ามี

	for {
		message, err := infrastructure.ReadInput(reader, infrastructure.PromptMessage) // อ่านข้อความจากผู้ใช้
		if err != nil {
			log.Println("Error:", err) // บันทึกข้อผิดพลาดถ้ามี
			continue
		}

		usecase.SendChatMessage(username, message)                                               // ส่งข้อความไปยังเซิร์ฟเวอร์
		fmt.Printf("[%s] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), username, message) // แสดงข้อความในแชท

		if Confirm(infrastructure.PromptLeaveChat) { // ถามผู้ใช้ว่าต้องการออกจากห้องแชทหรือไม่
			usecase.LeaveChat(username)                     // ออกจากห้องแชท
			fmt.Printf("%s has left the chat.\n", username) // แสดงข้อความเมื่อผู้ใช้ออกจากห้องแชท
			break
		}
	}
}

// Confirm ถามผู้ใช้เพื่อรับข้อมูลแบบใช่/ไม่ใช่
func Confirm(prompt string) bool {
	response, err := infrastructure.ReadInput(bufio.NewReader(os.Stdin), prompt) // อ่านข้อความจากผู้ใช้
	if err != nil {
		log.Println("Error reading input:", err) // บันทึกข้อผิดพลาดถ้ามี
		return false                             // คืนค่า false ถ้าอ่านข้อความไม่สำเร็จ
	}
	return strings.EqualFold(response, "y") // คืนค่า true ถ้าผู้ใช้ตอบ "y"
}
