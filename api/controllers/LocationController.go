package controllers

import (
	"matel/exceptions"
	"matel/models"
	"matel/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProvince creates a new Province.
func CreateProvince(c *gin.Context) {
	var province models.Province
	if err := c.ShouldBindJSON(&province); err != nil {
		exceptions.AppException(c, "Invalid input data")
		return
	}

	createdProvince, err := repository.CreateProvince(c, province)
	if err != nil {
		exceptions.AppException(c, "Failed to create province")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdProvince, "message": "Province created successfully"})
}

// GetProvince retrieves a single Province by its ID.
func GetProvince(c *gin.Context) {

	province, err := repository.GetAllProvinces(c)
	if err != nil {
		exceptions.AppException(c, "Province not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": province, "message": "Province found"})
}

// UpdateProvince updates an existing Province.
func UpdateProvince(c *gin.Context) {
	provinceID, _ := strconv.ParseUint(c.Param("province-id"), 10, 64)

	var province models.Province
	if err := c.ShouldBindJSON(&province); err != nil {
		exceptions.AppException(c, "Invalid input data")
		return
	}

	updatedProvince, err := repository.UpdateProvince(c, uint(provinceID), province)
	if err != nil {
		exceptions.AppException(c, "Failed to update province")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedProvince, "message": "Province updated successfully"})
}

// DeleteProvince deletes a Province by its ID.
func DeleteProvince(c *gin.Context) {
	provinceID, _ := strconv.ParseUint(c.Param("province-id"), 10, 64)

	err := repository.DeleteProvince(c, uint(provinceID))
	if err != nil {
		exceptions.AppException(c, "Failed to delete province")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Province deleted successfully"})
}

// GetKabupaten retrieves a single Kabupaten by its ID.
func GetKabupaten(c *gin.Context) {
	kabupatenID, _ := strconv.ParseUint(c.Param("kabupaten-id"), 10, 64)

	kabupaten, err := repository.GetAllKabupaten(c, uint(kabupatenID))
	if err != nil {
		exceptions.AppException(c, "Kabupaten not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kabupaten, "message": "Kabupaten found"})
}

// UpdateKabupaten updates an existing Kabupaten.
func UpdateKabupaten(c *gin.Context) {
	kabupatenID, _ := strconv.ParseUint(c.Param("kabupaten-id"), 10, 64)

	var kabupaten models.Kabupaten
	if err := c.ShouldBindJSON(&kabupaten); err != nil {
		exceptions.AppException(c, "Invalid input data")
		return
	}

	updatedKabupaten, err := repository.UpdateKabupaten(c, uint(kabupatenID), kabupaten)
	if err != nil {
		exceptions.AppException(c, "Failed to update Kabupaten")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedKabupaten, "message": "Kabupaten updated successfully"})
}

// DeleteKabupaten deletes a Kabupaten by its ID.
func DeleteKabupaten(c *gin.Context) {
	kabupatenID, _ := strconv.ParseUint(c.Param("kabupaten-id"), 10, 64)

	err := repository.DeleteKabupaten(c, uint(kabupatenID))
	if err != nil {
		exceptions.AppException(c, "Failed to delete Kabupaten")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kabupaten deleted successfully"})
}

// GetKecamatan retrieves a single Kecamatan by its ID.
func GetKecamatan(c *gin.Context) {
	kecamatanID, _ := strconv.ParseUint(c.Param("kecamatan-id"), 10, 64)

	kecamatan, err := repository.GetAllKecamatan(c, uint(kecamatanID))
	if err != nil {
		exceptions.AppException(c, "Kecamatan not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kecamatan, "message": "Kecamatan found"})
}

// UpdateKecamatan updates an existing Kecamatan.
func UpdateKecamatan(c *gin.Context) {
	kecamatanID, _ := strconv.ParseUint(c.Param("kecamatan-id"), 10, 64)

	var kecamatan models.Kecamatan
	if err := c.ShouldBindJSON(&kecamatan); err != nil {
		exceptions.AppException(c, "Invalid input data")
		return
	}

	updatedKecamatan, err := repository.UpdateKecamatan(c, uint(kecamatanID), kecamatan)
	if err != nil {
		exceptions.AppException(c, "Failed to update Kecamatan")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedKecamatan, "message": "Kecamatan updated successfully"})
}

// DeleteKecamatan deletes a Kecamatan by its ID.
func DeleteKecamatan(c *gin.Context) {
	kecamatanID, _ := strconv.ParseUint(c.Param("kecamatan-id"), 10, 64)

	err := repository.DeleteKecamatan(c, uint(kecamatanID))
	if err != nil {
		exceptions.AppException(c, "Failed to delete Kecamatan")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kecamatan deleted successfully"})
}
