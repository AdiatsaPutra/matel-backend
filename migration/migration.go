package main

import (
	config "motor/configs"
	"motor/models"
)

func main() {
	config.InitDB().AutoMigrate(&models.User{})
	config.InitDB().AutoMigrate(&models.Member{})
	config.InitDB().AutoMigrate(&models.Kecamatan{})
	config.InitDB().AutoMigrate(&models.Kabupaten{})
	config.InitDB().AutoMigrate(&models.Province{})
}
