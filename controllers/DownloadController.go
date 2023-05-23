package controllers

import (
	"archive/zip"
	"bytes"
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
	err := config.InitDB().Find(&data).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Info(data)

	// Create a new SQLite database file
	sqliteDB, err := gorm.Open(sqlite.Open("exported.db"), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clear existing data from the table
	err = sqliteDB.Exec("DELETE FROM m_leasing").Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// AutoMigrate your model in the SQLite database
	err = sqliteDB.AutoMigrate(&models.Leasing{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	batchSize := 100 // Set the desired batch size for insertion
	totalData := len(data)
	batchCount := totalData / batchSize

	// Insert data into the SQLite database in batches
	for i := 0; i <= batchCount; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > totalData {
			end = totalData
		}

		batch := data[start:end]
		err = sqliteDB.Create(&batch).Error
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Set the response headers for file download
	filepath := "exported.db"
	file, err := os.Open(filepath)
	if err != nil {
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
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create a new file in the zip archive.
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	header.Name = fileInfo.Name()

	writer, err := zipWriter.CreateHeader(header)
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
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=archive.zip")

	// Send the zip file to the user.
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}
