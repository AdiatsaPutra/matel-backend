package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	config "motor/configs"
	"motor/models"
	"motor/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
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
	_, err = file.WriteString("INSERT INTO m_leasing (nomorPolisi, noMesin, noRangka) VALUES\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	var leasings []models.LeasingToExport
	err = sourceDB.Table("m_leasing").Select("nomorPolisi, noMesin, noRangka").Find(&leasings).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch data from table: %v": err.Error()})
	}

	// Menulis data ke file
	for i, l := range leasings {
		// log.Printf("Writing data %d of %d\n", i+1, len(leasings))
		_, err = file.WriteString(fmt.Sprintf("('%s', '%s', '%s')", l.NomorPolisi, l.NoMesin, l.NoRangka))
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

	payloads.HandleSuccess(c, "Berhasil update database", "Berhasil", 200)
}
