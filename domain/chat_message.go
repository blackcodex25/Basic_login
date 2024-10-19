package domain

import (
	"time"
)

// ChatMessage แทนข้อความในห้องแชท
type ChatMessage struct {
	TimeStamp time.Time // เวลาที่ส่งข้อความ
	Sender    string    // ผู้ส่งข้อความ
	Message   string    // ข้อความ
}
