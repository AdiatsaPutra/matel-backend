package models

type Bank struct {
	ID    uint   `json:"id" gorm:"column:id"`
	Image string `json:"image" gorm:"column:image"`
	Bank  string `json:"bank" gorm:"column:bank"`
}

// custom tablename
func (e *Bank) TableName() string {
	return "m_bank"
}
