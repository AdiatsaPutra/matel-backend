package controllers

import (
	config "matel/configs"
	"matel/models"
	"matel/payloads"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLeasing(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")    // Get the page parameter from the query string
	limitStr := c.DefaultQuery("limit", "20") // Get the limit parameter from the query string
	searchQuery := c.Query("search")          // Get the search query from the query string

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid limit parameter"})
		return
	}

	db := config.InitDB()

	var leasing []models.Leasing
	var total int64

	query := db

	if searchQuery != "" {
		query = query.Where("leasing LIKE ? OR cabang LIKE ? OR nomorPolisi LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	query.Find(&leasing)

	if err := query.Model(&models.Leasing{}).Count(&total).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve leasing"})
		return
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(1000)
	query.Find(&leasing)

	data := make(map[string]interface{})
	data["total"] = total
	data["leasing"] = leasing

	payloads.HandleSuccess(c, data, "Leasing found", http.StatusOK)
}

