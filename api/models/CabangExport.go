package models

type CabangExport struct {
	Cabang          string `json:"nama_cabang" gorm:"column:nama_cabang"`
	NoHP            string `json:"no_hp" gorm:"column:no_hp"`
	LatestCreatedAt string `json:"latest_created_at" gorm:"column:latest_created_at"`
}
