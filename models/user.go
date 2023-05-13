package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"column:id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	IsAdmin   uint      `json:"is_admin"`
	Token     string    `json:"token"`
	gorm.Model
}

// custom tablename
func (e *User) TableName() string {
	return "m_users"
}
