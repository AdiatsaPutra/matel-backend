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
	result := config.InitDB().Raw("SELECT COUNT(*) FROM m_leasing WHERE deleted_at IS NULL").Scan(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return uint(count.Int64), nil

}

func GetLeasingChart(c *gin.Context) ([]models.LeasingChart, error) {
	var leasingChart []models.LeasingChart
	query := `SELECT m_kendaraan.cabang AS leasing_name, COUNT(*) AS count
	FROM m_kendaraan
	WHERE id < 99999999 AND deleted_at IS NULL
	GROUP BY cabang
	ORDER BY cabang;`
	result := config.InitDB().Raw(query).Scan(&leasingChart)

	if result.Error != nil {
		return leasingChart, result.Error
	}

	return leasingChart, nil

}

func GetKendaraanPerCabangTotal(c *gin.Context, leasing string, cabang string) (uint, error) {
	var count sql.NullInt64
	var kendaraan models.Kendaraan

	if cabang != "" {

		result := config.InitDB().Model(&kendaraan).Count(&count.Int64).Where("leasing = ?", leasing).Where("cabang = ?", cabang)

		if result.Error != nil {
			return 0, result.Error
		}

		return uint(count.Int64), nil
	}

	result := config.InitDB().Model(&kendaraan).Count(&count.Int64).Where("leasing = ?", leasing)

	if result.Error != nil {
		return 0, result.Error
	}

	return uint(count.Int64), nil

}

func GetKendaraanTotal(c *gin.Context) (uint, error) {
	var count sql.NullInt64
	result := config.InitDB().Raw(`SELECT COUNT(*)
	FROM m_kendaraan
	WHERE id < 99999999;`).Scan(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return uint(count.Int64), nil

}

func GetLeasingByID(c *gin.Context, LeasingID uint) (models.Kendaraan, error) {
	var leasing = models.Kendaraan{}
	result := config.InitDB().Unscoped().Where("id = ?", LeasingID).First(&leasing)

	if result.Error != nil {
		return leasing, result.Error
	}

	return leasing, nil

}

func UpdateSearched(c *gin.Context, LeasingID uint) error {
	var user models.Kendaraan

	err := config.InitDB().Model(&user).Where("id = ?", LeasingID).Update("searched", 1).Error

	if err != nil {
		return err
	}
	return nil

}

func GetLeasingByNopolHistory(c *gin.Context, UserID uint) ([]models.Kendaraan, error) {
	var user models.User
	var leasings []models.Kendaraan

	err := config.InitDB().Model(&user).Where("id = ?", UserID).First(&user).Error
	if err != nil {
		return nil, err
	}

	numbersViewed := user.NoPolHistory
	numbers := strings.Split(numbersViewed, ",")

	// Query the leasings with the given numbers
	err = config.InitDB().Model(&models.Kendaraan{}).Where("id IN (?)", numbers).Find(&leasings).Error
	if err != nil {
		return nil, err
	}

	return leasings, nil
}
