package models

import (
	"gorm.io/gorm"
)

type UserDetail struct {
	ID                uint   `json:"id" gorm:"column:id"`
	UserName          string `json:"username"`
	Password          string `json:"password"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	IsAdmin           uint   `json:"is_admin"`
	DeviceID          string `json:"device_id" gorm:"column:device_id"`
	ProvinceID        uint   `json:"province_id" gorm:"column:province_id"`
	ProvinceName      string `json:"province_name" gorm:"column:province_name"`
	KabupatenID       uint   `json:"kabupaten_id" gorm:"column:kabupaten_id"`
	KabupatenName     string `json:"kabupaten_name" gorm:"column:kabupaten_name"`
	KecamatanID       uint   `json:"kecamatan_id" gorm:"column:kecamatan_id"`
	KecamatanName     string `json:"kecamatan_name" gorm:"column:kecamatan_name"`
	Status            uint   `json:"status" gorm:"column:status"`
	SubscriptionMonth uint   `json:"subscription_month" gorm:"column:subscription_month"`
	StartSubscription string `json:"start_subscrition" gorm:"column:start_subscrition"`
	EndSubscription   string `json:"end_subscription" gorm:"column:end_subscription"`
	NoPolHistory      string `json:"nopol_history" gorm:"column:nopol_history"`
	Token             string `json:"token"`
	gorm.Model
}
