package models

import (
	"gorm.io/gorm"
)

type Kendaraan struct {
	ID          uint   `json:"id" gorm:"column:id"`
	Leasing     string `json:"leasing" gorm:"column:leasing"`
	Cabang      string `json:"cabang" gorm:"column:cabang"`
	CabangData  string `json:"cabangData" gorm:"column:cabangData"`
	NoKontrak   string `json:"no_kontrak" gorm:"column:noKontrak"`
	NamaDebitur string `json:"nama_debitur" gorm:"column:namaDebitur"`
	NomorPolisi string `json:"nomor_polisi" gorm:"column:nomorPolisi"`
	SisaHutang  string `json:"sisa_hutang" gorm:"column:sisaHutang"`
	Tipe        string `json:"tipe" gorm:"column:tipe"`
	Tahun       string `json:"tahun" gorm:"column:tahun"`
	NoRangka    string `json:"no_rangka" gorm:"column:noRangka"`
	NoMesin     string `json:"no_mesin" gorm:"column:noMesin"`
	PIC         string `json:"pic" gorm:"column:pic"`
	Status      uint   `json:"status" gorm:"column:status"`
	Versi       uint   `json:"versi" gorm:"column:versi"`
	gorm.Model
}

// custom tablename
func (e *Kendaraan) TableName() string {
	return "m_kendaraan"
}
