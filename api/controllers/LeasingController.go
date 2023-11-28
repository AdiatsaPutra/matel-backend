package controllers

import (
	"io"
	config "matel/configs"
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"matel/repository"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetKendaraan(c *gin.Context) {

	pageNumber, _ := strconv.Atoi(c.Query("page"))
	if pageNumber <= 0 {
		pageNumber = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit <= 0 && limit != -1 {
		limit = 10
	}

	query := config.InitDB().Model(&models.Kendaraan{})

	if leasing := c.Query("leasing"); leasing != "" {
		query = query.Where("leasing LIKE ?", ""+leasing+"")
	}

	if cabang := c.Query("cabang"); cabang != "" {
		query = query.Where("cabang LIKE ?", ""+cabang+"")
	}

	query = query.Order("created_at DESC")

	var kendaraans []models.Kendaraan
	offset := (pageNumber - 1) * limit

	var count int64
	result := query.Count(&count)
	if result.Error != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	if limit == -1 {
		result := query.Find(&kendaraans)
		if result.Error != nil {
			exceptions.AppException(c, "Something went wrong")
			return
		}
	} else {
		result := query.Offset(offset).Limit(limit).Find(&kendaraans)
		if result.Error != nil {
			exceptions.AppException(c, "Something went wrong")
			return
		}

	}

	if search := c.Query("search"); search != "" {
		query = query.Where("nomorPolisi LIKE ?", ""+search+"")
	}

	data := make(map[string]interface{})
	data["total"] = count
	data["kendaraan"] = kendaraans

	config.CloseDB(config.InitDB())

	payloads.HandleSuccess(c, data, "Kendaraan found", http.StatusOK)
}

func DeleteKendaraan(c *gin.Context) {
	leasing := c.Query("leasing")
	leasingID := c.Query("leasing_id")
	cabang := c.Query("cabang")

	leasingIDInt, _ := strconv.Atoi(leasingID)
	SetVersiCabang(c, uint(leasingIDInt), cabang, true)
	// SetVersiCabang(c, uint(leasingIDInt), cabang, true)

	var kendaraan models.Kendaraan

	deleteResult := config.InitDB().Where("leasing = ? AND cabang = ?", leasing, cabang).Delete(&kendaraan)
	if deleteResult.Error != nil {
		exceptions.AppException(c, "Failed to delete Kendaraan")
		return
	}
	payloads.HandleSuccess(c, "Success", "Success", http.StatusOK)
}

func DeleteKendaraanByID(c *gin.Context) {
	kendaraanID := c.Param("id")

	kendaraanIDInt, _ := strconv.Atoi(kendaraanID)
	// SetVersiCabang(c, uint(leasingIDInt), cabang, true)

	var kendaraan models.Kendaraan

	deleteResult := config.InitDB().Where("id = ?", kendaraanIDInt).Delete(&kendaraan)
	if deleteResult.Error != nil {
		exceptions.AppException(c, deleteResult.Error.Error())
		return
	}

	payloads.HandleSuccess(c, "Success", "Success", http.StatusOK)
}

func DeleteAllKendaraan(c *gin.Context) {
	deleteResult := config.InitDB().Exec("DELETE from m_kendaraan")
	if deleteResult.Error != nil {
		exceptions.AppException(c, "Failed to delete Kendaraan")
		return
	}

	payloads.HandleSuccess(c, "Success", "Success", http.StatusOK)
}

func DownloadTemplate(c *gin.Context) {
	filePath := "leasing-template.csv"
	file, err := os.Open(filePath)
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		exceptions.AppException(c, "Failed to download file")
		return
	}
}

func DownloadTemplateCabang(c *gin.Context) {
	filePath := "leasing-template-cabang.csv"
	file, err := os.Open(filePath)
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		exceptions.AppException(c, "Failed to download file")
		return
	}
}

func GetLeasing(c *gin.Context) {
	searchQuery := c.Query("search") // Get the search query from the query string

	var kendaraan []models.Kendaraan

	query := config.InitDB()

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
