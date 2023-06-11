package models

type LeasingToExport struct {
	ID string `gorm:"column:id"`
	NomorPolisi string `gorm:"column:nomorPolisi"`
	NoRangka    string `gorm:"column:noRangka"`
	NoMesin     string `gorm:"column:noMesin"`
}
