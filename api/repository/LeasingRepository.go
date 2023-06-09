package repository

import (
	"database/sql"
	config "matel/configs"
	"matel/models"

	"github.com/gin-gonic/gin"
)

func GetLeasingTotal(c *gin.Context) (uint, error) {
	var count sql.NullInt64
	result := config.InitDB().Raw("SELECT COUNT(*) FROM m_leasing").Scan(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return uint(count.Int64), nil

}

func GetLeasingByID(c *gin.Context, LeasingID uint) (models.Leasing, error) {
	var leasing = models.Leasing{}
	result := config.InitDB().First(&leasing, LeasingID)

	if result.Error != nil {
		return models.Leasing{}, result.Error
	}

	return leasing, nil

}

func UpdateSearched(c *gin.Context, LeasingID uint) error {
	var user models.Leasing

	err := config.InitDB().Model(&user).Where("id = ?", LeasingID).Update("searched", 1).Error

	if err != nil {
		return err
	}
	return nil

}
