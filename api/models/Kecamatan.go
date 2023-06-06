package models

import "gorm.io/gorm"

type Kecamatan struct {
	ID          uint   `json:"id" gorm:"column:id"`
	ProvinceID  uint   `json:"province_id" gorm:"column:province_id"`
	KabupatenID uint   `json:"kabupaten_id" gorm:"column:kabupaten_id"`
	Name        string `json:"name" gorm:"column:name"`
	gorm.Model
}

// custom tablename
func (e *Kecamatan) TableName() string {
	return "m_kecamatan"
}
