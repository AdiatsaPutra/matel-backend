package models

import "gorm.io/gorm"

type Home struct {
	ID             uint `json:"id" gorm:"column:id"`
	KendaraanTotal uint `json:"kendaraan_total" gorm:"column:kendaraan_total"`
	gorm.Model
}

// custom tablename
func (e *Home) TableName() string {
	return "m_home"
}
