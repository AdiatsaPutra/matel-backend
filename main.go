package main

import (
	"compress/gzip"
	"io"
	config "motor/configs"
	"motor/controllers"
	"motor/models"
	"motor/security"
	"net/http"
	"os"
	"path/filepath"

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
	r.GET("/export", exportHandler)

	r.GET("/member", security.AuthMiddleware(), controllers.GetAllMember)

	r.GET("/province", controllers.GetProvince)
	r.GET("/kabupaten/:province-id", controllers.GetKabupaten)
	r.GET("/kecamatan/:kabupaten-id", controllers.GetKecamatan)

	r.Run()
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
	filename := "exported.db"
	// c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	// c.Writer.Header().Set("Content-Type", "application/octet-stream")
	// c.File(filepath.Join(".", filename))
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(c.Writer)
	defer gzipWriter.Close()

	// Open the SQLite database file
	filePath := filepath.Join(".", filename)
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Copy the SQLite database file to the gzip writer
	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the content encoding to gzip
	c.Writer.Header().Set("Content-Encoding", "gzip")
}
