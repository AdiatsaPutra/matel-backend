package controllers

import (
	"matel/exceptions"
	"matel/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetProvince retrieves a single Province by its ID.
func GetProvince(c *gin.Context) {

	province, err := repository.GetAllProvinces(c)
	if err != nil {
		exceptions.AppException(c, "Province not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": province, "message": "Province found"})
}

// GetKabupaten retrieves a single Kabupaten by its ID.
func GetKabupaten(c *gin.Context) {
	kabupatenID, _ := strconv.ParseUint(c.Param("province-id"), 10, 64)

	logrus.Info(kabupatenID)

	kabupaten, err := repository.GetAllKabupaten(c, uint(kabupatenID))
	if err != nil {
		exceptions.AppException(c, "Kabupaten not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kabupaten, "message": "Kabupaten found"})
}

// GetKecamatan retrieves a single Kecamatan by its ID.
func GetKecamatan(c *gin.Context) {
	kecamatanID, _ := strconv.ParseUint(c.Param("kabupaten-id"), 10, 64)

	kecamatan, err := repository.GetAllKecamatan(c, uint(kecamatanID))
	if err != nil {
		exceptions.AppException(c, "Kecamatan not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kecamatan, "message": "Kecamatan found"})
}
