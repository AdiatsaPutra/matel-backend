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

func GetKendaraan(c *gin.Context) {
	pageNumber, _ := strconv.Atoi(c.Query("page"))
	if pageNumber <= 0 {
		pageNumber = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit <= 0 {
		limit = 10
	}

	query := config.InitDB().Model(&models.Kendaraan{})

	if leasing := c.Query("leasing"); leasing != "" {
		query = query.Where("leasing LIKE ?", "%"+leasing+"%")
	}

	if cabang := c.Query("cabang"); cabang != "" {
		query = query.Where("cabang LIKE ?", "%"+cabang+"%")
	}

	query = query.Order("created_at ASC")

	var kendaraans []models.Kendaraan
	offset := (pageNumber - 1) * limit

	result := query.Offset(offset).Limit(limit).Find(&kendaraans)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, kendaraans, "Kendaraan found", http.StatusOK)
}

func GetLeasing(c *gin.Context) {
	searchQuery := c.Query("search") // Get the search query from the query string

	var kendaraan []models.Kendaraan

	query := config.InitDB().Limit(100)

	if searchQuery != "" {
		query = query.Find(&kendaraan).Where("leasing LIKE ? OR cabang LIKE ? OR nomorPolisi LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	if err := query.Find(&kendaraan).Error; err != nil {
		return
	}

	payloads.HandleSuccess(c, kendaraan, "Leasing found", http.StatusOK)
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

func GetLeasingHistory(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)

	if UserID == 0 {
		exceptions.AppException(c, "Not authorized")
		return
	}

	leasing, err := repository.GetLeasingByNopolHistory(c, uint(UserID))

	if err != nil {
		exceptions.AppException(c, "Leasing tidak ditemukan")
		return
	}

	payloads.HandleSuccess(c, leasing, "Leasing found", http.StatusOK)
}
