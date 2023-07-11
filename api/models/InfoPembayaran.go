package models

import (
	"gorm.io/gorm"
)

type InfoPembayaran struct {
	ID         uint   `json:"id" gorm:"column:id"`
	Nama       string `json:"nama_bank" gorm:"column:nama_bank"`
	NoRekening string `json:"no_rekening" gorm:"column:no_rekening"`
	FotoBank   string `json:"foto_bank" gorm:"column:foto_bank"`
	gorm.Model
}

// custom tablename
func (e *InfoPembayaran) TableName() string {
	return "m_info_pembayaran"
}
