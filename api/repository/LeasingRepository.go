package repository

import (
	"database/sql"
	config "matel/configs"
	"matel/models"
	"strings"

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

func GetLeasingByNopolHistory(c *gin.Context, UserID uint) ([]models.Leasing, error) {
	var user models.User
	var leasings []models.Leasing

	err := config.InitDB().Model(&user).Where("id = ?", UserID).First(&user).Error
	if err != nil {
		return nil, err
	}

	numbersViewed := user.NoPolHistory
	numbers := strings.Split(numbersViewed, ",")

	// Query the leasings with the given numbers
	err = config.InitDB().Model(&models.Leasing{}).Where("id IN (?)", numbers).Find(&leasings).Error
	if err != nil {
		return nil, err
	}

	return leasings, nil
}
