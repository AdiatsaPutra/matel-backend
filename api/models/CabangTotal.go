package models

type CabangTotal struct {
	ID              int    `json:"id" gorm:"column:id"`
	Cabang          string `json:"nama_cabang" gorm:"column:nama_cabang"`
	LatestCreatedAt string `json:"latest_created_at" gorm:"column:latest_created_at"`
	KendaraanTotal  uint   `json:"kendaraan_total" gorm:"column:kendaraan_total"`
	NoHP            string `json:"no_hp" gorm:"column:no_hp"`
}
