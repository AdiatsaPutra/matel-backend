package main

import (
	config "motor/configs"
	"motor/models"
)

func main() {
	config.InitDB().AutoMigrate(&models.User{})
	config.InitDB().AutoMigrate(&models.Member{})
}
