package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID                uint   `json:"id" gorm:"column:id"`
	UserName          string `json:"username"`
	Password          string `json:"password"`
	PasswordToView    string `json:"password_to_view"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	IsAdmin           uint   `json:"is_admin"`
	DeviceID          string `json:"device_id" gorm:"column:device_id"`
	ProvinceID        uint   `json:"province_id" gorm:"column:province_id"`
	KabupatenID       uint   `json:"kabupaten_id" gorm:"column:kabupaten_id"`
	KecamatanID       uint   `json:"kecamatan_id" gorm:"column:kecamatan_id"`
	Status            uint   `json:"status" gorm:"column:status"`
	SubscriptionMonth uint   `json:"subscription_month" gorm:"column:subscription_month"`
	StartSubscription string `json:"start_subscrition" gorm:"column:start_subscrition"`
	EndSubscription   string `json:"end_subscription" gorm:"column:end_subscription"`
	NoPolHistory      string `json:"nopol_history" gorm:"column:nopol_history"`
	Token             string `json:"token"`
	gorm.Model
}

// custom tablename
func (e *User) TableName() string {
	return "m_users"
}
