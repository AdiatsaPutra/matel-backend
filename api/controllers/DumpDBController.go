package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	config "matel/configs"
	"matel/models"
	"matel/payloads"
	"os"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func DumpSQLHandler(c *gin.Context) {

	filepath := "export.sql"

	// Koneksi ke database sumber
	sourceDB := config.InitDB()

	db, err := sourceDB.DB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer db.Close()

	// Membuka file untuk menulis hasil dump
	file, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to create file: %v": err.Error()})
	}
	defer file.Close()

	// Menulis header SQL ke file
	_, err = file.WriteString("INSERT INTO m_kendaraan (id, nomorPolisi, noMesin, noRangka) VALUES\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	var leasings []models.LeasingToExport
	err = sourceDB.Table("m_kendaraan").Select("id, nomorPolisi, noMesin, noRangka").Find(&leasings).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch data from table: %v": err.Error()})
	}

	// Menulis data ke file
	for i, l := range leasings {
		// log.Printf("Writing data %d of %d\n", i+1, len(leasings))
		_, err = file.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s')", l.ID, l.NomorPolisi, l.NoMesin, l.NoRangka))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to write to file: %v": err.Error()})
		}

		if i < len(leasings)-1 {
			_, err = file.WriteString(",\n")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"failed to write to file: %v": err.Error()})
			}
		} else {
			_, err = file.WriteString(";")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"failed to write to file: %v": err.Error()})
			}
		}
	}

	fileSource, err := os.Open(filepath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer fileSource.Close()

	// Get the file information.
	fileInfo, err := fileSource.Stat()
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
	_, err = io.Copy(writer, fileSource)
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

	payloads.HandleSuccess(c, "Berhasil update database", "Berhasil", 200)
}

func UpdateSQLHandler(c *gin.Context) {

	filepath := "export_new_only.sql"

	// Get the date parameter from the request
	dateParam := c.Query("date")

	// Parse the date parameter into a time.Time object
	date, err := time.Parse("2006-01-02-15-04-05", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Koneksi ke database sumber
	sourceDB := config.InitDB()

	db, err := sourceDB.DB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer db.Close()

	// Membuka file untuk menulis hasil dump
	file, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to create file: %v": err.Error()})
	}
	defer file.Close()

	// Menulis header SQL ke file
	_, err = file.WriteString("INSERT INTO m_kendaraan (id, nomorPolisi, noMesin, noRangka) VALUES\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	var leasings []models.LeasingToExport
	err = sourceDB.Table("m_kendaraan").
		Select("id, nomorPolisi, noMesin, noRangka").
		Where("created_at >= ?", date).
		Find(&leasings).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch data from table": err.Error()})
		return
	}

	logrus.Info(leasings)

	// Menulis data ke file
	for i, l := range leasings {
		// log.Printf("Writing data %d of %d\n", i+1, len(leasings))
		_, err = file.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s')", l.ID, l.NomorPolisi, l.NoMesin, l.NoRangka))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to write to file: %v": err.Error()})
		}

		if i < len(leasings)-1 {
			_, err = file.WriteString(",\n")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"failed to write to file: %v": err.Error()})
			}
		} else {
			_, err = file.WriteString(";")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"failed to write to file: %v": err.Error()})
			}
		}
	}

	fileSource, err := os.Open(filepath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer fileSource.Close()

	// Get the file information.
	fileInfo, err := fileSource.Stat()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create a zip archive.
	zipFilePath := "archive_new_only.zip"
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
	_, err = io.Copy(writer, fileSource)
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

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilePath)
	c.Writer.Header().Set("Content-Type", "application/zip")

	c.File(zipFilePath)
}
