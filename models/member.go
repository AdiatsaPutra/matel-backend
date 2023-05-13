package models

import (
	"gorm.io/gorm"
)

type Member struct {
	ID          uint      `json:"id" gorm:"column:id"`
	Status      uint      `json:"status" gorm:"column:status"`
	UserID      uint      `json:"user_id" gorm:"column:user_id"`
	DeviceID    string    `json:"device_id" gorm:"column:device_id"`
	ProvinceID  uint      `json:"province_id" gorm:"column:province_id"`
	KabupatenID uint      `json:"kabupaten_id" gorm:"column:kabupaten_id"`
	KecamatanID uint      `json:"kecamatan_id" gorm:"column:kecamatan_id"`
	gorm.Model
}

// custom tablename
func (e *Member) TableName() string {
	return "m_member"
}
