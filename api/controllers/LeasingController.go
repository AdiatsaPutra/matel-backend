package controllers

import (
	config "matel/configs"
	"matel/models"
	"matel/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLeasing(c *gin.Context) {
	// pageStr := c.DefaultQuery("page", "1")    // Get the page parameter from the query string
	// limitStr := c.DefaultQuery("limit", "20") // Get the limit parameter from the query string
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

	var leasing []models.Leasing

	query := config.InitDB().Limit(100)

	if searchQuery != "" {
        query = query.Find(&leasing).Where("leasing LIKE ? OR cabang LIKE ? OR nomorPolisi LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
    }

	if err := query.Find(&leasing).Error; err != nil {
        return
    }

	payloads.HandleSuccess(c, leasing, "Leasing found", http.StatusOK)
}

