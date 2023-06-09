package controllers

import (
	"matel/exceptions"
	"matel/payloads"
	"matel/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProvince(c *gin.Context) {
	province, err := repository.GetProvince(c)

	if err != nil {
		exceptions.AppException(c, "Province not found")
		return
	}

	payloads.HandleSuccess(c, province, "Province found", http.StatusOK)
}

func GetKabupaten(c *gin.Context) {

	type kabupatenParam struct {
		KabupatenID uint
	}

	KabupatenID, _ := strconv.ParseUint(c.Param("province-id"), 10, 64)

	uintKabupatenID := uint(KabupatenID)

	var kabupaten = kabupatenParam{KabupatenID: uintKabupatenID}

	kabupatenResult, err := repository.GetKabupaten(c, kabupaten.KabupatenID)

	if err != nil {
		exceptions.AppException(c, "Kabupaten not found")
		return
	}

	payloads.HandleSuccess(c, kabupatenResult, "Kabupaten found", http.StatusOK)
}

func GetKecamatan(c *gin.Context) {

	type KecamatanParam struct {
		KabupatenID uint
	}

	KecamatanID, _ := strconv.ParseUint(c.Param("kabupaten-id"), 10, 64)

	uintKabupatenID := uint(KecamatanID)

	var kecamatan = KecamatanParam{KabupatenID: uintKabupatenID}

	kecamatanResult, err := repository.GetKecamatan(c, kecamatan.KabupatenID)

	if err != nil {
		exceptions.AppException(c, "Cant create member")
		return
	}

	payloads.HandleSuccess(c, kecamatanResult, "Kecamatan found", http.StatusOK)
}
