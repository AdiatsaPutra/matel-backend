package models

import (
	"gorm.io/gorm"
)

type InfoPembayaran struct {
	ID         uint   `json:"id" gorm:"column:id"`
	NoRekening string `json:"no_rekening" gorm:"column:no_rekening"`
	BankID     uint   `json:"bank_id" gorm:"column:bank_id"`
	gorm.Model
}

// custom tablename
func (e *InfoPembayaran) TableName() string {
	return "m_info_pembayaran"
}
