package controllers

import (
	config "matel/configs"
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBanks(c *gin.Context) {
	var banks []models.InfoPembayaran

	// Get query parameters
	searchQuery := c.Query("search")

	// Prepare database query
	db := config.InitDB()
	query := db
	if searchQuery != "" {
		// Add search condition to the query
		query = query.Where("name LIKE ?", "%"+searchQuery+"%")
	}

	// Retrieve banks from the database
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
