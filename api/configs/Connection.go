package config

import (
	"log"
	"net/url"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	val := url.Values{}
	val.Add("parseTime", "True")
	val.Add("loc", "Asia/Jakarta")

	dsn := ""

	if os.Getenv("GIN_MODE") == "release" {
		dsn = "root:root@tcp(db)/matel?charset=utf8mb4&parseTime=True&loc=Local"
	} else {
		dsn = "root:1Ultramilk!@tcp(127.0.0.1:3306)/motor?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Cannot connected database ", err)
		return nil
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil
	}

	err = sqlDB.Ping()

	if err != nil {
		return nil
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(time.Minute * 3)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 3)

	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal("Cannot get database connection: ", err)
		return
	}

	sqlDB.Close()
}
