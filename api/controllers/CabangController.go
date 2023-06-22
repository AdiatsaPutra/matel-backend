package controllers

import (
	"fmt"
	config "matel/configs"
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateCabang(c *gin.Context) {
	var payload models.Cabang
	if err := c.ShouldBindJSON(&payload); err != nil {
		exceptions.BadRequest(c)
		return
	}

	cabang := models.Cabang{
		LeasingID: payload.LeasingID,
		NamaCabang: payload.NamaCabang,
	}

	result := config.InitDB().Create(&cabang)
	logrus.Info(result.Error)
	logrus.Info(cabang)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, cabang, "Cabang created", 200)

}

func GetCabang(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 0 
	}

	search := c.Query("search")
	leasingID := c.Query("leasing_id")

	offset := (page - 1) * limit

	var cabang []models.Cabang
	var count int64

	query := config.InitDB().Model(&models.Cabang{})

	if search != "" {
		query = query.Where("nama_cabang LIKE ?", fmt.Sprintf("%%%s%%", search))
	}

	if leasingID != "" {
		query = query.Where("leasing_id = ?", leasingID)
	}

	query = query.Order("created_at DESC")
	query.Count(&count)

	if page == 0 && limit == 0 {
		query.Find(&cabang)
	} else {
		result := query.Offset(offset).Limit(limit).Find(&cabang)
		if result.Error != nil {
			exceptions.AppException(c, "Something went wrong")
			return
		}
	}

	data := make(map[string]interface{})
	data["cabang"] = cabang
	data["total"] = count

	payloads.HandleSuccess(c, data, "Data found", 200)
}

func UpdateCabang(c *gin.Context) {
	cabangID := c.Param("id")

	var payload models.Cabang
	if err := c.ShouldBindJSON(&payload); err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	var cabang models.Cabang
	result := config.InitDB().First(&cabang, cabangID)
	if result.Error != nil {
		payloads.HandleSuccess(c, "Leasing not found", "Success", 200)
		return
	}

	cabang.NamaCabang = payload.NamaCabang

	result = config.InitDB().Save(&cabang)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, cabang, "Success", 200)
}

func DeleteCabang(c *gin.Context) {

	cabangID := c.Param("id")

	var cabang models.Cabang
	result := config.InitDB().First(&cabang, cabangID)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	result = config.InitDB().Delete(&cabang)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, "Cabang deleted", "Success", 200)
}