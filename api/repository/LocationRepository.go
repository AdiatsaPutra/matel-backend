package repository

import (
	config "matel/configs"
	"matel/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetProvince(c *gin.Context) ([]models.Province, error) {
	var province []models.Province
	result := config.InitDB().Find(&province)

	if result.Error != nil {
		return province, result.Error
	}

	return province, nil

}

func GetKabupaten(c *gin.Context, provinceID uint) ([]models.Kabupaten, error) {
	var kabupaten []models.Kabupaten
	if provinceID == 0 {
		result := config.InitDB().Find(&kabupaten)
		logrus.Info(kabupaten)

		if result.Error != nil {
			return kabupaten, result.Error
		}

		return kabupaten, nil
	} else {
		result := config.InitDB().Where(&models.Kabupaten{ProvinceID: provinceID}).Find(&kabupaten)

		if result.Error != nil {
			return kabupaten, result.Error
		}

		return kabupaten, nil

	}

}

func GetKecamatan(c *gin.Context, KabupatenID uint) ([]models.Kecamatan, error) {
	var kecamatan []models.Kecamatan
	if KabupatenID == 0 {
		result := config.InitDB().Find(&kecamatan)
		if result.Error != nil {
			return kecamatan, result.Error
		}

		return kecamatan, nil
	} else {
		result := config.InitDB().Where(&models.Kecamatan{KabupatenID: KabupatenID}).Find(&kecamatan)
		if result.Error != nil {
			return kecamatan, result.Error
		}

		return kecamatan, nil
	}

}
