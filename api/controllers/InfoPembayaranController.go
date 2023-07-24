package controllers

import (
	config "matel/configs"
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBanks(c *gin.Context) {
	var banks []models.BankInfoPembayaran

	// Get query parameters
	searchQuery := c.Query("search")

	// Prepare database query
	db := config.InitDB()

	// Construct the SQL query with the required fields and the LEFT JOIN
	query := db.Table("m_info_pembayaran").
		Select("m_info_pembayaran.id, m_info_pembayaran.no_rekening, m_info_pembayaran.bank_id, m_bank.image, m_bank.bank, m_info_pembayaran.created_at, m_info_pembayaran.updated_at").
		Joins("LEFT JOIN m_bank ON m_info_pembayaran.bank_id = m_bank.id")

	if searchQuery != "" {
		// Add search condition to the query
		query = query.Where("m_bank.bank LIKE ?", "%"+searchQuery+"%")
	}

	// Retrieve banks from the database
	err := query.Find(&banks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, banks, "Banks found", 200)
}

func GetBankData(c *gin.Context) {
	var banks []models.Bank

	db := config.InitDB()
	query := db

	err := query.Find(&banks).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, banks, "Banks found", 200)
}

func GetBank(c *gin.Context) {
	id := c.Param("id")

	var bank models.InfoPembayaran
	err := config.InitDB().First(&bank, id).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	c.JSON(http.StatusOK, bank)
}

func CreateBank(c *gin.Context) {
	var bank models.InfoPembayaran
	err := c.BindJSON(&bank)
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	err = config.InitDB().Create(&bank).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	c.JSON(http.StatusCreated, bank)
}

func UpdateBank(c *gin.Context) {
	id := c.Param("id")

	var bank models.InfoPembayaran
	err := config.InitDB().First(&bank, id).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	err = c.BindJSON(&bank)
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	err = config.InitDB().Save(&bank).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	c.JSON(http.StatusOK, bank)
}

func DeleteBank(c *gin.Context) {
	id := c.Param("id")

	var bank models.InfoPembayaran
	err := config.InitDB().First(&bank, id).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	err = config.InitDB().Delete(&bank).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank deleted successfully"})
}
