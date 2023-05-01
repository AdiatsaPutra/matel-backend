package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"column:id"`
	UserName  string    `json:"username"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	DeviceId  string    `json:"device_id"`
	IsAdmin   uint      `json:"is_admin"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// custom tablename
func (e *User) TableName() string {
	return "m_users"
}
