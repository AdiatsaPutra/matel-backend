package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	config "matel/configs"
	"matel/models"
	"matel/payloads"
	"os"
	"strconv"

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

	// Get Cabang With Version
	var cabang []models.Cabang
	err = sourceDB.Table("m_cabang").
		Select("id, versi").
		Find(&cabang).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch data from table": err.Error()})
		return
	}

	// Insert cabangForm into m_cabang

	_, err = file.WriteString("INSERT INTO m_cabang (id_source, versi) VALUES\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	for i, cb := range cabang {
		idInt := strconv.Itoa(int(cb.ID))
		versi := strconv.Itoa(cb.Versi)
		_, err = file.WriteString(fmt.Sprintf("('%s', '%s')", idInt, versi))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to write to file: %v": err.Error()})
		}

		if i < len(cabang)-1 {
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

	_, err = file.WriteString("\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write delete query to file": err.Error()})
		return
	}

	// Menulis header SQL ke file
	_, err = file.WriteString("INSERT INTO m_kendaraan (id, cabang, nomorPolisi, noMesin, noRangka) VALUES\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	var leasings []models.LeasingToExport
	err = sourceDB.Table("m_kendaraan").Select("id, cabang, nomorPolisi, noMesin, noRangka").Find(&leasings).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch data from table: %v": err.Error()})
	}

	// Menulis data ke file
	for i, l := range leasings {
		// log.Printf("Writing data %d of %d\n", i+1, len(leasings))
		_, err = file.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')", l.ID, l.Cabang, l.NomorPolisi, l.NoMesin, l.NoRangka))
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

	// dateParam := c.Query("date")

	type CabangForm struct {
		ID    string `json:"id_source"`
		Versi int    `json:"versi"`
	}

	var cabangForm []CabangForm
	var cabangFormUnupdated []CabangForm

	if err := c.BindJSON(&cabangForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cabangFormUnupdated = append(cabangFormUnupdated, cabangForm...)

	// date, err := time.Parse("2006-01-02-15-04-05", dateParam)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
	// 	return
	// }

	sourceDB := config.InitDB()

	db, err := sourceDB.DB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer db.Close()

	file, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to create file: %v": err.Error()})
	}
	defer file.Close()

	// Get Cabang With Version
	var cabang []models.Cabang
	err = sourceDB.Table("m_cabang").
		Select("id, nama_cabang, versi").
		Find(&cabang).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch data from table": err.Error()})
		return
	}

	existingCabangMap := make(map[string]int)

	for _, cf := range cabangForm {
		existingCabangMap[cf.ID] = cf.Versi
	}

	// for _, cb := range cabang {
	// 	idStr := strconv.Itoa(int(cb.ID))
	// 	if _, ok := existingCabangMap[cb.NamaCabang]; !ok {
	// 		cabangForm = append(cabangForm, CabangForm{ID: idStr, Versi: cb.Versi})
	// 	} else {
	// 		existingCabangMap[cb.NamaCabang] = cb.Versi
	// 	}
	// }

	for i := range cabangForm {
		versi := existingCabangMap[cabangForm[i].ID]
		if versi == 0 {
			cabangForm[i].Versi = 1
		} else {
			cabangForm[i].Versi = versi
		}
	}

	_, err = file.WriteString("DELETE FROM m_cabang;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	_, err = file.WriteString("\n\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write delete query to file": err.Error()})
		return
	}

	// Insert cabangForm into m_cabang

	var comparedCabangForm []CabangForm
	for _, cf := range cabangForm {
		found := false
		for _, cfu := range cabangFormUnupdated {
			if cf.ID == cfu.ID {
				found = true
				break
			}
		}

		if !found {
			comparedCabangForm = append(comparedCabangForm, cf)
		}
	}

	_, err = file.WriteString("INSERT INTO m_cabang (id_source, versi) VALUES\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	for i, cb := range cabangForm {
		versi := strconv.Itoa(cb.Versi)
		_, err = file.WriteString(fmt.Sprintf("('%s', '%s')", cb.ID, versi))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to write to file: %v": err.Error()})
		}

		if i < len(cabangForm)-1 {
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

	_, err = file.WriteString("\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write delete query to file": err.Error()})
		return
	}

	for _, cc := range comparedCabangForm {
		_, err = file.WriteString(fmt.Sprintf("DELETE FROM m_kendaraan WHERE cabang = '%s';\n", cc.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to write delete query to file": err.Error()})
			return
		}
	}

	_, err = file.WriteString("\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	logrus.Info(comparedCabangForm)
	logrus.Info(cabangForm)
	logrus.Info(cabangFormUnupdated)

	for _, cc := range comparedCabangForm {
		var leasings []models.LeasingToExport
		err = sourceDB.Table("m_kendaraan").
			Select("id, cabang, nomorPolisi, noMesin, noRangka").
			Where("cabang_id = ?", cc.ID).
			Where("versi < ?", cc.Versi).
			Find(&leasings).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch data from table": err.Error()})
			return
		}

		_, err = file.WriteString("INSERT INTO m_kendaraan (id, cabang, nomorPolisi, noMesin, noRangka) VALUES\n")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
		}

		for i, l := range leasings {
			_, err = file.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')", l.ID, l.Cabang, l.NomorPolisi, l.NoMesin, l.NoRangka))
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

		_, err = file.WriteString("\n")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
		}
	}

	fileSource, err := os.Open(filepath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer fileSource.Close()

	fileInfo, err := fileSource.Stat()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	zipFilePath := "archive_new_only.zip"
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create ZIP file")
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	header.Name = fileInfo.Name()

	writer, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   filepath,
		Method: zip.Deflate,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, err = io.Copy(writer, fileSource)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = zipWriter.Close()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilePath)
	c.Writer.Header().Set("Content-Type", "application/zip")

	c.File(zipFilePath)
}
