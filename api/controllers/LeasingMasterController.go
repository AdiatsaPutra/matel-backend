package controllers

import (
	"fmt"
	config "matel/configs"
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateLeasing(c *gin.Context) {

	var payload models.Leasing
	if err := c.ShouldBindJSON(&payload); err != nil {
		exceptions.BadRequest(c)
		return
	}

	leasing := models.Leasing{
		NamaLeasing: payload.NamaLeasing,
		NamaPIC:     payload.NamaPIC,
		NoHPPIC:     payload.NoHPPIC,
	}

	result := config.InitDB().Create(&leasing)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, leasing, "Leasing created", 200)

}

func GetLeasingMaster(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	search := c.Query("search")

	offset := (page - 1) * limit

	var leasings []models.Leasing
	var count int64

	query := config.InitDB().Model(&models.Leasing{})

	if search != "" {
		query = query.Where("nama_leasing LIKE ?", fmt.Sprintf("%%%s%%", search))
	}

	query = query.Order("created_at ASC")
	
	query.Count(&count)

	result := query.Offset(offset).Limit(limit).Find(&leasings)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	data := make(map[string]interface{})
	data["leasing"] = leasings
	data["total"] = count

	payloads.HandleSuccess(c, data, "Data found", 200)

}

func UpdateLeasing(c *gin.Context) {
	leasingID := c.Param("id")

	var payload models.Leasing
	if err := c.ShouldBindJSON(&payload); err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	var leasing models.Leasing
	result := config.InitDB().First(&leasing, leasingID)
	if result.Error != nil {
		payloads.HandleSuccess(c, "Leasing not found", "Success", 200)
		return
	}

	leasing.NamaLeasing = payload.NamaLeasing
	leasing.NamaPIC = payload.NamaPIC
	leasing.NoHPPIC = payload.NoHPPIC

	result = config.InitDB().Save(&leasing)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, leasing, "Success", 200)
}

func DeleteLeasing(c *gin.Context) {

	leasingID := c.Param("id")

	var leasing models.Leasing
	result := config.InitDB().First(&leasing, leasingID)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	result = config.InitDB().Delete(&leasing)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, "Leasing deleted", "Success", 200)
}
