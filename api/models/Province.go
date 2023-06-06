package models

import "gorm.io/gorm"

type Province struct {
	ID   uint   `json:"id" gorm:"column:id"`
	Nama string `json:"name" gorm:"column:name"`
	gorm.Model
}

// custom tablename
func (e *Province) TableName() string {
	return "m_province"
}
