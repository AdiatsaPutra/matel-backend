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

	// Membuka file untuk menulis hasil dump
	file, err := os.Create("dump.sql")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()

	// Menulis header SQL ke file
	_, err = file.WriteString("INSERT INTO m_leasing (nomorPolisi, noMesin, noRangka) VALUES\n")
	if err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}

	var leasings []models.LeasingToExport
	err = sourceDB.Table("m_leasing").Select("nomorPolisi, noMesin, noRangka").Find(&leasings).Error
	if err != nil {
		log.Fatalf("failed to fetch data from table: %v", err)
	}

	// Menulis data ke file
	for i, l := range leasings {
		log.Printf("Writing data %d of %d\n", i+1, len(leasings))
		_, err = file.WriteString(fmt.Sprintf("('%s', '%s', '%s')", l.NomorPolisi, l.NoMesin, l.NoRangka))
		if err != nil {
			log.Fatalf("failed to write to file: %v", err)
		}

		if i < len(leasings)-1 {
			_, err = file.WriteString(",\n")
			if err != nil {
				log.Fatalf("failed to write to file: %v", err)
			}
		} else {
			_, err = file.WriteString(";")
			if err != nil {
				log.Fatalf("failed to write to file: %v", err)
			}
		}
	}

	fmt.Println("Table dumped successfully!")
}
