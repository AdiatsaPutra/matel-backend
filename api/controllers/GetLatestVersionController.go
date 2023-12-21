package controllers

import (
	"matel/exceptions"
	"matel/payloads"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetLatestVersion(c *gin.Context) {
	var body struct {
		Version string `json:"version" validate:"required"`
	}

	c.ShouldBindJSON(&body)

	if body.Version == "" {
		exceptions.AppException(c, "Please add version")
		return
	}

	data := make(map[string]interface{})

	latestVersion := "1.8.0"

	isMajorUpdate := strings.Split(body.Version, ".")[0] < strings.Split(latestVersion, ".")[0]
	isMinorUpdate := strings.Split(body.Version, ".")[1] < strings.Split(latestVersion, ".")[1]

	data["is_latest"] = body.Version == latestVersion
	data["is_force_update"] = isMajorUpdate || isMinorUpdate
	data["latest_version"] = latestVersion

	payloads.HandleSuccess(c, data, "Data found", http.StatusOK)

}
