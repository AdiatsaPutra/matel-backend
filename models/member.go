package models

import (
	"time"
)

type Member struct {
	ID        uint      `json:"id" gorm:"column:id"`
	Status    uint      `json:"status" gorm:"column:status"`
	UserID    uint      `json:"user_id" gorm:"column:user_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// custom tablename
func (e *Member) TableName() string {
	return "m_member"
}
