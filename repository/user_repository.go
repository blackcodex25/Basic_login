package repository

import (
	"Basic_login/domain"
	"errors"
	"log"
	"sync"
	"time"
)

// Error messages
var (
	// สร้าง errors.New และกำหนดค่าข้อความให้ตัวแปร ErrUserNotFound และ ErrUserExists
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

// สร้าง struct Repository ที่มีฟิลด์ชื่อ users ของ type []domain.User และ mutex ของ type sync.Mutex
type InMemoryUserRepository struct {
	mu            sync.RWMutex            // ใช้เพื่อควบคุมการเข้าถึงข้อมูลใน โครงสร้าง InMemoryUserRepository โดยให้การอ่านและเขียนพร้อมกันอย่างปลอดภัยในหลายๆ เธรด
	userIDCounter int64                   // เพิ่มตัวเลขเพื่อสร้าง ID ของผู้ใช้
	users         map[string]*domain.User // สร้าง map สําหรับเก็บข้อมูลผู้ใช้ โดยค่าจะเป็นตัวชี้ไปยังโครงสร้าง User
	userIDs       map[int64]*domain.User  // สร้าง map สำหรับเก็บผู้ใช้ตามรหัสประจำตัว
	chatRoom      *domain.ChatRoom        // ตัวแปรที่เก็บข้อมูลของห้องสนทนา
}

// ฟังก์ชัน NewInMemoryUserRepository พร้อมกับการกำหนดขนาดของ buffer สำหรับห้องสนทนา
func NewInMemoryUserRepository(bufferSize int) *InMemoryUserRepository {
	repo := &InMemoryUserRepository{ // repo เป็น pointer ไปยัง InMemoryUserRepository ใหม่
		// ซึ่งมีการสร้าง map สำหรับเก็บผู้ใช้ และรหัสประจำตัวผู้ใช้
		users:   make(map[string]*domain.User),
		userIDs: make(map[int64]*domain.User),
		// สร้างห้องสนทนาพร้อมกับขนาด buffer ที่กำหนด
		chatRoom: domain.NewChatRoom(bufferSize),
	}
	go repo.chatRoom.Run() // จะเริ่มการทำงานของฟังก์ชัน Run ใน goroutine ใหม่ เพื่อให้ห้องสนทนาสามารถทำงานได้ในพื้นหลัง
	return repo            // คืนค่า repo ซึ่งเป็น instance ของ InMemoryUserRepository
}

// ฟังก์ชัน GetByID ใช้ในการดึงข้อมูลผู้ใช้จาก InMemoryUserRepository ตามรหัสประจำตัว (ID)
// parameter id ใช้เพื่อระบุรหัสประจำตัวของผู้ใช้ที่ต้องการดึงข้อมูล
func (repo *InMemoryUserRepository) GetByID(id int64) (*domain.User, error) {
	repo.mu.RLock()         // เรียกใช้เพื่อทำการล็อกการอ่าน (Read Lock) ซึ่งช่วยให้มั่นใจว่าข้อมูลใน userIDs จะไม่ถูกเปลี่ยนแปลง
	defer repo.mu.RUnlock() // ใช้เพื่อปลดล็อกการอ่านเมื่อฟังก์ชันสิ้นสุด

	user, exists := repo.userIDs[id] // ตรวจสอบว่ามีผู้ใช้ตามรหัสประจำตัวที่ระบุหรือไม่
	if !exists {                     // ถ้าผู้ใช้ไม่พบ ฟังก์ชันจะคืนค่า nil และส่งคืนข้อผิดพลาด
		return nil, ErrUserNotFound
	}
	return user, nil // หากผู้ใช้พบ ฟังก์ชันจะคืนค่า pointer ไปยัง user ที่ค้นพบ
}

// ฟังก์ชัน Create ใช้ในการเพิ่มผู้ใช้ใหม่ลงใน InMemoryUserRepository
func (repo *InMemoryUserRepository) Create(user *domain.User) error {
	// ล็อกการเขียน
	repo.mu.Lock()         // ใช้เพื่อทำการล็อกการเขียน (write lock) เพื่อป้องกันการเข้าถึงข้อมูลพร้อมกันจากหลายเธรด
	defer repo.mu.Unlock() // ใช้เพื่อปลดล็อกการเขียนเมื่อฟังก์ชันสิ้นสุด

	// ตรวจสอบผู้ใช้
	if _, exists := repo.users[user.Username]; exists { // ตรวจสอบว่าผู้ใช้ที่มีชื่อผู้ใช้นีี้มีอยู่แล้วใน repository หรือไม่
		return ErrUserExists // ถ้ามีอยู่แล้ว ฟังก์ชันจะคืนค่าข้อผิดพลาด
	}

	// เพิ่มผู้ใช้ใหม่
	repo.userIDCounter++                                               // เพิ่มค่าตัวนับ ID ของผู้ใช้ใหม่
	user.ID = repo.userIDCounter                                       // กำหนดค่า ID ของผู้ใช้ใหม่ให้กับ user
	repo.users[user.Username] = user                                   // เพิ่มผู้ใช้ใหม่ลงใน users
	repo.userIDs[user.ID] = user                                       // เพิ่มผู้ใช้ใหม่ลงใน userIDs
	log.Printf("Created user: %s with ID: %d", user.Username, user.ID) // บันทึกการสร้างผู้ใช้ใหม่ใน log

	// เข้าร่วมในห้องสนทนา
	repo.chatRoom.Join <- user.Username // ส่งชื่อผู้ใช้ไปยัง channel ของห้องสนทนาเพื่อให้ผู้ใช้เข้าร่วม
	return nil                          // คืนค่า nil ถ้าการสร้างผู้ใช้สําเร็จ
}

// ฟังก์ชัน GetByUsername มีพารามิเตอร์ username ในการระบุชื่อผู้ใช้ และคืนค่า user และดึงข้อมูลจาก InMemoryUserRepository
// ตามชื่อผู้ใช้
func (repo *InMemoryUserRepository) GetByUsername(username string) (*domain.User, error) {
	repo.mu.RLock()         // ทำการล็อกการอ่าน เพื่อป้องกันไม่ให้มีการเปลี่ยนแปลงข้อมูลในขณะทำการอ่าน
	defer repo.mu.RUnlock() // ใช้เพื่อปลดล็อคเมื่อฟังก์ชันสิ้นสุด

	// ตรวจสอบผู้ใช้
	user, exists := repo.users[username] // ตรวจสอบว่าชื่อผู้ใช้มีอยู่ใน map ของ users หรือไม่
	if !exists {                         // ถ้าผู้ใช้ไม่พบ ฟังก์ชันจะคืนค่า nil และส่งคืนข้อผิดพลาด
		return nil, ErrUserNotFound
	}
	return user, nil // หากผู้ใช้พบ ฟังก์ชันจะคืนค่า user ไปยัง pointer domain.User
}

// ฟังก์ชัน GetAll ใช้ในการดึงข้อมูลผู้ใช้ทั้งหมดจาก InMemoryUserRepository และคืนค่า users และ error
func (repo *InMemoryUserRepository) GetAll() ([]*domain.User, error) {
	repo.mu.RLock()         // ทำการล็อกการอ่าน เพื่อป้องกันไม่ให้มีการเปลี่ยนแปลงข้อมูลในขณะทำการอ่าน
	defer repo.mu.RUnlock() // ใช้เพื่อปลดล็อคเมื่อฟังก์ชันสิ้นสุด

	// สร้าง slice สําหรับเก็บข้อมูลผู้ใช้
	users := make([]*domain.User, 0, len(repo.users)) // สร้าง slice สำหรับเก็บข้อมูลผู้ใช้ทั้งหมด โดยมีขนาดเริ่มต้นเท่ากับจำนวนผู้ใช้ใน repo

	// เพิ่มผู้ใช้ลงใน slice
	for _, user := range repo.users { // ใช้ for เพื่อวนลอบผู้ใช้ใน map users
		users = append(users, user) // เพิ่มผู้ใช้แต่ละคนลงใน slice
	}

	// ฟังก์ชันนี้จะคืนค่า slice ของผู้ใช้ทั้งหมดพร้อมกับ error
	return users, nil // คืนค่า slice ของผู้ใช้ทั้งหมด และส่งคืนค่า nil
}

// ฟังก์ชัน Update ใช้ในการปรับปรุงข้อมูลผู้ใช้ที่มีอยู่ใน InMemoryUserRepository
func (repo *InMemoryUserRepository) Update(user *domain.User) error {
	repo.mu.Lock()         // ทำการล็อกการอ่าน เพื่อป้องกันไม่ให้มีการเปลี่ยนแปลงข้อมูลในขณะทำการอ่าน
	defer repo.mu.Unlock() // ใช้เพื่อปลดล็อคเมื่อฟังก์ชันสิ้นสุด

	// ตรวจสอบผู้ใช้
	if _, exists := repo.users[user.Username]; !exists { // ตรวจสอบว่าผู้ใช้ที่ต้องการอัปเดตมีอยู่ใน repository หรือไม่
		return ErrUserNotFound // ถ้าผู้ใช้ไม่พบ ฟังก์ชันจะคืนค่า ErrUserNotFound
	}

	// อัปเดตข้อมูลผู้ใช้
	// จะอัปเดตข้อมูลผู้ใช้ใน repo.users และ repo.userIDs ด้วยค่าใหม่จากผู้ใช้ที่รับเข้ามา
	repo.users[user.Username] = user
	repo.userIDs[user.ID] = user

	// บันทึกข้อมูล
	log.Printf("Updated user: %s", user.Username) // บันทึกการปรับปรุงข้อมูลผู้ใช้ใน log

	return nil // คืนค่า nil ถ้าการปรับปรุงข้อมูลผู้ใช้สําเร็จ
}

// ฟังก์ชัน SendChatMessage ใช้ในการส่งข้อความไปยัง chat room กำหนดพารามิเตอร์ sender ชื่อผู้ส่ง และ message ข้อความที่ต้องการส่ง
// โดยดึงข้อมูลผู้ใช้จาก InMemoryUserRepository
func (repo *InMemoryUserRepository) SendChatMessage(sender, message string) {
	// สร้างโครงสร้างข้อความของแชทใหม่ โดยตั้งค่าฟิลด์
	chatMessage := domain.ChatMessage{
		Sender:    sender,     // ชื่อผู้ส่ง
		Message:   message,    // ข้อความที่ส่ง
		TimeStamp: time.Now(), // เวลาที่ส่งข้อความ
	}
	// ส่งข้อความไปยัง channels Messages ของห้องสนทนา ซึ่งจะทำให้ห้องสนทนาได้รับข้อความ
	repo.chatRoom.Messages <- chatMessage
}

// ฟังก์ชัน LeaveChat กำหนดพารามิเตอร์ username ใช้ในการนำผู้ใช้ออกจากห้องสนทนา (chat room)
func (repo *InMemoryUserRepository) LeaveChat(username string) {
	// พารามิเตอร์ username ชื่อผู้ใช้ที่ต้องการออกจากห้องสนทนา
	repo.chatRoom.Leave <- username // ส่งชื่อผู้ใช้ไปยัง channel Leave ของห้องสนทนา
}
