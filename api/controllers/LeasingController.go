package controllers

import (
	config "matel/configs"
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"matel/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLeasing(c *gin.Context) {
	searchQuery := c.Query("search") // Get the search query from the query string

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

func GetLeasingDetail(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)
	LeasingID := c.Param("id")

	if UserID == 0 {
		exceptions.AppException(c, "Not authorized")
		return
	}

	if LeasingID == "" {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	LeasingIDInt, err := strconv.Atoi(LeasingID)

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	leasing, err := repository.GetLeasingByID(c, uint(LeasingIDInt))

	if err != nil {
		exceptions.AppException(c, "Leasing tidak ditemukan")
		return
	}

	e := repository.AddSearchHistory(c, UserID, uint(LeasingIDInt))

	if e != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, leasing, "Leasing found", http.StatusOK)
}
