package models

import (
	"gorm.io/gorm"
)

type Cabang struct {
	ID         uint   `json:"id" gorm:"column:id"`
	LeasingID  uint   `json:"leasing_id" gorm:"column:leasing_id"`
	NamaCabang string `json:"nama_cabang" gorm:"column:nama_cabang"`
	NoHP       string `json:"no_hp" gorm:"column:no_hp"`
	Versi      int    `json:"versi" gorm:"column:versi"`
	gorm.Model
}

// custom tablename
func (e *Cabang) TableName() string {
	return "m_cabang"
}
