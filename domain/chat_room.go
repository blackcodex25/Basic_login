package domain

import (
	"log"
	"sync"
)

// ChatRoom แทนห้องแชทที่ผู้ใช้สามารถเข้าร่วม ออกจากห้อง และส่งข้อความได้
type ChatRoom struct {
	Messages chan ChatMessage    // channels สำหรับส่งข้อความ
	Join     chan string         // channels สำหรับผู้ใช้ที่เข้าร่วม
	Leave    chan string         // channels สำหรับผู้ใช้ที่ออกจากห้อง
	Users    map[string]struct{} // map สำหรับเก็บผู้ใช้
	mu       sync.Mutex          // Mutex สำหรับการเข้าถึง Users อย่างปลอดภัยในหลายเธรด
}

// NewChatRoom สร้าง ChatRoom ใหม่พร้อม channels ที่มีการบัฟเฟอร์
func NewChatRoom(bufferSize int) *ChatRoom {
	return &ChatRoom{
		Messages: make(chan ChatMessage, bufferSize), // channels สำหรับส่งข้อความ
		Join:     make(chan string, bufferSize),      // channels สำหรับผู้ใช้ที่เข้าร่วม
		Leave:    make(chan string, bufferSize),      // channels สำหรับผู้ใช้ที่ออกจากห้อง
		Users:    make(map[string]struct{}),          // map สำหรับเก็บผู้ใช้
	}
}

// Run เริ่มต้นห้องแชทและฟังข้อความ การเข้าร่วม และการออกจากห้อง
func (c *ChatRoom) Run() {
	for {
		select {
		case message := <-c.Messages: // รับข้อความจาก channels Messages
			go c.processMessage(message) // ประมวลผลข้อความพร้อมกัน
		case user := <-c.Join: // รับผู้ใช้จาก channels Join
			go c.processJoin(user) // ประมวลผลการเข้าร่วมของผู้ใช้พร้อมกัน
		case user := <-c.Leave: // รับผู้ใช้จาก channels Leave
			go c.processLeave(user) // ประมวลผลการออกจากห้องของผู้ใช้พร้อมกัน
		}
	}
}

// processMessage ประมวลผลและบันทึกข้อความที่ได้รับ
func (c *ChatRoom) processMessage(message ChatMessage) {
	log.Printf("[%s] %s: %s\n", message.TimeStamp.Format("2006-01-02 15:04:05"), message.Sender, message.Message) // บันทึกข้อความที่ได้รับ
}

// processJoin จัดการการเข้าร่วมของผู้ใช้ในห้องแชท
func (c *ChatRoom) processJoin(user string) {
	c.mu.Lock() // Lock เพื่อความปลอดภัยในการเข้าถึง Users เป็นไปอย่างปลอดภัยในหลายเธรด
	defer c.mu.Unlock()

	c.Users[user] = struct{}{}                // เก็บเฉพาะการมีอยู่ของผู้ใช้
	log.Printf("%s joined the chat.\n", user) // บันทึกการเข้าร่วมของผู้ใช้
}

// processLeave จัดการการออกจากห้องของผู้ใช้
func (c *ChatRoom) processLeave(user string) {
	c.mu.Lock() // Lock เพื่อความปลอดภัยในการเข้าถึง Users เป็นไปอย่างปลอดภัยในหลายเธรด
	defer c.mu.Unlock()

	delete(c.Users, user)                   // ลบผู้ใช้จาก map Users
	log.Printf("%s left the chat.\n", user) // บันทึกการออกจากห้องของผู้ใช้
}
