package controllers

import (
	config "matel/configs"
	"matel/models"
	"matel/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLeasing(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")    // Get the page parameter from the query string
	limitStr := c.DefaultQuery("limit", "20") // Get the limit parameter from the query string
	searchQuery := c.Query("search")          // Get the search query from the query string

	// page, err := strconv.Atoi(pageStr)
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "Invalid page parameter"})
	// 	return
	// }

	// limit, err := strconv.Atoi(limitStr)
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "Invalid limit parameter"})
	// 	return
	// }

	db := config.InitDB()

	var leasing []models.Leasing

	query := db

	if searchQuery != "" {
		query = query.Where("leasing LIKE ? OR cabang LIKE ? OR nomorPolisi LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	query.Find(&leasing).Limit(1000)

	data := make(map[string]interface{})
	data["leasing"] = leasing

	payloads.HandleSuccess(c, data, "Leasing found", http.StatusOK)
}

