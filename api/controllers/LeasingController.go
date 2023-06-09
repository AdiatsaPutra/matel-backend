package controllers

import (
	config "matel/configs"
	"matel/models"
	"matel/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLeasing(c *gin.Context) {
	searchQuery := c.Query("search")          // Get the search query from the query string

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

