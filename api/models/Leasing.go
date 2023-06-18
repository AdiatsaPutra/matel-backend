package models

import (
	"gorm.io/gorm"
)

type Leasing struct {
	ID          uint   `json:"id" gorm:"column:id"`
	NamaLeasing string `json:"nama_leasing" gorm:"column:nama_leasing"`
	NamaPIC     string `json:"nama_pic" gorm:"column:nama_pic"`
	NoHPPIC     string `json:"no_hp_pic" gorm:"column:no_hp_pic"`
	gorm.Model
}

// custom tablename
func (e *Leasing) TableName() string {
	return "m_leasing"
}
