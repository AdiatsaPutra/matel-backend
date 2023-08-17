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
		LeasingID:   payload.LeasingID,
		NamaCabang:  payload.NamaCabang,
		NoHP:        payload.NoHP,
		Versi:       2,
		VersiMaster: 1,
	}

	result := config.InitDB().Create(&cabang)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}
	config.CloseDB(config.InitDB())

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

func GetCabangWithTotal(c *gin.Context) {
	leasingID := c.Query("leasing_id")

	db := config.InitDB()

	var results []models.CabangTotal

	err := db.Raw(`SELECT
			c.id,
			c.nama_cabang,
			c.no_hp,
			l.nama_leasing,
			COALESCE(k.kendaraan_total, 0) AS kendaraan_total,
			COALESCE(k.latest_created_at, NULL) AS latest_created_at
		FROM
			m_cabang c
		LEFT JOIN
			m_leasing l ON c.leasing_id = l.id
		LEFT JOIN (
			SELECT
				cabang,
				cabang_id,
				COUNT(nomorPolisi) AS kendaraan_total,
				MAX(created_at) AS latest_created_at
			FROM
				m_kendaraan
			WHERE
				deleted_at IS NULL
			GROUP BY
				cabang, cabang_id
		) k ON c.id = k.cabang_id
		WHERE
			c.leasing_id = ? AND c.deleted_at IS NULL;
		`, leasingID).Scan(&results).Error

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}
	config.CloseDB(config.InitDB())

	payloads.HandleSuccess(c, results, "Data found", 200)
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

	// Store the original NamaCabang value for comparison
	originalNamaCabang := cabang.NamaCabang

	// Update the cabang record
	cabang.NamaCabang = payload.NamaCabang
	cabang.NoHP = payload.NoHP
	result = config.InitDB().Save(&cabang)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	// Update the m_kendaraan table
	kendaraanResult := config.InitDB().Model(&models.Kendaraan{}).Where("cabang = ?", originalNamaCabang).Update("cabang", cabang.NamaCabang)
	if kendaraanResult.Error != nil {
		exceptions.AppException(c, "Something went wrong while updating m_kendaraan")
		return
	}

	config.CloseDB(config.InitDB())

	payloads.HandleSuccess(c, cabang, "Success", 200)
}

func SetVersiCabang(c *gin.Context, LeasingID uint, CabangName string, Reset bool) {
	var cabang models.Cabang
	result := config.InitDB().Where("leasing_id = ? AND nama_cabang = ? AND deleted_at IS NULL", LeasingID, CabangName).Find(&cabang)
	if result.Error != nil {
		payloads.HandleSuccess(c, "Leasing not found", "Success", 200)
		return
	}

	cabang.VersiMaster = cabang.VersiMaster + 1

	if Reset {
		cabang.Versi = 1
	} else {
		// cabang.Versi = cabang.Versi + 1
	}

	result = config.InitDB().Save(&cabang)
	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return
	}

	// payloads.HandleSuccess(c, cabang, "Success", 200)
}

func GetCabangVersi(CabangName string) int {

	var cabang models.Cabang
	result := config.InitDB().Where("nama_cabang = ? AND deleted_at IS NULL", CabangName).Find(&cabang)
	if result.Error != nil {
		return 0
	}

	return cabang.Versi
}

func DeleteCabang(c *gin.Context) {
	cabangID := c.Param("id")
	queryParam := c.DefaultQuery("query", "")

	var cabang models.Cabang
	result := config.InitDB().First(&cabang, cabangID)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	logrus.Info("Log")
	logrus.Info(queryParam)

	if queryParam == "delete-cabang" {
		logrus.Info(queryParam)
		result = config.InitDB().Delete(&cabang)
		if result.Error != nil {
			exceptions.AppException(c, "Something went wrong")
			return
		}

	}

	cabang.VersiMaster = cabang.VersiMaster + 1
	result = config.InitDB().Save(&cabang)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	var kendaraan models.Kendaraan
	deleteKendaraanResult := config.InitDB().Where("cabang = ?", cabang.NamaCabang).Delete(&kendaraan)
	if deleteKendaraanResult.Error != nil {
		exceptions.AppException(c, "Failed to delete Kendaraan")
		return
	}

	config.CloseDB(config.InitDB())

	payloads.HandleSuccess(c, "Cabang deleted", "Success", 200)
}
