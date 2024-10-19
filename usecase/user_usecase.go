package usecase

import (
	"Basic_login/domain"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"golang.org/x/crypto/argon2"
)

// โครงสร้างข้อมูล Config สำหรับการกำหนดค่าที่ใช้สำหรับการเข้ารหัส
type Config struct {
	SaltLength   int    // ความยาวของ salt
	ArgonTime    uint32 // จำนวนรอบการเข้ารหัส
	ArgonMemory  uint32 // จํานวนหน่วยความจำใน kilobytes
	ArgonKeyLen  uint32 // ความยาวของ key
	ArgonThreads uint8  // จำนวนเธรด (threads) ที่จะใช้ในการประมวลผลแฮช เพื่อเพิ่มประสิทธิภาพในการคำนวณ
}

// ฟังก์ชัน DefaultConfig สำหรับส่งคืนค่าโครงสร้าง Config ที่มีการตั้งค่าพื้นฐานสำหรับ argon2
func DefaultConfig() *Config {
	return &Config{
		ArgonTime:    1,         // ตั้งค่าเป็น 1 รอบ
		ArgonMemory:  64 * 1024, // ตั้งค่าเป็น 64 * 1024 ซึ่งแปลว่าใช้หน่วยความจำ 64 KB ในการประมวลผล
		ArgonThreads: 4,         //  ตั้งค่าเป็น 4 ซึ่งหมายถึงจำนวนเธรดที่ใช้ในการประมวลผลแฮช เพื่อให้สามารถคำนวณได้เร็วขึ้นในระบบที่รองรับหลายเธรด
		ArgonKeyLen:  32,        // ตั้งค่าเป็น 32 ซึ่งระบุความยาวของคีย์ที่ได้จากการแฮชเป็น 32 bytes
		SaltLength:   16,        // ตั้งค่าเป็น 16 ซึ่งหมายถึงความยาวของ "salt" ที่จะใช้ในการแฮช
	}
}

// โครงสร้าง Constants สำหรับกำหนดค่าคงที่ ที่เกี่ยวข้องกับบทบาทผู้ใช้ และ ข้อความแสดงข้อผิดพลาดในระบบ
type Constants struct {
	RoleAdmin                string // ค่าที่ใช้สำหรับบทบาทผู้ดูแลระบบ
	RoleUser                 string // ค่าที่ใช้สำหรับบทบาทผู้ใช้
	ErrUserNotFound          error  // ข้อความข้อผิดพลาดที่แสดงเมื่อไม่พบผู้ใช้ในระบบ
	ErrInvalidPassword       error  // ข้อความข้อผิดพลาดที่แสดงเมื่อรหัสผ่านไม่ถูกต้อง
	ErrUsernameTooLong       error  // ข้อความข้อผิดพลาดที่แสดงเมื่อชื่อผู้ใช้ที่ป้อนยาวเกินกว่าที่กำหนด
	ErrUsernameAlreadyExists error  // ข้อความข้อผิดพลาดที่แสดงเมื่อมีการพยายามสร้างชื่อผู้ใช้ที่มีอยู่แล้วในระบบ
}

// ฟังก์ชัน NewConstants เพื่อสร้างและคืนค่าให้โครงสร้าง Constants ที่มีการกำหนดค่าคงที่สำหรับบทบาทของผู้ใช้และข้อความแสดงข้อผิดพลาด
func NewConstants() *Constants {
	return &Constants{
		RoleAdmin:                "admin",                                                         // ตั้งค่าเป็น "admin" เพื่อระบุบทบาทของผู้ดูแลระบบ
		RoleUser:                 "user",                                                          // ตั้งค่าเป็น "user" เพื่อระบุบทบาทของผู้ใช้ทั่วไป
		ErrUserNotFound:          errors.New("user not found"),                                    // ใช้ errors.New("user not found") เพื่อสร้างข้อผิดพลาดที่จะแจ้งว่าไม่พบผู้ใช้ในระบบ
		ErrInvalidPassword:       errors.New("invalid password"),                                  // ใช้ errors.New("invalid password") เพื่อสร้างข้อผิดพลาดที่จะแจ้งว่า รหัสผ่านที่ป้อนไม่ถูกต้อง
		ErrUsernameTooLong:       errors.New("username must be between 5 and 20 characters long"), // ใช้ errors.New("username must be between 5 and 20 characters long") เพื่อสร้างข้อผิดพลาดที่จะแจ้งว่าชื่อผู้ใช้ต้องมีความยาวระหว่าง 5 ถึง 20 ตัวอักษร
		ErrUsernameAlreadyExists: errors.New("the username is already taken"),                     // ใช้ errors.New("the username is already taken") เพื่อสร้างข้อผิดพลาดที่จะแจ้งว่าชื่อผู้ใช้ที่พยายามลงทะเบียนมีอยู่แล้วในระบบ
	}
}

// โครงสร้าง interface UserRepository ใช้สำหรับการดำเนินการกับผู้ใช้ในระบบ
type UserRepository interface {
	GetByID(id int64) (*domain.User, error)              // ฟังก์ชันนี้ใช้เพื่อดึงข้อมูลผู้ใช้จากฐานข้อมูลตาม id ที่ระบุ โดยจะคืนค่าผู้ใช้ (*domain.User) และข้อผิดพลาด (error) หากไม่พบผู้ใช้หรือเกิดข้อผิดพลาดในการดึงข้อมูล
	Create(user *domain.User) error                      // ฟังก์ชันนี้ใช้เพื่อสร้างผู้ใช้ใหม่ในฐานข้อมูล โดยรับพารามิเตอร์เป็นผู้ใช้ (*domain.User) และจะคืนค่าข้อผิดพลาดหากเกิดปัญหาในการสร้าง
	GetByUsername(username string) (*domain.User, error) // ฟังก์ชันนี้ใช้เพื่อดึงข้อมูลผู้ใช้จากฐานข้อมูลตามชื่อผู้ใช้ (username) โดยคืนค่าผู้ใช้และข้อผิดพลาดตามปกติ
	GetAll() ([]*domain.User, error)                     // ฟังก์ชันนี้ใช้เพื่อดึงข้อมูลผู้ใช้ทั้งหมดจากฐานข้อมูล โดยคืนค่าลิสต์ของผู้ใช้ ([]*domain.User) และข้อผิดพลาด
	Update(user *domain.User) error                      // ฟังก์ชันนี้ใช้เพื่อปรับปรุงข้อมูลผู้ใช้ที่มีอยู่ในฐานข้อมูล โดยรับพารามิเตอร์เป็นผู้ใช้และคืนค่าข้อผิดพลาดหากเกิดปัญหา
	SendChatMessage(sender, message string)              // ฟังก์ชันนี้ใช้สำหรับส่งข้อความแชทจากผู้ส่ง (sender) ไปยังข้อความ (message) ที่ระบุ
	LeaveChat(username string)                           // ฟังก์ชันนี้ใช้สำหรับให้ผู้ใช้ (username) ออกจากการแชท
}

// โครงสร้าง UserUsecase ใช้สำหรับการดำเนินการที่เกี่ยวข้องกับผู้ใช้ในระบบ
type UserUsecase struct { //  โครงสร้าง UserUsecase มีฟิลด์ต่างๆ เช่น UserRepo (interface UserRepository) สำหรับการเข้าถึงข้อมูลผู้ใช้, Config (*Config) สำหรับการกำหนดค่า, Logger (*log.Logger) สำหรับการเขียนล็อก, และ Constants (*Constants) สำหรับค่าคงที่ที่ใช้ในระบบ
	UserRepo  UserRepository // ฟิลด์สำหรับการเข้าถึงข้อมูลผู้ใช้
	Config    *Config        // ฟิลด์สำหรับการกำหนดค่า
	Logger    *log.Logger    // ฟิลด์สำหรับการเขียนล็อก
	Constants *Constants     // ฟิลด์สำหรับค่าคงที่ที่ใช้ในระบบ
}

// NewUserUsecase สร้างและคืนค่า UserUsecase ใหม่
func NewUserUsecase(repo UserRepository, config *Config, logger *log.Logger, constants *Constants) *UserUsecase {
	return &UserUsecase{
		UserRepo:  repo,      // กำหนดค่า UserRepo จากพารามิเตอร์ repo
		Config:    config,    // กำหนดค่า Config จากพารามิเตอร์ config
		Logger:    logger,    // กำหนดค่า Logger จากพารามิเตอร์ logger
		Constants: constants, // กำหนดค่า Constants จากพารามิเตอร์ constants
	}
}

// GetUserByID ดึงข้อมูลผู้ใช้ตาม ID
func (u *UserUsecase) GetUserByID(id int64) (*domain.User, error) {
	return u.UserRepo.GetByID(id) // เรียกใช้ฟังก์ชัน GetByID จาก UserRepo เพื่อดึงข้อมูลผู้ใช้
}

// Update ปรับปรุงข้อมูลผู้ใช้
func (u *UserUsecase) Update(user *domain.User) error {
	return u.UserRepo.Update(user) // เรียกใช้ฟังก์ชัน Update จาก UserRepo เพื่อปรับปรุงข้อมูลผู้ใช้
}

// SendChatMessage ส่งข้อความแชท
func (u *UserUsecase) SendChatMessage(sender, message string) {
	u.UserRepo.SendChatMessage(sender, message) // เรียกใช้ฟังก์ชัน SendChatMessage จาก UserRepo เพื่อส่งข้อความแชท
}

// LeaveChat ให้ผู้ใช้ (username) ออกจากการแชท
func (u *UserUsecase) LeaveChat(username string) {
	u.UserRepo.LeaveChat(username) // เรียกใช้ฟังก์ชัน LeaveChat จาก UserRepo เพื่อให้ผู้ใช้ออกจากการแชท
}

// CreateUser สร้างผู้ใช้ใหม่
func (u *UserUsecase) CreateUser(user *domain.User, role string) error {
	// ตรวจสอบชื่อผู้ใช้ หากไม่ถูกต้องให้คืนค่าข้อผิดพลาด
	if err := validateUsername(user.Username, u.Constants); err != nil {
		return err
	}

	user.Role = role                                                   // กำหนดบทบาทให้กับผู้ใช้
	hashedPassword, salt, err := HashPassword(user.Password, u.Config) // แฮชรหัสผ่านและสร้าง salt
	if err != nil {
		return err // หากเกิดข้อผิดพลาดในการแฮชคืนค่าข้อผิดพลาด
	}
	user.Password = hashedPassword // กำหนดรหัสผ่านที่แฮชแล้ว
	user.Salt = salt               // กำหนดค่า salt

	if err := u.UserRepo.Create(user); err != nil {
		if errors.Is(err, u.Constants.ErrUsernameAlreadyExists) {
			return u.Constants.ErrUsernameAlreadyExists // หากชื่อผู้ใช้มีอยู่แล้วให้คืนค่าข้อผิดพลาด
		}
		return err // คืนค่าข้อผิดพลาดอื่นๆ
	}

	u.Logger.Printf("User created: %s\n", user.Username) // บันทึกการสร้างผู้ใช้ในล็อก
	return nil                                           // คืนค่า nil หากสร้างผู้ใช้สำเร็จ
}

// Login ทำการเข้าสู่ระบบของผู้ใช้
func (u *UserUsecase) Login(username, password string) (*domain.User, error) {
	user, err := u.UserRepo.GetByUsername(username) // ดึงข้อมูลผู้ใช้จากฐานข้อมูลตามชื่อผู้ใช้
	if err != nil {
		return nil, u.Constants.ErrUserNotFound // หากไม่พบผู้ใช้ให้คืนค่าข้อผิดพลาด
	}

	if !ValidatePassword(password, user.Password, user.Salt, u.Config) {
		return nil, u.Constants.ErrInvalidPassword // หากรหัสผ่านไม่ถูกต้องให้คืนค่าข้อผิดพลาด
	}

	u.Logger.Printf("Login successful: %s with role: %s\n", user.Username, user.Role) // บันทึกการเข้าสู่ระบบของผู้ใช้ในล็อก
	return user, nil                                                                  // คืนค่าผู้ใช้และ nil หากเข้าสู่ระบบสำเร็จ
}

// HashPassword สร้างแฮชสำหรับรหัสผ่านที่เป็นข้อความธรรมดาที่ให้มา
func HashPassword(password string, config *Config) (string, []byte, error) {
	salt, err := generateSalt(config.SaltLength) // สร้าง salt สำหรับรหัสผ่าน
	if err != nil {
		return "", nil, err // หากเกิดข้อผิดพลาดในการสร้าง salt คืนค่าข้อผิดพลาด
	}
	hashedPassword := argon2.IDKey([]byte(password), salt, config.ArgonTime, config.ArgonMemory, config.ArgonThreads, config.ArgonKeyLen) // แฮชรหัสผ่านโดยใช้ argon2
	return base64.RawStdEncoding.EncodeToString(hashedPassword), salt, nil                                                                // คืนค่าแฮชที่เข้ารหัสเป็น base64 พร้อมกับ salt
}

// ValidatePassword เปรียบเทียบรหัสผ่านที่เป็นข้อความธรรมดากับรหัสผ่านที่แฮชเก็บไว้
func ValidatePassword(password, hashedPassword string, salt []byte, config *Config) bool {
	hash := argon2.IDKey([]byte(password), salt, config.ArgonTime, config.ArgonMemory, config.ArgonThreads, config.ArgonKeyLen) // แฮชรหัสผ่านที่ให้มา
	return hashedPassword == base64.RawStdEncoding.EncodeToString(hash)                                                         // คืนค่าผลลัพธ์ว่าแฮชที่เก็บไว้ตรงกับแฮชที่สร้างจากรหัสผ่านหรือไม่
}

// generateSalt สร้าง salt แบบสุ่ม
func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)               // สร้าง slice ของ byte ที่มีความยาวตามที่กำหนด
	if _, err := rand.Read(salt); err != nil { // อ่านข้อมูลแบบสุ่มลงใน salt
		return nil, err // หากเกิดข้อผิดพลาดในการอ่าน จะคืนค่าข้อผิดพลาด
	}
	return salt, nil // คืนค่า salt ที่สร้างขึ้น
}

// validateUsername ตรวจสอบชื่อผู้ใช้ว่าตรงตามข้อกำหนดหรือไม่
func validateUsername(username string, constants *Constants) error {
	if len(username) < 5 || len(username) > 20 { // ตรวจสอบความยาวของชื่อผู้ใช้
		return constants.ErrUsernameTooLong // หากชื่อผู้ใช้ยาวเกินไปหรือสั้นเกินไปให้คืนค่าข้อผิดพลาด
	}
	return nil // คืนค่า nil หากชื่อผู้ใช้ถูกต้อง
}
