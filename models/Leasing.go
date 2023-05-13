package models

import (
	"gorm.io/gorm"
)

type Leasing struct {
	ID          uint   `json:"id" gorm:"column:id"`
	Leasing     string `json:"leasing" gorm:"column:leasing"`
	Cabang      string `json:"cabang" gorm:"column:cabang"`
	NoKontrak   string `json:"no_kontrak" gorm:"column:no_kontrak"`
	NamaDebitur string `json:"nama_debitur" gorm:"column:nama_debitur"`
	NomorPolisi string `json:"nomor_polisi" gorm:"column:nomor_polisi"`
	SisaHutang  uint   `json:"sisa_hutang" gorm:"column:sisa_hutang"`
	Tipe        string `json:"tipe" gorm:"column:tipe"`
	Tahun       string `json:"tahun" gorm:"column:tahun"`
	NoRangka    string `json:"no_rangka" gorm:"column:no_rangka"`
	NoMesin     string `json:"no_mesin" gorm:"column:no_mesin"`
	PIC         string `json:"pic" gorm:"column:pic"`
	Status      uint   `json:"status" gorm:"column:status"`
	gorm.Model
}

// custom tablename
func (e *Leasing) TableName() string {
	return "m_leasing"
}
