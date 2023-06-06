package models

import "gorm.io/gorm"

type Kabupaten struct {
	ID         uint   `json:"id" gorm:"column:id"`
	ProvinceID uint   `json:"province_id" gorm:"column:province_id"`
	Name       string `json:"name" gorm:"column:name"`
	gorm.Model
}

// custom tablename
func (e *Kabupaten) TableName() string {
	return "m_kabupaten"
}
