package auth

import "time"

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleManager   Role = "manager"
	RoleUser      Role = "user"
	RoleReception Role = "reception"
)

type User struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Role      Role      `json:"role" gorm:"default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
