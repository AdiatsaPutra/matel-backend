package controllers

import (
	"encoding/csv"
	"fmt"
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

type LeasingResponse struct {
	total   int64
	leasing []models.Leasing
}

func GetLeasing(c *gin.Context) {
	// leasing, err := repository.GetLeasing(c)
	pageStr := c.DefaultQuery("page", "1")    // Get the page parameter from the query string
	limitStr := c.DefaultQuery("limit", "20") // Get the limit parameter from the query string

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

	var leasing []models.Leasing
	offset := (page - 1) * limit

	// Query the leasing table with pagination
	config.InitDB().Offset(offset).Limit(limit).Find(&leasing)
	var total int64
	if err := config.InitDB().Model(&models.Leasing{}).Count(&total).Error; err != nil {
		fmt.Println("Failed to retrieve total:", err)
		return
	}

	var data = LeasingResponse{
		total:   total,
		leasing: leasing,
	}

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
