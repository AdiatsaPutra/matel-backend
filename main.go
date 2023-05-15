package main

import (
	"archive/zip"
	"bytes"
	"io"
	config "motor/configs"
	"motor/controllers"
	"motor/models"
	"motor/security"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/profil", security.AuthMiddleware(), controllers.GetProfile)

	r.GET("/leasing", controllers.GetLeasing)
	r.POST("/upload-leasing", controllers.UploadLeasing)
	r.GET("/export", handleZipDownload)

	r.GET("/member", security.AuthMiddleware(), controllers.GetAllMember)

	r.GET("/province", controllers.GetProvince)
	r.GET("/kabupaten/:province-id", controllers.GetKabupaten)
	r.GET("/kecamatan/:kabupaten-id", controllers.GetKecamatan)

	r.Run()
}

func handleZipDownload(c *gin.Context) {
	// Open the file to be zipped.
	filePath := "exported.db"
	file, err := os.Open(filePath)
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

func exportHandler(c *gin.Context) {

	// Retrieve all data from the table
	var data []models.Leasing
	err := config.InitDB().Find(&data).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new SQLite database file
	sqliteDB, err := gorm.Open(sqlite.Open("exported.db"), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// defer sqliteDB.Close()

	// AutoMigrate your model in the SQLite database
	err = sqliteDB.AutoMigrate(&models.Leasing{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert data into the SQLite database
	sqliteDB.Create(&data)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

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
