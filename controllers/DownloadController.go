package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"sync"

	config "motor/configs"
	"motor/models"
	"net/http"
	"os"

	"motor/payloads"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ExportHandler(c *gin.Context) {
	// Retrieve all data from the table
	var data []models.Leasing
	err := config.InitDB().Select("nomorPolisi, noRangka, noMesin").Limit(10000).Find(&data).Error
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

	// // Clear existing data from the table
	// err = sqliteDB.Exec("DELETE FROM m_leasing").Error
	// if err != nil {
	// 	logrus.Error(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// AutoMigrate your model in the SQLite database
	err = sqliteDB.AutoMigrate(&models.LeasingToExport{})
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	batchSize := 10 // Set the desired batch size for insertion
	totalData := len(data)
	batchCount := totalData / batchSize

	// Convert data to the desired struct with selected fields
	var leasingData []models.LeasingToExport
	for _, d := range data {
		nomorPolisi := d.NomorPolisi
		if nomorPolisi == "" {
			nomorPolisi = " "
		}

		noRangka := d.NoRangka
		if noRangka == "" {
			noRangka = " "
		}

		noMesin := d.NoMesin
		if noMesin == "" {
			noMesin = " "
		}

		leasingData = append(leasingData, models.LeasingToExport{
			NomorPolisi: nomorPolisi,
			NoRangka:    noRangka,
			NoMesin:     noMesin,
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

	payloads.HandleSuccess(c, "Berhasil update database", "Berhasil", 200)

}

func DumpDBHandler(c *gin.Context) {
	sourceDB := config.InitDB()

	// Membuat channel untuk mengirim hasil dumping
	resultCh := make(chan error)

	// Membuat WaitGroup untuk menunggu selesainya semua worker
	var wg sync.WaitGroup

	// Jumlah worker yang akan digunakan
	workerCount := 5

	// Menginisialisasi worker
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(sourceDB, resultCh, &wg)
	}

	// Menunggu selesainya semua worker
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Mengumpulkan hasil dumping
	var dumpError error
	for result := range resultCh {
		if result != nil {
			dumpError = result
		}
	}

	if dumpError != nil {
		c.AbortWithError(http.StatusInternalServerError, dumpError)
		fmt.Println("Failed to dump database:", dumpError)
	} else {
		fmt.Println("Database dumped successfully")

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

		payloads.HandleSuccess(c, "Berhasil update database", "Berhasil", 200)
	}
}

func worker(sourceDB *gorm.DB, resultCh chan<- error, wg *sync.WaitGroup) {
	// Koneksi ke database tujuan
	destinationDB, err := gorm.Open(sqlite.Open("exported.db"), &gorm.Config{})
	if err != nil {
		resultCh <- err
		wg.Done()
		return
	}

	rawSourceDB, err := sourceDB.DB()
	if err != nil {
		resultCh <- err
		wg.Done()
		return
	}

	rawDestinationDB, err := destinationDB.DB()
	if err != nil {
		resultCh <- err
		wg.Done()
		return
	}

	defer rawSourceDB.Close()
	defer rawDestinationDB.Close()

	// Migrate skema di database tujuan berdasarkan struktur skema dari database sumber
	err = destinationDB.AutoMigrate(&models.LeasingToExport{})
	if err != nil {
		resultCh <- err
		wg.Done()
		return
	}

	// Salin data dari tabel sumber ke tabel tujuan
	var mLeasings []models.Leasing
	sourceDB.Select("nomorPolisi, noMesin, noRangka").Find(&mLeasings)
	for _, mLeasing := range mLeasings {
		exportedMLeasing := models.LeasingToExport{
			NomorPolisi: mLeasing.NomorPolisi,
			NoMesin:     mLeasing.NoMesin,
			NoRangka:    mLeasing.NoRangka,
		}
		destinationDB.Create(&exportedMLeasing)
	}

	wg.Done()
	resultCh <- nil
}

func ExportHandlerNew(c *gin.Context) {

	sourceDB := config.InitDB()

	destinationDB, err := gorm.Open(sqlite.Open("exported.db"), &gorm.Config{})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Copy data to destination database
	batchSize := 100000
	offset := 0
	for {
		var sourceData []models.Leasing
		result := sourceDB.Offset(offset).Limit(batchSize).Find(&sourceData)
		if result.Error != nil {
			fmt.Println("Failed to retrieve users:", result.Error)
			break
		}

		// Copy data to destination database
		for _, sourceUser := range sourceData {
			destinationUser := models.LeasingToExport{NomorPolisi: sourceUser.NomorPolisi, NoMesin: sourceUser.NoMesin, NoRangka: sourceUser.NoRangka}
			result := destinationDB.Create(&destinationUser)
			if result.Error != nil {
				fmt.Println("Failed to copy user:", result.Error)
				break
			}
		}

		// Check if all data has been processed
		if len(sourceData) < batchSize {
			break
		}

		offset += batchSize
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

	payloads.HandleSuccess(c, "Berhasil update database", "Berhasil", 200)
}

func DownloadLeasing(c *gin.Context) {
	zipFilePath := "archive.zip"

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilePath)
	c.Writer.Header().Set("Content-Type", "application/zip")

	c.File(zipFilePath)

	// Remove file after sending it to the user.
	// os.Remove(zipFilePath)

}
