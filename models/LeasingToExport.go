package models

type LeasingToExport struct {
	NomorPolisi string `gorm:"column:nomorPolisi"`
	NoRangka    string `gorm:"column:noRangka"`
	NoMesin     string `gorm:"column:noMesin"`
}
