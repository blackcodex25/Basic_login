package main

import (
	"Basic_login/controllers"
	"Basic_login/repository"
	"Basic_login/usecase"
	"log"
)

func main() {
	// userRepo สร้าง instance ของ repository ข้อมูลผู้ใช้ในหน่วยความจำ (In-Memory)
	userRepo := repository.NewInMemoryUserRepository(1000) // สร้าง repository สำหรับเก็บข้อมูลผู้ใช้ในหน่วยความจำ
	// userUsecase สร้าง instance ของ use case สำหรับจัดการกับผู้ใช้ โดยใช้ repository และการตั้ง
	userUsecase := usecase.NewUserUsecase(
		userRepo,                // ใช้สำหรับเก็บข้อมูลผู้ใช้ในหน่วยความจำ
		usecase.DefaultConfig(), // ใช้สำหรับตั้งค่า Argon2 ในการเข้ารหัส
		log.Default(),           // ใช้สำหรับการบันทึกข้อมูล (logging) ในระบบ ในกรณีที่เกิดข้อผิดพลาด
		usecase.NewConstants(),  // ใช้สำหรับตั้งค่าค่าคงที่
	)

	// เรียกฟังก์ชัน CreateUser จาก controllers เพื่อสร้างผู้ใช้ใหม่หากมีข้อผิดพลาดจะถูกล็อก
	err := controllers.CreateUser(userUsecase) // สร้างผู้ใช้ใหม่
	if err != nil {
		log.Fatalf("Failed to create user: %v\n", err) // หากเกิดข้อผิดพลาดในการสร้างผู้ใช้ให้ล็อกข้อผิดพลาด
	}

	// เริ่มการสนทนาผ่านฟังก์ชัน StartChat ซึ่งใช้ userUsecase สำหรับจัดการกับผู้ใช้
	controllers.StartChat(userUsecase) // เริ่มการสนทนาสำหรับผู้ใช้
}
