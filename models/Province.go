package models

type Province struct {
	ID   uint   `json:"id" gorm:"column:id"`
	Nama string `json:"name" gorm:"column:name"`
}

// custom tablename
func (e *Province) TableName() string {
	return "m_province"
}
