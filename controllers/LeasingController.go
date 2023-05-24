package controllers

import (
	"encoding/csv"
	"io"
	config "motor/configs"
	"motor/exceptions"
	"motor/models"
	"motor/payloads"
	"net/http"
	"strconv"
	"strings"

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
	query = query.Offset(offset).Limit(limit)
	query.Find(&leasing)

	data := make(map[string]interface{})
	data["total"] = total
	data["leasing"] = leasing

	payloads.HandleSuccess(c, data, "Leasing found", http.StatusOK)
}

func UploadLeasing(c *gin.Context) {
	// Retrieve the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		exceptions.BadRequest(c, "Masukkan data valid")
		return
	}

	// Open the uploaded file
	csvFile, err := file.Open()
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}
	defer csvFile.Close()

	// Create a CSV reader
	reader := csv.NewReader(csvFile)

	// Skip the header row
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	// Process each row of the CSV file
	for {
		// Read each row
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			exceptions.AppException(c, "Something went wrong")
			return
		}

		replacer := strings.NewReplacer(",", "")

		sisaHutang, err := strconv.Atoi(replacer.Replace(string(row[5])))
		if err != nil {
			exceptions.AppException(c, "Something went wrong")
			return
		}

		status, err := strconv.Atoi(row[11])
		if err != nil {
			exceptions.AppException(c, "Something went wrong")
			return
		}

		leasing := models.Leasing{
			Leasing:     row[0],
			Cabang:      row[1],
			NoKontrak:   row[2],
			NamaDebitur: row[3],
			NomorPolisi: row[4],
			SisaHutang:  uint(sisaHutang),
			Tipe:        row[6],
			Tahun:       row[7],
			NoRangka:    row[8],
			NoMesin:     row[9],
			PIC:         row[10],
			Status:      uint(status),
		}

		err = config.InitDB().Create(&leasing).Error
		if err != nil {
			exceptions.AppException(c, "Something went wrong")
			return
		}
	}

	payloads.HandleSuccess(c, "Berhasil mengupload", "Berhasil", 200)
}
