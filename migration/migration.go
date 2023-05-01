package main

import (
	"motor/configs"
	"motor/models"
)

func init() {
	configs.ConnectDatabase()
}

func main() {
	configs.DB.AutoMigrate(&models.User{})
	configs.DB.AutoMigrate(&models.Member{})
}
