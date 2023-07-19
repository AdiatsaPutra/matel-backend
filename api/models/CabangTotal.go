package models

type CabangTotal struct {
	Cabang          string `json:"nama_cabang" gorm:"column:nama_cabang"`
	LatestCreatedAt string `json:"latest_created_at" gorm:"column:latest_created_at"`
	KendaraanTotal  uint   `json:"kendaraan_total" gorm:"column:kendaraan_total"`
}
