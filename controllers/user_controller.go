package controllers

import (
	"Basic_login/domain"
	"Basic_login/infrastructure"
	"Basic_login/usecase"
	"bufio"
	"fmt"
	"os"
)

// CreateUser สร้างผู้ใช้ใหม่
func CreateUser(usecase *usecase.UserUsecase) error {
	reader := bufio.NewReader(os.Stdin) // สร้าง reader สำหรับอ่านข้อมูลจาก stdin

	userSet, err := GetExistingUsers(usecase) // ดึงผู้ใช้ที่มีอยู่แล้ว
	if err != nil {
		return err
	}

	username, password, role, err := infrastructure.ReadUserInput(reader) //  อ่านข้อมูลผู้ใช้
	if err != nil {
		return err
	}

	if _, exists := userSet[username]; exists { // ตรวจสอบว่าชื่อผู้ใช้มีอยูู่แล้วหรือไม่
		return fmt.Errorf("username '%s' already exists", username) // ถ้ามีอยู่แล้วให้คืนค่าข้อผิดพลาด
	}

	if err := usecase.CreateUser(&domain.User{Username: username, Password: password}, role); err != nil { // สร้างผู้ใช้ใหม่
		return fmt.Errorf("error creating user: %w", err) // คืนค่าข้อผิดพลาด ถ้าสร้างไม่สำเร็จ
	}

	fmt.Println("User created successfully.") // แสดงข้อความเมื่อสร้างผู้ใช้สำเร็จ
	return nil
}

// GetExistingUsers ดึงผู้ใช้ที่มีอยู่ทั้งหมด
func GetExistingUsers(usecase *usecase.UserUsecase) (map[string]struct{}, error) {
	existingUsers, err := usecase.UserRepo.GetAll() // ดึงผู้ใช้ทั้งหมดจาก UserRepo
	if err != nil {
		return nil, err // คืนค่าข้อผิดพลาด ถ้าดึงไม่สำเร็จ
	}

	userSet := make(map[string]struct{}, len(existingUsers)) // สร้าง map สำหรับเก็บชื่อผู้ใช้
	for _, user := range existingUsers {
		userSet[user.Username] = struct{}{} // เก็บเฉพาะการมีอยู่ของชื่อผู้ใช้
	}

	return userSet, nil // คืนค่า map ผู้ใช้ที่มีอยู่แล้ว
}
