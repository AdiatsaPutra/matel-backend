package repository

import (
	config "motor/configs"
	"motor/exceptions"
	"motor/models"

	"github.com/gin-gonic/gin"
)

func GetProvince(c *gin.Context) ([]models.Province, error) {
	var province []models.Province
	result := config.InitDB().Find(&province)

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return province, result.Error
	}

	return province, nil

}

func GetKabupaten(c *gin.Context, provinceID uint) ([]models.Kabupaten, error) {
	var kabupaten []models.Kabupaten
	result := config.InitDB().Where(&models.Kabupaten{ProvinceID: provinceID}).Find(&kabupaten)

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return kabupaten, result.Error
	}

	return kabupaten, nil

}

func GetKecamatan(c *gin.Context, kecamatanID uint) ([]models.Kecamatan, error) {
	var kecamatan []models.Kecamatan
	result := config.InitDB().Where(&models.Kecamatan{ProvinceID: kecamatanID}).Find(&kecamatan)

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return kecamatan, result.Error
	}

	return kecamatan, nil

}
