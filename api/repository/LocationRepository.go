package repository

import (
	config "matel/configs"
	"matel/models"

	"github.com/gin-gonic/gin"
)

// CreateProvince creates a new Province.
func CreateProvince(c *gin.Context, province models.Province) (models.Province, error) {
	result := config.InitDB().Create(&province)
	if result.Error != nil {
		return province, result.Error
	}
	return province, nil
}

// UpdateProvince updates an existing Province.
func UpdateProvince(c *gin.Context, id uint, province models.Province) (models.Province, error) {
	var updatedProvince models.Province
	result := config.InitDB().Model(&updatedProvince).Where("id = ?", id).Updates(province)
	if result.Error != nil {
		return updatedProvince, result.Error
	}
	return updatedProvince, nil
}

// DeleteProvince deletes a Province by its ID.
func DeleteProvince(c *gin.Context, id uint) error {
	result := config.InitDB().Delete(&models.Province{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetKabupatenByID retrieves a single Kabupaten by its ID.
func GetKabupatenByID(c *gin.Context, id uint) (models.Kabupaten, error) {
	var kabupaten models.Kabupaten
	result := config.InitDB().First(&kabupaten, id)
	if result.Error != nil {
		return kabupaten, result.Error
	}
	return kabupaten, nil
}

// CreateKabupaten creates a new Kabupaten.
func CreateKabupaten(c *gin.Context, kabupaten models.Kabupaten) (models.Kabupaten, error) {
	result := config.InitDB().Create(&kabupaten)
	if result.Error != nil {
		return kabupaten, result.Error
	}
	return kabupaten, nil
}

// UpdateKabupaten updates an existing Kabupaten.
func UpdateKabupaten(c *gin.Context, id uint, kabupaten models.Kabupaten) (models.Kabupaten, error) {
	var updatedKabupaten models.Kabupaten
	result := config.InitDB().Model(&updatedKabupaten).Where("id = ?", id).Updates(kabupaten)
	if result.Error != nil {
		return updatedKabupaten, result.Error
	}
	return updatedKabupaten, nil
}

// DeleteKabupaten deletes a Kabupaten by its ID.
func DeleteKabupaten(c *gin.Context, id uint) error {
	result := config.InitDB().Delete(&models.Kabupaten{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetKecamatanByID retrieves a single Kecamatan by its ID.
func GetKecamatanByID(c *gin.Context, id uint) (models.Kecamatan, error) {
	var kecamatan models.Kecamatan
	result := config.InitDB().First(&kecamatan, id)
	if result.Error != nil {
		return kecamatan, result.Error
	}
	return kecamatan, nil
}

// CreateKecamatan creates a new Kecamatan.
func CreateKecamatan(c *gin.Context, kecamatan models.Kecamatan) (models.Kecamatan, error) {
	result := config.InitDB().Create(&kecamatan)
	if result.Error != nil {
		return kecamatan, result.Error
	}
	return kecamatan, nil
}

// UpdateKecamatan updates an existing Kecamatan.
func UpdateKecamatan(c *gin.Context, id uint, kecamatan models.Kecamatan) (models.Kecamatan, error) {
	var updatedKecamatan models.Kecamatan
	result := config.InitDB().Model(&updatedKecamatan).Where("id = ?", id).Updates(kecamatan)
	if result.Error != nil {
		return updatedKecamatan, result.Error
	}
	return updatedKecamatan, nil
}

// DeleteKecamatan deletes a Kecamatan by its ID.
func DeleteKecamatan(c *gin.Context, id uint) error {
	result := config.InitDB().Delete(&models.Kecamatan{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAllProvinces retrieves all Provinces.
func GetAllProvinces(c *gin.Context) ([]models.Province, error) {
	var provinces []models.Province
	result := config.InitDB().Find(&provinces)

	if result.Error != nil {
		return provinces, result.Error
	}

	return provinces, nil
}

// GetAllKabupaten retrieves all Kabupatens or filtered by ProvinceID.
func GetAllKabupaten(c *gin.Context, provinceID uint) ([]models.Kabupaten, error) {
	var kabupaten []models.Kabupaten
	if provinceID == 0 {
		result := config.InitDB().Find(&kabupaten)

		if result.Error != nil {
			return kabupaten, result.Error
		}

		return kabupaten, nil
	} else {
		result := config.InitDB().Where("province_id = ?", provinceID).Find(&kabupaten)

		if result.Error != nil {
			return kabupaten, result.Error
		}

		return kabupaten, nil

	}
}

// GetAllKecamatan retrieves all Kecamatans or filtered by KabupatenID.
func GetAllKecamatan(c *gin.Context, kabupatenID uint) ([]models.Kecamatan, error) {
	var kecamatan []models.Kecamatan
	if kabupatenID == 0 {
		result := config.InitDB().Find(&kecamatan)
		if result.Error != nil {
			return kecamatan, result.Error
		}

		return kecamatan, nil
	} else {
		result := config.InitDB().Where("kabupaten_id = ?", kabupatenID).Find(&kecamatan)
		if result.Error != nil {
			return kecamatan, result.Error
		}

		return kecamatan, nil
	}
}
