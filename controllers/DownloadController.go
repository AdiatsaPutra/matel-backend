package controllers

import (
	"archive/zip"
	"io"
	config "motor/configs"
	"motor/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ExportHandler(c *gin.Context) {
	// Retrieve all data from the table
	var data []models.Leasing
	err := config.InitDB().Select("nomorPolisi, noRangka, noMesin").Find(&data).Error
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No data found"})
		return
	}

	// Create a new SQLite database file
	sqliteDB, err := gorm.Open(sqlite.Open("exported.db"), &gorm.Config{})
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// defer sqliteDB.Close()

	// Clear existing data from the table
	err = sqliteDB.Exec("DELETE FROM m_leasing").Error
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// AutoMigrate your model in the SQLite database
	err = sqliteDB.AutoMigrate(&models.LeasingToExport{})
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	batchSize := 100 // Set the desired batch size for insertion
	totalData := len(data)
	batchCount := totalData / batchSize

	// Convert data to the desired struct with selected fields
	var leasingData []models.LeasingToExport
	for _, d := range data {
		leasingData = append(leasingData, models.LeasingToExport{
			NomorPolisi: d.NomorPolisi,
			NoRangka:    d.NoRangka,
			NoMesin:     d.NoMesin,
		})
	}

	// Insert data into the SQLite database in batches
	for i := 0; i <= batchCount; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > totalData {
			end = totalData
		}

		batch := leasingData[start:end]
		err = sqliteDB.Create(&batch).Error
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Set the response headers for file download
	filepath := "exported.db"
	file, err := os.Open(filepath)
	if err != nil {
		logrus.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	// Get the file information.
	fileInfo, err := file.Stat()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create a zip archive.
	zipFilePath := "archive.zip"
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create ZIP file")
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Create a new file in the zip archive.
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	header.Name = fileInfo.Name()

	writer, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   filepath,
		Method: zip.Deflate, // Menggunakan metode kompresi Deflate
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Copy the file content to the zip archive.
	_, err = io.Copy(writer, file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Close the zip archive.
	err = zipWriter.Close()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set the response headers.
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilePath)
	c.Writer.Header().Set("Content-Type", "application/zip")

	// Send the zip file to the user.
	c.File(zipFilePath)

	// Remove file after sending it to the user.
	os.Remove(zipFilePath)

}
