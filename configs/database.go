package configs

import (
	"log"
	// "os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	dsn := "root:1Ultramilk!@tcp(127.0.0.1:3306)/motor?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Cant connect to database")
	}

	log.Println("Success to connect")
}
