package main

import (
	config "matel/configs"
	"matel/models"
)

func main() {
	config.InitDB().AutoMigrate(&models.User{})
	config.InitDB().AutoMigrate(&models.Kecamatan{})
	config.InitDB().AutoMigrate(&models.Kabupaten{})
	config.InitDB().AutoMigrate(&models.Province{})
	config.InitDB().AutoMigrate(&models.Kendaraan{})
	config.InitDB().AutoMigrate(&models.Leasing{})
	config.InitDB().AutoMigrate(&models.Cabang{})
}
