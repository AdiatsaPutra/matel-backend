package models

import (
	"gorm.io/gorm"
)

type UserChange struct {
	ID              uint   `json:"id" gorm:"column:id"`
	UserID          uint   `json:"user_id" gorm:"column:user_id"`
	UnupdatedStatus uint   `json:"unupdated_status" gorm:"column:unupdated_status"`
	UpdatedStatus   uint   `json:"updated_status" gorm:"column:updated_status"`
	TimeUpdated     string `json:"start_subscrition" gorm:"column:start_subscrition"`
	gorm.Model
}

// custom tablename
func (e *UserChange) TableName() string {
	return "m_users_change"
}
