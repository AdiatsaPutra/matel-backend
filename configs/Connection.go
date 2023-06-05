package config

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	val := url.Values{}
	val.Add("parseTime", "True")
	val.Add("loc", "Asia/Jakarta")

	dsn := "root:1Ultramilk!@tcp(127.0.0.1:3306)/motor?charset=utf8mb4&parseTime=True&loc=Local"
	// dsn := "root:root@tcp(167.172.69.241:3306)/matel?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Cannot connected database ", err)
		return nil
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal("Cannot connected database ", err)
		return nil
	}

	err = sqlDB.Ping()

	if err != nil {
		log.Fatal("Request Timeout ", err)
		return nil
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(time.Minute * 3)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 3)

	log.Info("Connected Database")

	return db
}
