package domain

type User struct {
	Salt     []byte
	Username string
	Password string
	Role     string // เช่น "admin", "user", "guest"
	ID       int64
}
