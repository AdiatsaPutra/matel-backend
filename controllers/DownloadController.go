package controllers

import (
	"compress/gzip"
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
	err = sqliteDB.AutoMigrate(&models.Leasing{})
	if err != nil {
		logrus.Error(err)
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

	// Create a new gzipped file
	gzippedFile, err := os.Create("archive.gz")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer gzippedFile.Close()

	// Create a new gzip writer
	gzipWriter := gzip.NewWriter(gzippedFile)
	defer gzipWriter.Close()

	// Copy the contents of the file to the gzip writer
	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Flush the gzip writer to ensure all data is written
	gzipWriter.Flush()

	// open archived file
	archivedFile, err := os.Open("archive.gz")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer archivedFile.Close()

	// set filename and extension
	fileName := "archive.gz"
	contentType := "application/octet-stream"

	// set header response
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.Writer.Header().Set("Content-Type", contentType)

	c.File(fileName)

	// Send the compressed zip file to the user.
	// c.Data(http.StatusOK, "application/gzip", buf.Bytes())
}
