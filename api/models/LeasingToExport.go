package models

type LeasingToExport struct {
	ID          string `gorm:"column:id"`
	CabangID    string `gorm:"column:cabang_id"`
	Cabang      string `gorm:"column:cabang"`
	NomorPolisi string `gorm:"column:nomorPolisi"`
	NoRangka    string `gorm:"column:noRangka"`
	NoMesin     string `gorm:"column:noMesin"`
	Versi       string `gorm:"column:versi"`
}
