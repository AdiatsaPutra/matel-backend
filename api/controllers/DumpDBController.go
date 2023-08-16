package controllers

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"log"
	config "matel/configs"
	"matel/models"
	"matel/payloads"
	"os"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

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
		Select("id, versi, versi_master").
		Find(&cabang).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch data from table": err.Error()})
		return
	}

	// Insert cabangForm into m_cabang

	_, err = file.WriteString("INSERT INTO m_cabang (id_source, versi, versi_master) VALUES\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	for i, cb := range cabang {
		idInt := strconv.Itoa(int(cb.ID))
		versi := strconv.Itoa(cb.Versi)
		versiMaster := strconv.Itoa(cb.VersiMaster)
		_, err = file.WriteString(fmt.Sprintf("('%s', '%s', '%s')", idInt, versi, versiMaster))
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
	_, err = file.WriteString("INSERT INTO m_kendaraan (id_source, cabang_id, nomorPolisi, noMesin, noRangka) VALUES\n")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to write header to file: %v": err.Error()})
	}

	var leasings []models.LeasingToExport
	err = sourceDB.Table("m_kendaraan").Select("id, cabang_id, nomorPolisi, noMesin, noRangka").Where("deleted_at IS NULL").Find(&leasings).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to fetch data from table: %v": err.Error()})
	}

	// Menulis data ke file
	for i, l := range leasings {
		// log.Printf("Writing data %d of %d\n", i+1, len(leasings))
		_, err = file.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')", l.ID, l.CabangID, l.NomorPolisi, l.NoMesin, l.NoRangka))
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

type Cabang struct {
	ID          int `json:"id"`
	Versi       int `json:"versi"`
	VersiMaster int `json:"versi_master"`
}

type MKendaraan struct {
	ID          int    `json:"id"`
	CabangID    int    `json:"cabang_id"`
	NomorPolisi string `json:"nomorPolisi"`
	NoRangka    string `json:"noRangka"`
	NoMesin     string `json:"noMesin"`
}

type Item struct {
	IDSource    int `json:"id_source"`
	Versi       int `json:"versi"`
	VersiMaster int `json:"versi_master"`
}

func compareData(apiData []Item, dbData []Cabang) []map[string]interface{} {
	results := []map[string]interface{}{}

	for _, dbItem := range dbData {
		dbID, dbVersi, dbVersiMaster := dbItem.ID, dbItem.Versi, dbItem.VersiMaster
		found := false

		for _, apiItem := range apiData {
			if apiItem.IDSource == dbID {
				if apiItem.Versi != dbVersi || apiItem.VersiMaster != dbVersiMaster {
					var status string
					if dbVersiMaster > apiItem.VersiMaster {
						status = "Perbedaan versi master"
					} else if dbVersi > apiItem.Versi {
						status = "Perbedaan versi"
					}

					compareResult := map[string]interface{}{
						"id_source":    dbID,
						"versi":        dbVersi,
						"versi_master": dbVersiMaster,
						"status":       status,
					}
					results = append(results, compareResult)
				}
				found = true
				break
			}
		}

		if !found {
			result := map[string]interface{}{
				"id_source":    dbID,
				"versi":        dbVersi,
				"versi_master": dbVersiMaster,
				"status":       "Cabang tidak ada dalam request API",
			}
			results = append(results, result)
		}
	}

	return results
}

func getMKendaraanByCabang(cabangID int) ([]MKendaraan, error) {
	query := fmt.Sprintf("SELECT id, cabang_id, nomorPolisi, noRangka, noMesin FROM m_kendaraan WHERE cabang_id = %d", cabangID)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MKendaraan
	for rows.Next() {
		var kendaraan MKendaraan
		err := rows.Scan(&kendaraan.ID, &kendaraan.CabangID, &kendaraan.NomorPolisi, &kendaraan.NoRangka, &kendaraan.NoMesin)
		if err != nil {
			return nil, err
		}
		results = append(results, kendaraan)
	}

	return results, nil
}

func getMKendaraanByCabangVersi(cabangID int, versi int) ([]MKendaraan, error) {
	query := fmt.Sprintf("SELECT id, cabang_id, nomorPolisi, noRangka, noMesin FROM m_kendaraan WHERE cabang_id = %d AND versi > %d", cabangID, versi)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MKendaraan
	for rows.Next() {
		var kendaraan MKendaraan
		err := rows.Scan(&kendaraan.ID, &kendaraan.CabangID, &kendaraan.NomorPolisi, &kendaraan.NoRangka, &kendaraan.NoMesin)
		if err != nil {
			return nil, err
		}
		results = append(results, kendaraan)
	}

	return results, nil
}

func createSQLFile(compareResults []map[string]interface{}, mKendaraanData []MKendaraan, dbData []Cabang) {
	sqlStatements := []string{}

	sqlStatements = append(sqlStatements, "DELETE FROM m_cabang;\n")
	sqlStatements = append(sqlStatements, "INSERT INTO m_cabang (id_source, versi, versi_master) VALUES")
	for idx, dbItem := range dbData {
		statement := fmt.Sprintf("(%d, %d, %d)", dbItem.ID, dbItem.Versi, dbItem.VersiMaster)
		if idx == len(dbData)-1 {
			statement += ";\n"
		} else {
			statement += ","
		}
		sqlStatements = append(sqlStatements, statement)
	}

	for _, result := range compareResults {
		if status, ok := result["status"].(string); ok && status == "Perbedaan versi master" {
			sqlStatements = append(sqlStatements, fmt.Sprintf("DELETE FROM m_kendaraan WHERE cabang_id = %d;\n", result["id_source"].(int)))
			break
		}
	}

	sqlStatements = append(sqlStatements, "INSERT INTO m_kendaraan (id_source, cabang_id, nomor_polisi, no_rangka, no_mesin) VALUES")
	for idx, kendaraan := range mKendaraanData {
		statement := fmt.Sprintf("(%d, %d, '%s', '%s', '%s')", kendaraan.ID, kendaraan.CabangID, kendaraan.NomorPolisi, kendaraan.NoRangka, kendaraan.NoMesin)
		if idx == len(mKendaraanData)-1 {
			statement += ";\n"
		} else {
			statement += ","
		}
		sqlStatements = append(sqlStatements, statement)
	}

	sqlFile, err := os.Create("output.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlFile.Close()

	for _, statement := range sqlStatements {
		sqlFile.WriteString(statement + "\n")
	}
}

func UpdateSQLHandler(c *gin.Context) {
	var items []Item
	if err := c.ShouldBindJSON(&items); err != nil {
		logrus.Info(items)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := db.Query("SELECT id, versi, versi_master FROM m_cabang")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var dbData []Cabang
	for rows.Next() {
		var cabang Cabang
		err := rows.Scan(&cabang.ID, &cabang.Versi, &cabang.VersiMaster)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		dbData = append(dbData, cabang)
	}

	compareResults := compareData(items, dbData)
	mKendaraanData := []MKendaraan{}

	for _, result := range compareResults {
		switch status := result["status"].(string); {
		case status == "Cabang tidak ada dalam request API":
			kendaraanData, err := getMKendaraanByCabang(result["id_source"].(int))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			mKendaraanData = append(mKendaraanData, kendaraanData...)
		case status == "Perbedaan versi master":
			kendaraanData, err := getMKendaraanByCabang(result["id_source"].(int))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			mKendaraanData = append(mKendaraanData, kendaraanData...)
		case status == "Perbedaan versi":
			versi := result["versi"].(int)
			for _, item := range items {
				if item.IDSource == result["id_source"].(int) {
					versi = item.Versi
					break
				}
			}
			kendaraanData, err := getMKendaraanByCabangVersi(result["id_source"].(int), versi)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			mKendaraanData = append(mKendaraanData, kendaraanData...)
		}
	}

	createSQLFile(compareResults, mKendaraanData, dbData)

	c.JSON(http.StatusOK, gin.H{
		"compare_results":  compareResults,
		"m_kendaraan_data": mKendaraanData,
	})

}
