package models

import (
	"gorm.io/gorm"
)

type BankInfoPembayaran struct {
	ID         uint   `json:"id" gorm:"column:id"`
	Nama       string `json:"nama_bank" gorm:"column:nama_bank"`
	NoRekening string `json:"no_rekening" gorm:"column:no_rekening"`
	BankID     uint   `json:"bank_id" gorm:"column:bank_id"`
	Image      string `json:"image" gorm:"column:image"`
	Bank       string `json:"bank" gorm:"column:bank"`
	gorm.Model
}
