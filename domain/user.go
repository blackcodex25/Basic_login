package domain

// โครงสร้างข้อมูล User
type User struct {
	Salt     []byte // ใช้สำหรับเก็บค่า salt ที่ใช้ในการ hash รหัสผ่าน
	Username string // เก็บชื่อผู้ใช้งาน
	Password string // เก็บรหัสผ่าน
	Role     string // บทบาทผู้ใช้ เช่น "admin", "user", "guest"
	ID       int64  // รหัสประจำตัวผู้ใช้ การระบุผู้ใช้ในระบบ
}

//  effectively
