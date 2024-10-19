package infrastructure

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	// คำถามสำหรับการป้อนข้อมูลของผู้ใช้
	promptUsername  = "Enter username: "                            // คำถามสำหรับการป้อนชื่อผู้ใช้
	promptPassword  = "Enter password: "                            // คำถามสำหรับการป้อนรหัสผ่าน
	promptRole      = "Enter role (admin/user): "                   // คำถามสำหรับการป้อนบทบาท
	promptContinue  = "Do you want to create another user? (y/n): " // คำถามสำหรับการสร้างผู้ใช้ใหม่
	PromptMessage   = "Enter message: "                             // คำถามสำหรับการป้อนข้อความ
	PromptLeaveChat = "Do you want to leave the chat? (y/n): "      // คำถามสำหรับการออกจากห้องสนทนา

	// บทบาท
	roleAdmin = "admin" // บทบาทผู้ดูแลระบบ
	roleUser  = "user"  // บทบาทผู้ใช้ทั่วไป
)

// คำถามสำหรับการป้อนข้อมูลของผู้ใช้
var Prompts = map[string]string{
	"username": promptUsername, // ชื่อผู้ใช้
	"password": promptPassword, // รหัสผ่าน
	"role":     promptRole,     // บทบาท
}

// ReadInput จัดการการป้อนข้อมูลของผู้ใช้
func ReadInput(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)                     // แสดงคำถามให้ผู้ใช้ป้อน
	input, err := reader.ReadString('\n') // อ่านข้อมูลที่ผู้ใช้ป้อนจนถึงบรรทัดใหม่
	if err != nil {
		return "", err // หากเกิดข้อผิดพลาดในการอ่าน คืนค่าข้อผิดพลาด
	}
	return strings.TrimSpace(input), nil // คืนค่าข้อมูลที่ถูกตัดช่องว่างที่ผู้ใช้ป้อน
}

// ReadUserInput รวบรวมข้อมูลผู้ใช้
func ReadUserInput(reader *bufio.Reader) (string, string, string, error) {
	username, err := ReadInput(reader, Prompts["username"]) // อ่านชื่อผู้ใช้
	if err != nil {
		return "", "", "", err // หากเกิดข้อผิดพลาดในการอ่านชื่อผู้ใช้คืนค่าข้อผิดพลาด
	}

	password, err := ReadInput(reader, Prompts["password"]) // อ่านรหัสผ่าน
	if err != nil {
		return "", "", "", err // หากเกิดข้อผิดพลาดในการอ่านรหัสผ่านคืนค่าข้อผิดพลาด
	}

	role, err := ReadRole(reader) // อ่านบทบาท
	if err != nil {
		return "", "", "", err // หากเกิดข้อผิดพลาดในการอ่านบทบาทคืนค่าข้อผิดพลาด
	}

	return username, password, role, nil // คืนค่าข้อมูลชื่อผู้ใช้, รหัสผ่าน, และบทบาท
}

// ReadRole ตรวจสอบบทบาทของผู้ใช้
func ReadRole(reader *bufio.Reader) (string, error) {
	role, err := ReadInput(reader, Prompts["role"]) // อ่านบทบาทจากผู้ใช้
	if err != nil {
		return "", err // หากเกิดข้อผิดพลาดในการอ่านคืนค่าข้อผิดพลาด
	}

	if !IsValidRole(role) { // ตรวจสอบว่าบทบาทที่ป้อนถูกต้องหรือไม่
		return "", fmt.Errorf("invalid role: %s. Please enter 'admin' or 'user'", role) // หากบทบาทไม่ถูกต้องให้คืนค่าข้อผิดพลาด
	}
	return role, nil // คืนค่าบทบาทที่ถูกต้อง
}

// IsValidRole ตรวจสอบว่าบทบาทที่ป้อนถูกต้องหรือไม่
func IsValidRole(role string) bool {
	return role == roleAdmin || role == roleUser // คืนค่าจริงหากบทบาทตรงกับ roleAdmin หรือ roleUser
}
