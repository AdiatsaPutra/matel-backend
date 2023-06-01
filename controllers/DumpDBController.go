package controllers

import (
	// "database/sql"
	"fmt"
	"log"
	"os"

	config "motor/configs"
	"motor/models"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/go-sql-driver/mysql"
	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"
)

func DumpSQLHandler(c *gin.Context) {
	// Koneksi ke database sumber
	sourceDB := config.InitDB()

	db, err := sourceDB.DB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer db.Close()

	// Membuat file untuk menyimpan hasil dump
	file, err := os.Create("exported.sql")
	if err != nil {
		log.Fatal("Failed to create exported SQL file:", err)
	}
	defer file.Close()

	// Menginisialisasi string untuk menyimpan hasil dump
	var dumpSQL string

	// Mengambil hasil dump tabel menggunakan GORM
	err = sourceDB.Model(&models.Leasing{}).Select("nomorPolisi, noMesin, noRangka").Find(&[]models.LeasingToExport{}).Error
	if err != nil {
		log.Fatal("Failed to get dump table statement:", err)
	}

	// Menyimpan hasil dump ke dalam file
	_, err = file.WriteString(dumpSQL)
	if err != nil {
		log.Fatal("Failed to write dump to file:", err)
	}

	fmt.Println("Database dumped successfully")
}
